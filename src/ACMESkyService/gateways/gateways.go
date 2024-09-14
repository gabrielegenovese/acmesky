package gateways

import (
	"acmesky/gateways/flight_subscription"
	"acmesky/gateways/flight_matcher"

	"github.com/gin-gonic/gin"
)

func Listen() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	v1 := r.Group("/api/v1")
	{
		flight_subscription.Listen(v1)
		flight_matcher.Listen(v1)
	}

	r.Run("0.0.0.0:8090")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
