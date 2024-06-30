package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	dbClient "flightcompany/dao/db"
	"flightcompany/gateways/bookings"
	"flightcompany/gateways/flights"
	// _ "flightcompany/docs"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			FlightCompany Swagger API
//	@version		1.0
//	@description	This is an API specification for the FlightCompany services.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	MIT
//	@license.url	https://opensource.org/license/mit

//	@host		localhost:8091
//	@BasePath	/

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	godotenv.Load(".env")
	dbClient.InitDB()
	defer dbClient.CloseClient()

	router := gin.Default()
	flights.Listen(router)
	bookings.Listen(router)
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8091")
}
