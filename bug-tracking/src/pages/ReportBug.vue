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

      <!-- Assignee Selection -->
      <div class="form-group">
        <label class="form-label">Assign To</label>
        <select v-model="bug.assignedTo" class="form-control">
          <option value="">Unassigned</option>
          <option v-for="user in developers" :key="user" :value="user">{{ user }}</option>
        </select>
      </div>

      <!-- Submit Button -->
      <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
        {{ isSubmitting ? 'Submitting...' : 'Submit Bug' }}
      </button>
    </form>
  </div>
</template>

<script>
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

export default {
  setup() {
    const authStore = useAuthStore();
    const router = useRouter();
    return { authStore, router };
  },
  data() {
    return {
      bug: {
        title: '',
        description: '',
        priority: 'Low',
        status: 'open',
        reporter: '',
        assignedTo: null
      },
      developers: ['Alice', 'Bob', 'Charlie'], // Example developer names
      isSubmitting: false
    };
  },
  methods: {
    async submitBug() {
      if (this.isSubmitting) return;

      this.isSubmitting = true;
      try {
        // Set the reporter as the current user
        this.bug.reporter = this.authStore.currentUser.name;

        const response = await fetch('http://localhost:8080/api/bugs', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.authStore.token}`
          },
          body: JSON.stringify(this.bug)
        });

        if (!response.ok) {
          throw new Error('Failed to submit bug');
        }

        const result = await response.json();
        alert('Bug submitted successfully!');
        
        // Reset form
        this.bug = {
          title: '',
          description: '',
          priority: 'Low',
          status: 'open',
          reporter: '',
          assignedTo: null
        };

        // Navigate back to bug list
        this.router.push('/bugs');
      } catch (error) {
        console.error('Error submitting bug:', error);
        alert('Failed to submit bug. Please try again.');
      } finally {
        this.isSubmitting = false;
      }
    }
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
  