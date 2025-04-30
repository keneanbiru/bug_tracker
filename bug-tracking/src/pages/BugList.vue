<template>
  <MainLayout>
    <div class="bug-list">
      <div v-if="!authStore.isAuthenticated" class="error-message">
        Please log in to view bugs.
      </div>
      <div v-else>
        <div class="page-header">
          <div class="header-content">
            <h1>Bug List</h1>
            <router-link v-if="authStore.canReportBugs" to="/report" class="btn btn-primary">
              <span class="icon">+</span>
              Report New Bug
            </router-link>
          </div>
        </div>

        <div class="filters-container">
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
                <option value="open">Open</option>
                <option value="in-progress">In Progress</option>
                <option value="resolved">Resolved</option>
                <option value="closed">Closed</option>
              </select>
            </div>
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
            <h2>Update Bug Status</h2>
            <div class="form-group">
              <label for="status">New Status:</label>
              <select 
                id="status" 
                v-model="selectedStatus" 
                class="form-control"
                data-testid="status-select"
              >
                <option value="OPEN">Open</option>
                <option value="IN_PROGRESS">In Progress</option>
                <option value="RESOLVED">Resolved</option>
                <option value="CLOSED">Closed</option>
              </select>
            </div>
            <div class="modal-actions">
              <button 
                @click="closeStatusModal" 
                class="btn btn-secondary"
                data-testid="cancel-status-update"
              >
                Cancel
              </button>
              <button 
                @click="updateStatus" 
                class="btn btn-primary"
                data-testid="confirm-status-update"
                :disabled="!selectedStatus"
              >
                Update
              </button>
            </div>
          </div>
        </div>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <div v-if="bugStore.loading" class="loading">
          <div class="spinner"></div>
          <p>Loading bugs...</p>
        </div>
        <div v-else-if="bugStore.error" class="error">
          {{ bugStore.error }}
        </div>
        <div v-else class="bugs-grid">
          <div v-for="bug in filteredBugs" :key="bug.id" class="bug-card">
            <div class="bug-header">
              <h3>{{ bug.title }}</h3>
              <span :class="['status', bug.status ? bug.status.toLowerCase() : '']">{{ formatStatus(bug.status || 'OPEN') }}</span>
            </div>
            <p class="description">{{ bug.description }}</p>
            <div class="bug-footer">
              <div class="assignee">
                <span v-if="bug.assigned_to">Assigned to: {{ bug.assigned_to.name }}</span>
                <span v-else class="unassigned">Unassigned</span>
              </div>
              <div class="actions">
                <router-link :to="'/bugs/' + bug.id" class="btn btn-secondary btn-sm">
                  View Details
                </router-link>
                <button
                  v-if="authStore.canAssignBugs && !bug.assigned_to"
                  @click="assignBug(bug.id)"
                  class="btn btn-primary"
                  data-testid="assign-button"
                >
                  Assign
                </button>
                <button
                  v-if="authStore.canEditBugStatus && bug.assigned_to?.id === authStore.currentUser?.id"
                  @click="updateBugStatus(bug.id)"
                  class="btn btn-primary btn-sm"
                  data-testid="update-status-button"
                >
                  Update Status
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useBugStore } from '../stores/bug';
import api from '../services/api';

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
const errorMessage = ref('');

// Fetch bugs and developers when component mounts
onMounted(async () => {
  console.log('BugList component mounted');
  console.log('Auth token:', authStore.token);
  
  if (!authStore.isAuthenticated) {
    console.log('User not authenticated');
    return;
  }
  
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
    const response = await api.get('/auth/developers');
    developers.value = response.data.map(dev => ({
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
    filtered = filtered.filter((bug) => bug.status.toLowerCase() === statusFilter.value.toLowerCase());
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
  console.log('Assign button clicked for bug:', bugId);
  selectedBugId.value = bugId;
  showAssignModal.value = true;
};

const confirmAssignBug = async () => {
  if (!selectedDeveloper.value || !selectedBugId.value) return;

  console.log('Confirming assign with:', {
    bugId: selectedBugId.value,
    developerId: selectedDeveloper.value
  });

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
  console.log('Opening status update modal for bug:', bugId);
  const bug = bugStore.bugs.find(b => b.id === bugId);
  console.log('Found bug:', bug);
  
  if (bug) {
    // Ensure status is in the correct format
    const validStatuses = ['OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED'];
    selectedStatus.value = validStatuses.includes(bug.status) ? bug.status : 'OPEN';
    selectedBugId.value = bugId;
    console.log('Current bug status:', bug.status);
    showStatusModal.value = true;
  } else {
    console.error('Bug not found:', bugId);
    alert('Bug not found. Please refresh the page and try again.');
  }
};

const updateStatus = async () => {
  if (!selectedStatus.value || !selectedBugId.value) return;

  try {
    console.log('Updating status:', {
      bugId: selectedBugId.value,
      status: selectedStatus.value,
      currentBug: bugStore.bugs.find(b => b.id === selectedBugId.value)
    });

    // Validate status before sending
    const validStatuses = ['OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED'];
    if (!validStatuses.includes(selectedStatus.value)) {
      throw new Error('Invalid status value');
    }

    const result = await bugStore.updateBugStatus(selectedBugId.value, selectedStatus.value);
    console.log('Status update result:', result);

    if (result) {
      closeStatusModal();
    } else {
      throw new Error('Failed to update status');
    }
  } catch (error) {
    console.error('Error updating status:', error);
    const errorMessage = error.response?.data?.message || 
                        error.message || 
                        'Failed to update bug status. Please try again.';
    alert(errorMessage);
  }
};

const closeStatusModal = () => {
  console.log('Closing status update modal');
  showStatusModal.value = false;
  selectedStatus.value = '';
  selectedBugId.value = null;
};

const formatStatus = (status) => {
  if (!status) return 'Unknown';
  
  const statusMap = {
    'OPEN': 'Open',
    'IN_PROGRESS': 'In Progress',
    'RESOLVED': 'Resolved',
    'CLOSED': 'Closed',
    'open': 'Open',
    'in_progress': 'In Progress',
    'resolved': 'Resolved',
    'closed': 'Closed'
  };
  
  return statusMap[status] || status;
};
</script>

<style scoped>
.bug-list {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h1 {
  font-size: 2.5rem;
  color: #2d3748;
  margin: 0;
}

.filters-container {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.filters {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.search {
  flex: 1;
}

.status-filter {
  width: 200px;
}

.form-control {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.form-control:focus {
  outline: none;
  border-color: #4299e1;
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.2);
}

.bugs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.bug-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.bug-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.bug-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.bug-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #2d3748;
  flex: 1;
  margin-right: 1rem;
}

.status {
  padding: 0.5rem 1rem;
  border-radius: 999px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: capitalize;
}

.status.open {
  background-color: #fff3cd;
  color: #856404;
}

.status.in_progress {
  background-color: #cce5ff;
  color: #004085;
}

.status.resolved {
  background-color: #d4edda;
  color: #155724;
}

.status.closed {
  background-color: #e2e8f0;
  color: #4a5568;
}

.description {
  color: #4a5568;
  margin-bottom: 1.5rem;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
}

.bug-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}

.assignee {
  font-size: 0.875rem;
  color: #718096;
}

.unassigned {
  color: #e53e3e;
  font-style: italic;
}

.actions {
  display: flex;
  gap: 0.35rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  border: none;
}

/* Add new styles for small buttons */
.btn-sm {
  padding: 0.35rem 0.6rem;
  font-size: 0.75rem;
  min-width: auto;
  border-radius: 6px;
}

.btn-primary {
  background-color: #4299e1;
  color: white;
}

.btn-primary:hover {
  background-color: #3182ce;
}

.btn-secondary {
  background-color: #e2e8f0;
  color: #4a5568;
}

.btn-secondary:hover {
  background-color: #cbd5e0;
}

.icon {
  margin-right: 0.5rem;
  font-size: 1.25rem;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #718096;
}

.loading .spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top: 4px solid #4299e1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

.error {
  text-align: center;
  padding: 2rem;
  color: #e53e3e;
  background: #fff5f5;
  border-radius: 8px;
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
  padding: 2rem;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-content h2,
.modal-content h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #2d3748;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

select.form-control {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='%234a5568' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  background-size: 1.25rem;
  padding-right: 2.5rem;
}
</style>
  