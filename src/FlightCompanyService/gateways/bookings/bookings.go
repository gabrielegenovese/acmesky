package bookings

import "github.com/gin-gonic/gin"

func rest_bookFlights(ctx *gin.Context) {

}

func rest_unbookFlight(ctx *gin.Context) {

}

func rest_confirmBooking(ctx *gin.Context) {

}

func Listen(router *gin.Engine) {

	bookings := router.Group("/bookings")
	{
		bookings.POST("/reserve", rest_bookFlights)
		bookings.DELETE("/reserve/:bookingID", rest_unbookFlight)
		bookings.POST("/confirm/:bookingID", rest_confirmBooking)
	}
}
