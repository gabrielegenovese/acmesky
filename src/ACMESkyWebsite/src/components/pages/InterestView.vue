<template>
  <div class="grid h-[80vh] place-items-center">
      <div class="text-4xl">Welcome to ACMESky!</div>
      <div class="text-2xl">Tell us some of your travel preferences!</div>
      <form>
        <Combobox v-model="selectedOriginAirport">
          <label for="airport-origin">Origin Airport</label>
          <ComboboxInput 
            id="airport-origin"
            @change="queryOrigin = $event.target.value" 
            :displayValue="(airport) => airport.name"
          />
          <ComboboxOptions>
            <ComboboxOption
              v-for="airport in filteredOriginAirports"
              :key="airport.airport_id"
              :value="airport"
            >
              {{ airport.name }}
            </ComboboxOption>
          </ComboboxOptions>
        </Combobox>
        
        <Combobox v-model="selectedDestAirport">
          <label for="airport-dest">Destination Airport</label>
          <ComboboxInput 
            id="airport-dest"
            @change="queryDest = $event.target.value"
            :displayValue="(airport) => airport.name"
          />
          <ComboboxOptions>
            <ComboboxOption
              v-for="airport in filteredDestAirports"
              :key="airport.airport_id"
              :value="airport"
            >
              {{ airport.name }}
            </ComboboxOption>
          </ComboboxOptions>
        </Combobox>
      </form>
  </div>
</template>

<script setup>
import { ref, computed, watchEffect} from 'vue'
import {
  Combobox,
  ComboboxInput,
  ComboboxOptions,
  ComboboxOption,
} from '@headlessui/vue'
import {getAirports} from '../../api/Airports'
const selectedOriginAirport = ref(null)
const selectedDestAirport = ref(null)
const queryOrigin = ref('')
const queryDest = ref('')
const filteredOriginAirports =ref([]);
const filteredDestAirports =ref([]);

watchEffect(async () => {
  let airports = await getAirports(queryOrigin.value);
  filteredOriginAirports.value = airports;
})


watchEffect(async () => {
  let airports = await getAirports(queryDest.value);
  filteredDestAirports.value = airports;
})

</script>
