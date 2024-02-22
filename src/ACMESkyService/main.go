package main

import (
	"github.com/gin-gonic/gin"

	"acmesky/gateways/flight_subscription"
	travelPrefWorker "acmesky/workers/travel_preference"
)

func main() {

	workers := travelPrefWorker.RegisterWorkers()
	defer travelPrefWorker.UnregisterWorkers(workers)

	router := gin.Default()
	flight_subscription.Listen(router)
	router.Run("localhost:8080")
}
