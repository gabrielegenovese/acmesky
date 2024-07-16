<script setup lang="ts">
import ChatPreview from "../components/ChatPreview.vue";
import Chat from "../components/Chat.vue";
const chats = [
	{
		username: "ACMESky",
		profileImage: "airplane",
		previewText: "Lorem ipsum",
		selected: true,
	},
	{
		username: "Alice",
		profileImage: "woman",
		previewText: "Hi, news for the airplane tickets?",
	},
	{
		username: "Bob",
		profileImage: "man",
		previewText: "saturday night beer?",
	},
];
</script>

<template>
	<main class="flex flex-row">
		<div class="basis-1/8 flex flex-col gap-4 border-r-2 p-4">
			<h1 class="flex flex-row items-center gap-2 text-xl font-bold">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="size-8"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 0 1-2.555-.337A5.972 5.972 0 0 1 5.41 20.97a5.969 5.969 0 0 1-.474-.065 4.48 4.48 0 0 0 .978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25Z"
					/>
				</svg>
				Prontogram
			</h1>
			<input
				type="search"
				class="rounded-xl"
				placeholder="Search or start new chat"
			/>
			<ChatPreview
				v-for="chat in chats"
				:username="chat.username"
				:profileImage="'/images/' + chat.profileImage + '.jpg'"
				:previewText="chat.previewText"
				:selected="chat.selected"
			/>
		</div>
		<Chat
			username="ACMESKy"
			profileImage="/images/airplane.jpg"
			class="basis-7/8 grow"
		/>
	</main>
</template>
