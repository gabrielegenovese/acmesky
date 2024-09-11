package prontogram

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

type Message struct {
	Sender         util.User `json:"sender"`
	ReceiverUserId string    `json:"receiverUserId"`
	Content        string    `json:"content"`
}

func SendNewMessageHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	util.OfferSelected[variables["offerCode"].(string)] = make(chan struct{})
	messageRest := Message{
		Sender:         util.ProntogramUser,
		ReceiverUserId: variables["prontogramID"].(string),
		Content:        "OfferCode: " + variables["offerCode"].(string),
	}
	data, _ := json.Marshal(messageRest)
	http.Post(os.Getenv("PRONTOGRAM_API")+"/api/users/prontogram/messages", "application/json", bytes.NewReader(data))
	message, _ := util.ZbClient.NewPublishMessageCommand().MessageName("MessageNewMessage").CorrelationKey("correlation").VariablesFromMap(variables)
	_, err = message.Send(util.Ctx)

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
