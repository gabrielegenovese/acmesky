package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"acmesky/gateways/flight_subscription"
	dbClient "acmesky/repository/db"
	travelPrefWorker "acmesky/workers/travel_preference"
)

func main() {
	godotenv.Load(".env")

	dbClient.InitDB()
	workers := travelPrefWorker.RegisterWorkers()
	defer travelPrefWorker.UnregisterWorkers(workers)

	router := gin.Default()
	flight_subscription.Listen(router)
	router.Run("localhost:8080")
}
