package acmesky

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"workers/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/gin-gonic/gin"
)

func SaveInterestHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	// Wait for confirm
	<-util.InterestSaved[variables["correlationKey"].(string)]
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)

	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func UpdateFlightListHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	http.Get(os.Getenv("ACMESKY_SERVICE_API") + "/api/v1/updateFlights")

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)

	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func NewInterestSaved(c *gin.Context) {
	id := c.Param("id")
	close(util.InterestSaved[id])
}

func SendNewOffertHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	response, err := http.Get(os.Getenv("ACMESKY_SERVICE_API") + "/api/v1/offers")
	if err != nil {
		util.FailJob(client, job)
	}

	if response.StatusCode == 200 {
		offers := make([]util.Offer, 0)
		json.NewDecoder(response.Body).Decode(&offers)
		for _, offer := range offers {
			variables := make(map[string]interface{})
			variables["prontogramId"] = offer.TravelPreference.ProntogramID
			variables["offerCode"] = fmt.Sprintf("%d", offer.OfferCode)
			message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageNewOffert").CorrelationKey("correlation").VariablesFromMap(variables)
			_, err = message.Send(util.Ctx)
			if err != nil {
				util.FailJob(client, job)
				return
			}
		}
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func bookFlight(flightID string, prontogramID string, seats uint) (string, float64, error) {
	fmt.Println("Depart func: ", flightID)
	flightIDint, _ := strconv.ParseInt(flightID, 10, 64)
	bookingRequest := util.FlightBooking{
		FlightID: flightIDint,
		CustomerName: prontogramID,
		CustomerSurname: prontogramID,
		SeatsCount: int(seats),
	}
	data, _ := json.Marshal(bookingRequest)
	response, _ := http.Post(os.Getenv("FLIGHT_COMPANY_API") + "/bookings/", "application/json", bytes.NewReader(data))
	if response.StatusCode == 200 {
		var booking util.FlightBooking
		json.NewDecoder(response.Body).Decode(&booking)
		response, _ := http.Get(os.Getenv("FLIGHT_COMPANY_API") + "/flight/" + flightID)
		var flight util.Flight
		json.NewDecoder(response.Body).Decode(&flight)
		fmt.Println("Flight ", flight)
		totalPrice := flight.FlightPrice * float64(seats)
		return fmt.Sprintf("%d", booking.BookingID), totalPrice, nil
	} else {
		return "", 0, errors.New("Booking error")
	}
}

func BookFlightHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	response, err := http.Get(os.Getenv("ACMESKY_SERVICE_API") + "/api/v1/offer/" + variables["offerCode"].(string))
	if err != nil {
		util.FailJob(client, job)
	}

	if response.StatusCode == 200 {
		var offer util.Offer
		json.NewDecoder(response.Body).Decode(&offer)
		fmt.Println("Depart: ", offer.DepartFlight.FlightID)
		departBooking, departPrice, err := bookFlight(offer.DepartFlight.FlightID, offer.TravelPreference.ProntogramID, offer.TravelPreference.SeatsCount)
		if err != nil {
			util.BuyResults[variables["offerCode"].(string)] <- util.BuyResult{Success: false}
			util.FailJob(client, job)
			return
		}
		variables["departBooking"] = departBooking
		returnBooking, returnPrice, err := bookFlight(offer.ReturnFlight.FlightID, offer.TravelPreference.ProntogramID, offer.TravelPreference.SeatsCount)
		if err != nil {
			util.BuyResults[variables["offerCode"].(string)] <- util.BuyResult{Success: false}
			util.FailJob(client, job)
			return
		}
		variables["totalPrice"] = departPrice + returnPrice
		variables["flightBooked"] = true
		variables["returnBooking"] = returnBooking
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func SendNewPaymentHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageNewPayment").CorrelationKey(variables["offerCode"].(string)).VariablesFromMap(variables)
	_, err = message.Send(util.Ctx)
	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func UnbookFlightHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	httpClient := &http.Client{}
	req, err := http.NewRequest("DELETE", os.Getenv("ACMESKY_SERVICE_API") + "/" + variables["departBooking"].(string), nil)
	_, err = httpClient.Do(req)

	if err != nil {
		util.FailJob(client, job)
		return
	}

	req, err = http.NewRequest("DELETE", os.Getenv("ACMESKY_SERVICE_API") + "/" + variables["returnBooking"].(string), nil)
	_, err = httpClient.Do(req)

	if err != nil {
		util.FailJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func SendUnbookFlightHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageUnbookFlight").CorrelationKey(variables["offerCode"].(string)).VariablesFromMap(variables)
	_, err = message.Send(util.Ctx)
	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		util.FailJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}
