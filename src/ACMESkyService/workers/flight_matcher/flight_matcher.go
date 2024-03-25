package flightMatcher

import (
	entities "acmesky/entities"
	flightsRepo "acmesky/repository/flights"
	travelPreferenceRepo "acmesky/repository/travel_preference"
	zbSingleton "acmesky/workers"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	zeebeEntities "github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

func RegisterWorkers() []worker.JobWorker {
	client := *zbSingleton.GetInstance()
	workers := []worker.JobWorker{
		client.
			NewJobWorker().
			JobType("loadTravelPreferences").
			Handler(HandleLoadTravelPreferences).
			Open(),
		client.
			NewJobWorker().
			JobType("fetchAndStoreFlightsByTravelPreference").
			Handler(HandleFetchFlightsByTravelPreference).
			Open(),
		client.
			NewJobWorker().
			JobType("storeFlights").
			Handler(HandleStoreFlights).
			Open(),
		client.
			NewJobWorker().
			JobType("searchMatchByTravelPreference").
			Handler(HandleSearchMatchByTravelPreference).
			Open(),
	}
	return workers
}

func HandleLoadTravelPreferences(client worker.JobClient, job zeebeEntities.Job) {

	fmt.Println("Getting customers' Travel Preferences")
	prefs, err := travelPreferenceRepo.GetAllCustomerFlightPreferencesNotOutdated()

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	if err != nil {
		// cant get preference -> fails job

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

	fmt.Println("Got customers' Travel Preferences")
	command, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"prefs": prefs,
		})
	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on complete job with key [%d]: [%s]", job.Key, err))
	}
	_, err = command.Send(ctx)
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
	var flightCompanyID int64

	_flightCompanyID, hasFlightCompanyID := vars["flight_company_id"]

	if hasFlightCompanyID {
		flightCompanyID = _flightCompanyID.(int64)
	} else {
		flightCompanyID = 1
	}

	pref := entities.CustomerFlightSubscriptionRequestFromMap(vars["pref"].(map[string]interface{}))
	fmt.Printf("Fething preference using flight company ID %d\n", flightCompanyID)
	flights, fetchErr := FetchFlightsByCompanyID(pref, flightCompanyID)

	if fetchErr == nil {
		fmt.Printf("Storing %d fetched flights\n", len(flights))
		if len(flights) > 0 {
			dbErr = flightsRepo.AddFlights(flights)
		}
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

				if err == nil {
					_, err = command.Send(ctx)
				}

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
		fmt.Printf("Fetch or Store error: %s\n", err.Error())
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

	fmt.Printf("Store successfull of %d items\n", len(flights))
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

func HandleStoreFlights(client worker.JobClient, job zeebeEntities.Job) {

	vars, err := job.GetVariablesAsMap()
	if err != nil {
		return
	}
	var dbErr error
	var flights []entities.Flight
	flightsRaw := vars["flights"].([]map[string]interface{})

	for i := 0; i < len(flightsRaw); i++ {
		flights = append(flights, entities.FlightFromMapFromMap(flightsRaw[i]))
	}

	if len(flights) > 0 {
		dbErr = flightsRepo.AddFlights(flights)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	if dbErr != nil {
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
			log.Println(dbErr)
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

func HandleSearchMatchByTravelPreference(client worker.JobClient, job zeebeEntities.Job) {

}
