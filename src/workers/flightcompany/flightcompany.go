package flightcompany

import (
	"context"
	"log"
	"workers/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/gin-gonic/gin"
)

type Flight struct {
	FlightID             string  `json:"flight_id" example:"1"`
	AirportOriginID      string  `json:"airport_origin_id" example:"5"`
	AirportDestinationID string  `json:"airport_destination_id" example:"20"`
	DepartDatetime       string  `json:"depart_datetime" example:"2024-01-01 14:00:00"`
	ArrivalDatetime      string  `json:"arrival_datetime" example:"2024-01-10 05:30:00"`
	FlightPrice          float64 `json:"flight_price" example:"199.99"`
	AvailableSeats       uint    `json:"available_seats_count" example:"2"`
}

func NewFlight(c *gin.Context) {
	var newFlight Flight
	if err := c.BindJSON(&newFlight); err != nil {
		log.Printf("Bind error: %s", err)
		return
	}
	// create a new process instance
	variables := make(map[string]interface{})
	variables["FlightID"] = newFlight.FlightID
	variables["AirportOriginID"] = newFlight.AirportOriginID
	variables["AirportDestinationID"] = newFlight.AirportDestinationID
	variables["DepartDatetime"] = newFlight.DepartDatetime
	variables["ArrivalDatetime"] = newFlight.ArrivalDatetime
	variables["FlightPrice"] = newFlight.FlightPrice
	variables["AvailableSeats"] = newFlight.AvailableSeats
	variables["correlationKey"] = "correlation"
	request, err := util.ZbClient.NewCreateInstanceCommand().BPMNProcessId("FlightCompany").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	result, err := request.Send(util.Ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Result: ", result.String())
}

func SendNewFlightHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageNewFlight").CorrelationKey("correlation").VariablesFromMap(variables)
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
