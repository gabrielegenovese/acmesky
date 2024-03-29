package flightMatcher

import (
	entities "acmesky/entities"
	flightsRepo "acmesky/repository/flights"
	travelPreferenceRepo "acmesky/repository/travel_preference"
	zbSingleton "acmesky/workers"
	zeebeUtils "acmesky/workers/utils/zeebe_utils"
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
			JobType("findSolutionsByTravelPreference").
			Handler(HandleFindSolutionsByTravelPreference).
			Open(),
		client.
			NewJobWorker().
			JobType("prepareOffersForCustomer").
			Handler(HandlePrepareOfferForCustomer).
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
		if zeebeUtils.Handle_BP_fail_allow_retry(client, job, err, 10*time.Second) {
			return
		}

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

type fetchFlightsParameters struct {
	Pref      entities.CustomerFlightSubscription `json:"pref,omitempty"`
	CompanyID int64                               `json:"flight_company_id,omitempty"`
}

func HandleFetchFlightsByTravelPreference(client worker.JobClient, job zeebeEntities.Job) {
	var dbErr error
	var fetchParams fetchFlightsParameters
	err := job.GetVariablesAs(&fetchParams)

	if err != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelFn()
		_, _ = client.NewFailJobCommand().JobKey(job.Key).Retries(0).ErrorMessage(err.Error()).Send(ctx)
		return
	}

	if fetchParams.CompanyID == 0 {
		fetchParams.CompanyID = 1
	}

	fmt.Printf("Fething preference using flight company ID %d\n", fetchParams.CompanyID)
	flights, fetchErr := FetchFlightsByCompanyID(fetchParams.Pref.CustomerFlightSubscriptionRequest, fetchParams.CompanyID)

	if fetchErr == nil {
		fmt.Printf("Storing %d fetched flights\n", len(flights))
		if len(flights) > 0 {
			dbErr = flightsRepo.AddFlights(flights)
		}
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	if dbErr != nil || fetchErr != nil {
		if zeebeUtils.Handle_BP_fail_allow_retry(client, job, dbErr, 5*time.Second) {
			return
		} else if zeebeUtils.Handle_BP_fail_allow_retry(client, job, fetchErr, 10*time.Second) {
			return
		} else if fetchErr != nil {
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

type storeFlightsParameters struct {
	Flights   []entities.Flight `json:"flights"`
	CompanyID int64             `json:"flight_company_id,omitempty"`
}

func HandleStoreFlights(client worker.JobClient, job zeebeEntities.Job) {
	var dbErr error
	var storeParams storeFlightsParameters

	err := job.GetVariablesAs(&storeParams)

	if err != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		_, _ = client.NewFailJobCommand().JobKey(job.Key).Retries(0).ErrorMessage(err.Error()).Send(ctx)
		return
	}

	if len(storeParams.Flights) > 0 {
		dbErr = flightsRepo.AddFlights(storeParams.Flights)
	}

	if zeebeUtils.Handle_BP_fail_allow_retry(client, job, dbErr, 5*time.Second) {
		return
	} else if dbErr != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelFn()
		_, err := client.
			NewFailJobCommand().
			JobKey(job.Key).
			Retries(0).
			ErrorMessage(dbErr.Error()).
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
			"flights": storeParams.Flights,
		})

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on creating complete job with key [%d]: [%s]", job.Key, err))
		return
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()
	_, err = command.Send(ctx)

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on complete job with key [%d]: [%s]", job.Key, err))
	}

}

type findSolutionsParameters struct {
	Pref entities.CustomerFlightSubscription `json:"pref,omitempty"`
}

func HandleFindSolutionsByTravelPreference(client worker.JobClient, job zeebeEntities.Job) {

	var findParams findSolutionsParameters
	err := job.GetVariablesAs(&findParams)

	if err != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		_, _ = client.NewFailJobCommand().JobKey(job.Key).Retries(0).ErrorMessage(err.Error()).Send(ctx)
		return
	}
	var dbErr error
	var solutions []entities.Solution
	fmt.Printf("getting pref\n")

	fmt.Printf("Searching solutions\n")
	solutions, dbErr = flightsRepo.GetSolutionsFromPreference(findParams.Pref.CustomerFlightSubscriptionRequest)
	fmt.Printf("Got %d solutions\n", len(solutions))

	if zeebeUtils.Handle_BP_fail_allow_retry(client, job, dbErr, 5*time.Second) {
		return
	} else if dbErr != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelFn()

		_, err := client.
			NewFailJobCommand().
			JobKey(job.Key).
			Retries(0).
			RetryBackoff(5 * time.Second).
			ErrorMessage(dbErr.Error()).
			Send(ctx)

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, err))
		} else {
			log.Println(dbErr)
		}
		return
	}

	fmt.Printf("Found %d solutions\n", len(solutions))
	command, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"solutions": solutions,
		})

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on creating complete job with key [%d]: [%s]", job.Key, err))
		return
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()
	_, err = command.Send(ctx)

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on complete job with key [%d]: [%s]", job.Key, err))
	}
}

type prepareOffersParameters struct {
	Pref     entities.CustomerFlightSubscription `json:"pref,omitempty"`
	Solution entities.Solution                   `json:"solution,omitempty"`
}

func HandlePrepareOfferForCustomer(client worker.JobClient, job zeebeEntities.Job) {
	var offer entities.ReservedOffer
	var prepareOffersParams prepareOffersParameters
	fmt.Printf("HandlePrepareOfferForCustomer")

	err := job.GetVariablesAs(&prepareOffersParams)

	if err != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		_, _ = client.NewFailJobCommand().JobKey(job.Key).Retries(0).ErrorMessage(err.Error()).Send(ctx)
		return
	}

	fmt.Printf("Preparing offer for TravelPreference %d\n", prepareOffersParams.Pref.TravelPreferenceID)
	flights, dbErr := flightsRepo.GetFlight(
		[]string{
			prepareOffersParams.Solution.DepartFlight.FlightID,
			prepareOffersParams.Solution.ReturnFlight.FlightID,
		},
		[]int64{
			prepareOffersParams.Solution.DepartFlight.FlightCompanyID,
			prepareOffersParams.Solution.ReturnFlight.FlightCompanyID,
		},
	)
	if zeebeUtils.Handle_BP_fail_allow_retry(client, job, dbErr, 5*time.Second) {
		return
	} else if dbErr == nil {
		var offerCode int64
		fmt.Printf("Preparing offer ... \n")
		var totalPrice float32 = 0
		for _, f := range flights {
			totalPrice += float32(f.FlightPrice)
		}
		offerCode, dbErr = travelPreferenceRepo.AddReservedOffer(prepareOffersParams.Pref.TravelPreferenceID, totalPrice, flights)
		if dbErr == nil {
			offer, dbErr = travelPreferenceRepo.GetRservedOffer(offerCode)
		}
	}

	if zeebeUtils.Handle_BP_fail_allow_retry(client, job, dbErr, 5*time.Second) {
		return
	} else if dbErr != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelFn()

		_, err := client.
			NewFailJobCommand().
			JobKey(job.Key).
			Retries(0).
			RetryBackoff(5 * time.Second).
			ErrorMessage(dbErr.Error()).
			Send(ctx)

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, err))
		} else {
			log.Println(dbErr)
		}
		return
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	command, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"offer": offer,
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
