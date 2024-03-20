package entities

type FlightBooking struct {
	BookingID       int64  `json:"booking_id"`
	FlightID        int64  `json:"flight_id"`
	CustomerName    string `json:"customer_name"`
	CustomerSurname string `json:"customer_surname"`
	SeatsCount      int    `json:"passengers_count"`
}
