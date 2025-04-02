<template>
  <div class="bug-details-container">
    <div class="bug-details-header">
      <h1>Bug Details</h1>
      <button @click="$router.back()" class="back-button">Back to List</button>
    </div>

    <div v-if="bug" class="bug-details-content">
      <div class="bug-info">
        <div class="bug-title-section">
          <h2>{{ bug.title }}</h2>
          <span :class="['status-badge', bug.status]">{{ bug.status }}</span>
        </div>

        <div class="bug-meta">
          <div class="meta-item">
            <span class="meta-label">Reported by:</span>
            <span class="meta-value">{{ bug.reporter }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">Assigned to:</span>
            <span class="meta-value">{{ bug.assignee || 'Unassigned' }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">Severity:</span>
            <span class="meta-value">{{ bug.severity }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">Reported on:</span>
            <span class="meta-value">{{ bug.createdAt }}</span>
          </div>
        </div>

        <div class="bug-description">
          <h3>Description</h3>
          <p>{{ bug.description }}</p>
        </div>

        <div v-if="canUpdateBug" class="bug-actions">
          <div class="form-group">
            <label for="status">Update Status</label>
            <select id="status" v-model="newStatus" class="status-select">
              <option value="open">Open</option>
              <option value="in-progress">In Progress</option>
              <option value="resolved">Resolved</option>
            </select>
          </div>

          <div v-if="authStore.isAdmin" class="form-group">
            <label for="assignee">Assign To</label>
            <select id="assignee" v-model="newAssignee" class="assignee-select">
              <option value="">Select developer</option>
              <option v-for="developer in developers" :key="developer.id" :value="developer.id">
                {{ developer.name }}
              </option>
            </select>
          </div>

          <button @click="updateBug" class="update-button">Update Bug</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// Mock data - replace with API call
const bug = ref(null)
const developers = ref([
  { id: 1, name: 'Mike Johnson' },
  { id: 2, name: 'Sarah Wilson' }
])

const newStatus = ref('')
const newAssignee = ref('')

const canUpdateBug = computed(() => {
  if (!bug.value) return false
  if (authStore.isAdmin) return true
  if (authStore.isDeveloper && bug.value.assignee === authStore.currentUser?.name) return true
  return false
})

onMounted(async () => {
  try {
    // TODO: Replace with actual API call
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    bug.value = {
      id: route.params.id,
      title: 'Login button not working',
      description: 'Users cannot log in using the login button on the homepage. The button appears to be unresponsive when clicked.',
      status: 'open',
      severity: 'high',
      reporter: 'John Doe',
      assignee: null,
      createdAt: '2024-04-02'
    }
    newStatus.value = bug.value.status
    newAssignee.value = bug.value.assignee
  } catch (error) {
    console.error('Error fetching bug details:', error)
    // TODO: Show error message to user
  }
})

const updateBug = async () => {
  try {
    // TODO: Replace with actual API call
    const updateData = {
      status: newStatus.value,
      assignee: newAssignee.value
    }

    console.log('Updating bug:', updateData)
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Update local state
    bug.value.status = newStatus.value
    bug.value.assignee = newAssignee.value

    // TODO: Show success message
  } catch (error) {
    console.error('Error updating bug:', error)
    // TODO: Show error message to user
  }
}
</script>

<style scoped>
.bug-details-container {
  max-width: 1000px;
  margin: 0 auto;
}

.bug-details-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.back-button {
  background-color: #e2e8f0;
  color: #4a5568;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.back-button:hover {
  background-color: #cbd5e0;
}

.bug-details-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.bug-title-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.bug-title-section h2 {
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

.bug-meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 1rem;
  background-color: #f8fafc;
  border-radius: 4px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.meta-label {
  font-size: 0.875rem;
  color: #64748b;
}

.meta-value {
  font-weight: 500;
  color: #1e293b;
}

.bug-description {
  margin-bottom: 2rem;
}

.bug-description h3 {
  margin-bottom: 1rem;
  color: #2c3e50;
}

.bug-description p {
  color: #4a5568;
  line-height: 1.6;
}

.bug-actions {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e2e8f0;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4a5568;
}

.status-select,
.assignee-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  margin-bottom: 1rem;
}

.update-button {
  background-color: #4CAF50;
  color: white;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.update-button:hover {
  background-color: #45a049;
}
</style>
  