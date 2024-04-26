package flights

import (
	flightsDAO "flightcompany/dao/impl/flights"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Get a list of Flights matching the filter
//
//	@Summary		Get all filtered flights
//	@Description	Get a list of all available flights over thr provided filters
//	@Tags			Flight
//	@Produce		json
//	@Param			origin_airport			query	string	true	"ID of origin airport from GET /airports"
//	@Param			dest_airport			query	string	true	"ID of destination airport from GET /airports"
//	@Param			passengers_count		query	number	true	"number > 0 of min passengers seats count that need be available"
//	@Param			start_range_datetime	query	string	true	"string rappresenting the datetime for start filter range"
//	@Param			end_range_datetime		query	string	true	"string rappresenting the datetime for end filter range"
//	@Success		200						{array}	entities.Flight
//	@Failure		400
//	@Failure		500
//	@Router			/flights [get]
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

	flights, err := flightsDAO.GetFlights(
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
