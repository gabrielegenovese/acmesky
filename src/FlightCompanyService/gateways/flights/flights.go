package flights

import (
	flightsRepo "flightcompany/repository/flights"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func rest_getFlights(ctx *gin.Context) {

	var originAirportID string = ctx.Query("origin_airport")
	var destinationAirportID string = ctx.Query("dest_airport")
	passengersCount, err1 := strconv.Atoi(ctx.Query("passengers_count"))
	departDateTime, err2 := time.Parse(time.DateOnly, ctx.Query("depart_date"))

	if err1 != nil || err2 != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	passengersCount = max(1, passengersCount)

	flights, err := flightsRepo.GetFlights(
		originAirportID,
		destinationAirportID,
		departDateTime,
		passengersCount,
	)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.IndentedJSON(http.StatusOK, flights)
	}

}

func Listen(router *gin.Engine) {

	router.GET("/flights", rest_getFlights)
}
