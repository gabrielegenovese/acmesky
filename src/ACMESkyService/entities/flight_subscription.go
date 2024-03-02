package entities

type CustomerFlightSubscription struct {
	ProntogramID         string `json:"customer_prontogram_id"`
	AirportOriginID      string `json:"airport_id_origin"`
	AirportDestinationID string `json:"airport_id_destination"`
	// start travel date range in ISO 8601 format (with timezone UTC)
	DateStartISO8601 string `json:"travel_date_start"`
	// end travel date range in ISO 8601 format (with timezone UTC)
	DateEndISO8601 string  `json:"travel_date_end"`
	Budget         float32 `json:"travel_max_price"`
	SeatsCount     uint    `json:"travel_seats_count"`
}

func CustomerFlightSubscriptionFromMap(m map[string]interface{}) CustomerFlightSubscription {
	v := CustomerFlightSubscription{
		ProntogramID:         m["customer_prontogram_id"].(string),
		AirportOriginID:      m["airport_id_origin"].(string),
		AirportDestinationID: m["airport_id_destination"].(string),
		DateStartISO8601:     m["travel_date_start"].(string),
		DateEndISO8601:       m["travel_date_end"].(string),
		Budget:               float32(m["travel_max_price"].(float64)),
		SeatsCount:           uint(m["travel_seats_count"].(float64)),
	}
	return v
}

func (v CustomerFlightSubscription) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"airport_id_origin":      v.AirportOriginID,
		"airport_id_destination": v.AirportDestinationID,
		"travel_date_start":      v.DateStartISO8601,
		"travel_date_end":        v.DateEndISO8601,
		"travel_max_price":       v.Budget,
		"customer_prontogram_id": v.ProntogramID,
		"travel_seats_count":     v.SeatsCount,
	}
	return m
}
