package prontogram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"workers/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

type Message struct {
	Sender         util.User `json:"sender"`
	ReceiverUserId string    `json:"receiverUserId"`
	Content        string    `json:"content"`
}

func SendNewMessageHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()
	var messageParams util.SendMessageParams

	err := job.GetVariablesAs(&messageParams)
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}
	offer := messageParams.Offer
	offerCode := messageParams.OfferCode
	response, err := http.Get(os.Getenv("ACMESKY_SERVICE_API") + "/airports")
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}
	var airports []util.Airport
	json.NewDecoder(response.Body).Decode(&airports)

	departFlight := offer.DepartFlight
	returnFlight := offer.ReturnFlight
	var departAirport util.Airport
	var returnAirport util.Airport
	for _, airport := range airports {
		if airport.AirportID == departFlight.AirportOriginID {
			departAirport = airport
		} else if airport.AirportID == returnFlight.AirportOriginID {
			returnAirport = airport
		}
	}
	datetimeDepart, _ := time.Parse(time.DateOnly, offer.TravelPreference.DateStartISO8601)
	datetimeReturn, _ := time.Parse(time.DateOnly, offer.TravelPreference.DateEndISO8601)

	messageContent := fmt.Sprintf("ACMESKY found a flight travel offer with return flight for you (%v seats) from %v to %v between %s and %s within your budget of %v %v\n"+
		"We offer the following flights at the price of %v %v :\n"+
		" - Depart flight from %v at %v will arrive at %v \n"+
		" - Return flight from %v at %v will return at %v. \n"+
		"Use the code %v within 24h to purchase this offer on our portal at this reserved price !\n",
		offer.TravelPreference.SeatsCount, departAirport.City, returnAirport.City, datetimeDepart.Format(time.DateOnly), datetimeReturn.Format(time.DateOnly), offer.TravelPreference.Budget, "Euro",
		offer.TotalPrice, "Euro",
		departAirport.Name, departFlight.DepartDatetime, departFlight.ArrivalDatetime,
		returnAirport.Name, returnFlight.DepartDatetime, returnFlight.ArrivalDatetime,
		offerCode,
	)

	util.OfferSelected[offerCode] = make(chan struct{})
	messageRest := Message{
		Sender:         util.ProntogramUser,
		ReceiverUserId: offer.TravelPreference.ProntogramID,
		Content:        messageContent,
	}
	data, _ := json.Marshal(messageRest)
	http.Post(os.Getenv("PRONTOGRAM_API")+"/api/users/prontogram/messages", "application/json", bytes.NewReader(data))

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromObject(messageParams)
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
