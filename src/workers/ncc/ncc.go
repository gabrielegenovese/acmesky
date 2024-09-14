package ncc

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"workers/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

func NearestNCCHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)

	response, err := http.Get(os.Getenv("NCC_API") + "/getNCC")
	var nccs []util.NCC
	minDistance := 99999
	var nearestNCC util.NCC
	json.NewDecoder(response.Body).Decode(&nccs)
	for _, ncc := range nccs {
		response, _ := http.Get(os.Getenv("DISTANCE_API") + "/distance?from=" + ncc.Location + "&to=" + variables["address"].(string))
		var distance util.DistanceResBody
		json.NewDecoder(response.Body).Decode(&distance)
		if distance.Value < minDistance {
			minDistance = distance.Value
			nearestNCC = ncc
		}
	}
	variables["nearestNCC"] = nearestNCC.Id

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

func BookNCCHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)

	bookingRequest := util.NCCBookingRequest{
		NccId: variables["nearestNCC"].(string),
		Name: variables["prontogramId"].(string),
		Date: "",//TODO
		
	}
	data, _ := json.Marshal(bookingRequest)
	response, err := http.Post(os.Getenv("NCC_API") + "/book", "application/json", bytes.NewReader(data))
	if response.StatusCode == 200 {
		variables["nccBooked"] = true
	} else {
		variables["nccBooked"] = false
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

func SendNCCBookingHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	util.NCCResponses[variables["paymentId"].(string)] <- util.NCCResponse{
		Success: variables["nccBooked"].(bool),
		NearestNCC: variables["nearestNCC"].(string),
	}
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageNCCBooking").CorrelationKey(variables["paymentId"].(string)).VariablesFromMap(variables)
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
