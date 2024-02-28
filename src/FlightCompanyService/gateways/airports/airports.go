package flight_subscription

import (
	airportsRepo "flightcompany/repository/airports"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getAlbums responds with the list of all albums as JSON.
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
