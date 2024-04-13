package flights

import (
	"acmesky/dao/entities"
	flight_company "acmesky/services/flights/flight_company"
)

func FetchFlightsByCompanyID(pref entities.CustomerFlightSubscriptionRequest, flightCompanyID int64) ([]entities.Flight, error) {
	switch flightCompanyID {
	case 1:
		{
			return flight_company.FetchFlights(pref)
		}
	// case x:
	//	{
	//		return strategyFetchFlightFrom_X(pref)
	//	}
	default:
		{
			return flight_company.FetchFlights(pref)
		}
	}
}
