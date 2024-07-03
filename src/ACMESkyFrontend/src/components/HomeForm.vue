<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { Airport } from "../types/Airport";
import SubscribePreferenceDialog from "../components/SubscribePreferenceDialog.vue";

const departCity = ref<string>();
const landCity = ref<string>();
const departAirportID = ref<string>();
const landAirportID = ref<string>();

// already normalized
const departDate = ref<Date>(new Date());
const returnDate = ref<Date>();
const budget = ref<number>();
const seatsCount = ref<number>(1);

const airports = ref([] as Airport[]);
function openDialog() {
	// departDate.value = new Date(departDate.value)
	if (departDate.value && returnDate.value && departDate.value >= returnDate.value ) {
		alert(
			"You cannot travel in time: your departure must before your return",
		);
	} else {
		const airport_origin = airports.value.find(airport => airport.city == departCity.value);
		const airport_destination = airports.value.find(airport => airport.city == landCity.value);

		if (!airport_origin) {
			alert("Select a valid depart airport");
		}
		else if ( !airport_destination ) {
			alert("Select a valid land airport");
		} 
		else {
			departAirportID.value = airport_origin.airport_id;
			landAirportID.value = airport_destination.airport_id;
			
			(<HTMLDialogElement>(
				document.getElementById("subscribePreferenceDialog")
			)).showModal();
		}
	}
}
onMounted(() => {
	fetch(import.meta.env.VITE_SKY_SERVICE_API + "/airports")
		.then((response) => response.json())
		.then((json) => {
			airports.value = json;
		});
});
</script>

<template>
	<div
		class="relative flex h-screen w-full items-center bg-cover bg-fixed bg-center p-8"
		style="background-image: url(/images/airplane-bg.jpg)"
	>
		<div class="flex w-96 flex-col gap-4 rounded-xl bg-white p-4">
			<h1 class="text-center text-3xl font-bold drop-shadow-xl">
				Find your best journey
			</h1>
			<form class="flex flex-col gap-4" @submit.prevent="openDialog">
				<div class="flex flex-col gap-2">
					<label for="when">When you want to relax?</label>
					<div class="flex flex-row gap-2">
						<input
							id="dateStart"
							type="date"
							class="grow rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
							required
							:value="departDate.toISOString().slice(0,10)"
							@input="departDate = new Date($event as any)"
						/>
						<input
							id="dateEnd"
							type="date"
							class="grow rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
							required
							:value="returnDate?.toISOString().slice(0,10)"
							@input="returnDate = new Date($event as any)"
						/>
					</div>
				</div>
				<div class="flex flex-col gap-2">

					<label for="travel_origin">Where to take off ?</label>
					<input
						id="travel_origin"
						type="search"
						list="airports_depart"
						name="travel_origin"
						autocomplete="off"
						class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
						v-model="departCity"
						required
					/>
					<datalist id="airports_depart">
						<option
							v-for="airport in airports"
							:key="`airport-d-${airport.airport_id}`"
							:value="airport.city"
						>
							{{ airport.city }} ({{ airport.name }})
						</option>
					</datalist>
				</div>
				<div class="flex flex-col gap-2">
					<label for="travel_destination">Where to land ?</label>
					<input
						id="travel_destination"
						type="search"
						list="airports_land"
						name="travel_destination"
						autocomplete="off"
						class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
						v-model="landCity"
						required
					/>
					<datalist id="airports_land">
						<option
							v-for="airport in airports"
							:key="`airport-l-${airport.airport_id}`"
							:value="airport.city"
						>
							{{ airport.city }} ({{ airport.name }})
						</option>
					</datalist>
				</div>
				<div class="flex flex-col gap-2">
					<label for="price">How cheap?</label>
					<input
						id="price"
						type="number"
						name="price"
						class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
						v-model="budget"
						required
					/>
				</div>
				<div class="flex flex-col gap-2">
					<label for="seatsCount">How many seats?</label>
					<input
						id="seatsCount"
						type="number"
						min="1"
						step="1"
						name="seatsCount"
						class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
						v-model="seatsCount"
						required
					/>
				</div>
				<SubscribePreferenceDialog
					v-bind:budget="budget"
					v-bind:departDate="(departDate as Date)"
					v-bind:returnDate="(returnDate as Date)"
					v-bind:departAirportID="departAirportID"
					v-bind:landAirportID="landAirportID"
					v-bind:seatsCount="seatsCount"
				/>
				<input
					class="rounded bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
					type="submit"
					value="Send me the best offers"
				/>
			</form>
		</div>
	</div>
</template>
