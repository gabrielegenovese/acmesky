package airportsDAO

import (
	dbClient "acmesky/dao/db"
	"acmesky/dao/entities"
	"database/sql"
	"fmt"
)

func GetAirports(query string) ([]entities.Airport, error) {
	db := dbClient.GetInstance()
	var airports []entities.Airport
	var rows *sql.Rows
	var err error

	if db == nil {
		fmt.Println("ERROR NIL")
	}
	if len(query) > 0 {
		rows, err = db.Query("SELECT * FROM Airports WHERE Name LIKE %?% OR City LIKE %?% ORDER BY Name ASC", query)
	} else {
		rows, err = db.Query("SELECT * FROM Airports ORDER BY Name ASC")
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

func GetAirportsById(ids []string) ([]entities.Airport, error) {
	db := dbClient.GetInstance()
	var airports []entities.Airport
	var rows *sql.Rows
	var err error

	if db == nil {
		fmt.Println("ERROR NIL")
	}

	if len(ids) > 0 {
		sqlStr := "SELECT * FROM Airports WHERE AirportID IN ("
		vals := []interface{}{}
		for _, id := range ids {
			sqlStr += "?,"
			vals = append(vals, id)
		}

		//trim the last ,
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		sqlStr += ")"

		//prepare the statement
		rows, err = db.Query(sqlStr, vals...)
	} else {
		rows, err = db.Query("SELECT * FROM Airports ORDER BY Name ASC")
	}

	if err != nil {
		return nil, fmt.Errorf("airportsByIDs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var airport entities.Airport
		if err := rows.Scan(&airport.AirportID, &airport.Name, &airport.City); err != nil {
			return nil, fmt.Errorf("airportsByIDs: %v", err)
		}

		airports = append(airports, airport)
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return airports, nil
		}
		return nil, fmt.Errorf("airportsByIDs: %v", err)
	}
	return airports, nil
}
