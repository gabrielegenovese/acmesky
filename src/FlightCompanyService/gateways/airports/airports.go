package flight_subscription

import (
	airportsDAO "flightcompany/dao/impl/airports"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get a list of all available airports
//
//	@Summary		Get all airports
//	@Description	Get a list of all available airports
//	@Tags			Airports
//	@Produce		json
//	@Param			query	query	string	false	"Search string to return airports which names contains this substring"
//	@Success		200		{array}	entities.Airport
//	@Failure		500
//	@Router			/airports [get]
func rest_getAirports(ctx *gin.Context) {

	searchQuery := ctx.Query("query")
	airports, err := airportsDAO.GetAirports(searchQuery)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, airports)
	} else {
		ctx.IndentedJSON(http.StatusOK, airports)
	}

}

func Listen(router *gin.Engine) {

	router.GET("/airports", rest_getAirports)
}
