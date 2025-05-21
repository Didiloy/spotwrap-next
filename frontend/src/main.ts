import { createApp } from "vue";
import App from "./App.vue";
import Home from "./views/Home.vue";
import Settings from "./views/Settings.vue";
import Search from "./views/Search.vue";
import Subscriptions from "./views/Subscriptions.vue";
import ArtistDetails from "./views/details/ArtistDetails.vue";
import AlbumDetails from "./views/details/AlbumDetails.vue";
import TrackDetails from "./views/details/TrackDetails.vue";

import "./style.css";
import { createMemoryHistory, createRouter } from "vue-router";
import { createI18n } from "vue-i18n";
import messages from "./i18n/i18n.json";
import { createPinia } from "pinia";

const routes = [
  { path: "/", component: Home },
  { path: "/settings", component: Settings },
  { path: "/search/:term?", name: "search", component: Search },
  { path: "/subscriptions", component: Subscriptions },
  { path: "/artist/:id", component: ArtistDetails },
  { path: "/album/:id", component: AlbumDetails },
  { path: "/track/:id", component: TrackDetails },
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

const pinia = createPinia();

createApp(App).use(i18n).use(router).use(pinia).mount("#app");
