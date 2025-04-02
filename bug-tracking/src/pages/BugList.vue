<template>
  <div class="bug-list">
    <div class="header">
      <h1>Bug List</h1>
      <router-link v-if="authStore.canReportBugs" to="/bugs/new" class="btn btn-primary">
        Report New Bug
      </router-link>
    </div>

    <div class="filters">
      <div class="search">
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Search bugs..."
          class="form-control"
        />
      </div>
      <div class="status-filter">
        <select v-model="statusFilter" class="form-control">
          <option value="">All Statuses</option>
          <option value="open">Open</option>
          <option value="in-progress">In Progress</option>
          <option value="resolved">Resolved</option>
          <option value="closed">Closed</option>
        </select>
      </div>
    </div>

    <!-- Developer Selection Modal -->
    <div v-if="showAssignModal" class="modal-overlay">
      <div class="modal-content">
        <h3>Assign Bug to Developer</h3>
        <div class="form-group">
          <label for="developer">Select Developer:</label>
          <select 
            id="developer" 
            v-model="selectedDeveloper" 
            class="form-control"
          >
            <option value="">Choose a developer...</option>
            <option 
              v-for="dev in developers" 
              :key="dev.id" 
              :value="dev.id"
            >
              {{ dev.name }}
            </option>
          </select>
        </div>
        <div class="modal-actions">
          <button 
            @click="showAssignModal = false" 
            class="btn btn-secondary"
          >
            Cancel
          </button>
          <button 
            @click="confirmAssignBug" 
            class="btn btn-primary"
            :disabled="!selectedDeveloper"
          >
            Assign
          </button>
        </div>
      </div>
    </div>

    <div class="bugs-grid">
      <div v-for="bug in filteredBugs" :key="bug.id" class="bug-card">
        <div class="bug-header">
          <h3>{{ bug.title }}</h3>
          <span :class="['status', bug.status]">{{ bug.status }}</span>
        </div>
        <p class="description">{{ bug.description }}</p>
        <div class="bug-footer">
          <div class="assignee">
            <span v-if="bug.assignee">Assigned to: {{ bug.assignee }}</span>
            <span v-else>Unassigned</span>
          </div>
          <div class="actions">
            <router-link :to="'/bugs/' + bug.id" class="btn btn-secondary">
              View Details
            </router-link>
            <button
              v-if="authStore.canAssignBugs && !bug.assignee"
              @click="assignBug(bug.id)"
              class="btn btn-primary"
            >
              Assign
            </button>
            <button
              v-if="authStore.canEditBugStatus && bug.assignee === authStore.currentUser?.name"
              @click="updateBugStatus(bug.id)"
              class="btn btn-primary"
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
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();
const searchQuery = ref('');
const statusFilter = ref('');
const showAssignModal = ref(false);
const selectedDeveloper = ref('');
const developers = ref([]);
const selectedBugId = ref(null);

// Fetch developers when component mounts
onMounted(async () => {
  if (authStore.canAssignBugs) {
    await fetchDevelopers();
  }
});

// Mock data - replace with actual API calls
const bugs = ref([
  {
    id: 1,
    title: 'Login button not working',
    description: 'The login button on the homepage is not responding to clicks',
    status: 'open',
    assignee: null
  },
  // Add more mock bugs...
]);

const filteredBugs = computed(() => {
  let filtered = bugs.value;
  
  // Apply search filter
  if (searchQuery.value) {
    filtered = filtered.filter(bug => 
      bug.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      bug.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
  }
  
  // Apply status filter
  if (statusFilter.value) {
    filtered = filtered.filter(bug => bug.status === statusFilter.value);
  }

  // Filter by user role
  if (authStore.currentUser?.role === 'developer') {
    // Developers only see bugs assigned to them
    filtered = filtered.filter(bug => bug.assignee === authStore.currentUser.name);
  }

  return filtered;
});

const fetchDevelopers = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/users/developers', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch developers');
    }
    
    developers.value = await response.json();
  } catch (error) {
    console.error('Error fetching developers:', error);
    // You might want to show an error message to the user here
  }
};

const assignBug = (bugId) => {
  selectedBugId.value = bugId;
  showAssignModal.value = true;
};

const confirmAssignBug = async () => {
  if (!selectedDeveloper.value || !selectedBugId.value) return;

  try {
    const response = await fetch(`http://localhost:8080/api/bugs/${selectedBugId.value}/assign`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({
        developerId: selectedDeveloper.value
      })
    });

    if (!response.ok) {
      throw new Error('Failed to assign bug');
    }

    // Update the local bug data
    const updatedBug = await response.json();
    const bugIndex = bugs.value.findIndex(b => b.id === selectedBugId.value);
    if (bugIndex !== -1) {
      bugs.value[bugIndex] = updatedBug;
    }

    // Reset the modal state
    showAssignModal.value = false;
    selectedDeveloper.value = '';
    selectedBugId.value = null;
  } catch (error) {
    console.error('Error assigning bug:', error);
    // You might want to show an error message to the user here
  }
};

const updateBugStatus = async (bugId) => {
  // Implement status update logic
  console.log('Updating status for bug:', bugId);
};
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
  align-items: center;
  margin-bottom: 1rem;
}

.status {
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.875rem;
  text-transform: capitalize;
}

.status.open {
  background-color: #ffebee;
  color: #c62828;
}

.status.in-progress {
  background-color: #fff3e0;
  color: #ef6c00;
}

.status.resolved {
  background-color: #e8f5e9;
  color: #2e7d32;
}

.status.closed {
  background-color: #eceff1;
  color: #455a64;
}

.description {
  color: #666;
  margin-bottom: 1rem;
}

.bug-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.assignee {
  color: #666;
  font-size: 0.875rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  text-decoration: none;
  font-size: 0.875rem;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.form-control {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-content h3 {
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #666;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
  