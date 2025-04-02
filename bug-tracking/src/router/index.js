import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";

const routes = [
  {
    path: "/",
    redirect: "/dashboard",
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("../pages/Dashboard.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/bugs",
    name: "BugList",
    component: () => import("../pages/BugList.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/bugs/new",
    name: "ReportBug",
    component: () => import("../pages/ReportBug.vue"),
    meta: {
      requiresAuth: true,
      roles: ["developer", "manager", "admin"],
    },
  },
  {
    path: "/bugs/:id",
    name: "BugDetails",
    component: () => import("../pages/BugDetails.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("../pages/Login.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/register",
    name: "Register",
    component: () => import("../pages/Register.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("../pages/NotFound.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const allowedRoles = to.meta.roles;

  // If route requires auth and user is not authenticated
  if (requiresAuth && !authStore.isAuthenticated) {
    next("/login");
    return;
  }

  // If route has role restrictions
  if (allowedRoles && !allowedRoles.includes(authStore.role)) {
    next("/dashboard");
    return;
  }

  // If user is authenticated and tries to access login/register
  if (!requiresAuth && authStore.isAuthenticated) {
    next("/dashboard");
    return;
  }

  next();
});

export default router;
