package entities

type Flight struct {
	FlightID             string  `json:"flight_id"`
	AirportOriginID      string  `json:"airport_origin_id"`
	AirportDestinationID string  `json:"airport_destination_id"`
	DepartDatetime       string  `json:"depart_datetime"`
	ArrivalDatetime      string  `json:"arrival_datetime"`
	FlightPrice          float64 `json:"flight_price"`
	AvailableSeats       uint    `json:"available_seats_count"`
}
