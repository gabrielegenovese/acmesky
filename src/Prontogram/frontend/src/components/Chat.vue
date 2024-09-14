<script setup lang="ts">
import { ref } from "vue";
import MessageBubble from "./MessageBubble.vue";
const props = defineProps({
	username: String,
	profileImage: String,
});
const messages = ref([]);
function getMessages() {
	if (localStorage.userId) {
		fetch(
			import.meta.env.VITE_PRONTOGRAM_SERVICE_API +
				"/users/" +
				localStorage.userId +
				"/getMessages",
			{
				method: "POST",
				headers: {
					Accept: "application/json",
					"Content-Type": "application/json",
				},
				body: JSON.stringify({
					userId: localStorage.userId,
					sid: localStorage.sid,
				}),
			},
		)
			.then((response) => response.json())
			.then((json) => {
				if (json.messages) {
					messages.value = json.messages.reverse();
				}
			});
	}
}
setInterval(getMessages, 3000);
</script>

<template>
	<div class="flex h-screen flex-col">
		<div class="flex flex-row items-center gap-4 bg-sky-200 p-4">
			<img class="w-16 rounded-full" :src="props.profileImage" />
			<div class="flex flex-col">
				<h2 class="text-lg font-bold">{{ props.username }}</h2>
				<p>Online</p>
			</div>
		</div>
		<div
			class="flex grow flex-col overflow-y-auto"
			style="background-image: url(/images/doodle.png)"
		>
			<div class="flex grow flex-col-reverse content-end overflow-y-auto">
				<MessageBubble v-for="message in messages" time="13:00">{{
					// @ts-ignore
					message.content
				}}</MessageBubble>
				<MessageBubble time="13:00"
					>Lorem ipsum dolor sit amet, consectetur adipiscing elit,
					sed do eiusmod tempor incididunt ut labore et dolore magna
					aliqua. Ut enim ad minim veniam, quis nostrud exercitation
					ullamco laboris nisi ut aliquip ex ea commodo consequat.
					Duis aute irure dolor in reprehenderit in voluptate velit
					esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
					occaecat cupidatat non proident, sunt in culpa qui officia
					deserunt mollit anim id est laborum.</MessageBubble
				>
				<MessageBubble time="12:50" v-bind="props"
					>Lorem ipsum dolor sit amet, consectetur adipiscing elit,
					sed do eiusmod tempor incididunt ut labore et dolore magna
					aliqua. Ut enim ad minim veniam, quis nostrud exercitation
					ullamco laboris nisi ut aliquip ex ea commodo consequat.
					Duis aute irure dolor in reprehenderit in voluptate velit
					esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
					occaecat cupidatat non proident, sunt in culpa qui officia
					deserunt mollit anim id est laborum.
				</MessageBubble>
				<MessageBubble time="11:40" v-bind="props"
					>Lorem ipsum dolor sit amet, consectetur adipiscing elit,
					sed do eiusmod tempor incididunt ut labore et dolore magna
					aliqua. Ut enim ad minim veniam, quis nostrud exercitation
					ullamco laboris nisi ut aliquip ex ea commodo consequat.
					Duis aute irure dolor in reprehenderit in voluptate velit
					esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
					occaecat cupidatat non proident, sunt in culpa qui officia
					deserunt mollit anim id est laborum.
				</MessageBubble>
				<MessageBubble time="11:32"
					>Lorem ipsum dolor sit amet, consectetur adipiscing elit,
					sed do eiusmod tempor incididunt ut labore et dolore magna
					aliqua. Ut enim ad minim veniam, quis nostrud exercitation
					ullamco laboris nisi ut aliquip ex ea commodo consequat.
					Duis aute irure dolor in reprehenderit in voluptate velit
					esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
					occaecat cupidatat non proident, sunt in culpa qui officia
					deserunt mollit anim id est laborum.</MessageBubble
				>
				<MessageBubble time="10:20" v-bind="props"
					>Lorem ipsum dolor sit amet, consectetur adipiscing elit,
					sed do eiusmod tempor incididunt ut labore et dolore magna
					aliqua. Ut enim ad minim veniam, quis nostrud exercitation
					ullamco laboris nisi ut aliquip ex ea commodo consequat.
					Duis aute irure dolor in reprehenderit in voluptate velit
					esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
					occaecat cupidatat non proident, sunt in culpa qui officia
					deserunt mollit anim id est laborum.
				</MessageBubble>
			</div>
			<div class="flex flex-row p-4">
				<input
					class="grow rounded-l-xl"
					type="text"
					placeholder="Message"
				/>
				<button
					class="rounded-r-xl bg-sky-600 px-4 text-white hover:bg-sky-500"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						class="size-6"
					>
						<title>Send message</title>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5"
						/>
					</svg>
				</button>
			</div>
		</div>
	</div>
</template>
