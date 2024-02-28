package main

import (
	"github.com/gin-gonic/gin"

	"flightcompany/gateways/flights"
	dbClient "flightcompany/repository/db"
)

func main() {

	dbClient.InitDB()

	router := gin.Default()
	flights.Listen(router)
	router.Run("localhost:8081")
}
