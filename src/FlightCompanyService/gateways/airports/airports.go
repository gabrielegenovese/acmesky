package flight_subscription

import (
	airportsRepo "flightcompany/repository/airports"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * Return a JSON list of Airports
 * Method: GET
 * Query parameters:
 * - query: a string to return airports which names contains this substring
 */
func rest_getAirports(ctx *gin.Context) {

	searchQuery := ctx.Query("query")
	airports, err := airportsRepo.GetAirports(searchQuery)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, airports)
	} else {
		ctx.IndentedJSON(http.StatusOK, airports)
	}

}

func Listen(router *gin.Engine) {

	router.GET("/airports", rest_getAirports)
}
