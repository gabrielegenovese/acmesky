
import { Airport } from "../models/Airport";

export async function getAirports(query) {

    const url = new URL(`${import.meta.env.VITE_API_BASEURL}/airports`);
    
    if( query && query.length > 0)
        url.searchParams.set("query", query);

    const response = await fetch(
        url,
        {
            method: 'GET',
        }
    )
    return await response.json() as Airport[];
}
