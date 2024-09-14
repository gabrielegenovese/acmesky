<script setup lang="ts">
import { onMounted, ref } from "vue";
import Flight from "../components/Flight.vue";
import NCCCard from "../components/NCCCard.vue";
import type { Airport } from "../types/Airport";
import type { Payment } from "../types/Payment";
import type { Offer } from "../types/Offer";
import type { NCC } from "../types/NCC";

const airports = ref([] as Airport[]);
const payment = ref({} as Payment);
const offer = ref({} as Offer);
const nearestNCC = ref({} as NCC);
const originAirport = ref({} as Airport);
const address = ref("");

function print() {
	window.print();
}

async function searchNCC() {
	await fetch(import.meta.env.VITE_SKY_SERVICE_API + "/airports")
		.then((response) => response.json())
		.then((json) => {
			airports.value = json;
		});
	await fetch(
		import.meta.env.VITE_FLIGHT_COMPANY_API +
			"/flight/" +
			offer.value.depart_flight.flight_id,
	)
		.then((response) => response.json())
		.then((json) => {
			originAirport.value = <Airport>(
				airports.value.find(
					(e) => e.airport_id == json.airport_origin_id,
				)
			);
		});
	fetch(import.meta.env.VITE_WORKERS_API + "/searchNCC", {
		method: "POST",
		headers: {
			Accept: "application/json",
			"Content-Type": "application/json",
		},
		body: JSON.stringify({
			paymentId: payment.value.id,
			prontogramId: offer.value.travel_preference.customer_prontogram_id,
			address: address.value,
			airportAddress: originAirport.value.city, // TODO Address
		}),
	})
		.then((response) => response.json())
		.then((json) => {
			nearestNCC.value = json[0]; // TODO find the nearest
		});
}

onMounted(async () => {
	const params = new URLSearchParams(document.location.search);
	await fetch(
		import.meta.env.VITE_BANK_API + "/payment/" + params.get("payment"),
	)
		.then((response) => response.json())
		.then((json) => {
			payment.value = json;
		});
	fetch(
		import.meta.env.VITE_SKY_SERVICE_API +
			"/offer/" +
			payment.value.description,
	)
		.then((response) => response.json())
		.then((json) => {
			offer.value = json;
		});
});
</script>

<template>
	<div class="mx-auto flex max-w-prose flex-col gap-4 p-16 text-justify">
		<div v-if="payment.paid" class="text-center">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="mx-auto size-12 rounded-full bg-green-500 p-2 text-white"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="m4.5 12.75 6 6 9-13.5"
				/>
			</svg>
			<h1 class="mx-auto text-3xl">Payment completed</h1>
			<p>Paid: {{ payment.amount }} €</p>
		</div>
		<div v-else class="text-center">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="mx-auto size-12 rounded-full bg-red-500 p-2 text-white"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M6 18 18 6M6 6l12 12"
				/>
			</svg>
			<h1 class="mx-auto text-3xl">Payment error</h1>
			<a
				v-if="payment.description"
				:href="'offer?offerCode=' + payment.description"
				class="text-sky-600 hover:text-sky-500"
				>Retry</a
			>
		</div>
	</div>
	<div class="mx-auto flex max-w-screen-lg flex-col gap-4 p-4 text-justify">
		<Flight
			v-if="offer.depart_flight"
			:seats="offer.travel_preference.travel_seats_count"
			:flight="offer.depart_flight.flight_id"
		/>
		<Flight
			v-if="offer.return_flight"
			:seats="offer.travel_preference.travel_seats_count"
			:flight="offer.return_flight.flight_id"
		/>
	</div>
	<div class="mx-auto flex max-w-prose">
		<button
			class="mx-auto flex flex-row gap-2 rounded bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
			@click="print"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="size-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M6.72 13.829c-.24.03-.48.062-.72.096m.72-.096a42.415 42.415 0 0 1 10.56 0m-10.56 0L6.34 18m10.94-4.171c.24.03.48.062.72.096m-.72-.096L17.66 18m0 0 .229 2.523a1.125 1.125 0 0 1-1.12 1.227H7.231c-.662 0-1.18-.568-1.12-1.227L6.34 18m11.318 0h1.091A2.25 2.25 0 0 0 21 15.75V9.456c0-1.081-.768-2.015-1.837-2.175a48.055 48.055 0 0 0-1.913-.247M6.34 18H5.25A2.25 2.25 0 0 1 3 15.75V9.456c0-1.081.768-2.015 1.837-2.175a48.041 48.041 0 0 1 1.913-.247m10.5 0a48.536 48.536 0 0 0-10.5 0m10.5 0V3.375c0-.621-.504-1.125-1.125-1.125h-8.25c-.621 0-1.125.504-1.125 1.125v3.659M18 10.5h.008v.008H18V10.5Zm-3 0h.008v.008H15V10.5Z"
				/>
			</svg>
			Print
		</button>
	</div>
	<div
		v-if="payment.amount > 10"
		class="mx-auto flex max-w-prose flex-col gap-2 p-4 text-center"
	>
		<!-- FIXME -->
		<p class="text-xl">
			You can book a NCC for free!<br />What's your address?
		</p>
		<input
			id="address"
			v-model="address"
			type="text"
			class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
		/>
		<button
			@click="searchNCC"
			class="rounded bg-sky-600 p-2 text-white hover:bg-sky-500"
		>
			Book
		</button>
		<NCCCard v-if="nearestNCC" :ncc="nearestNCC" />
	</div>
</template>
