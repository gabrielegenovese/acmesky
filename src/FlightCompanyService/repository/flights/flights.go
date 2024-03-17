package flightsRepo

import (
	"database/sql"
	"flightcompany/entities"
	dbClient "flightcompany/repository/db"
	"fmt"
	"time"
)

func GetFlights(airportOriginID string, airportDestID string, startDatetime time.Time, endDatetime time.Time, minPassengers int) ([]entities.Flight, error) {
	db := dbClient.GetInstance()
	var flights []entities.Flight = []entities.Flight{}
	var rows *sql.Rows
	var err error
	var query string

	if db == nil {
		fmt.Println("ERROR NIL")
	}

	rows, err = db.Query(
		"SELECT F.FlightID, AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime, FlightPrice, LeftSeats"+
			" FROM Flights F INNER JOIN FlightCurrentSeats FCS ON F.FlightID = FCS.FlightID"+
			" WHERE (AirportOriginID = ? and AirportDestinationID = ?)"+
			" AND (LeftSeats >= ?)"+
			" AND (DepartDatetime >= STR_TO_DATE(?, '%Y-%m-%d %H:%i:%S') and DepartDatetime <= STR_TO_DATE(?, '%Y-%m-%d %H:%i:%S') )"+
			" ORDER BY DepartDatetime ASC",
		airportOriginID, airportDestID,
		minPassengers,
		startDatetime.UTC().Format(time.DateTime),
		endDatetime.UTC().Format(time.DateTime),
	)

	if err != nil {
		return nil, fmt.Errorf("flightsByQuery %q: %v", query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var flight entities.Flight
		if err := rows.Scan(&flight.FlightID, &flight.AirportOriginID, &flight.AirportDestinationID, &flight.DepartDatetime, &flight.ArrivalDatetime, &flight.FlightPrice, &flight.AvailableSeats); err != nil {
			return nil, fmt.Errorf("flightsByQuery %q: %v", query, err)
		}

		flights = append(flights, flight)
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return flights, nil
		}
		return nil, fmt.Errorf("flightsByQuery %q: %v", query, err)
	}
	return flights, nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
