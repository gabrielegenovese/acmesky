<script setup lang="ts">
function openUseOfferCodeDialog() {
	(<HTMLDialogElement>(
		document.getElementById("useOfferCodeDialog")
	)).showModal();
}
</script>

<template>
	<header
		class="fixed z-10 flex h-14 w-full flex-wrap justify-between backdrop-blur"
	>
		<RouterLink
			to="/"
			class="flex flex-row items-center gap-2 px-5 drop-shadow-sm hover:bg-orange-500"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="size-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5"
				/>
			</svg>
			ACMESky
		</RouterLink>
		<button
			class="px-5 hover:bg-orange-500 md:hidden"
			onclick="document.querySelector('nav').classList.toggle('hidden')"
		>
			<svg
				class="h-12 w-12"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
			>
				<title>Toggle navigation</title>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
				/>
			</svg>
		</button>
		<nav
			class="flex hidden w-full flex-col gap-8 text-center md:flex md:w-auto md:flex-row"
		>
			<RouterLink
				to="/how"
				class="my-auto text-lg font-bold hover:text-orange-500"
				>How it works</RouterLink
			>
			<RouterLink
				to="/about"
				class="my-auto text-lg font-bold hover:text-orange-500"
				>About us</RouterLink
			>
			<RouterLink
				to="/contact"
				class="my-auto text-lg font-bold hover:text-orange-500"
				>Contact</RouterLink
			>
		</nav>
		<button
			class="mx-2 my-2 rounded-full border border-2 bg-sky-500 px-4 font-bold text-white hover:bg-sky-400 hover:drop-shadow"
			@click="openUseOfferCodeDialog"
		>
			Use offer code
		</button>
	</header>
</template>
