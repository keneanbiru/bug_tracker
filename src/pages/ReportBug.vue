<template>
  <div class="report-bug">
    <div class="header">
      <h1>Report New Bug</h1>
    </div>

    <form @submit.prevent="submitBug" class="bug-form">
      <div class="form-group">
        <label for="title">Title</label>
        <input
          id="title"
          v-model="bugData.title"
          type="text"
          class="form-control"
          required
        >
      </div>

      <div class="form-group">
        <label for="description">Description</label>
        <textarea
          id="description"
          v-model="bugData.description"
          class="form-control"
          rows="4"
          required
        ></textarea>
      </div>

      <div class="form-group">
        <label for="priority">Priority</label>
        <select
          id="priority"
          v-model="bugData.priority"
          class="form-control"
          required
        >
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
          <option value="critical">Critical</option>
        </select>
      </div>

      <div class="form-actions">
        <router-link to="/bugs" class="btn btn-secondary">Cancel</router-link>
        <button 
          type="submit" 
          class="btn btn-primary"
          :disabled="isSubmitting"
        >
          {{ isSubmitting ? 'Submitting...' : 'Submit Bug' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useBugStore } from '../stores/bug'

const router = useRouter()
const bugStore = useBugStore()

const bugData = ref({
  title: '',
  description: '',
  priority: 'low'
})

const isSubmitting = ref(false)

const submitBug = async () => {
  if (isSubmitting.value) return
  
  isSubmitting.value = true
  try {
    await bugStore.createBug(bugData.value)
    router.push('/bugs')
  } catch (error) {
    console.error('Error submitting bug:', error)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.report-bug {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  margin-bottom: 2rem;
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
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}
</style> 