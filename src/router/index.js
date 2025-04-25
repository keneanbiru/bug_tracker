import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/bugs'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../pages/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/bugs',
    name: 'BugList',
    component: () => import('../pages/BugList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/bugs/new',
    name: 'ReportBug',
    component: () => import('../pages/ReportBug.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['developer', 'manager', 'admin']
    }
  },
  {
    path: '/bugs/:id',
    name: 'BugDetails',
    component: () => import('../pages/BugDetails.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiredRoles = to.meta.roles

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }

  if (requiredRoles && !requiredRoles.includes(authStore.role)) {
    next('/bugs')
    return
  }

  if (to.path === '/login' && authStore.isAuthenticated) {
    next('/bugs')
    return
  }

  next()
})

export default router 