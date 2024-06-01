package main

import (
	"github.com/joho/godotenv"

	dbClient "acmesky/dao/db"
	"acmesky/gateways"
	"acmesky/services/notification/prontogram"
	zbSingleton "acmesky/workers"
	flightMatcher "acmesky/workers/flight_matcher"
	travelPrefWorker "acmesky/workers/travel_preference"
)

//	@title			ACMESky Swagger API
//	@version		1.0
//	@description	This is an API specification for the ACMESKy services.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	MIT
//	@license.url	https://opensource.org/license/mit

//	@host		localhost:8090
//	@BasePath	/api/v1

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	godotenv.Load(".env")

	dbClient.InitDB()
	workers := travelPrefWorker.RegisterWorkers()
	workers = append(workers, flightMatcher.RegisterWorkers()...)
	defer zbSingleton.UnregisterWorkers(workers)

	prontogram.Init()

	gateways.Listen()
}
