<script setup lang="ts">
import { onMounted } from "vue";
function login(event: Event) {
	const formData = new FormData(event.target as HTMLFormElement);
	fetch(
		import.meta.env.VITE_PRONTOGRAM_SERVICE_API +
			"/auth/" +
			formData.get("username") +
			"/login",
		{
			method: "POST",
			headers: {
				Accept: "application/json",
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				userId: formData.get("username"),
				password: formData.get("password"),
			}),
		},
	).then((response) => {
		if (response.status == 200) {
			response.json().then((json) => {
				(<HTMLDialogElement>(
					document.getElementById("loginDialog")
				)).close();
				localStorage.userId = json.userId;
				localStorage.sid = json.sid;
			});
		} else if (response.status == 401) {
			alert("Incorrect password");
		} else if (response.status == 404) {
			alert("User not found");
		}
	});
}
onMounted(() => {
	if (!localStorage.userId) {
		(<HTMLDialogElement>document.getElementById("loginDialog")).showModal();
	}
});
</script>

<template>
	<dialog id="loginDialog" class="w-1/3 rounded p-4 backdrop:bg-gray-900/75">
		<form class="flex flex-col" @submit.prevent="login">
			<label for="username">Username</label>
			<input
				type="text"
				id="username"
				class="mb-2 rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
				name="username"
				required
			/>
			<label for="password">Password</label>
			<input
				type="password"
				id="password"
				class="mb-2 rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
				name="password"
				required
			/>
			<a href="/create" class="text-sky-600 hover:text-sky-500"
				>Don't have an account? Create one</a
			>
			<input
				class="mx-auto rounded-xl bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
				type="submit"
				value="Login"
			/>
		</form>
	</dialog>
</template>
