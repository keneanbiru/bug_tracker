<template>
  <div class="bugs-container">
    <div class="bugs-header">
      <h1>Bug List</h1>
      <router-link v-if="authStore.isAdmin" to="/bugs/new" class="new-bug-button">
        Report New Bug
      </router-link>
    </div>

    <div class="bugs-filters">
      <select v-model="statusFilter" class="filter-select">
        <option value="">All Status</option>
        <option value="open">Open</option>
        <option value="in-progress">In Progress</option>
        <option value="resolved">Resolved</option>
      </select>

      <select v-if="authStore.isAdmin" v-model="assigneeFilter" class="filter-select">
        <option value="">All Assignees</option>
        <option v-for="developer in developers" :key="developer.id" :value="developer.id">
          {{ developer.name }}
        </option>
      </select>
    </div>

    <div class="bugs-list">
      <div v-for="bug in filteredBugs" :key="bug.id" class="bug-card">
        <div class="bug-header">
          <h3>{{ bug.title }}</h3>
          <span :class="['status-badge', bug.status]">{{ bug.status }}</span>
        </div>
        <p class="bug-description">{{ bug.description }}</p>
        <div class="bug-footer">
          <div class="bug-meta">
            <span>Reported by: {{ bug.reporter }}</span>
            <span v-if="bug.assignee">Assigned to: {{ bug.assignee }}</span>
          </div>
          <div class="bug-actions">
            <router-link :to="'/bugs/' + bug.id" class="view-button">View Details</router-link>
            <button 
              v-if="canUpdateBug(bug)" 
              @click="updateBugStatus(bug)"
              class="update-button"
            >
              Update Status
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

// Mock data - replace with API calls
const bugs = ref([
  {
    id: 1,
    title: 'Login button not working',
    description: 'Users cannot log in using the login button on the homepage',
    status: 'open',
    reporter: 'John Doe',
    assignee: null
  },
  {
    id: 2,
    title: 'Database connection error',
    description: 'Getting connection timeout when trying to access the database',
    status: 'in-progress',
    reporter: 'Jane Smith',
    assignee: 'Mike Johnson'
  }
])

const developers = ref([
  { id: 1, name: 'Mike Johnson' },
  { id: 2, name: 'Sarah Wilson' }
])

const statusFilter = ref('')
const assigneeFilter = ref('')

const filteredBugs = computed(() => {
  return bugs.value.filter(bug => {
    const matchesStatus = !statusFilter.value || bug.status === statusFilter.value
    const matchesAssignee = !assigneeFilter.value || bug.assignee === assigneeFilter.value
    return matchesStatus && matchesAssignee
  })
})

const canUpdateBug = (bug) => {
  if (authStore.isAdmin) return true
  if (authStore.isDeveloper && bug.assignee === authStore.currentUser?.name) return true
  return false
}

const updateBugStatus = (bug) => {
  // TODO: Implement status update logic
  console.log('Updating bug status:', bug.id)
}
</script>

<style scoped>
.bugs-container {
  max-width: 1200px;
  margin: 0 auto;
}

.bugs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.new-bug-button {
  background-color: #4CAF50;
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  text-decoration: none;
  transition: background-color 0.2s;
}

.new-bug-button:hover {
  background-color: #45a049;
}

.bugs-filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.filter-select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  min-width: 200px;
}

.bugs-list {
  display: grid;
  gap: 1.5rem;
}

.bug-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.bug-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.bug-header h3 {
  margin: 0;
  color: #2c3e50;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-badge.open {
  background-color: #fef3c7;
  color: #92400e;
}

.status-badge.in-progress {
  background-color: #dbeafe;
  color: #1e40af;
}

.status-badge.resolved {
  background-color: #dcfce7;
  color: #166534;
}

.bug-description {
  color: #4a5568;
  margin-bottom: 1.5rem;
}

.bug-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.bug-meta {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  color: #718096;
  font-size: 0.875rem;
}

.bug-actions {
  display: flex;
  gap: 1rem;
}

.view-button,
.update-button {
  padding: 0.5rem 1rem;
  border-radius: 4px;
  text-decoration: none;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.view-button {
  background-color: #e2e8f0;
  color: #4a5568;
}

.view-button:hover {
  background-color: #cbd5e0;
}

.update-button {
  background-color: #4299e1;
  color: white;
  border: none;
}

.update-button:hover {
  background-color: #3182ce;
}
</style> 