package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var readyClose = make(chan struct{})

var api = "http://localhost:8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("ZEEBE_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}

	// deploy process
	ctx := context.Background()
	response, err := zbClient.NewDeployResourceCommand().AddResourceFile("../ncc.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response: ", response.String())

	// create a new process instance
	variables := make(map[string]interface{})
	variables["name"] = "mario"
	variables["price"] = "100"
	variables["location"] = "Roma"

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("ncc").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	result, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result: ", result.String())

	jobWorkerAddNCC := zbClient.NewJobWorker().JobType("addNCC").Handler(addNCC).Open()
	jobWorkerGetNCC := zbClient.NewJobWorker().JobType("getNCC").Handler(getNCC).Open()
	<-readyClose
	jobWorkerAddNCC.Close()
	jobWorkerAddNCC.AwaitClose()
	jobWorkerGetNCC.Close()
	jobWorkerGetNCC.AwaitClose()
}

func addNCC(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	json_data, err := json.Marshal(variables)

	resp, _ := http.Post(api+"/addNCC", "application/json", bytes.NewBuffer(json_data))
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println(res)
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

func getNCC(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	resp, _ := http.Get(api + "/getNCC")
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println(res)
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
