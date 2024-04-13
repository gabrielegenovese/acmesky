package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	dbClient "acmesky/dao/db"
	"acmesky/gateways/flight_subscription"
	"acmesky/services/notification/prontogram"
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

	prontogram.Init()

	router := gin.Default()
	flight_subscription.Listen(router)
	router.Run("localhost:8090")
}
