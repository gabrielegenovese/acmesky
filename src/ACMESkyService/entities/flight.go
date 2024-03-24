package entities

type Flight struct {
	FlightID             string  `json:"flight_id"`
	FlightCompanyID      int64   `json:"flight_company_id"`
	AirportOriginID      string  `json:"airport_origin_id"`
	AirportDestinationID string  `json:"airport_destination_id"`
	DepartDatetime       string  `json:"depart_datetime"`
	ArrivalDatetime      string  `json:"arrival_datetime"`
	FlightPrice          float32 `json:"flight_price"`
	AvailableSeats       uint    `json:"available_seats_count"`
}

func FlightFromMapFromMap(m map[string]interface{}) Flight {
	v := Flight{
		FlightID:             m["flight_id"].(string),
		FlightCompanyID:      int64(m["flight_company_id"].(float64)),
		AirportOriginID:      m["airport_id_origin"].(string),
		AirportDestinationID: m["airport_id_destination"].(string),
		DepartDatetime:       m["depart_datetime"].(string),
		ArrivalDatetime:      m["arrival_datetime"].(string),
		FlightPrice:          float32(m["travel_max_price"].(float64)),
		AvailableSeats:       uint(m["travel_seats_count"].(float64)),
	}
	return v
}
