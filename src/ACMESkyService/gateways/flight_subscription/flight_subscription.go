package flight_subscription

import (
	"acmesky/dao/entities"
	airportsDAO "acmesky/dao/impl/airports"
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

type ResError struct {
	Error string `json:"error"`
}

//	@Summary		Get all airports
//	@Description	Get a list of all available airports
//	@Tags			airport
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	false	"Search query"
//	@Success		200		{array}		[]entities.Airport
//	@Failure		500		{object}	ResError
//	@Router			/airports [get]
func rest_getAirports(ctx *gin.Context) {

	searchQuery := ctx.Query("query")
	airports, err := airportsDAO.GetAirports(searchQuery)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, ResError{Error: err.Error()})
	} else {
		ctx.IndentedJSON(http.StatusOK, airports)
	}

}

// Create a Travel Preference and subscribe to notification service (a sort of Newsletter like servive)
//	@Summary		Add a Travel preference
//	@Description	Add a Travel preference and subscribe to notification service
//	@Tags			travel_preference
//	@Accept			json
//	@Param			body	body		entities.CustomerFlightSubscriptionRequest	true	"subscription data"	SchemaExample(testUser)
//	@Success		200		{object}	string
//	@Failure		400		{object}	ResError
//	@Failure		500		{object}	ResError
//	@Router			/subscribe [post]
func rest_subscribeTravelPreference(context *gin.Context) {
	var newSubRequest entities.CustomerFlightSubscriptionRequest

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
		context.IndentedJSON(http.StatusInternalServerError, ResError{Error: err.Error()})
	} else {
		// waiting result
		outVars := <-result

		if _, hasError := outVars["errorCode"]; hasError {
			context.IndentedJSON(http.StatusInternalServerError, ResError{Error: "bpmn has got an error"})
		} else {
			context.Status(http.StatusOK)
		}
	}
	chanBPRepo.UnsetContext(bpk_uuid.String())
}

func bpmn_NotifyReceivedTravelPreference(zBClient zbc.Client, bpk string, newSubRequest entities.CustomerFlightSubscriptionRequest) (*pb.PublishMessageResponse, error) {

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

func Listen(router *gin.RouterGroup) {

	router.GET("/airports", rest_getAirports)
	router.POST("/subscribe", rest_subscribeTravelPreference)
}
