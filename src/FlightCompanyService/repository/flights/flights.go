package flightsRepo

import (
	"database/sql"
	"flightcompany/entities"
	dbClient "flightcompany/repository/db"
	"fmt"
	"time"
)

func GetFlights(airportOriginID string, airportDestID string, departDatetime time.Time, minPassengers int) ([]entities.Flight, error) {
	db := dbClient.GetInstance()
	var flights []entities.Flight
	var rows *sql.Rows
	var err error
	var query string

	if db == nil {
		fmt.Println("ERROR NIL")
	}

	var nextDate = departDatetime
	nextDate = nextDate.Add(time.Hour * 24).Truncate(time.Hour * 24)

	// TODO: edit such that query should consider results on same hour
	rows, err = db.Query(
		"SELECT * FROM Flights"+
			" WHERE (AirportOriginID == ? and AirportDestinationID == ?) "+
			" AND (AvailableSeats >= ?) "+
			" AND (DepartDatime >= ? and DepartDatime <= DATE(?))"+
			" ORDER BY DepartDatetime ASC",
		airportOriginID, airportDestID,
		minPassengers,
		departDatetime.UTC().Format(time.DateTime), nextDate.UTC().Format(time.DateTime),
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
