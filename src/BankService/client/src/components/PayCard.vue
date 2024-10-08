<script setup lang="ts">
import { ref, watch } from "vue";
import type { Payment } from "../types/Payment";
import { useRoute } from "vue-router";

const BASEURL = import.meta.env.VITE_BASEURL;

const route = useRoute();
const data = ref<Payment | null>(null);
const id = route.params.id;
watch(() => route.params.id, fetchData, { immediate: true });

async function fetchData() {
  const ENDPOINT = BASEURL + "/payment/" + id;
  let res = await fetch(ENDPOINT);
  const response = await res.json();
  if (response.res) {
    alert("Error: " + response.res);
    return;
  }
  data.value = response as Payment;
  if (route.query.redirecturi) {
    if (!data.value.id) {
      window.location.replace(route.query.redirecturi as string);
    }
    if (data.value.paid) {
      console.log("already payed");
      window.location.replace(route.query.redirecturi as string);
    }
  } else {
    data.value = null;
  }
}

const pay = async () => {
  const ENDPOINT = BASEURL + "/payment/pay/" + id;
  let res = await fetch(ENDPOINT, {
    method: "POST",
  });
  const data = await res.json();
  console.log(data);
  if (data.res == "Done") {
      window.location.replace(route.query.redirecturi as string);
  } else {
    alert("Error: " + data.res);
  }
};
</script>

<template>
  <div
    v-if="data"
    class="min-w-screen min-h-screen bg-gray-200 flex items-center justify-center px-5 pb-10 pt-16"
  >
    <div
      class="w-full mx-auto rounded-lg bg-white shadow-lg p-5 text-gray-700"
      style="max-width: 600px"
    >
      <div class="w-full pt-1 pb-5">
        <div
          class="bg-indigo-500 text-white overflow-hidden rounded-full w-20 h-20 -mt-16 mx-auto shadow-lg flex justify-center items-center"
        >
          <img class="p-2 rounded-full" src="../assets/acme.jpg" />
        </div>
      </div>
      <div class="mb-10">
        <h1 class="text-center font-bold text-xl uppercase">ACME BANK PAYMENT SERVICE</h1>
      </div>
      <div class="text-lg mb-2">
        <p class="font-bold">Info</p>
        <p class="">Amount: {{ data.amount }}</p>
        <p class="">User: {{ data.user }}</p>
        <p class="">Description: {{ data.description }}</p>
      </div>
      <div class="mb-3 flex -mx-2">
        <div class="px-2">
          <label for="type1" class="flex items-center cursor-pointer">
            <input
              type="radio"
              class="form-radio h-5 w-5 text-indigo-500"
              name="type"
              id="type1"
              checked
            />
            <img
              src="https://leadershipmemphis.org/wp-content/uploads/2020/08/780370.png"
              class="h-8 ml-3"
            />
          </label>
        </div>
        <div class="px-2">
          <label for="type2" class="flex items-center cursor-pointer">
            <input type="radio" class="form-radio h-5 w-5 text-indigo-500" name="type" id="type2" />
            <img
              src="https://www.sketchappsources.com/resources/source-image/PayPalCard.png"
              class="h-8 ml-3"
            />
          </label>
        </div>
      </div>
      <div class="mb-3">
        <label class="font-bold text-sm mb-2 ml-1">Name on card</label>
        <div>
          <input
            class="w-full px-3 py-2 mb-1 border-2 border-gray-200 rounded-md focus:outline-none focus:border-indigo-500 transition-colors"
            placeholder="John Smith"
            type="text"
          />
        </div>
      </div>
      <div class="mb-3">
        <label class="font-bold text-sm mb-2 ml-1">Card number</label>
        <div>
          <input
            class="w-full px-3 py-2 mb-1 border-2 border-gray-200 rounded-md focus:outline-none focus:border-indigo-500 transition-colors"
            placeholder="0000 0000 0000 0000"
            type="text"
          />
        </div>
      </div>
      <div class="mb-3 -mx-2 flex items-end">
        <div class="px-2 w-1/2">
          <label class="font-bold text-sm mb-2 ml-1">Expiration date</label>
          <div>
            <select
              class="form-select w-full px-3 py-2 mb-1 border-2 border-gray-200 rounded-md focus:outline-none focus:border-indigo-500 transition-colors cursor-pointer"
            >
              <option value="01">01 - January</option>
              <option value="02">02 - February</option>
              <option value="03">03 - March</option>
              <option value="04">04 - April</option>
              <option value="05">05 - May</option>
              <option value="06">06 - June</option>
              <option value="07">07 - July</option>
              <option value="08">08 - August</option>
              <option value="09">09 - September</option>
              <option value="10">10 - October</option>
              <option value="11">11 - November</option>
              <option value="12">12 - December</option>
            </select>
          </div>
        </div>
        <div class="px-2 w-1/2">
          <select
            class="form-select w-full px-3 py-2 mb-1 border-2 border-gray-200 rounded-md focus:outline-none focus:border-indigo-500 transition-colors cursor-pointer"
          >
            <option value="2024">2024</option>
            <option value="2025">2025</option>
            <option value="2026">2026</option>
            <option value="2027">2027</option>
            <option value="2028">2028</option>
            <option value="2029">2029</option>
            <option value="2030">2030</option>
            <option value="2031">2031</option>
            <option value="2032">2032</option>
            <option value="2033">2033</option>
            <option value="2034">2034</option>
            <option value="2035">2035</option>
            <option value="2036">2036</option>
            <option value="2037">2037</option>
            <option value="2038">2038</option>
            <option value="2039">2039</option>
          </select>
        </div>
      </div>
      <div class="mb-10">
        <label class="font-bold text-sm mb-2 ml-1">Security code</label>
        <div>
          <input
            class="w-32 px-3 py-2 mb-1 border-2 border-gray-200 rounded-md focus:outline-none focus:border-indigo-500 transition-colors"
            placeholder="000"
            type="text"
          />
        </div>
      </div>
      <div>
        <button
          @click="pay"
          class="block w-full max-w-xs mx-auto bg-indigo-500 hover:bg-indigo-700 focus:bg-indigo-700 text-white rounded-lg px-3 py-3 font-semibold"
        >
          PAY NOW
        </button>
      </div>
    </div>
  </div>

  <dib v-else>Payment not exist or redirecturi not inserted😢</dib>
</template>

<style scoped>
.form-radio {
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  display: inline-block;
  vertical-align: middle;
  background-origin: border-box;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  flex-shrink: 0;
  border-radius: 100%;
  border-width: 2px;
}

.form-radio:checked {
  background-image: url("data:image/svg+xml,%3csvg viewBox='0 0 16 16' fill='white' xmlns='http://www.w3.org/2000/svg'%3e%3ccircle cx='8' cy='8' r='3'/%3e%3c/svg%3e");
  border-color: transparent;
  background-color: currentColor;
  background-size: 100% 100%;
  background-position: center;
  background-repeat: no-repeat;
}

@media not print {
  .form-radio::-ms-check {
    border-width: 1px;
    color: transparent;
    background: inherit;
    border-color: inherit;
    border-radius: inherit;
  }
}

.form-radio:focus {
  outline: none;
}

.form-select {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%23a0aec0'%3e%3cpath d='M15.3 9.3a1 1 0 0 1 1.4 1.4l-4 4a1 1 0 0 1-1.4 0l-4-4a1 1 0 0 1 1.4-1.4l3.3 3.29 3.3-3.3z'/%3e%3c/svg%3e");
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  background-repeat: no-repeat;
  padding-top: 0.5rem;
  padding-right: 2.5rem;
  padding-bottom: 0.5rem;
  padding-left: 0.75rem;
  font-size: 1rem;
  line-height: 1.5;
  background-position: right 0.5rem center;
  background-size: 1.5em 1.5em;
}

.form-select::-ms-expand {
  color: #a0aec0;
  border: none;
}

@media not print {
  .form-select::-ms-expand {
    display: none;
  }
}

@media print and (-ms-high-contrast: active), print and (-ms-high-contrast: none) {
  .form-select {
    padding-right: 0.75rem;
  }
}
</style>
