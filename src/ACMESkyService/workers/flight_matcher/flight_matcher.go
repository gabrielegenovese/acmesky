package flightMatcher

import (
	entities "acmesky/entities"
	travelPreferenceRepo "acmesky/repository/travel_preference"
	zbSingleton "acmesky/workers"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	zeebeEntities "github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func RegisterWorkers() []worker.JobWorker {
	client := *zbSingleton.GetInstance()
	workers := []worker.JobWorker{
		client.
			NewJobWorker().
			JobType("loadTravelPreferencesAndDispatchBP").
			Handler(HandleLoadTravelPreferencesAndDispatchBP).
			Open(),
		client.
			NewJobWorker().
			JobType("fetchFlightsByTravelPreference").
			Handler(HandleFetchFlightsByTravelPreference).
			Open(),
		client.
			NewJobWorker().
			JobType("searchMatchByTravelPreference").
			Handler(HandleSearchMatchByTravelPreference).
			Open(),
	}
	return workers
}

func HandleLoadTravelPreferencesAndDispatchBP(client worker.JobClient, job zeebeEntities.Job) {
	zbClient := *zbSingleton.GetInstance()

	prefs, err := travelPreferenceRepo.GetAllCustomerFlightPreferencesNotOutdated()

	if err != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelFn()
		_, err := client.
			NewFailJobCommand().
			JobKey(job.Key).
			Retries(job.GetRetries() - 1).
			RetryBackoff(5 * time.Second).
			ErrorMessage(err.Error()).
			Send(ctx)

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, err))
		} else {
			log.Println(err)
		}
		return
	}

	for _, _pref := range prefs {
		go func(zbClient *zbc.Client, client worker.JobClient, job zeebeEntities.Job, pref entities.CustomerFlightSubscription) {
			command, err := (*zbClient).NewCreateInstanceCommand().
				BPMNProcessId("aaaa").
				LatestVersion().
				VariablesFromObject(pref)

			if err != nil {
				log.Println(fmt.Errorf("[BPMNERROR] error on creating BPMN process from job key [%d]: [%s]", job.Key, err))
				return
			}

			ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFn()

			_, err = command.Send(ctx)
			if err != nil {
				log.Println(fmt.Errorf("[BPMNERROR] error on send created BPMN process from job key [%d]: [%s]", job.Key, err))
			}
		}(&zbClient, client, job, _pref)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	_, err = client.NewCompleteJobCommand().
		JobKey(job.Key).
		Send(ctx)

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on complete job with key [%d]: [%s]", job.Key, err))
	}
}

func HandleFetchFlightsByTravelPreference(client worker.JobClient, job zeebeEntities.Job) {
	var dbErr error

	vars, err := job.GetVariablesAsMap()
	if err != nil {
		return
	}

	pref := entities.CustomerFlightSubscriptionRequestFromMap(vars)

	flights, fetchErr := fetchFlightsForPreference(pref)

	if fetchErr == nil {
		// TODO save flights to DB
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	if dbErr != nil || fetchErr != nil {
		if fetchErr != nil {
			if strings.Contains(fetchErr.Error(), "HTTP_ERROR") {
				command, err := client.
					NewThrowErrorCommand().
					JobKey(job.Key).
					ErrorCode("HTTP_ERROR").
					ErrorMessage(fetchErr.Error()).
					VariablesFromMap(map[string]interface{}{
						"errorCode": "HTTP_ERROR",
						"errorMsg":  fetchErr.Error(),
					})

				_, err = command.Send(ctx)

				if err != nil {
					log.Println(fmt.Errorf("[BPMNERROR] error on throwing error on job with key [%d]: [%s]", job.Key, err))
				} else {
					log.Println(fetchErr)
				}
				return
			}
		}

		// fail as unhandled if we are here
		if fetchErr != nil {
			err = fetchErr
		} else {
			err = dbErr
		}
		_, err := client.
			NewFailJobCommand().
			JobKey(job.Key).
			Retries(job.GetRetries() - 1).
			RetryBackoff(30 * time.Second).
			ErrorMessage(err.Error()).
			Send(ctx)

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, err))
		} else {
			log.Println(fetchErr)
		}
		return
	}

	command, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"flights": flights,
		})

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on creating complete job with key [%d]: [%s]", job.Key, err))
		return
	}
	_, err = command.Send(ctx)

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on complete job with key [%d]: [%s]", job.Key, err))
	}
}

func fetchFlightsForPreference(pref entities.CustomerFlightSubscriptionRequest) ([]entities.Flight, error) {

	var FLIGHT_COMPANY_ADDRESS string = os.Getenv("FLIGHT_COMPANY_ADDRESS")
	var flights []entities.Flight = []entities.Flight{}

	departDateRangeEnd, parseErr := time.Parse(time.DateTime, pref.DateStartISO8601)
	if parseErr != nil {
		return flights, fmt.Errorf("PARSE_ERROR:" + parseErr.Error())
	}
	departDateRangeEnd = departDateRangeEnd.Add(24 * time.Hour).Truncate(24 * time.Hour)

	res, err := http.Get(
		"http://" + FLIGHT_COMPANY_ADDRESS + "/flights?" + url.Values{
			"origin_airport":       {pref.AirportOriginID},
			"dest_airport":         {pref.AirportDestinationID},
			"passengers_count":     {strconv.FormatUint(uint64(pref.SeatsCount), 10)},
			"start_range_datetime": {pref.DateStartISO8601},
			"end_range_datetime":   {departDateRangeEnd.UTC().Format(time.DateTime)},
		}.Encode())
	defer res.Body.Close()

	if err != nil {
		return flights, fmt.Errorf("CONNECTION_ERROR:" + err.Error())
	}

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return flights, fmt.Errorf("HTTP_ERROR:" + res.Status)
	}

	decodeErr := json.NewDecoder(res.Body).Decode(&flights)

	if err != nil {
		return flights, fmt.Errorf("PARSE_ERROR:" + decodeErr.Error())
	}

	return flights, nil
}

func HandleSearchMatchByTravelPreference(client worker.JobClient, job zeebeEntities.Job) {

}
