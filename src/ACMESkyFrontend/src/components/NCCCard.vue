<script setup lang="ts">
import type { NCC } from "../types/NCC";

const price = defineModel();

const props = defineProps<{
	ncc: NCC;
}>();
</script>

<template>
	<div
		class="mx-auto flex w-48 flex-col gap-2 rounded-xl border border-gray-500 p-4 shadow"
	>
		<p class="text-xl font-bold">{{ props.ncc.name }}</p>
		<p>{{ props.ncc.location }}</p>
		<div class="mx-auto flex flex-row gap-2">
			<span class="line-through">{{ props.ncc.price }} €</span>
			<span class="text-red-500">Free!</span>
		</div>
	</div>
</template>
