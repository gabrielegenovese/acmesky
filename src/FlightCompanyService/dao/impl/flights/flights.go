package flightsDAO

import (
	"database/sql"
	dbClient "flightcompany/dao/db"
	"flightcompany/dao/entities"
	"fmt"
	"time"
)

func AddFlight(newFlight entities.Flight) (int64, error) {
	db := dbClient.GetInstance()
	if db == nil {
		fmt.Println("ERROR NIL")
	}
	result, err := db.Exec(
		"INSERT INTO Flights" +
    		" (AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime, AvailableSeats, FlightPrice)" +
			" VALUES" +
			"(?, ?, ?, ?, ?, ?)",
		newFlight.AirportOriginID,
		newFlight.AirportDestinationID,
		newFlight.DepartDatetime,
		newFlight.ArrivalDatetime,
		newFlight.FlightPrice,
		newFlight.AvailableSeats,
	)

	if err != nil {
		return 0, fmt.Errorf("addFlight: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
        return 0, fmt.Errorf("addFlight: %v", err)
    }
    return id, nil
}

func GetFlight(flightId int64) (entities.Flight, error) {
	db := dbClient.GetInstance()
	var flight entities.Flight
	var rows *sql.Rows
	var err error
	var query string

	if db == nil {
		fmt.Println("ERROR NIL")
	}

	rows, err = db.Query(
		"SELECT F.FlightID, AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime, FlightPrice, LeftSeats"+
			" FROM Flights F INNER JOIN FlightCurrentSeats FCS ON F.FlightID = FCS.FlightID"+
			" WHERE F.FlightID = ?",
		flightId,
	)

	if err != nil {
		return flight, fmt.Errorf("flightsByQuery %q: %v", query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var flight entities.Flight
		if err := rows.Scan(&flight.FlightID, &flight.AirportOriginID, &flight.AirportDestinationID, &flight.DepartDatetime, &flight.ArrivalDatetime, &flight.FlightPrice, &flight.AvailableSeats); err != nil {
			return flight, fmt.Errorf("flightsByQuery %q: %v", query, err)
		}
		return flight, nil
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return flight, nil
		}
		return flight, fmt.Errorf("flightsByQuery %q: %v", query, err)
	}
	return flight, nil
}

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

func GetAllFlights() ([]entities.Flight, error) {
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
			" FROM Flights F INNER JOIN FlightCurrentSeats FCS ON F.FlightID = FCS.FlightID",
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

func ReserveFlightBooking(booking entities.FlightBooking) (int64, error) {
	db := dbClient.GetInstance()

	result, err := db.Exec(
		"INSERT INTO FlightBookings (FlightID, SeatsCount, ReservationFlightPrice, CustomerName, CustomerSurname) "+
			" SELECT B.FlightID, B.SeatsCount, F.FlightPrice, B.CustomerName, B.CustomerSurname"+
			" FROM Flights F INNER JOIN ("+
			"   SELECT (?) as FlightID, (?) as SeatsCount, (?) as CustomerName, (?) as CustomerSurname, (?) as BuyerID"+
			" ) B ON F.FlightID = B.FlightID",
		booking.FlightID,
		booking.SeatsCount,
		booking.CustomerName,
		booking.CustomerSurname,
		nil,
	)

	if err != nil {
		return 0, fmt.Errorf("[DBERROR] addBookingReservation: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[DBERROR] addBookingReservation: %v", err)
	}
	return id, nil
}

func ConfirmFlightBooking(bookingID int64) error {
	db := dbClient.GetInstance()

	result, err := db.Exec(
		"UPDATE FlightBookings"+
			" SET BoughtDatetime = NOW()"+
			" WHERE BookingID = ? AND BoughtDatetime IS NULL",
		bookingID,
	)

	if err != nil {
		return fmt.Errorf("[DBERROR] setBookingConfirm: %v", err)
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[DBERROR] setBookingConfirm: %v", err)
	} else if rowsCount < 1 {
		return fmt.Errorf("NOT FOUND %d", bookingID)
	}
	return nil
}

func RemoveFlightBooking(bookingID int64) error {
	db := dbClient.GetInstance()

	result, err := db.Exec(
		"DELETE FROM FlightBookings"+
			" WHERE BookingID = ? AND BoughtDatetime IS NULL",
		bookingID,
	)

	if err != nil {
		return fmt.Errorf("[DBERROR] removeBooking: %v", err)
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[DBERROR] removeBooking: %v", err)
	} else if rowsCount < 1 {
		return fmt.Errorf("NOT FOUND %d", bookingID)
	}
	return nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
