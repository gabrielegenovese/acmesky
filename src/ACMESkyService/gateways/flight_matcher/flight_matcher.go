package flight_matcher

import (
	"acmesky/dao/entities"
	flightsDAO "acmesky/dao/impl/flights"
	travelPreferenceDAO "acmesky/dao/impl/travel_preference"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResError struct {
	Error string `json:"error"`
}

type Offer struct {
	TravelPreference         entities.CustomerFlightSubscription `json:"travel_preference"`
	OfferCode                int64   `json:"offer_code"`
	TotalPrice               float32 `json:"offer_price"`
	DepartFlight             entities.Flight `json:"depart_flight"`
	ReturnFlight             entities.Flight `json:"return_flight"`
}

func rest_offers(ctx *gin.Context) {
	// Get customer travel preferences
	prefs, err := travelPreferenceDAO.GetAllCustomerFlightPreferencesNotOutdated()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResError{Error: err.Error()})
	}
	var offers []Offer
	for _, pref := range prefs {
		// Find unoffered flight solutions matching with travel preference
		solutions, fetchErr := flightsDAO.GetSolutionsFromPreference(pref.CustomerFlightSubscriptionRequest)
		if fetchErr != nil {
			ctx.JSON(http.StatusInternalServerError, ResError{Error: fetchErr.Error()})
		}
		for _, solution := range solutions {
			// Prepare Offer bundle for each solution
			flights := []entities.Flight{solution.DepartFlight, solution.ReturnFlight}
			totalPrice := 0.0
			for _, f := range flights {
				totalPrice += f.FlightPrice * float64(pref.SeatsCount)
			}

			offerCode, dbErr := travelPreferenceDAO.AddReservedOffer(pref.TravelPreferenceID, float32(totalPrice), []entities.Flight{solution.DepartFlight, solution.ReturnFlight})
			if dbErr == nil {
				offer, dbErr := travelPreferenceDAO.GetReservedOffer(offerCode)
				if dbErr == nil {
					offers = append(offers, Offer{
						TravelPreference: pref,
						OfferCode: offer.OfferCode,
						TotalPrice: offer.TotalPrice,
						DepartFlight: solution.DepartFlight,
						ReturnFlight: solution.ReturnFlight,
					})
				}
			}
		}
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResError{Error: err.Error()})
	}
	if len(offers) == 0 {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.IndentedJSON(http.StatusOK, offers)
	}
}

func rest_updateFlights(ctx *gin.Context) {
	response, _ := http.Get(os.Getenv("FLIGHT_COMPANY_API") + "/allFlights")
	flights := make([]entities.Flight, 0)
	json.NewDecoder(response.Body).Decode(&flights)
	err := flightsDAO.AddFlights(flights)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, ResError{Error: err.Error()})
	} else {
		ctx.IndentedJSON(http.StatusOK, flights)
	}
}

func rest_offerId(ctx *gin.Context) {
	offerCode, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	offer, err := travelPreferenceDAO.GetReservedOffer(offerCode)
	offerBundle, err := travelPreferenceDAO.GetOfferBundle(offerCode)
	travelPreference, err := travelPreferenceDAO.GetTravelPreference(offer.TravelPreferenceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResError{Error: err.Error()})
	} else {
		ctx.IndentedJSON(http.StatusOK, Offer{
						TravelPreference: travelPreference,
						OfferCode: offer.OfferCode,
						TotalPrice: offer.TotalPrice,
						DepartFlight: offerBundle[0],
						ReturnFlight: offerBundle[1],
					})
	}
}

func Listen(router *gin.RouterGroup) {
	router.GET("/updateFlights", rest_updateFlights)
	router.GET("/offers", rest_offers)
	router.GET("/offer/:id", rest_offerId)
}
