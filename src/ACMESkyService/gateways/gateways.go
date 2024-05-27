package gateways

import (
	"acmesky/gateways/flight_subscription"

	"github.com/gin-gonic/gin"
)

// @title           ACMESky Swagger API
// @version         1.0
// @description     This is an API specification for the ACMESKy services.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/license/mit

// @host      localhost:8090
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func Listen() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	v1 := r.Group("/api/v1")
	{
		flight_subscription.Listen(v1)
	}

	r.Run("localhost:8090")
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
