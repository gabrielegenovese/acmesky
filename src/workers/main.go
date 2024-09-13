package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"workers/acmesky"
	"workers/bank"
	"workers/flightcompany"
	"workers/ncc"
	"workers/prontogram"
	"workers/user"
	"workers/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var readyClose = make(chan struct{})

var api = "http://localhost:8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// ncc_backend_api, ok := os.LookupEnv("NCC_BACKEND_API")
	// if !ok {
	// 	ncc_backend_api = "http://localhost"
	// }
	resp, err := http.Post(os.Getenv("PRONTOGRAM_API")+"/api/auth/prontogram/login", "application/json", bytes.NewBuffer([]byte(`{"userId": "prontogram", "password": "sicurissima"}`)))
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&util.ProntogramUser)
	log.Println(util.ProntogramUser)

	router := gin.Default()
	router.POST("/newInterest", user.NewInterest)
	router.GET("/buyOffer/:id", user.BuyOffer)
	router.POST("/searchNCC", user.SearchNCC)
	router.GET("/newInterestSaved/:id", acmesky.NewInterestSaved)
	router.POST("/newFlight", flightcompany.NewFlight)
	router.GET("/paymentCompleted/:id", bank.PaymentCompleted)

	util.ZbClient, err = zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("ZEEBE_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}
	// deploy process
	util.Ctx = context.Background()
	response, err := util.ZbClient.NewDeployResourceCommand().AddResourceFile("./acmesky.bpmn").Send(util.Ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response: ", response.String())

	util.ZbClient.NewJobWorker().JobType("NewInterest").Handler(user.NewInterestHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SendNewInterest").Handler(user.SendNewInterestHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SelectOffer").Handler(user.SelectOfferHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SendBookFlight").Handler(user.SendBookFlightHandler).Open()
	util.ZbClient.NewJobWorker().JobType("CalculateDistance").Handler(user.CalculateDistanceHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SendSearchNCC").Handler(user.SendSearchNCCHandler).Open()

	util.ZbClient.NewJobWorker().JobType("SaveInterest").Handler(acmesky.SaveInterestHandler).Open()
	util.ZbClient.NewJobWorker().JobType("UpdateFlightList").Handler(acmesky.UpdateFlightListHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SendNewOffert").Handler(acmesky.SendNewOffertHandler).Open()
	util.ZbClient.NewJobWorker().JobType("BookFlight").Handler(acmesky.BookFlightHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SendNewPayment").Handler(acmesky.SendNewPaymentHandler).Open()

	util.ZbClient.NewJobWorker().JobType("CreateNewPayment").Handler(bank.CreateNewPaymentHandler).Open()
	util.ZbClient.NewJobWorker().JobType("ExecutePayment").Handler(bank.ExecutePaymentHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SendReceipt").Handler(bank.SendReceiptHandler).Open()

	util.ZbClient.NewJobWorker().JobType("SendNewMessage").Handler(prontogram.SendNewMessageHandler).Open()

	util.ZbClient.NewJobWorker().JobType("SendNewFlight").Handler(flightcompany.SendNewFlightHandler).Open()

	util.ZbClient.NewJobWorker().JobType("NearestNCC").Handler(ncc.NearestNCCHandler).Open()
	util.ZbClient.NewJobWorker().JobType("BookNCC").Handler(ncc.BookNCCHandler).Open()
	util.ZbClient.NewJobWorker().JobType("SendNCCBooking").Handler(ncc.SendNCCBookingHandler).Open()

	util.ZbClient.NewJobWorker().JobType("ReceiveMessage").Handler(receiveMessageHandler).Open()
	util.ZbClient.NewJobWorker().JobType("genericHandler").Handler(genericHandler).Open()
	router.Run(":8080")
	// jobWorkerAddNCC := zbClient.NewJobWorker().JobType("addNCC").Handler(addNCC).Open()
	// jobWorkerGetNCC := zbClient.NewJobWorker().JobType("getNCC").Handler(getNCC).Open()
	// <-readyClose
	// jobWorkerAddNCC.Close()
	// jobWorkerAddNCC.AwaitClose()
	// jobWorkerGetNCC.Close()
	// jobWorkerGetNCC.AwaitClose()
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
