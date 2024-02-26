package travelPreferenceRepo

import (
	"acmesky/entities"
	dbClient "acmesky/repository/db"
	"fmt"
)

func AddCustomerSubscribtionPreference(pref entities.CustomerFlightSubscription) (int64, error) {

	db, err := dbClient.GetInstance()
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(
		"INSERT INTO TravelPreference (Budget, AirportOriginID, AirportDestinationID, TravelDateStart, TravelDateEnd, ProntogramID) VALUES (?, ?, ?, ?, ?, ?)",
		pref.Budget, pref.AirportOriginID, pref.AirportDestinationID, pref.DateStartISO8601, pref.DateEndISO8601, pref.ProntogramID)
	if err != nil {
		return 0, fmt.Errorf("addCustomerTravelPreference: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addCustomerTravelPreference: %v", err)
	}
	return id, nil
}
