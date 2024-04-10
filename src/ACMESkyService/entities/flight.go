package entities

type Flight struct {
	FlightID             string  `json:"flight_id,omitempty"`
	FlightCompanyID      int64   `json:"flight_company_id,omitempty"`
	AirportOriginID      string  `json:"airport_origin_id,omitempty"`
	AirportDestinationID string  `json:"airport_destination_id,omitempty"`
	DepartDatetime       string  `json:"depart_datetime,omitempty"`
	ArrivalDatetime      string  `json:"arrival_datetime,omitempty"`
	FlightPrice          float64 `json:"flight_price,omitempty"`
	AvailableSeats       uint    `json:"available_seats_count,omitempty"`
}

func FlightFromMapFromMap(m map[string]interface{}) Flight {
	v := Flight{
		FlightID:             m["flight_id"].(string),
		FlightCompanyID:      int64(m["flight_company_id"].(float64)),
		AirportOriginID:      m["airport_id_origin"].(string),
		AirportDestinationID: m["airport_id_destination"].(string),
		DepartDatetime:       m["depart_datetime"].(string),
		ArrivalDatetime:      m["arrival_datetime"].(string),
		FlightPrice:          m["travel_max_price"].(float64),
		AvailableSeats:       uint(m["travel_seats_count"].(float64)),
	}
	return v
}
