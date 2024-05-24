import { createApp } from 'vue'
import { createWebHistory, createRouter } from 'vue-router'
import './style.css'
import App from './App.vue'
import WelcomeView from "./components/pages/WelcomeView.vue";
import InterestView from "./components/pages/InterestView.vue";
import BuyView from "./components/pages/BuyView.vue";

const routes = [
    { path: '/', component: WelcomeView },
    { path: '/interests', component: InterestView },
    { path: '/buy', component: BuyView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

createApp(App)
    .use(router)
    .mount('#app')
