package main

import (
	"github.com/gin-gonic/gin"

	"acmesky/gateways/flight_subscription"
	dbClient "acmesky/repository/db"
	travelPrefWorker "acmesky/workers/travel_preference"
)

func main() {

	dbClient.InitDB()
	workers := travelPrefWorker.RegisterWorkers()
	defer travelPrefWorker.UnregisterWorkers(workers)

	router := gin.Default()
	flight_subscription.Listen(router)
	router.Run("localhost:8080")
}
