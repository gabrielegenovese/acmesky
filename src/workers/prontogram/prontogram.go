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
	response, err := http.Get(os.Getenv("ACMESKY_SERVICE_API") + "/api/v1/airports")
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}
	airports := make([]util.Airport, 0)
	err = json.NewDecoder(response.Body).Decode(&airports)
	if err != nil {
		fmt.Printf("Decode Error: %+v\n", err)
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}
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

	datetimeDepart, _ := time.Parse(time.DateTime, offer.TravelPreference.DateStartISO8601)
	datetimeReturn, _ := time.Parse(time.DateTime, offer.TravelPreference.DateEndISO8601)

	messageContent := fmt.Sprintf("<p>ACMESKY found a flight travel offer with return flight for you (%v seats) from <em>%s</em> to <em>%s</em> between <time>%s</time> and <time>%s</time> within your budget of <strong>%v %v</strong><br>"+
		"We offer the following flights at the price of <strong>%v %v</strong> :<ul>"+
		"<li>Depart flight from <em>%s</em> at <time>%v</time> will arrive at <time>%s</time></li>"+
		"<li>Return flight from <em>%s</em> at <time>%v</time> will return at <time>%s</time></li>"+
		"</ul></p><p>Use the code <strong><code>%v</code></strong> within 24h to purchase this offer on our portal at this reserved price !</p>",
		offer.TravelPreference.SeatsCount, departAirport.City, returnAirport.City, datetimeDepart.Format(time.DateOnly), datetimeReturn.Format(time.DateOnly), offer.TravelPreference.Budget, "&euro;",
		offer.TotalPrice, "&euro;",
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
