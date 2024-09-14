package bank

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
	"github.com/gin-gonic/gin"
)

func CreateNewPaymentHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)

	paymentRequest := util.PaymentReq{
		User: variables["prontogramId"].(string),
		Description: variables["offerCode"].(string),
		Amount: variables["totalPrice"].(float64),
	}
	data, _ := json.Marshal(paymentRequest)
	response, err := http.Post(os.Getenv("BANK_API") + "/payment/new", "application/json", bytes.NewReader(data))
	var payment util.Payment
	json.NewDecoder(response.Body).Decode(&payment)
	if err != nil {
		util.BuyResults[variables["offerCode"].(string)] <- util.BuyResult{Success: false}
		util.FailJob(client, job)
		return
	}
	variables["paymentId"] = payment.ID.String()
	util.PaymentCompleted[payment.ID.String()] = make(chan struct{})
	util.BuyResults[variables["offerCode"].(string)] <- util.BuyResult{Success: true, PaymentLink: payment.Link, PaymentID: payment.ID.String()}

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

func ExecutePaymentHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	// Wait for confirm
	<-util.PaymentCompleted[variables["paymentId"].(string)]
	variables["paymentCompleted"] = true
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

func PaymentCompleted(c *gin.Context) {
	id := c.Param("id")
	close(util.PaymentCompleted[id])
}

func SendReceiptHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		util.FailJob(client, job)
		return
	}

	log.Println("Start job", jobKey, "of type", job.Type)
	message, err := util.ZbClient.NewPublishMessageCommand().MessageName("MessageReceipt").CorrelationKey(variables["paymentId"].(string)).VariablesFromMap(variables)
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
