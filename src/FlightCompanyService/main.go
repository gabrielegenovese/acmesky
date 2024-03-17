package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"flightcompany/gateways/flights"
	dbClient "flightcompany/repository/db"
)

func main() {

	godotenv.Load(".env")
	dbClient.InitDB()

	router := gin.Default()
	flights.Listen(router)
	router.Run("localhost:8079")
}
