<template>
  <div class="bug-list">
    <div class="header">
      <h1>Bug Tracker</h1>
      <router-link 
        v-if="authStore.canReportBugs" 
        to="/bugs/new" 
        class="btn btn-primary"
      >
        Report New Bug
      </router-link>
    </div>

    <div class="filters">
      <div class="search">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search bugs..." 
          class="form-control"
        >
      </div>
      <select v-model="statusFilter" class="form-control status-filter">
        <option value="">All Statuses</option>
        <option value="open">Open</option>
        <option value="in-progress">In Progress</option>
        <option value="resolved">Resolved</option>
      </select>
    </div>

    <div v-if="bugStore.loading" class="loading">
      Loading bugs...
    </div>
    <div v-else-if="bugStore.error" class="error">
      {{ bugStore.error }}
    </div>
    <div v-else class="bugs-grid">
      <div v-for="bug in filteredBugs" :key="bug.id" class="bug-card">
        <div class="bug-header">
          <h3>{{ bug.title }}</h3>
          <span :class="['status', bug.status]">{{ bug.status }}</span>
        </div>
        <p class="description">{{ bug.description }}</p>
        <div class="priority">Priority: {{ bug.priority }}</div>
        <div class="meta">
          <div>Reported by: {{ bug.reported_by.name }}</div>
          <div v-if="bug.assigned_to">
            Assigned to: {{ bug.assigned_to.name }}
          </div>
          <div v-else>Unassigned</div>
        </div>
        <div class="actions">
          <button 
            v-if="authStore.canAssignBugs && bug.status !== 'resolved'" 
            @click="openAssignModal(bug)"
            class="btn btn-secondary"
          >
            Assign
          </button>
          <button 
            v-if="canUpdateStatus(bug)" 
            @click="openStatusModal(bug)"
            class="btn btn-primary"
          >
            Update Status
          </button>
        </div>
      </div>
    </div>

    <!-- Assignment Modal -->
    <div v-if="showAssignModal" class="modal">
      <div class="modal-content">
        <h3>Assign Bug</h3>
        <select v-model="selectedDeveloper" class="form-control">
          <option value="">Select Developer</option>
          <option 
            v-for="dev in developers" 
            :key="dev.id" 
            :value="dev.id"
          >
            {{ dev.name }}
          </option>
        </select>
        <div class="modal-actions">
          <button @click="closeAssignModal" class="btn btn-secondary">Cancel</button>
          <button @click="assignBug" class="btn btn-primary">Assign</button>
        </div>
      </div>
    </div>

    <!-- Status Update Modal -->
    <div v-if="showStatusModal" class="modal">
      <div class="modal-content">
        <h3>Update Status</h3>
        <select v-model="selectedStatus" class="form-control">
          <option value="open">Open</option>
          <option value="in-progress">In Progress</option>
          <option value="resolved">Resolved</option>
        </select>
        <div class="modal-actions">
          <button @click="closeStatusModal" class="btn btn-secondary">Cancel</button>
          <button @click="updateStatus" class="btn btn-primary">Update</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useBugStore } from '../stores/bug'

const authStore = useAuthStore()
const bugStore = useBugStore()

const searchQuery = ref('')
const statusFilter = ref('')
const showAssignModal = ref(false)
const showStatusModal = ref(false)
const selectedDeveloper = ref('')
const selectedStatus = ref('')
const selectedBug = ref(null)
const developers = ref([])

// Fetch bugs and developers on component mount
onMounted(async () => {
  await bugStore.fetchBugs()
  if (authStore.canAssignBugs) {
    await fetchDevelopers()
  }
})

const filteredBugs = computed(() => {
  let bugs = bugStore.bugs

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    bugs = bugs.filter(bug => 
      bug.title.toLowerCase().includes(query) ||
      bug.description.toLowerCase().includes(query)
    )
  }

  if (statusFilter.value) {
    bugs = bugs.filter(bug => bug.status === statusFilter.value)
  }

  return bugs
})

const fetchDevelopers = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/users/developers', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    if (!response.ok) throw new Error('Failed to fetch developers')
    developers.value = await response.json()
  } catch (error) {
    console.error('Error fetching developers:', error)
  }
}

const canUpdateStatus = (bug) => {
  if (!authStore.canEditBugStatus) return false
  if (authStore.isDeveloper) {
    return bug.assigned_to?.id === authStore.currentUser?.id
  }
  return true
}

const openAssignModal = (bug) => {
  selectedBug.value = bug
  showAssignModal.value = true
}

const closeAssignModal = () => {
  showAssignModal.value = false
  selectedDeveloper.value = ''
  selectedBug.value = null
}

const assignBug = async () => {
  if (!selectedDeveloper.value || !selectedBug.value) return
  
  try {
    await bugStore.assignBug(selectedBug.value.id, selectedDeveloper.value)
    closeAssignModal()
  } catch (error) {
    console.error('Error assigning bug:', error)
  }
}

const openStatusModal = (bug) => {
  selectedBug.value = bug
  selectedStatus.value = bug.status
  showStatusModal.value = true
}

const closeStatusModal = () => {
  showStatusModal.value = false
  selectedStatus.value = ''
  selectedBug.value = null
}

const updateStatus = async () => {
  if (!selectedStatus.value || !selectedBug.value) return
  
  try {
    await bugStore.updateBugStatus(selectedBug.value.id, selectedStatus.value)
    closeStatusModal()
  } catch (error) {
    console.error('Error updating status:', error)
  }
}
</script>

<style scoped>
.bug-list {
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.search {
  flex: 1;
}

.status-filter {
  width: 200px;
}

.bugs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.bug-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.bug-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 1rem;
}

.status {
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

.status.open {
  background-color: #fff3cd;
  color: #856404;
}

.status.in-progress {
  background-color: #cce5ff;
  color: #004085;
}

.status.resolved {
  background-color: #d4edda;
  color: #155724;
}

.description {
  margin-bottom: 1rem;
}

.priority {
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.meta {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 1rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  width: 100%;
  max-width: 400px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
}
</style> 