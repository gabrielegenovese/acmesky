package entities

type Flight struct {
	FlightID             string  `json:"flight_id" example:"1"`
	AirportOriginID      string  `json:"airport_origin_id" example:"5"`
	AirportDestinationID string  `json:"airport_destination_id" example:"20"`
	DepartDatetime       string  `json:"depart_datetime" example:"2024-01-01 14:00:00"`
	ArrivalDatetime      string  `json:"arrival_datetime" example:"2024-01-10 05:30:00"`
	FlightPrice          float64 `json:"flight_price" example:"199.99"`
	AvailableSeats       uint    `json:"available_seats_count" example:"2"`
}
