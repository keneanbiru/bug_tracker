import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import { createPinia } from "pinia";
import "./assets/main.css"; // Import our main CSS file
import Dashboard from "./pages/Dashboard.vue";
import BugList from "./pages/BugList.vue";
import ReportBug from "./pages/ReportBug.vue";
import BugDetails from "./pages/BugDetails.vue";
import Login from "./pages/Login.vue";
import NotFound from "./pages/NotFound.vue";

// Define routes
const routes = [
  { path: "/", component: Dashboard },
  { path: "/bugs", component: BugList },
  { path: "/report", component: ReportBug },
  { path: "/bug/:id", component: BugDetails },
  { path: "/login", component: Login },
  { path: "/:pathMatch(.*)*", component: NotFound },
];

// Create Pinia instance
const pinia = createPinia();

// Create and mount Vue app
createApp(App).use(router).use(pinia).mount("#app");
