package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/joho/godotenv"
)

var readyClose = make(chan struct{})

func main() {
	godotenv.Load()
	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("ZEEBE_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}

	// deploy process
	ctx := context.Background()
	response, err := zbClient.NewDeployResourceCommand().AddResourceFile("../bank.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response: ", response.String())

	// create a new process instance
	variables := make(map[string]interface{})
	variables["user"] = "mario"
	variables["amount"] = 100
	variables["description"] = "asd"

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("bank").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	result, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result: ", result.String())

	jobWorkerCreatePayment := zbClient.NewJobWorker().JobType("create-payment").Handler(createPayment).Open()
	jobWorkerPay := zbClient.NewJobWorker().JobType("pay").Handler(pay).Open()
	<-readyClose
	jobWorkerCreatePayment.Close()
	jobWorkerCreatePayment.AwaitClose()
	jobWorkerPay.Close()
	jobWorkerPay.AwaitClose()
}

func createPayment(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	json_data, _ := json.Marshal(variables)

	resp, _ := http.Post("http://localhost:3000/payment/new", "application/json", bytes.NewBuffer(json_data))
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println(res)
	variables["id"] = res["id"]
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Processing order:", variables["id"])
	log.Println("Name:", variables["user"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func pay(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	json_data, _ := json.Marshal(variables)

	resp, _ := http.Post("http://localhost:3000/payment/pay/"+variables["id"].(string), "application/json", bytes.NewBuffer(json_data))
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println("Pay result: ", res)

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
	close(readyClose)
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}
