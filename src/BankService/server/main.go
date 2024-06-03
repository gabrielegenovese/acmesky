package main

import (
	"bank/api"
	"bank/util"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Listen  string `toml:"listen"`
	BaseURL string `toml:"base_url"`
	DbURI   string `toml:"db_uri" required:"true"`
}

var (
	// Default config values
	config = Config{
		Listen:  "0.0.0.0:8001",
		BaseURL: "http://localhost:8001",
	}
)

//	@title			Bank Service API
//	@version		1.0
//	@description	This is a minimal microservice to act as a bank.
//	@contact.name	Gabriele Genovese
//	@contact.email	gabriele.genovese2@studio.unibo.it
//	@license.name	GPLv2
//	@license.url	https://www.gnu.org/licenses/old-licenses/gpl-2.0.html
//	@Host			localhost:8001
//	@BasePath		/
func main() {
	router := gin.Default()

	err := loadConfig()
	if err != nil {
		_ = fmt.Errorf("load config error %w", err)
		os.Exit(1)
	}

	err = util.ConnectDb(config.DbURI)
	if err != nil {
		os.Exit(1)
	}
	db := util.GetDb()
	err = db.AutoMigrate(&api.Payment{})
	if err != nil {
		_ = fmt.Errorf("connet to db %w", err)
		os.Exit(1)
	}

	router.Use(CORSMiddleware())

	router.PUT("/payment/new", api.NewPayment)
	router.POST("/payment/pay/:id", api.PayPaymentById)
	router.GET("/payment/:id", api.GetPaymentById)
	router.DELETE("/payment/:id", api.DelPaymentById)

	router.Run(config.Listen)
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

func loadConfig() (err error) {
	file, err := os.Open("config.toml")
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	_, err = toml.NewDecoder(file).Decode(&config)
	if err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close config file: %w", err)
	}

	return nil
}
