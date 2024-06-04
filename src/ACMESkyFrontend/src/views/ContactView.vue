<template>
	<div class="mx-auto flex max-w-prose flex-col gap-4 p-16 text-justify">
		<h1 class="mx-auto text-3xl">Contact</h1>
		<p>You can contact us for any needs.</p>
		<form class="flex flex-col gap-2">
			<div class="flex flex-col gap-1">
				<label for="name">Your name:</label>
				<input
					type="text"
					id="name"
					name="name"
					class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
					required
				/>
			</div>
			<div class="flex flex-col gap-1">
				<label for="mail">Your e-mail address:</label>
				<input
					type="email"
					id="mail"
					name="mail"
					class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
					required
				/>
			</div>
			<div class="flex flex-col gap-1">
				<label for="message">Your message:</label>
				<textarea
					id="message"
					name="message"
					class="rounded focus:border-gray-500 focus:ring-0 focus:drop-shadow-md"
					required
				></textarea>
			</div>
			<input
				type="submit"
				class="rounded bg-sky-600 p-2 text-white hover:bg-sky-500 hover:drop-shadow-xl"
				value="Send message"
			/>
		</form>
	</div>
</template>
