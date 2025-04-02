<template>
  <div v-if="hasAccess">
    <slot></slot>
  </div>
  <div v-else class="access-denied">
    <h2>Access Denied</h2>
    <p>You don't have permission to access this page.</p>
    <router-link to="/dashboard" class="btn btn-primary">Go to Dashboard</router-link>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useAuthStore } from '../stores/auth';

const props = defineProps({
  allowedRoles: {
    type: Array,
    required: true
  }
});

const authStore = useAuthStore();
const hasAccess = computed(() => props.allowedRoles.includes(authStore.role));
</script>

<style scoped>
.access-denied {
  text-align: center;
  padding: 2rem;
}

.access-denied h2 {
  color: #dc3545;
  margin-bottom: 1rem;
}

.access-denied p {
  color: #666;
  margin-bottom: 1.5rem;
}

.btn-primary {
  background-color: #007bff;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  text-decoration: none;
  display: inline-block;
}

.btn-primary:hover {
  background-color: #0056b3;
}
</style> 