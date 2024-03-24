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

	departDateRangeEnd, parseErr := time.Parse(time.DateTime, pref.DateStartISO8601)
	if parseErr != nil {
		return flights, fmt.Errorf("PARSE_ERROR:" + parseErr.Error())
	}
	departDateRangeEnd = departDateRangeEnd.Add(24 * time.Hour).Truncate(24 * time.Hour)

	res, err := http.Get(
		"http://" + FLIGHT_COMPANY_ADDRESS + "/flights?" + url.Values{
			"origin_airport":       {pref.AirportOriginID},
			"dest_airport":         {pref.AirportDestinationID},
			"passengers_count":     {strconv.FormatUint(uint64(pref.SeatsCount), 10)},
			"start_range_datetime": {pref.DateStartISO8601},
			"end_range_datetime":   {departDateRangeEnd.UTC().Format(time.DateTime)},
		}.Encode())

	if err != nil {
		return flights, fmt.Errorf("CONNECTION_ERROR:" + err.Error())
	}
	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return flights, fmt.Errorf("HTTP_ERROR:" + res.Status)
	}

	decodeErr := json.NewDecoder(res.Body).Decode(&flights)

	if err != nil {
		return flights, fmt.Errorf("PARSE_ERROR:" + decodeErr.Error())
	}
	for i := 0; i < len(flights); i++ {
		flights[i].FlightCompanyID = FlightCompanyID
	}
	return flights, nil
}
