export interface Flight {
	flight_id: string;
	airport_origin_id: string;
	airport_destination_id: string;
	depart_datetime: string;
	arrival_datetime: string;
	flight_price: number;
	available_seats_count: number;
}
