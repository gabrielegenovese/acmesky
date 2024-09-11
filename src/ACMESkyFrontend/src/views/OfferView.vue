<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import Flight from "../components/Flight.vue";
import type { Offer } from "../types/Offer";
import type { Payment } from "../types/Payment";

const offer = ref({} as Offer);
const departFlightPrice = ref(0);
const returnFlightPrice = ref(0);
const totalPrice = computed(() => {
	return departFlightPrice.value + returnFlightPrice.value;
});

function buy() {
	fetch(import.meta.env.VITE_WORKERS_API + "/buyOffer/" + offer.value.offer_code)
		.then((response) => response.json())
		.then((json) => {
			if (json.success) {
			window.location.href =
				json.payment_link + window.origin + "/receipt?payment=" + json.payment_id;
			} else {
				alert("Buy error");
			}
		});
}

onMounted(() => {
	const params = new URLSearchParams(document.location.search);
	fetch(
		import.meta.env.VITE_SKY_SERVICE_API +
			"/offer/" +
			params.get("offerCode"),
	)
		.then((response) => response.json())
		.then((json) => {
			offer.value = json;
		});
});
</script>

<template>
	<div class="mx-auto flex max-w-prose flex-col gap-4 p-16 text-justify">
		<h1 class="mx-auto text-3xl">Checkout</h1>
		<p v-if="offer.travel_preference">
			Welcome {{ offer.travel_preference.customer_prontogram_id }}, here's
			your flights
		</p>
	</div>
	<div class="mx-auto flex max-w-screen-lg flex-col gap-4 p-4 text-justify">
		<Flight
			v-if="offer.depart_flight"
			v-model="departFlightPrice"
			:seats="offer.travel_preference.travel_seats_count"
			:flight="offer.depart_flight.flight_id"
		/>
		<Flight
			v-if="offer.return_flight"
			v-model="returnFlightPrice"
			:seats="offer.travel_preference.travel_seats_count"
			:flight="offer.return_flight.flight_id"
		/>
	</div>
	<div class="mx-auto flex max-w-prose">
		<button
			class="mx-auto rounded bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
			@click="buy"
		>
			Buy now ({{ totalPrice }} €)
		</button>
	</div>
</template>
