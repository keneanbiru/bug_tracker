<template>
  <div class="new-bug-container">
    <h1>Report New Bug</h1>
    <div v-if="bugStore.error" class="error-message">
      {{ bugStore.error }}
    </div>
    <form @submit.prevent="handleSubmit" class="bug-form">
      <div class="form-group">
        <label for="title">Title</label>
        <input 
          type="text" 
          id="title" 
          v-model="title" 
          required 
          placeholder="Enter bug title"
        >
      </div>

      <div class="form-group">
        <label for="description">Description</label>
        <textarea 
          id="description" 
          v-model="description" 
          required 
          placeholder="Describe the bug in detail"
          rows="4"
        ></textarea>
      </div>

      <div class="form-group">
        <label for="severity">Severity</label>
        <select id="severity" v-model="severity" required>
          <option value="">Select severity</option>
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
          <option value="critical">Critical</option>
        </select>
      </div>

      <div class="form-group" v-if="authStore.isAdmin">
        <label for="assignee">Assign To</label>
        <select id="assignee" v-model="assignee">
          <option value="">Select developer</option>
          <option v-for="developer in developers" :key="developer.id" :value="developer.id">
            {{ developer.name }}
          </option>
        </select>
      </div>

      <div class="form-actions">
        <button type="button" @click="$router.back()" class="cancel-button">Cancel</button>
        <button type="submit" class="submit-button">Submit Bug</button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useBugStore } from '../stores/bug'

const router = useRouter()
const authStore = useAuthStore()
const bugStore = useBugStore()

// Mock data - replace with API call
const developers = ref([
  { id: 1, name: 'Mike Johnson' },
  { id: 2, name: 'Sarah Wilson' }
])

const title = ref('')
const description = ref('')
const severity = ref('')
const assignee = ref('')

const handleSubmit = async () => {
  try {
    const bugData = {
      title: title.value,
      description: description.value,
      severity: severity.value,
      assignee: assignee.value,
      reporter: authStore.currentUser?.name,
      status: 'open'
    }

    await bugStore.createBug(bugData)
    router.push('/bugs')
  } catch (err) {
    console.error('Error submitting bug:', err)
  }
}
</script>

<style scoped>
.new-bug-container {
  max-width: 800px;
  margin: 0 auto;
}

.new-bug-container h1 {
  margin-bottom: 2rem;
  color: #2c3e50;
}

.bug-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4a5568;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.cancel-button,
.submit-button {
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.cancel-button {
  background-color: #e2e8f0;
  color: #4a5568;
  border: none;
}

.cancel-button:hover {
  background-color: #cbd5e0;
}

.submit-button {
  background-color: #4CAF50;
  color: white;
  border: none;
}

.submit-button:hover {
  background-color: #45a049;
}

.error-message {
  background-color: #ffebee;
  color: #c62828;
  padding: 0.75rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  text-align: center;
}
</style> 