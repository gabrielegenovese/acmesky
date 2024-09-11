export interface Offer {
	travel_preference: {
		customer_prontogram_id: string;
		airport_id_origin: string;
		airport_id_destination: string;
		travel_date_start: string;
		travel_date_end: string;
		travel_max_price: number;
		travel_seats_count: number;
		travel_preference_id: number;
	};
	offer_code: number;
	offer_price: number;
	depart_flight: {
		flight_id: string;
		flight_company_id: number;
	};
	return_flight: {
		flight_id: string;
		flight_company_id: number;
	};
}
