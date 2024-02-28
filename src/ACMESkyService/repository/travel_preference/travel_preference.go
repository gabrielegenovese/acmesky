package travelPreferenceRepo

import (
	"acmesky/entities"
	dbClient "acmesky/repository/db"
	"fmt"
)

func AddCustomerSubscribtionPreference(pref entities.CustomerFlightSubscription) (int64, error) {

	db := dbClient.GetInstance()

	result, err := db.Exec(
		"INSERT INTO TravelPreferences (AirportOriginID, AirportDestinationID, TravelDateStart, TravelDateEnd, SeatsCount, Budget, ProntogramID) VALUES (?, ?, ?, ?, ?, ?, ?)",
		pref.AirportOriginID, pref.AirportDestinationID, pref.DateStartISO8601, pref.DateEndISO8601, pref.SeatsCount, pref.Budget, pref.ProntogramID)
	if err != nil {
		return 0, fmt.Errorf("[DBERROR] addCustomerTravelPreference: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[DBERROR] addCustomerTravelPreference: %v", err)
	}
	return id, nil
}
