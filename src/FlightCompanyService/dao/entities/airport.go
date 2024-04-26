package entities

type Airport struct {
	AirportID string `json:"airport_id" example:"5"`
	Name      string `json:"name" example:"Aeroporto di Bologna-Guglielmo Marconi"`
	City      string `json:"city" example:"Bologna"`
}
