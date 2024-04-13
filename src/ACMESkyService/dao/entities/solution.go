package entities

type Solution struct {
	DepartFlight Flight `json:"depart_flight"`
	ReturnFlight Flight `json:"return_flight"`
}
