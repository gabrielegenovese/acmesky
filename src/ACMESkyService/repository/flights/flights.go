package flightsRepo

import (
	"acmesky/entities"
	dbClient "acmesky/repository/db"
	"fmt"
)

func AddFlights(flights []entities.Flight) error {

	db := dbClient.GetInstance()

	sqlStr := "INSERT INTO Flights" +
		" (CompanyFlightID, CompanyID, CompanyFlightPrice, AvailableSeats, AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime)" +
		" VALUES "
	vals := []interface{}{}

	for _, f := range flights {
		sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, f.FlightID, f.FlightCompanyID, f.FlightPrice, f.AvailableSeats, f.AirportOriginID, f.FlightPrice, f.AirportDestinationID, f.DepartDatetime, f.ArrivalDatetime)
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
