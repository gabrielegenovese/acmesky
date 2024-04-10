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

func AddReservedOffer(travelPreferenceId int64, totalOfferPrice float32, flights []entities.Flight) (int64, error) {
	db := dbClient.GetInstance()

	transaction, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v", err)
	}

	resultOffer, err := transaction.Exec(
		"INSERT INTO ReservedOffers (TravelPreferenceID, TotalOfferPrice) VALUES (?, ?)",
		travelPreferenceId, totalOfferPrice,
	)
	if err != nil {
		abortErr := transaction.Rollback()
		if abortErr != nil {
			return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v %v", abortErr, err)
		}
		return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v", err)
	}

	offerCode, err := resultOffer.LastInsertId()
	if err != nil {
		abortErr := transaction.Rollback()
		if abortErr != nil {
			return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v %v", abortErr, err)
		}
		return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v", err)
	}

	sqlQueryBundle := "INSERT INTO OffersBundles (OfferCode, CompanyFlightID, CompanyID) VALUES "
	vals := []interface{}{}
	for _, f := range flights {
		sqlQueryBundle += "(?, ?, ?),"
		vals = append(vals, offerCode, f.FlightID, f.FlightCompanyID)
	}
	//trim the last ,
	sqlQueryBundle = sqlQueryBundle[0 : len(sqlQueryBundle)-1]

	_, err = transaction.Exec(
		sqlQueryBundle,
		vals...,
	)
	if err != nil {
		abortErr := transaction.Rollback()
		if abortErr != nil {
			return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v %v", abortErr, err)
		}
		return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v", err)
	}

	err = transaction.Commit()
	if err != nil {
		abortErr := transaction.Rollback()
		if abortErr != nil {
			return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v %v", abortErr, err)
		}
		return 0, fmt.Errorf("[DBERROR] AddReservedOffer: %v", err)
	}

	return offerCode, nil
}

func GetReservedOffer(offerCode int64) (entities.ReservedOffer, error) {
	db := dbClient.GetInstance()
	var reservedOffer entities.ReservedOffer

	if db == nil {
		fmt.Println("ERROR NIL")
	}

	rows, err := db.Query(
		"SELECT OfferCode, TravelPreferenceID, TotalOfferPrice, DATE_FORMAT(StartReservationDatetime, '%Y-%m-%d %H:%i:%S'), DATE_FORMAT(EndReservationDatetime, '%Y-%m-%d %H:%i:%S')"+
			" FROM ReservedOffers"+
			" WHERE OfferCode = ?",
		offerCode,
	)

	if err != nil {
		return reservedOffer, fmt.Errorf("GetReservedOffer: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&reservedOffer.OfferCode, &reservedOffer.TravelPreferenceID, &reservedOffer.TotalPrice, &reservedOffer.StartReservationDatetime, &reservedOffer.EndReservationDatetime); err != nil {
			return reservedOffer, fmt.Errorf("GetReservedOffer: %v", err)
		}
		return reservedOffer, nil
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return reservedOffer, fmt.Errorf("GetReservedOffer: NOT_FOUND")
		}
		return reservedOffer, fmt.Errorf("GetReservedOffer: %v", err)
	}
	return reservedOffer, nil
}
