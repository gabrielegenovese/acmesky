<script setup lang="ts">
import { RouterView } from "vue-router";
import Header from "./components/Header.vue";
import UseOfferCodeDialog from "./components/UseOfferCodeDialog.vue";
</script>

<template>
	<Header />
	<RouterView />
	<UseOfferCodeDialog />
	<footer class="flex flex-col p-4 text-gray-600">
		<div class="mx-auto">Copyright 2024 ACMESky</div>
		<div class="mx-auto">P. IVA 00000000</div>
	</footer>
</template>
