import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";

// Import route components
import Dashboard from "../pages/Dashboard.vue";
import Login from "../pages/Login.vue";
import Register from "../pages/Register.vue";
import BugList from "../pages/BugList.vue";
import ReportBug from "../pages/ReportBug.vue";
import BugDetails from "../pages/BugDetails.vue";
import NotFound from "../pages/NotFound.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: "/register",
    name: "Register",
    component: Register,
    meta: { requiresGuest: true }
  },
  {
    path: "/bugs",
    name: "BugList",
    component: BugList,
    meta: { requiresAuth: true }
  },
  {
    path: "/report",
    name: "ReportBug",
    component: ReportBug,
    meta: { requiresAuth: true }
  },
  {
    path: "/bugs/:id",
    name: "BugDetail",
    component: BugDetails,
    meta: { requiresAuth: true }
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: NotFound
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest);

  if (requiresAuth && !authStore.isAuthenticated) {
    next("/login");
  } else if (requiresGuest && authStore.isAuthenticated) {
    next("/");
  } else {
    next();
  }
});

export default router;
