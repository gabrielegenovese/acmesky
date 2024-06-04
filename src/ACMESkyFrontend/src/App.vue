<script setup lang="ts">
import { RouterView } from "vue-router";
import Header from "./components/Header.vue";
import UseCouponDialog from "./components/UseCouponDialog.vue";
</script>

<template>
	<Header />
	<RouterView />
	<UseCouponDialog />
	<footer class="flex flex-col p-4 text-gray-600">
		<div class="mx-auto">Copyright 2024 ACMESky</div>
		<div class="mx-auto">P. IVA 00000000</div>
	</footer>
</template>
