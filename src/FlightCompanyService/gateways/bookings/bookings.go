package bookings

import (
	"flightcompany/dao/entities"
	flightsDAO "flightcompany/dao/impl/flights"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Book a flight
//
//	@Summary		Book a flight
//	@Description	book a flight which can be confirmed or unbooked
//	@Tags			Bookings
//	@Produce		json
//	@Param			origin_airport			query	string					true	"ID of origin airport from GET /airports"
//	@Param			dest_airport			query	string					true	"ID of destination airport from GET /airports"
//	@Param			passengers_count		query	number					true	"number > 0 of min passengers seats count that need be available"
//	@Param			start_range_datetime	query	string					true	"string rappresenting the datetime for start filter range"
//	@Param			end_range_datetime		query	string					true	"string rappresenting the datetime for end filter range"
//	@Success		200						{array}	entities.FlightBooking	"Flight has been booked"
//	@Failure		400						"If provided booking is in wrong format"
//	@Failure		500						"If any error occurred"
//	@Router			/bookings/ [post]
func rest_bookFlights(ctx *gin.Context) {
	var bookingRequest entities.FlightBooking
	var bookingRespose entities.FlightBooking
	var err error

	if err := ctx.BindJSON(&bookingRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	bookingRespose.BookingID, err = flightsDAO.ReserveFlightBooking(bookingRequest)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.IndentedJSON(http.StatusOK, bookingRespose)
	}
}

// Unbook a flight
//
//	@Summary		Unbook a flight
//	@Description	Unbook a booked flight
//	@Tags			Bookings
//	@Produce		json
//	@Param			bookingID	path	string	true	"A booking ID"
//	@Success		200			"Flight has been unbooked"
//	@Failure		400			"If booking has not been provided"
//	@Failure		404			"If booking doesn't exist"
//	@Failure		500			"If any error occurred and flight has not been unbooked"
//	@Router			/bookings/{bookingID} [delete]
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

	err = flightsDAO.RemoveFlightBooking(bookingID)

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

// Confirm a booking
//
//	@Summary		Confirm a booking
//	@Description	Confirm / pay a booking of flight which was not unconfirmed yet
//	@Tags			Bookings
//	@Produce		json
//	@Param			bookingID	path	string	true	"A booking ID"
//	@Success		200			"Booking has been confirmed"
//	@Failure		400			"If booking has not been provided"
//	@Failure		404			"If booking doesn't exist"
//	@Failure		500			"If any error occurred and flight has not been confirmed"
//	@Router			/bookings/{bookingID}/confirm [post]
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

	err = flightsDAO.ConfirmFlightBooking(bookingID)

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
