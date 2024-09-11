<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { Airport } from "../types/Airport";
import type { Flight } from "../types/Flight";

const price = defineModel();

const props = defineProps<{
	flight: string;
	seats: number;
}>();
const completeDateTime = new Intl.DateTimeFormat("en-GB", {
	day: "numeric",
	month: "long",
	year: "numeric",
	hour: "numeric",
	minute: "numeric",
});

const airports = ref([] as Airport[]);
const flight = ref({} as Flight);
const originAirport = ref({} as Airport);
const destinationAirport = ref({} as Airport);

onMounted(async () => {
	await fetch(import.meta.env.VITE_SKY_SERVICE_API + "/airports")
		.then((response) => response.json())
		.then((json) => {
			airports.value = json;
		});
	fetch(import.meta.env.VITE_FLIGHT_COMPANY_API + "/flight/" + props.flight)
		.then((response) => response.json())
		.then((json) => {
			flight.value = json;
			originAirport.value = <Airport>(
				airports.value.find(
					(e) => e.airport_id == flight.value.airport_origin_id,
				)
			);
			destinationAirport.value = <Airport>(
				airports.value.find(
					(e) => e.airport_id == flight.value.airport_destination_id,
				)
			);
			price.value = flight.value.flight_price * props.seats;
		});
});
</script>

<template>
	<div
		v-if="flight.flight_id"
		class="flex flex-col gap-2 rounded-xl border border-gray-500 p-4 shadow"
	>
		<div class="flex w-full flex-row items-center justify-around">
			<div class="flex basis-1/2 flex-col">
				<p>From</p>
				<p v-if="originAirport.airport_id" class="font-bold">
					{{ originAirport.name }}
				</p>
				<p v-if="originAirport.airport_id" class="font-bold">
					{{ originAirport.city }}
				</p>
				<p>
					{{
						completeDateTime.format(
							new Date(flight.depart_datetime),
						)
					}}
				</p>
			</div>
			<div>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="size-12"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5"
					/>
				</svg>
			</div>
			<div class="flex basis-1/2 flex-col text-right">
				<p>To</p>
				<p class="font-bold">{{ destinationAirport.name }}</p>
				<p class="font-bold">{{ destinationAirport.city }}</p>
				<p>
					{{
						completeDateTime.format(
							new Date(flight.arrival_datetime),
						)
					}}
				</p>
			</div>
		</div>
		<div class="flex flex-row justify-around">
			<p>Seats: {{ props.seats }}</p>
			<p>
				Price: <span class="font-bold">{{ price }} €</span>
			</p>
		</div>
	</div>
</template>
