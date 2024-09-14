<script setup lang="ts">
function signup(event: Event) {
	const formData = new FormData(event.target as HTMLFormElement);
	fetch(
		import.meta.env.VITE_PRONTOGRAM_SERVICE_API +
			"/users/" +
			formData.get("username"),
		{
			method: "POST",
			headers: {
				Accept: "application/json",
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				userId: formData.get("username"),
				display_name: formData.get("name"),
				password: formData.get("password"),
			}),
		},
	).then((response) => {
		console.log(response);
		if (response.status == 200) {
			alert("User created");
			window.location.href = "/";
		} else if (response.status == 403) {
			alert("User already exists");
		}
	});
}
</script>

<template>
	<div class="mx-auto max-w-md p-4">
		<h1 class="text-center text-xl font-bold">Create account</h1>
		<form class="flex flex-col" @submit.prevent="signup">
			<label for="name">Name</label>
			<input id="name" name="name" type="text" class="mb-2 rounded-xl" />
			<label for="username">Username</label>
			<input
				type="text"
				id="username"
				class="mb-2 rounded-xl"
				name="username"
				required
			/>
			<label for="password">Password</label>
			<input
				type="password"
				id="password"
				class="mb-2 rounded-xl"
				name="password"
				required
			/>
			<input
				type="submit"
				value="Create account"
				class="rounded-xl bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
			/>
		</form>
	</div>
</template>
