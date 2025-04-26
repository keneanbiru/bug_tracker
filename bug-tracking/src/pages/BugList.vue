<template>
  <div class="bug-list">
    <div class="header">
      <h1>Bug List</h1>
      <router-link v-if="authStore.canReportBugs" to="/report" class="btn btn-primary">
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
          <option value="all">All Statuses</option>
          <option value="OPEN">Open</option>
          <option value="IN_PROGRESS">In Progress</option>
          <option value="RESOLVED">Resolved</option>
          <option value="CLOSED">Closed</option>
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

    <!-- Status Update Modal -->
    <div v-if="showStatusModal" class="modal-overlay">
      <div class="modal-content">
        <h3>Update Bug Status</h3>
        <div class="form-group">
          <label for="status">Select Status:</label>
          <select 
            id="status" 
            v-model="selectedStatus" 
            class="form-control"
          >
            <option value="">Choose a status...</option>
            <option value="OPEN">Open</option>
            <option value="IN_PROGRESS">In Progress</option>
            <option value="RESOLVED">Resolved</option>
            <option value="CLOSED">Closed</option>
          </select>
        </div>
        <div class="modal-actions">
          <button 
            @click="showStatusModal = false" 
            class="btn btn-secondary"
          >
            Cancel
          </button>
          <button 
            @click="confirmUpdateStatus" 
            class="btn btn-primary"
            :disabled="!selectedStatus"
          >
            Update
          </button>
        </div>
      </div>
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
          <span :class="['status', bug.status.toLowerCase()]">{{ formatStatus(bug.status) }}</span>
        </div>
        <p class="description">{{ bug.description }}</p>
        <div class="bug-footer">
          <div class="assignee">
            <span v-if="bug.assigned_to">Assigned to: {{ bug.assigned_to.name }}</span>
            <span v-else class="unassigned">Unassigned</span>
          </div>
          <div class="actions">
            <router-link :to="'/bugs/' + bug.id" class="btn btn-secondary">
              View Details
            </router-link>
            <button
              v-if="authStore.canAssignBugs && !bug.assigned_to"
              @click="assignBug(bug.id)"
              class="btn btn-primary"
              data-test="assign-button"
            >
              Assign
            </button>
            <button
              v-if="authStore.canEditBugStatus && (bug.assigned_to?.id === authStore.currentUser?.id || authStore.isManager || authStore.isAdmin)"
              @click="updateBugStatus(bug.id)"
              class="btn btn-primary"
              data-test="update-status-button"
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
import { useBugStore } from '../stores/bug';

const authStore = useAuthStore();
const bugStore = useBugStore();
const searchQuery = ref('');
const statusFilter = ref('');
const showAssignModal = ref(false);
const selectedDeveloper = ref('');
const developers = ref([]);
const selectedBugId = ref(null);
const showStatusModal = ref(false);
const selectedStatus = ref('');

// Fetch bugs and developers when component mounts
onMounted(async () => {
  console.log('BugList component mounted');
  console.log('Auth token:', authStore.token);
  
  try {
    console.log('Fetching bugs...');
    await bugStore.fetchBugs();
    console.log('Bugs fetched:', bugStore.bugs);
    console.log('Loading state:', bugStore.loading);
    console.log('Error state:', bugStore.error);
    
    if (authStore.canAssignBugs) {
      await fetchDevelopers();
    }
  } catch (error) {
    console.error('Error in onMounted:', error);
  }
});

const fetchDevelopers = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/auth/developers', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch developers');
    }
    
    const data = await response.json();
    developers.value = data.map(dev => ({
      id: dev.id,
      name: dev.name
    }));
  } catch (error) {
    console.error('Error fetching developers:', error);
    alert('Failed to fetch developers. Please try again.');
  }
};

const filteredBugs = computed(() => {
  let filtered = [...bugStore.bugs];
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(
      (bug) =>
        bug.title.toLowerCase().includes(query) ||
        bug.description.toLowerCase().includes(query)
    );
  }

  // Filter by status
  if (statusFilter.value && statusFilter.value !== "all") {
    filtered = filtered.filter((bug) => bug.status === statusFilter.value);
  }

  // Filter by user role
  if (authStore.isDeveloper) {
    // Developers see bugs assigned to them OR unassigned bugs
    filtered = filtered.filter(
      (bug) =>
        !bug.assigned_to ||
        (bug.assigned_to && bug.assigned_to.id === authStore.currentUser?.id)
    );
  }

  return filtered;
});

const assignBug = (bugId) => {
  selectedBugId.value = bugId;
  showAssignModal.value = true;
};

const confirmAssignBug = async () => {
  if (!selectedDeveloper.value || !selectedBugId.value) return;

  try {
    await bugStore.assignBug(selectedBugId.value, selectedDeveloper.value);
    showAssignModal.value = false;
    selectedDeveloper.value = '';
    selectedBugId.value = null;
  } catch (error) {
    console.error('Error assigning bug:', error);
    alert(error.message || 'Failed to assign bug. Please try again.');
  }
};

const updateBugStatus = (bugId) => {
  selectedBugId.value = bugId;
  showStatusModal.value = true;
};

const confirmUpdateStatus = async () => {
  if (!selectedStatus.value || !selectedBugId.value) return;

  try {
    // Find the bug in the store
    const bugIndex = bugStore.bugs.findIndex(bug => bug.id === selectedBugId.value);
    if (bugIndex === -1) {
      throw new Error('Bug not found');
    }

    // Map the status to the format expected by the backend
    const statusMap = {
      'OPEN': 'open',
      'IN_PROGRESS': 'in-progress',
      'RESOLVED': 'resolved',
      'CLOSED': 'closed'
    };
    
    const mappedStatus = statusMap[selectedStatus.value] || selectedStatus.value;
    
    // Make API call to update the bug status in the database
    const response = await fetch(`http://localhost:8080/api/bugs/${selectedBugId.value}/status`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({ status: mappedStatus })
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to update bug status');
    }
    
    // Get the updated bug from the response
    const updatedBug = await response.json();
    
    // Update the bug in the store with the response data
    bugStore.bugs[bugIndex] = updatedBug;
    
    console.log('Bug status updated in database:', updatedBug);
    
    // Close the modal and reset the form
    showStatusModal.value = false;
    selectedStatus.value = '';
    selectedBugId.value = null;
    
    // Show success message
    alert('Bug status updated successfully!');
  } catch (error) {
    console.error('Error updating bug status:', error);
    alert(error.message || 'Failed to update bug status. Please try again.');
  }
};

const formatStatus = (status) => {
  switch (status) {
    case "OPEN":
      return "Open";
    case "IN_PROGRESS":
      return "In Progress";
    case "RESOLVED":
      return "Resolved";
    case "CLOSED":
      return "Closed";
    default:
      return status;
  }
};
</script>

<style scoped>
.bug-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filters {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
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
  gap: 20px;
}

.bug-card {
  background: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.bug-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.bug-header h3 {
  margin: 0;
  font-size: 1.2em;
}

.status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.9em;
}

.status.open {
  background-color: #ffebee;
  color: #c62828;
}

.status.in_progress {
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
  margin-bottom: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.bug-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.assignee {
  font-size: 0.9em;
  color: #666;
}

.unassigned {
  color: #999;
  font-style: italic;
}

.actions {
  display: flex;
  gap: 10px;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}

.error {
  text-align: center;
  padding: 20px;
  color: #c62828;
  background: #ffebee;
  border-radius: 4px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 100%;
  max-width: 500px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}
</style>
  