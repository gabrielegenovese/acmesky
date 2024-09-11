<script setup lang="ts"></script>

<template>
	<dialog
		id="useOfferCodeDialog"
		class="w-1/3 rounded p-4 backdrop:bg-gray-900/75"
	>
		<form class="flex flex-col gap-2" action="offer">
			<label for="prontogram-id">Enter your offer code</label>
			<input
				type="text"
				id="offerCode"
				class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
				name="offerCode"
				required
			/>
			<input
				class="mx-auto rounded bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
				type="submit"
				value="Use offer code"
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
