<script setup lang="ts">
const props = defineProps({
	username: String,
	profileImage: String,
	previewText: String,
	selected: Boolean,
});
</script>

<template>
	<div
		class="flex flex-row items-center gap-2 rounded-xl p-2"
		:class="[props.selected ? 'bg-sky-300' : '']"
	>
		<img class="size-16 rounded-full" :src="props.profileImage" />
		<div class="flex flex-col overflow-hidden whitespace-nowrap">
			<p class="font-bold">{{ props.username }}</p>
			<p class="w-full overflow-hidden text-ellipsis">
				{{ props.previewText }}
			</p>
		</div>
	</div>
</template>
