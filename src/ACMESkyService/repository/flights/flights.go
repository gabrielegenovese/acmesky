package flightsRepo

import (
	"acmesky/entities"
	dbClient "acmesky/repository/db"
	"database/sql"
	"fmt"
	"time"
)

func AddFlights(flights []entities.Flight) error {

	db := dbClient.GetInstance()

	sqlStr := "REPLACE INTO Flights" +
		" (CompanyFlightID, CompanyID, PassengerFlightPrice, AvailableSeats, AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime)" +
		" VALUES "
	vals := []interface{}{}

	for _, f := range flights {
		sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, f.FlightID, f.FlightCompanyID, f.FlightPrice, f.AvailableSeats, f.AirportOriginID, f.AirportDestinationID, f.DepartDatetime, f.ArrivalDatetime)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	//prepare the statement
	stmt, err := db.Prepare(sqlStr)

	if err != nil {
		return fmt.Errorf("[DBERROR] addFlights: %v", err)
	}

	//format all vals at once
	_, err = stmt.Exec(vals...)

	if err != nil {
		return fmt.Errorf("[DBERROR] addFlights: %v", err)
	}

	return nil
}

func GetSolutionsFromPreference(pref entities.CustomerFlightSubscriptionRequest) ([]entities.Solution, error) {
	db := dbClient.GetInstance()
	var solutions []entities.Solution = []entities.Solution{}
	var rows *sql.Rows
	var err error

	if db == nil {
		fmt.Println("ERROR NIL")
	}

	departDatetimeEnd, parseErrDepart := time.Parse(time.DateTime, pref.DateStartISO8601)
	returnDatetimeEnd, parseErrReturn := time.Parse(time.DateTime, pref.DateEndISO8601)
	if parseErrDepart != nil || parseErrReturn != nil {
		return solutions, fmt.Errorf("PARSE_ERROR: %s - %s", parseErrDepart.Error(), parseErrReturn.Error())
	}
	// these are the next days to depart and return dates truncated to the 0:00 of the next day
	departDatetimeEnd = departDatetimeEnd.Add(24 * time.Hour).Truncate(24 * time.Hour)
	returnDatetimeEnd = returnDatetimeEnd.Add(24 * time.Hour).Truncate(24 * time.Hour)

	rows, err = db.Query(
		"SELECT DEPART_F.CompanyFlightID, DEPART_F.CompanyID, RETURN_F.CompanyFlightID, RETURN_F.CompanyID"+
			" FROM"+
			" (SELECT CompanyFlightID, CompanyID, PassengerFlightPrice, AvailableSeats, AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime FROM Flights WHERE AirportOriginID = ? and AirportDestinationID = ? AND AvailableSeats >= ? AND (STR_TO_DATE(?, '%Y-%m-%d %H:%i:%S') <= DepartDatetime AND DepartDatetime < STR_TO_DATE(?, '%Y-%m-%d %H:%i:%S') ) ) AS DEPART_F"+
			" INNER JOIN"+
			" (SELECT CompanyFlightID, CompanyID, PassengerFlightPrice, AvailableSeats, AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime FROM Flights WHERE AirportOriginID = ? and AirportDestinationID = ? AND AvailableSeats >= ? AND (STR_TO_DATE(?, '%Y-%m-%d %H:%i:%S') <= DepartDatetime AND DepartDatetime < STR_TO_DATE(?, '%Y-%m-%d %H:%i:%S') ) ) AS RETURN_F"+
			" ON DEPART_F.AirportDestinationID = RETURN_F.AirportOriginID"+
			" WHERE ( (DEPART_F.PassengerFlightPrice * ?) + (RETURN_F.PassengerFlightPrice * ?) ) < ? AND RETURN_F.DepartDatetime > DEPART_F.ArrivalDatetime"+
			" ORDER BY DEPART_F.DepartDatetime, RETURN_F.ArrivalDatetime",
		pref.AirportOriginID, pref.AirportDestinationID, pref.SeatsCount, pref.DateStartISO8601, departDatetimeEnd.UTC().Format(time.DateTime),
		pref.AirportDestinationID, pref.AirportOriginID, pref.SeatsCount, pref.DateEndISO8601, returnDatetimeEnd.UTC().Format(time.DateTime),
		pref.SeatsCount, pref.SeatsCount, pref.Budget,
	)

	if err != nil {
		return nil, fmt.Errorf("flightsByPreference: %v", err)
	}
	defer rows.Close()

	var solution entities.Solution
	for rows.Next() {
		var departFlight entities.Flight
		var returnFlight entities.Flight
		if err := rows.Scan(&departFlight.FlightID, &departFlight.FlightCompanyID, &returnFlight.FlightID, &returnFlight.FlightCompanyID); err != nil {
			return nil, fmt.Errorf("flightsByPreference: %v", err)
		}
		solution = entities.Solution{
			DepartFlight: departFlight,
			ReturnFlight: returnFlight,
		}
		solutions = append(solutions, solution)
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return solutions, nil
		}
		return nil, fmt.Errorf("flightsByQuery: %v", err)
	}
	return solutions, nil
}

/*
SELECT DEPART_F.FlightID, DEPART_F.FlightCompanyID, RETURN_F.FlightID, RETURN_F.FlightCompanyID,
FROM
	(SELECT FlightID, FlightCompanyID, FlightPrice, AvailableSeats, AirportOriginID, AirportDestinationID FROM Flights WHERE AirportOriginID = ? and AirportDestinationID = ? AND AvailableSeats >= ? AND DepartDatetime >= ? ) AS DEPART_F
	INNER JOIN
	(SELECT FlightID, FlightCompanyID, FlightPrice, AvailableSeats, AirportOriginID, AirportDestinationID FROM Flights WHERE AirportOriginID = ? and AirportDestinationID = ? AND AvailableSeats >= ? AND DepartDatetime >= ?) AS RETURN_F
	ON DEPART_F.AirportDestinationID = RETURN_F.AirportOriginID
WHERE (DEPART_F.FlightPrice+RETURN_F.FlightPrice) < ? AND RETURN_F.DepartDatetime > DEPART_F.ArrivalDatetime
ORDER BY DEPART_F.DepartDatetime, RETURN_F.ArrivalDatetime
*/
