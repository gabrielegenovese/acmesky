package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"acmesky/gateways/flight_subscription"
	dbClient "acmesky/repository/db"
	zbSingleton "acmesky/workers"
	flightMatcher "acmesky/workers/flight_matcher"
	travelPrefWorker "acmesky/workers/travel_preference"
)

func main() {
	godotenv.Load(".env")

	dbClient.InitDB()
	workers := travelPrefWorker.RegisterWorkers()
	workers = append(workers, flightMatcher.RegisterWorkers()...)
	defer zbSingleton.UnregisterWorkers(workers)

	router := gin.Default()
	flight_subscription.Listen(router)
	router.Run("localhost:8090")
}
