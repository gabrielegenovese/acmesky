package flights

import (
	airportsRepo "flightcompany/repository/airports"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getAlbums responds with the list of all albums as JSON.
func rest_getFlights(ctx *gin.Context) {

	searchQuery := ctx.Query("query")
	airports, err := airportsRepo.GetAirports(searchQuery)

	var originAirportID string = ctx.Query("origin_airport")
	var destinationAirportID string = ctx.Query("dest_airport")

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, airports)
	} else {
		ctx.IndentedJSON(http.StatusOK, airports)
	}

}

func Listen(router *gin.Engine) {

	router.GET("/flights", rest_getFlights)
}
