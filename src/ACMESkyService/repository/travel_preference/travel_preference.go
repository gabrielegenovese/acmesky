package travelPreferenceRepo

import (
	"acmesky/entities"
	dbClient "acmesky/repository/db"
	"database/sql"
	"fmt"
)

func AddCustomerSubscribtionPreference(pref entities.CustomerFlightSubscriptionRequest) (int64, error) {

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

func GetAllCustomerFlightPreferencesNotOutdated() ([]entities.CustomerFlightSubscription, error) {
	db := dbClient.GetInstance()
	var prefs []entities.CustomerFlightSubscription = []entities.CustomerFlightSubscription{}
	var rows *sql.Rows
	var err error
	var query string

	if db == nil {
		fmt.Println("ERROR NIL")
	}

	rows, err = db.Query(
		"SELECT TravelPreferenceID, AirportOriginID, AirportDestinationID, DATE_FORMAT(TravelDateStart, '%Y-%m-%d %H:%i:%S'), DATE_FORMAT(TravelDateEnd, '%Y-%m-%d %H:%i:%S'), SeatsCount, Budget, ProntogramID" +
			" FROM TravelPreferences" +
			" WHERE (TravelPreferenceID NOT IN (SELECT DISTINCT TravelPreferenceID FROM ReservedOffers))" +
			" ORDER BY TravelDateStart ASC",
	)

	if err != nil {
		return nil, fmt.Errorf("GetAllPreferences %q: %v", query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var pref entities.CustomerFlightSubscription
		if err := rows.Scan(&pref.TravelPreferenceID, &pref.AirportOriginID, &pref.AirportDestinationID, &pref.DateStartISO8601, &pref.DateEndISO8601, &pref.SeatsCount, &pref.Budget, &pref.ProntogramID); err != nil {
			return nil, fmt.Errorf("GetAllPreferences %q: %v", query, err)
		}

		prefs = append(prefs, pref)
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return prefs, nil
		}
		return nil, fmt.Errorf("GetAllPreferences %q: %v", query, err)
	}
	return prefs, nil
}
