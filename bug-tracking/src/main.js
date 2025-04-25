import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./assets/main.css"; // Import our main CSS file

const app = createApp(App);
const pinia = createPinia();

// Use plugins
app.use(pinia);
app.use(router);

// Mount app
app.mount("#app");
