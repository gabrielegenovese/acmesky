<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
	departDate: Date | undefined,
	returnDate: Date | undefined,
	departAirportID: string | undefined,
	landAirportID: string | undefined,
	budget: number | undefined,
	seatsCount: number | undefined
}>();
const prontogramID = ref(null);

function subscribePreference() {
	fetch(
		import.meta.env.VITE_SKY_SERVICE_API + "/subscribe",
		{
			method: 'POST',
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				customer_prontogram_id: prontogramID.value,
				airport_id_origin: props.departAirportID,
				airport_id_destinarion: props.landAirportID,
				travel_date_start: props.departDate?.toISOString(),
				travel_date_end: props.returnDate?.toISOString(),
				travel_max_price: props.budget,
				travel_seats_count: props.seatsCount,
			})
		}
	)
}

</script>

<template>
	<dialog
		id="subscribePreferenceDialog"
		class="w-1/3 rounded p-4 backdrop:bg-gray-900/75"
	>
		<form class="flex flex-col gap-2" @submit.prevent="subscribePreference()">
			<label for="prontogram-id">What is your Prontogram ID? We will notify you the best offers</label>
			<input
				type="text"
				id="prontogram-id"
				class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
				name="prontogram-id"
				v-model="prontogramID"
				required
			/>
			<input
				class="mx-auto rounded bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
				type="submit"
				value="Subscribe"
			/>
		</form>
		<form method="dialog" class="absolute right-10 top-0">
			<button class="fixed m-2 rounded-full bg-gray-200">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="h-6"
				>
					<title>Close</title>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		</form>
	</dialog>
</template>
