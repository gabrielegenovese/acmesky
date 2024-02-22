package airportsRepo

import (
	"acmesky/entities"
	dbClient "acmesky/repository/db"
	"database/sql"
	"fmt"
)

func GetAirports(query string) ([]entities.Airport, error) {
	db, err := dbClient.GetInstance()
	var airports []entities.Airport
	if err != nil {
		return airports, nil
	}
	var rows *sql.Rows

	if len(query) > 0 {
		rows, err = db.Query("SELECT * FROM Airports WHERE Name LIKE %?% OR City LIKE %?% ORDER BY Name ASC", query)
	} else {
		rows, err = db.Query("SELECT * FROM Airports ORDER BY Name ASC", query)
	}

	if err != nil {
		return nil, fmt.Errorf("airportsByQuery %q: %v", query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var airport entities.Airport
		if err := rows.Scan(&airport.AirportID, &airport.Name, &airport.City); err != nil {
			return nil, fmt.Errorf("airportsByQuery %q: %v", query, err)
		}

		airports = append(airports, airport)
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return airports, nil
		}
		return nil, fmt.Errorf("airportsByQuery %q: %v", query, err)
	}
	return airports, nil
}
