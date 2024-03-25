package flightMatcher

import (
	"acmesky/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func FetchFlightsByCompanyID(pref entities.CustomerFlightSubscriptionRequest, flightCompanyID int64) ([]entities.Flight, error) {
	switch flightCompanyID {
	case 1:
		{
			return strategyFetchFlightFrom_FlightCompany(pref)
		}
	// case x:
	//	{
	//		return strategyFetchFlightFrom_X(pref)
	//	}
	default:
		{
			return strategyFetchFlightFrom_FlightCompany(pref)
		}
	}
}

func strategyFetchFlightFrom_FlightCompany(pref entities.CustomerFlightSubscriptionRequest) ([]entities.Flight, error) {
	var FLIGHT_COMPANY_ADDRESS string = os.Getenv("FLIGHT_COMPANY_ADDRESS")
	var FlightCompanyID int64 = 1
	var flights []entities.Flight = []entities.Flight{}

	departDateRangeEnd, parseErrDepart := time.Parse(time.DateTime, pref.DateStartISO8601)
	returnDateRangeEnd, parseErrReturn := time.Parse(time.DateTime, pref.DateEndISO8601)
	if parseErrDepart != nil || parseErrReturn != nil {
		return flights, fmt.Errorf("PARSE_ERROR: %s - %s", parseErrDepart.Error(), parseErrReturn.Error())
	}
	// these are the next days to depart and return dates truncated to the 0:00 of the next day
	departDateRangeEnd = departDateRangeEnd.Add(24 * time.Hour).Truncate(24 * time.Hour)
	returnDateRangeEnd = returnDateRangeEnd.Add(24 * time.Hour).Truncate(24 * time.Hour)
	resDepart, errDepart := http.Get(
		"http://" + FLIGHT_COMPANY_ADDRESS + "/flights?" + url.Values{
			"origin_airport":       {pref.AirportOriginID},
			"dest_airport":         {pref.AirportDestinationID},
			"passengers_count":     {strconv.FormatUint(uint64(pref.SeatsCount), 10)},
			"start_range_datetime": {pref.DateStartISO8601},
			"end_range_datetime":   {departDateRangeEnd.UTC().Format(time.DateTime)},
		}.Encode())

	resReturn, errReturn := http.Get(
		"http://" + FLIGHT_COMPANY_ADDRESS + "/flights?" + url.Values{
			"origin_airport":       {pref.AirportDestinationID},
			"dest_airport":         {pref.AirportOriginID},
			"passengers_count":     {strconv.FormatUint(uint64(pref.SeatsCount), 10)},
			"start_range_datetime": {pref.DateEndISO8601},
			"end_range_datetime":   {returnDateRangeEnd.UTC().Format(time.DateTime)},
		}.Encode())

	if errDepart != nil || errReturn != nil {
		return flights, fmt.Errorf("CONNECTION_ERROR: %s, %s", errDepart.Error(), errReturn.Error())
	}
	var responses []*http.Response = []*http.Response{resDepart, resReturn}
	var tmpFlights []entities.Flight = []entities.Flight{}
	for _, res := range responses {
		defer res.Body.Close()

		if !(200 <= res.StatusCode && res.StatusCode < 300) {
			return flights, fmt.Errorf("HTTP_ERROR:" + res.Status)
		}

		decodeErr := json.NewDecoder(res.Body).Decode(&tmpFlights)

		if decodeErr != nil {
			return flights, fmt.Errorf("PARSE_ERROR:" + decodeErr.Error())
		}

		flights = append(flights, tmpFlights...)
		tmpFlights = []entities.Flight{}
	}

	for i := 0; i < len(flights); i++ {
		flights[i].FlightCompanyID = FlightCompanyID
	}

	return flights, nil
}
