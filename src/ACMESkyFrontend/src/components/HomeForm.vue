<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { Airport } from "../types/Airport";
const whereValue = ref("");
const airports = ref([] as Airport[]);
function openDialog() {
	const dateStart = (<HTMLInputElement>document.getElementById("dateStart"))
		.value;
	const dateEnd = (<HTMLInputElement>document.getElementById("dateEnd"))
		.value;
	if (dateStart >= dateEnd) {
		alert(
			"You cannot travel in time: your departure must before your return",
		);
	} else if (
		!airports.value.some((airport) => airport.city == whereValue.value)
	) {
		alert("Select a valid airport");
	} else {
		(<HTMLDialogElement>(
			document.getElementById("subscribePreferenceDialog")
		)).showModal();
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
						/>
						<input
							id="dateEnd"
							type="date"
							class="grow rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
							required
						/>
					</div>
				</div>
				<div class="flex flex-col gap-2">
					<label for="where">Where to take off??</label>
					<input
						id="where"
						type="search"
						list="airports"
						name="where"
						autocomplete="off"
						class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
						v-model="whereValue"
						required
					/>
					<datalist id="airports">
						<option
							v-for="airport in airports"
							:key="`airport-${airport.airport_id}`"
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
						required
					/>
				</div>
				<input
					class="rounded bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
					type="submit"
					value="Send me the best offers"
				/>
			</form>
		</div>
	</div>
</template>
