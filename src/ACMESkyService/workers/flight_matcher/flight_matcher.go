package flightMatcher

import (
	entities "acmesky/dao/entities"
	airportsDAO "acmesky/dao/impl/airports"
	flightsDAO "acmesky/dao/impl/flights"
	travelPreferenceDAO "acmesky/dao/impl/travel_preference"
	"acmesky/services/flights"
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
		client.
			NewJobWorker().
			JobType("notifyReservedOffer").
			Handler(HandleNotifyReservedOffer).
			Open(),
	}
	return workers
}

func HandleLoadTravelPreferences(client worker.JobClient, job zeebeEntities.Job) {

	fmt.Println("Getting customers' Travel Preferences without offers")
	prefs, err := travelPreferenceDAO.GetAllCustomerFlightPreferencesNotOutdated()

	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	if err != nil {
		if zeebeUtils.Handle_BP_fail_allow_retry(client, job, err, 10*time.Second) {
			return
		}

		// cant get preference -> fails job

		_, errCmd := client.
			NewFailJobCommand().
			JobKey(job.Key).
			Retries(job.GetRetries() - 1).
			RetryBackoff(5 * time.Second).
			ErrorMessage(err.Error()).
			Send(ctx)

		if errCmd != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, errCmd))
		} else {
			log.Println(err)
		}
		return
	}

	fmt.Printf("Got %v customers' Travel Preferences\n", len(prefs))
	command, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"prefs": prefs,
		})
	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on complete job with key [%d]: [%s]", job.Key, err))
		return
	}
	_, err = command.Send(ctx)
	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] error on complete job with key [%d]: [%s]", job.Key, err))
		return
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
	flights, fetchErr := flights.FetchFlightsByCompanyID(fetchParams.Pref.CustomerFlightSubscriptionRequest, fetchParams.CompanyID)

	if fetchErr == nil {
		fmt.Printf("Storing %d fetched flights\n", len(flights))
		if len(flights) > 0 {
			dbErr = flightsDAO.AddFlights(flights)
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
		dbErr = flightsDAO.AddFlights(storeParams.Flights)
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

	fmt.Printf("Searching solutions for preference %v\n", findParams.Pref.TravelPreferenceID)
	solutions, dbErr = flightsDAO.GetSolutionsFromPreference(findParams.Pref.CustomerFlightSubscriptionRequest)
	fmt.Printf("Found %d solutions for preference %v\n", len(solutions), findParams.Pref.TravelPreferenceID)

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

	fmt.Printf("Found %d solutions for preference %v\n", len(solutions), findParams.Pref.TravelPreferenceID)
	/*for _, s := range solutions {
		fmt.Printf("Requested from %v to %v betwween %s and %s with budget < %v\n"+
			"Find Solution: \n"+
			"\tDepart: (%v,%v) from %v to %v betwween %s and %s with price %v\n"+
			"\tReturn: (%v,%v) from %v to %v betwween %s and %s with price %v\n",
			findParams.Pref.AirportOriginID, findParams.Pref.AirportDestinationID, findParams.Pref.DateStartISO8601, findParams.Pref.DateEndISO8601, findParams.Pref.Budget,
			s.DepartFlight.FlightCompanyID, s.DepartFlight.FlightID, s.DepartFlight.AirportOriginID, s.DepartFlight.AirportDestinationID, s.DepartFlight.DepartDatetime, s.DepartFlight.ArrivalDatetime, s.DepartFlight.FlightPrice,
			s.ReturnFlight.FlightCompanyID, s.ReturnFlight.FlightID, s.ReturnFlight.AirportOriginID, s.ReturnFlight.AirportDestinationID, s.ReturnFlight.DepartDatetime, s.ReturnFlight.ArrivalDatetime, s.ReturnFlight.FlightPrice,
		)
	}*/
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
	var airports []entities.Airport
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
	flights, dbErr := flightsDAO.GetFlight(
		[]string{
			prepareOffersParams.Solution.DepartFlight.FlightID,
			prepareOffersParams.Solution.ReturnFlight.FlightID,
		},
		[]int64{
			prepareOffersParams.Solution.DepartFlight.FlightCompanyID,
			prepareOffersParams.Solution.ReturnFlight.FlightCompanyID,
		},
	)
	prepareOffersParams.Solution.DepartFlight = flights[0]
	prepareOffersParams.Solution.ReturnFlight = flights[1]

	if dbErr == nil {
		airports, dbErr = airportsDAO.GetAirportsById([]string{
			flights[0].AirportOriginID,
			flights[1].AirportOriginID,
		})
	}

	if zeebeUtils.Handle_BP_fail_allow_retry(client, job, dbErr, 5*time.Second) {
		return
	} else if dbErr == nil {
		var offerCode int64
		fmt.Printf("Preparing offer for pref %v ...\n", prepareOffersParams.Pref.TravelPreferenceID)
		var totalPrice float32 = 0.
		for _, f := range flights {
			totalPrice += float32(f.FlightPrice) * float32(prepareOffersParams.Pref.SeatsCount)
		}
		offerCode, dbErr = travelPreferenceDAO.AddReservedOffer(prepareOffersParams.Pref.TravelPreferenceID, totalPrice, flights)
		if dbErr == nil {
			offer, dbErr = travelPreferenceDAO.GetReservedOffer(offerCode)
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

	fmt.Printf("Requested from %v to %v betwween %s and %s with budget < %v\n"+
		"Prepared offer (%v) for %v money since %v up to %v with the following flights:\n"+
		"\tDepart: (%v,%v) from %v to %v\n"+
		"\tReturn: (%v,%v) from %v to %v\n",
		prepareOffersParams.Pref.AirportOriginID, prepareOffersParams.Pref.AirportDestinationID, prepareOffersParams.Pref.DateStartISO8601, prepareOffersParams.Pref.DateEndISO8601, prepareOffersParams.Pref.Budget,
		offer.OfferCode, offer.TotalPrice, offer.StartReservationDatetime, offer.EndReservationDatetime,
		flights[0].FlightCompanyID, flights[0].FlightID, flights[0].AirportOriginID, flights[0].AirportDestinationID,
		flights[1].FlightCompanyID, flights[1].FlightID, flights[1].AirportOriginID, flights[1].AirportDestinationID,
	)

	output := offerData{
		Offer:         offer,
		Solution:      prepareOffersParams.Solution,
		DepartAirport: airports[0],
		ReturnAirport: airports[1],
	}
	command, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"offerData": output,
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

type offerData struct {
	Offer         entities.ReservedOffer `json:"offer,omitempty"`
	Solution      entities.Solution      `json:"solution,omitempty"`
	DepartAirport entities.Airport       `json:"departOriginAirport,omitempty"`
	ReturnAirport entities.Airport       `json:"returnOriginAirport,omitempty"`
}

type notifyOfferParameters struct {
	OfferData offerData                           `json:"offerData,omitempty"`
	Pref      entities.CustomerFlightSubscription `json:"pref,omitempty"`
}

func HandleNotifyReservedOffer(client worker.JobClient, job zeebeEntities.Job) {
	var input notifyOfferParameters
	fmt.Printf("HandleNotifyReservedOffer\n")

	err := job.GetVariablesAs(&input)

	if err != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		_, _ = client.NewFailJobCommand().JobKey(job.Key).Retries(0).ErrorMessage(err.Error()).Send(ctx)
		return
	}

	departAirport := input.OfferData.DepartAirport
	returnAirport := input.OfferData.ReturnAirport
	departFlight := input.OfferData.Solution.DepartFlight
	returnFlight := input.OfferData.Solution.ReturnFlight
	offerEndDatetime, _ := time.Parse(time.DateTime, input.OfferData.Offer.EndReservationDatetime)

	title := fmt.Sprintf("New ACMESKY travel offer from %s to %s until %s",
		departAirport.City, returnAirport.City, input.OfferData.Offer.EndReservationDatetime,
	)
	body := fmt.Sprintf("ACMESKY found a flight travel offer with return flight for you (%v seats) from %v to %v between %s and %s within your budget of %v %v\n"+
		"We offer the following flights at the price of %v %v :\n"+
		"- Depart from %v at %v will arrive at %v;\n"+
		"- Return from %v at %v will return at %v.\n"+
		"Use the code %v until %v on %v to purchase this offer on our portal at this reserved price !",
		input.Pref.SeatsCount, departAirport.City, returnAirport.City, input.Pref.DateStartISO8601, input.Pref.DateEndISO8601, input.Pref.Budget, "€",
		input.OfferData.Offer.TotalPrice, "€",
		departAirport.Name, departFlight.DepartDatetime, departFlight.ArrivalDatetime,
		returnAirport.Name, returnFlight.DepartDatetime, returnFlight.ArrivalDatetime,
		input.OfferData.Offer.OfferCode, offerEndDatetime.Format(time.TimeOnly), offerEndDatetime.Format(time.DateOnly),
	)

	messageRes, errSends := NotifyCustomer(input.Pref.CustomerFlightSubscriptionRequest, Notification{
		subject: title,
		content: body,
	})
	var errSend error
	if len(messageRes) < 1 {
		errSend = fmt.Errorf("[NotifyCustomer] no notification sent: %v", errSends)
	} else if len(errSends) > 0 {
		errSend = errSends[0]
	}

	if zeebeUtils.Handle_BP_fail_allow_retry(client, job, errSend, 5*time.Second) {
		return
	} else if errSend != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFn()

		_, err := client.
			NewFailJobCommand().
			JobKey(job.Key).
			Retries(0).
			RetryBackoff(5 * time.Second).
			ErrorMessage(errSend.Error()).
			Send(ctx)

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, err))
		} else {
			log.Println(errSend)
		}
		return
	}
	message := messageRes[0]
	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	command, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"notification_id": message.Id,
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
