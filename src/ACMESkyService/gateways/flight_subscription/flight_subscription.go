package flight_subscription

import (
	"acmesky/entities"
	airportsRepo "acmesky/repository/airports"
	zbSingleton "acmesky/workers"
	chanBPRepo "acmesky/workers/utils/channel_bp_repository"

	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func rest_subscribeTravelPreference(context *gin.Context) {
	var newSubRequest entities.CustomerFlightSubscription

	if err := context.BindJSON(&newSubRequest); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	zbClient := *zbSingleton.GetInstance()
	// Business Process Key
	bpk_uuid := uuid.New()
	chanBPRepo.SetContext(bpk_uuid.String())
	result := chanBPRepo.GetContext(bpk_uuid.String())

	_, err := bpmn_NotifyReceivedTravelPreference(zbClient, bpk_uuid.String(), newSubRequest)

	if err != nil {
		context.Status(http.StatusInternalServerError)
	} else {
		// waiting result
		outVars := <-result

		if _, hasError := outVars["errorCode"]; hasError {
			context.Status(http.StatusInternalServerError)
		} else {
			context.Status(http.StatusOK)
		}
	}
	chanBPRepo.UnsetContext(bpk_uuid.String())
}

func bpmn_NotifyReceivedTravelPreference(zBClient zbc.Client, bpk string, newSubRequest entities.CustomerFlightSubscription) (*pb.PublishMessageResponse, error) {

	vars := newSubRequest.ToMap()
	vars["bpk"] = bpk

	command, err := zBClient.NewPublishMessageCommand().
		MessageName("Message_ReceivedTravelSubscription").
		CorrelationKey("bpk").
		VariablesFromMap(vars)

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
