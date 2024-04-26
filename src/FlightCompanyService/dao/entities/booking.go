package entities

type FlightBooking struct {
	BookingID       int64  `json:"booking_id" example:"101"`
	FlightID        int64  `json:"flight_id" example:"2"`
	CustomerName    string `json:"customer_name" example:"Mario"`
	CustomerSurname string `json:"customer_surname" example:"Rossi"`
	SeatsCount      int    `json:"passengers_count" example:"2"`
}
