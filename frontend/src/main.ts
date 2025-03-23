import { createApp } from "vue";
import App from "./App.vue";
import Home from "./views/Home.vue";
import Settings from "./views/Settings.vue";
import Search from "./views/Search.vue";
import Subscriptions from "./views/Subscriptions.vue";
import "./style.css";
import { createMemoryHistory, createRouter } from "vue-router";
import { createI18n } from "vue-i18n";
import messages from "./i18n/i18n.json";

const routes = [
  { path: "/", component: Home },
  { path: "/settings", component: Settings },
  { path: "/search", component: Search },
  { path: "/subscriptions", component: Subscriptions },
];

const router = createRouter({
  history: createMemoryHistory(),
  routes,
});

const i18n = createI18n({
  locale: "fr",
  fallbackLocale: "en",
  messages: messages,
});

createApp(App).use(i18n).use(router).mount("#app");
