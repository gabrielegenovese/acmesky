package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"workers/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerFlightSubscriptionRequest struct {
	ProntogramID         string `json:"customer_prontogram_id"`
	AirportOriginID      string `json:"airport_id_origin"`
	AirportDestinationID string `json:"airport_id_destination"`
	// start travel date range in ISO 8601 format (with timezone UTC)
	DateStartISO8601 string `json:"travel_date_start"`
	// end travel date range in ISO 8601 format (with timezone UTC)
	DateEndISO8601 string  `json:"travel_date_end"`
	Budget         float32 `json:"travel_max_price"`
	SeatsCount     uint    `json:"travel_seats_count"`
}

func NewInterest(c *gin.Context) {
	var newInterest CustomerFlightSubscriptionRequest
	if err := c.BindJSON(&newInterest); err != nil {
		log.Printf("Bind error: %s", err)
		return
	}
	// create a new process instance
	variables := make(map[string]interface{})
	variables["prontogramId"] = newInterest.ProntogramID
	variables["AirportOriginID"] = newInterest.AirportOriginID
	variables["AirportDestinationID"] = newInterest.AirportDestinationID
	variables["DateStartISO8601"] = newInterest.DateStartISO8601
	variables["DateEndISO8601"] = newInterest.DateEndISO8601
	variables["Budget"] = newInterest.Budget
	variables["SeatsCount"] = newInterest.SeatsCount
	uuid := uuid.NewString()
	variables["correlationKey"] = uuid
	util.InterestSaved[uuid] = make(chan struct{})

	request, err := util.ZbClient.NewCreateInstanceCommand().BPMNProcessId("User").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	result, err := request.Send(util.Ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Result: ", result.String())
	c.Data(http.StatusOK, "text/plain", []byte(uuid))
}

func NewInterestHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
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

func SendNewInterestHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageNewInterest").CorrelationKey(variables["correlationKey"].(string)).VariablesFromMap(variables)
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

func SelectOfferHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	// Wait for use
	<-util.OfferSelected[variables["offerCode"].(string)]
	response, _ := http.Get(os.Getenv("ACMESKY_SERVICE_API") + "/api/v1/offer/" + variables["offerCode"].(string))

	variables["validOffer"] = response.StatusCode == 200

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

func BuyOffer(c *gin.Context) {
	id := c.Param("id")
	close(util.OfferSelected[id])
	util.BuyResults[id] = make(chan util.BuyResult)
	buyResult := <- util.BuyResults[id]
	c.IndentedJSON(http.StatusOK, buyResult)
}

func SendBookFlightHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageBookFlight").CorrelationKey(variables["offerCode"].(string)).VariablesFromMap(variables)
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

func CalculateDistanceHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	// Wait for use
	util.NCCSearchRequests[variables["paymentId"].(string)] = make(chan util.NCCSearchRequest)
	nccSearchRequest := <-util.NCCSearchRequests[variables["paymentId"].(string)]
	response, _ := http.Get(os.Getenv("DISTANCE_API") + "/distance?from=" + nccSearchRequest.Address + "&to=" + nccSearchRequest.AirportAddress)
	var distance util.DistanceResBody
	json.NewDecoder(response.Body).Decode(&distance)
	variables["distance"] = distance.Value
	variables["address"] = nccSearchRequest.Address

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

func SearchNCC(c *gin.Context) {
	var nccRequest util.NCCSearchRequest
	if err := c.BindJSON(&nccRequest); err != nil {
		log.Printf("Bind error: %s", err)
		return
	}
	util.NCCSearchRequests[nccRequest.PaymentId] <- nccRequest
	ncc := <- util.NCCResponses[nccRequest.PaymentId]
	c.IndentedJSON(http.StatusOK, ncc)
}

func SendSearchNCCHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageSearchNCC").CorrelationKey(variables["offerCode"].(string)).VariablesFromMap(variables)
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
