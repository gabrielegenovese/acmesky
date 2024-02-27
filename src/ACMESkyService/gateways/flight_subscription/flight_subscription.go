package flight_subscription

import (
	"acmesky/entities"
	airportsRepo "acmesky/repository/airports"
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

// getAlbums responds with the list of all albums as JSON.
func rest_getAirports(ctx *gin.Context) {

	searchQuery := ctx.Query("query")
	airports, err := airportsRepo.GetAirports(searchQuery)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, airports)
	} else {
		ctx.IndentedJSON(http.StatusOK, airports)
	}

}

func rest_subscribeTravelPreference(ctx *gin.Context) {
	var newSubRequest entities.CustomerFlightSubscription

	if err := ctx.BindJSON(&newSubRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	zbClient := *zbSingleton.GetInstance()

	_, err := bpmn_NotifyReceivedTravelPreference(zbClient, newSubRequest)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.Status(http.StatusOK)
	}

}

func bpmn_NotifyReceivedTravelPreference(zBClient zbc.Client, newSubRequest entities.CustomerFlightSubscription) (*pb.PublishMessageResponse, error) {

	command, err := zBClient.NewPublishMessageCommand().
		MessageName("Message_ReceivedTravelSubscription").
		CorrelationKey("").
		VariablesFromMap(newSubRequest.ToMap())

	if err != nil {
		log.Println(fmt.Errorf("failed to create process instance command for message [%+v]", newSubRequest))
		return nil, err
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	result, err := command.Send(ctx)

	if err != nil {
		log.Println(fmt.Errorf("error on saving preference: %s", err))
	} else {
		log.Printf("notifed we received [%+v] with key %d", newSubRequest, result.GetKey())
	}

	return result, err
}

func Listen(router *gin.Engine) {

	router.GET("/airports", rest_getAirports)
	router.PUT("/subscribe", rest_subscribeTravelPreference)
}
