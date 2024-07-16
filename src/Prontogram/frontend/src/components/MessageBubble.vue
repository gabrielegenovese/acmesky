<script setup lang="ts">
const props = defineProps({
	username: String,
	profileImage: String,
	time: String,
});
</script>

<template>
	<div
		class="m-4 flex w-1/3 flex-row gap-2"
		:class="[!props.username ? 'place-self-end' : '']"
	>
		<img class="h-8 self-end rounded-full" :src="props.profileImage" />
		<div
			class="grow flex-col gap-1 rounded-xl p-2"
			:class="[props.username ? 'bg-gray-100' : 'bg-sky-100']"
		>
			<div v-if="props.username" class="font-bold">
				{{ props.username }}
			</div>
			<slot></slot>
			<div class="text-right text-gray-500">{{ props.time }}</div>
		</div>
	</div>
</template>
