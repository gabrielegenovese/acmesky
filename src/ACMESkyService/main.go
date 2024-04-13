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

func main() {
	godotenv.Load(".env")

	dbClient.InitDB()
	workers := travelPrefWorker.RegisterWorkers()
	workers = append(workers, flightMatcher.RegisterWorkers()...)
	defer zbSingleton.UnregisterWorkers(workers)

	prontogram.Init()

	gateways.Listen()
}
