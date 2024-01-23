package main

import (
	"github.com/gin-gonic/gin"

	"acmesky/gateways/flight_subscription"
	"acmesky/workers/travel_preference"
)

func main() {

	workers := travel_preference.RegisterWorkers()
	defer travel_preference.UnregisterWorkers(workers)

	router := gin.Default()
	flight_subscription.Listen(router)
	router.Run("localhost:8080")
}
