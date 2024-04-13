package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	dbClient "flightcompany/dao/db"
	"flightcompany/gateways/bookings"
	"flightcompany/gateways/flights"
)

func main() {

	godotenv.Load(".env")
	dbClient.InitDB()

	router := gin.Default()
	flights.Listen(router)
	bookings.Listen(router)

	router.Run("localhost:8091")
}
