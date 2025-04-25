<template>
  <div class="layout">
    <header class="header">
      <div class="logo">
        <router-link to="/">Bug Tracker</router-link>
      </div>
      <nav class="nav">
        <router-link to="/" class="nav-link">Dashboard</router-link>
        <router-link to="/bugs" class="nav-link">Bugs</router-link>
        <router-link v-if="authStore.canReportBugs" to="/report" class="nav-link">Report Bug</router-link>
      </nav>
      <div class="user-menu">
        <span class="user-role">{{ authStore.currentUser?.name }} ({{ authStore.role }})</span>
        <button @click="handleLogout" class="logout-button">Logout</button>
      </div>
    </header>
    <main class="main">
      <slot></slot>
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #2c3e50;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: white;
}

.logo a {
  color: white;
  text-decoration: none;
  font-size: 1.5rem;
  font-weight: bold;
}

.nav {
  display: flex;
  gap: 1.5rem;
}

.nav-link {
  color: #ecf0f1;
  text-decoration: none;
  padding: 0.5rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: #34495e;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-role {
  color: #ecf0f1;
}

.logout-button {
  background-color: #e74c3c;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.logout-button:hover {
  background-color: #c0392b;
}

.main {
  flex: 1;
  padding: 2rem;
  background-color: #f5f5f5;
}
</style>
  