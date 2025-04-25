<template>
  <div class="bug-form-container">
    <h2 class="page-title">Report a New Bug</h2>

    <form @submit.prevent="submitBug" class="bug-form">
      <!-- Bug Title -->
      <div class="form-group">
        <label class="form-label">Title</label>
        <input v-model="bug.title" type="text" class="form-control" required />
      </div>

      <!-- Bug Description -->
      <div class="form-group">
        <label class="form-label">Description</label>
        <textarea v-model="bug.description" class="form-control" rows="4" required></textarea>
      </div>

      <!-- Priority Selection -->
      <div class="form-group">
        <label class="form-label">Priority</label>
        <select v-model="bug.priority" class="form-control">
          <option value="Low">Low</option>
          <option value="Medium">Medium</option>
          <option value="High">High</option>
          <option value="Critical">Critical</option>
        </select>
      </div>

      <!-- Assignee Selection - Only for managers and admins -->
      <div v-if="authStore.canAssignBugs" class="form-group">
        <label class="form-label">Assign To</label>
        <select v-model="bug.assignedTo" class="form-control">
          <option value="">Unassigned</option>
          <option v-for="dev in developers" :key="dev.id" :value="dev.id">
            {{ dev.name }}
          </option>
        </select>
      </div>

      <!-- Submit Button -->
      <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
        {{ isSubmitting ? 'Submitting...' : 'Submit Bug' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const bug = ref({
  title: '',
  description: '',
  priority: 'low',
  status: 'open',
  reporter: '',
  assignedTo: null
});

const developers = ref([]);
const isSubmitting = ref(false);

// Fetch developers list if user can assign bugs
onMounted(async () => {
  if (authStore.canAssignBugs) {
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
    }
  }
});

const submitBug = async () => {
  if (isSubmitting.value) return;

  isSubmitting.value = true;
  try {
    // Create a new object with only the required fields
    const bugData = {
      title: bug.value.title,
      description: bug.value.description,
      priority: bug.value.priority.toLowerCase()
    };

    const response = await fetch('http://localhost:8080/api/bugs', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify(bugData)
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to submit bug');
    }

    const result = await response.json();
    alert('Bug submitted successfully!');
    
    // Reset form
    bug.value = {
      title: '',
      description: '',
      priority: 'low',
      status: 'open',
      reporter: '',
      assignedTo: null
    };

    // Navigate back to bug list
    router.push('/bugs');
  } catch (error) {
    console.error('Error submitting bug:', error);
    alert(error.message || 'Failed to submit bug. Please try again.');
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped>
.bug-form-container {
  max-width: 42rem;
  margin: 0 auto;
  padding: 2rem;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.bug-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-label {
  display: block;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.form-control {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #0056b3;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
  