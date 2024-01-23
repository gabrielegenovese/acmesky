package flight_subscription

import (
	"acmesky/entities"
	zbSingleton "acmesky/workers"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/gin-gonic/gin"
)

var airports = []entities.Airport{
	{ID: "1", Name: "Aeroporto di Bologna-Guglielmo Marconi", City: "Bologna"},
	{ID: "2", Name: "Aeroporto del Salento", City: "Brindisi"},
}

// getAlbums responds with the list of all albums as JSON.
func getAirports(ctx *gin.Context) {
	// var client zbc.Client
	// client = *zbSingleton.GetInstance()

	ctx.IndentedJSON(http.StatusOK, airports)
}

func subscribeTravelPreference(ctx *gin.Context) {
	var newSubRequest entities.CustomerFlightSubscription

	if err := ctx.BindJSON(&newSubRequest); err != nil {
		log.Println(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", newSubRequest)

	zbClient := *zbSingleton.GetInstance()

	_, err := receivedTravelPreference(zbClient, newSubRequest)

	if err != nil {
		fmt.Println(fmt.Errorf("Error on saving preference: %s", err))
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.Status(http.StatusOK)
	}

}

func receivedTravelPreference(zBClient zbc.Client, newSubRequest entities.CustomerFlightSubscription) (*pb.PublishMessageResponse, error) {

	command, err := zBClient.NewPublishMessageCommand().
		MessageName("Message_ReceivedTravelSubscription").
		CorrelationKey("").
		VariablesFromMap(newSubRequest.ToMap())

	if err != nil {
		log.Println(fmt.Errorf("failed to create process instance command for message [%s]", newSubRequest))
		return nil, err
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	result, err := command.Send(ctx)

	if err == nil {
		log.Printf("notifed we received [%+v] with key %d", newSubRequest, result.GetKey())
	}

	return result, err
}

func Listen(router *gin.Engine) {

	router.GET("/airports", getAirports)
	router.PUT("/subscribe", subscribeTravelPreference)
}
