package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

var doneSecondStep = make(map[string]chan struct{})
var readyClose = make(chan struct{})

var api = "http://localhost:8080"

var zbClient zbc.Client
var ctx context.Context

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// ncc_backend_api, ok := os.LookupEnv("NCC_BACKEND_API")
	// if !ok {
	// 	ncc_backend_api = "http://localhost"
	// }
	router := gin.Default()
	router.POST("/newInterest", newInterest)
	router.POST("/secondStep/:id", secondStep)

	zbClient, err = zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("ZEEBE_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}
	// deploy process
	ctx = context.Background()
	response, err := zbClient.NewDeployResourceCommand().AddResourceFile("../../resources/acmesky.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response: ", response.String())


	zbClient.NewJobWorker().JobType("NewInterest").Handler(newInterestHandler).Open()
	zbClient.NewJobWorker().JobType("SecondStep").Handler(secondStepHandler).Open()
	zbClient.NewJobWorker().JobType("Timer").Handler(timerHandler).Open()
	zbClient.NewJobWorker().JobType("SendMessage").Handler(sendMessageHandler).Open()
	zbClient.NewJobWorker().JobType("ReceiveMessage").Handler(receiveMessageHandler).Open()
	zbClient.NewJobWorker().JobType("genericHandler").Handler(genericHandler).Open()
	router.Run(":8080")
	// jobWorkerAddNCC := zbClient.NewJobWorker().JobType("addNCC").Handler(addNCC).Open()
	// jobWorkerGetNCC := zbClient.NewJobWorker().JobType("getNCC").Handler(getNCC).Open()
	// <-readyClose
	// jobWorkerAddNCC.Close()
	// jobWorkerAddNCC.AwaitClose()
	// jobWorkerGetNCC.Close()
	// jobWorkerGetNCC.AwaitClose()
}

func newInterest(c *gin.Context) {
	var newInterest CustomerFlightSubscriptionRequest
	if err := c.BindJSON(&newInterest); err != nil {
		log.Printf("Bind error: %s", err)
		return
	}
	// create a new process instance
	variables := make(map[string]interface{})
	userId := string(rand.IntN(100))
	fmt.Println(newInterest)
	variables["AirportOriginID"] = newInterest.AirportOriginID
	variables["prontogramId"] = newInterest.ProntogramID
	variables["correlationKey"] = "correlation"
	doneSecondStep[userId] = make(chan struct{})
	fmt.Println("User id: ", userId);

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("NewClient").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	result, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result: ", result.String())
}

func secondStep(c *gin.Context) {
	userId := c.Param("id")
	close(doneSecondStep[userId])
}

func genericHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
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

func newInterestHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
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

func secondStepHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	fmt.Println("userIdHandler: ", variables["userId"]);
	<- doneSecondStep[variables["userId"].(string)];
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
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

func timerHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
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

func sendMessageHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := zbClient.NewPublishMessageCommand().MessageName("msgName").CorrelationKey("correlation").VariablesFromMap(variables)
	_, err = message.Send(ctx);
	message, err = zbClient.NewPublishMessageCommand().MessageName("msgTest").CorrelationKey("correlation").VariablesFromMap(variables)
	_, err = message.Send(ctx);
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
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

func receiveMessageHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
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

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}
