import { createRouter, createWebHistory } from 'vue-router'
import PayCard from "./components/PayCard.vue";

const routes = [
    { path: '/pay/:id', component: PayCard },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})