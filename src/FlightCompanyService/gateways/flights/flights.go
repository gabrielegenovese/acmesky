package flights

import (
	"bytes"
	"encoding/json"
	"flightcompany/dao/entities"
	flightsDAO "flightcompany/dao/impl/flights"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Add a flight
//
//	@Summary		Add a flight
//	@Description	Add a flight
//	@Tags			Flight
//	@Accept			json
//	@Produce		json
//	@Param			request			body	entities.Flight	true	"The new flight"
//	@Success		200						{object}	entities.Flight
//	@Failure		400
//	@Failure		500
//	@Router			/flights [post]
func rest_postFlights(ctx *gin.Context) {

	var newFlight entities.Flight
	if err := ctx.BindJSON(&newFlight); err != nil {
		log.Printf("Bind error: %s", err)
		return
	}
	id, err := flightsDAO.AddFlight(newFlight)
	if err != nil {
		log.Printf("Add error: %s", err)
		ctx.Status(http.StatusInternalServerError)
	} else {
		newFlight.FlightID = fmt.Sprintf("%d", id)
		data, _ := json.Marshal(newFlight)
		http.Post(os.Getenv("WORKERS_API") + "/newFlight", "application/json", bytes.NewReader(data))
		ctx.IndentedJSON(http.StatusOK, newFlight)
	}
}

// Get a list of Flights matching the filter
//
//	@Summary		Get all filtered flights
//	@Description	Get a list of all available flights over the provided filters
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

// Get a list of Flights
//
//	@Summary		Get all flights
//	@Description	Get a list of all available flights
//	@Tags			Flight
//	@Produce		json
//	@Success		200						{array}	entities.Flight
//	@Failure		400
//	@Failure		500
//	@Router			/flights [get]
func rest_getAllFlights(ctx *gin.Context) {
	flights, err := flightsDAO.GetAllFlights()

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.IndentedJSON(http.StatusOK, flights)
	}
}

// Get a specified Flight
//
//	@Summary		Get a flight
//	@Description	Get the flight with the specified ID
//	@Tags			Flight
//	@Produce		json
//	@Param			id			path	string	true	"ID of the flight"
//	@Success		200						{object}	entities.Flight
//	@Failure		400
//	@Failure		500
//	@Router			/flight/{id} [get]
func rest_getFlight(ctx *gin.Context) {
	flightId, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	flight, err := flightsDAO.GetFlight(flightId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.IndentedJSON(http.StatusOK, flight)
	}
}
func Listen(router *gin.Engine) {

	router.POST("/flights", rest_postFlights)
	router.GET("/flights", rest_getFlights)
	router.GET("/allFlights", rest_getAllFlights)
	router.GET("/flight/:id", rest_getFlight)
}
