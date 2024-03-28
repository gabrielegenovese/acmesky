package entities

type OfferPart struct {
	OfferCode int64  `json:"offer_code"`
	FlightID  string `json:"flight_id"`
	CompanyID int64  `json:"company_id"`
}

type ReservedOffer struct {
	OfferCode          int64   `json:"offer_code"`
	TotalPrice         float32 `json:"offer_price"`
	TravelPreferenceID int64   `json:"travel_preference_id"`
}
