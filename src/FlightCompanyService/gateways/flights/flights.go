package flights

import (
	flightsRepo "flightcompany/repository/flights"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*
 * Return a JSON list of Flight entities
 * Method: GET
 * Query Parameters:
 * - (Required) origin_airport: ID of origin airport from GET /airports
 * - (Required) dest_airport: ID of destination airport from GET /airports
 * - (Required) passengers_count: number > 0 of min passengers seats count that need be available
 * - (Required) start_range_datetime:  string rappresenting the datetime for start filter range
 * - (Required) end_range_datetime: string rappresenting the datetime for end filter range
 */
func rest_getFlights(ctx *gin.Context) {

	var originAirportID string = ctx.Query("origin_airport")
	var destinationAirportID string = ctx.Query("dest_airport")
	passengersCount, err1 := strconv.Atoi(ctx.Query("passengers_count"))
	startRangeDatetime, err2 := time.Parse(time.DateTime, ctx.Query("start_range_datetime"))
	endRangeDatetime, err3 := time.Parse(time.DateTime, ctx.Query("end_range_datetime"))

	if err1 != nil || err2 != nil || err3 != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	passengersCount = max(1, passengersCount)

	flights, err := flightsRepo.GetFlights(
		originAirportID,
		destinationAirportID,
		startRangeDatetime,
		endRangeDatetime,
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
