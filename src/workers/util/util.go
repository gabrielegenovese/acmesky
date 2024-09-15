package util

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/google/uuid"
)

type User struct {
	UserId string `json:"userId"`
	Sid    string `json:"sid"`
}

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

type CustomerFlightSubscriptionRequest struct {
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

type CustomerFlightSubscription struct {
	CustomerFlightSubscriptionRequest
	TravelPreferenceID int64 `json:"travel_preference_id"`
}

type Offer struct {
	TravelPreference CustomerFlightSubscription `json:"travel_preference"`
	OfferCode        int64                      `json:"offer_code"`
	TotalPrice       float32                    `json:"offer_price"`
	DepartFlight     Flight                     `json:"depart_flight"`
	ReturnFlight     Flight                     `json:"return_flight"`
}

type FlightBooking struct {
	BookingID       int64  `json:"booking_id" example:"101"`
	FlightID        int64  `json:"flight_id" example:"2"`
	CustomerName    string `json:"customer_name" example:"Mario"`
	CustomerSurname string `json:"customer_surname" example:"Rossi"`
	SeatsCount      int    `json:"passengers_count" example:"2"`
}

type Airport struct {
	AirportID string `json:"airport_id"`
	Name      string `json:"name"`
	City      string `json:"city"`
}

type Payment struct {
	ID          uuid.UUID    `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" gorm:"index"`
	User        string       `json:"user"`
	Description string       `json:"description"`
	Amount      float64      `json:"amount"`
	Link        string       `json:"link"`
	Paid        bool         `json:"paid"`
}

type PaymentReq struct {
	User        string  `json:"user"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type BuyResult struct {
	PaymentLink string `json:"payment_link"`
	PaymentID   string `json:"payment_id"`
	Success     bool   `json:"success"`
}

type NCCSearchRequest struct {
	PaymentId      string `json:"paymentId"`
	ProntogramId   string `json:"prontogramId"`
	Address        string `json:"address"`
	AirportAddress string `json:"airportAddress"`
}

type NCCBookingRequest struct {
	NccId string `xml:"nccId" json:"nccId"`
	Name  string `xml:"name" json:"name"`
	Date  string `xml:"date" json:"date"` // Format: dd/MM/yyyy kk:mm:ss
}

type NCC struct {
	Id       string `xml:"id" json:"id"`
	Name     string `xml:"name" json:"name"`
	Price    string `xml:"price" json:"price"`
	Location string `xml:"location" json:"location"`
}

type NCCResponse struct {
	NearestNCC string `ncc:"id"`
	Success    bool   `json:"success"`
}

type DistanceResBody struct {
	Distance string `json:"distance"`
	Value    int    `json:"value"`
	Status   string `json:"status"`
}

type SendMessageParams struct {
	Offer     Offer  `json:"offer"`
	OfferCode string `json:"offerCode" example:"2"`
}

var ZbClient zbc.Client
var Ctx context.Context
var InterestSaved = make(map[string]chan struct{})
var OfferSelected = make(map[string]chan struct{})
var BuyResults = make(map[string]chan BuyResult)
var PaymentCompleted = make(map[string]chan struct{})
var NCCSearchRequests = make(map[string]chan NCCSearchRequest)
var NCCResponses = make(map[string]chan NCCResponse)
var ProntogramUser User

func FailJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}
