package main

import (
	_ "acmeskyncc/docs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tiaguinho/gosoap"
)

type nccRequest struct {
	Name     string `xml:"name" json:"name"`
	Price    string `xml:"price" json:"price"`
	Location string `xml:"location" json:"location"`
}

type ncc struct {
	Id       string `xml:"id" json:"id"`
	Name     string `xml:"name" json:"name"`
	Price    string `xml:"price" json:"price"`
	Location string `xml:"location" json:"location"`
}

type nccList struct {
	Nccs []ncc `xml:"nccs"`
}

type bookingRequest struct {
	NccId string `xml:"nccId" json:"nccId"`
	Name  string `xml:"name" json:"name"`
	Date  string `xml:"date" json:"date"` // Format: dd/MM/yyyy kk:mm:ss
}

type booking struct {
	Success string `xml:"success" json:"success"`
}

var (
	soap *gosoap.Client
)

// add creates a new NCC
//
//	@Summary		New NCC
//	@Description	Create a new NCC
//	@Tags			ncc
//	@Accept			json
//	@Produce		json
//	@Param			request	body		nccRequest	true	"NCC"
//	@Success		200		{object}	ncc
//	@Failure		400
//	@Router			/addNCC [post]
func add(c *gin.Context) {
	var newNCC nccRequest
	if err := c.BindJSON(&newNCC); err != nil {
		log.Printf("Bind error: %s", err)
		return
	}
	params := gosoap.Params{
		"name":     newNCC.Name,
		"price":    newNCC.Price,
		"location": newNCC.Location,
	}
	res, err := soap.Call("add", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}
	var ncc ncc
	err = res.Unmarshal(&ncc)
	c.IndentedJSON(http.StatusCreated, ncc)
}

// get returns all the NCCs
//
//	@Summary		Get all NCCs
//	@Description	Get all the NCCs
//	@Tags			ncc
//	@Produce		json
//	@Success		200	{object}	[]ncc
//	@Failure		400
//	@Router			/getNCC [get]
func get(c *gin.Context) {
	res, err := soap.Call("get", gosoap.Params{})
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}
	var nccs nccList
	res.Unmarshal(&nccs)
	c.IndentedJSON(http.StatusOK, nccs.Nccs)
}

// getId returns a specified NCC
//
//	@Summary		Get a specified NCC
//	@Description	Get the NCC with the specified ID
//	@Tags			ncc
//	@Produce		json
//	@Param			id	path		string	true	"NCC ID"
//	@Success		200	{object}	ncc
//	@Failure		400
//	@Router			/getNCC/{id} [get]
func getId(c *gin.Context) {
	id := c.Param("id")
	res, err := soap.Call("getId", gosoap.Params{
		"value": id,
	})
	if err != nil {
		// Assume is due to NCCNotFound Exception
		// Jolie send 500 status with fault responses
		// but gosoap ignores the responses with status code >= 400
		// https://github.com/tiaguinho/gosoap/issues/86
		c.JSON(http.StatusNotFound, "NCC not found")
	} else {
		var ncc ncc
		err = res.Unmarshal(&ncc)
		if err != nil {
			c.JSON(http.StatusNotFound, "")
		}
		c.IndentedJSON(http.StatusOK, ncc)
	}
}

// book permits to book a specified NCC
//
//	@Summary		Book a NCC
//	@Description	Book a NCC
//	@Tags			booking
//	@Produce		json
//	@Param			request	body		bookingRequest	true	"Booking request"
//	@Success		200		{object}	booking
//	@Failure		400
//	@Router			/book [post]
func book(c *gin.Context) {
	var newBooking bookingRequest
	if err := c.BindJSON(&newBooking); err != nil {
		log.Printf("Bind error: %s", err)
		return
	}
	params := gosoap.Params{
		"nccId": newBooking.NccId,
		"name":  newBooking.Name,
		"date":  newBooking.Date,
	}
	res, err := soap.Call("book", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}
	var booking booking
	res.Unmarshal(&booking)
	c.IndentedJSON(http.StatusOK, booking)
}

// Middleware REST<->SOAP for NCC microservice
//
//	@title			NCC API
//	@version		1.0
//	@description	Manage and book NCCs.
//	@host			localhost:8080
//	@BasePath		/
func main() {
	ncc_backend_api, ok := os.LookupEnv("NCC_BACKEND_API")
	if !ok {
		ncc_backend_api = "http://localhost"
	}
	router := gin.Default()
	httpClient := &http.Client{}
	soapClient, err := gosoap.SoapClient(ncc_backend_api+":8001/WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}
	soap = soapClient
	router.POST("/addNCC", add)
	router.GET("/getNCC", get)
	router.GET("/getNCC/:id", getId)
	router.POST("/book", book)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8089")
}
