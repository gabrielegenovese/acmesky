package bookings

import (
	"flightcompany/entities"
	flightsRepo "flightcompany/repository/flights"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
 * Book a flight
 * Path: /bookings/
 * Method: POST
 * Return:
 * - 200 if flight has been booked
 * - 400 if body input is in wrong format
 * - 500 otherwise if errors occurred
 * Body:
 * - case 200: JSON with booking_id for this booking
 */
func rest_bookFlights(ctx *gin.Context) {
	var bookingRequest entities.FlightBooking
	var bookingRespose entities.FlightBooking
	var err error

	if err := ctx.BindJSON(&bookingRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	bookingRespose.BookingID, err = flightsRepo.ReserveFlightBooking(bookingRequest)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.IndentedJSON(http.StatusOK, bookingRespose)
	}
}

/*
 * unbook a flight booking which has not been confirmed yet
 * PATH: /bookings/:booking_id
 * Method: DELETE
 * Return:
 * - 200 if flight has been unbooked
 * - 400 if body input is in wrong format
 * - 404 if provided booking is not unconfirmed / doesnt existing
 * - 500 otherwise if errors occurred
 */
func rest_unbookFlight(ctx *gin.Context) {
	strBookingID, hasBookingID := ctx.Params.Get("bookingID")

	if !hasBookingID {
		ctx.Status(http.StatusBadRequest)
		return
	}

	bookingID, err := strconv.ParseInt(strBookingID, 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = flightsRepo.RemoveFlightBooking(bookingID)

	if err != nil {
		if strings.Contains(err.Error(), "NOT FOUND") {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
	} else {
		ctx.Status(http.StatusOK)
	}
}

/*
 * Confirm / pay a booking of flight which was not unconfirmed yet
 * PATH: /bookings/:booking_id/confirm
 * Method: POST
 * Return:
 * - 200 if flight has been confirmed
 * - 400 if body input is in wrong format
 * - 404 if provided booking is not unconfirmed / doesnt existing
 * - 500 otherwise if errors occurred
 */
func rest_confirmBooking(ctx *gin.Context) {
	strBookingID, hasBookingID := ctx.Params.Get("bookingID")

	if !hasBookingID {
		ctx.Status(http.StatusBadRequest)
		return
	}

	bookingID, err := strconv.ParseInt(strBookingID, 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = flightsRepo.ConfirmFlightBooking(bookingID)

	if err != nil {
		if strings.Contains(err.Error(), "NOT FOUND") {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
	} else {
		ctx.Status(http.StatusOK)
	}
}

func Listen(router *gin.Engine) {

	bookings := router.Group("/bookings")
	{
		bookings.POST("/", rest_bookFlights)
		bookings.DELETE("/:bookingID", rest_unbookFlight)
		bookings.POST("/:bookingID/confirm", rest_confirmBooking)
	}
}
