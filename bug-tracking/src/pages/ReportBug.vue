<template>
  <MainLayout>
    <div class="report-bug-container">
      <div class="report-bug-header">
        <h1>Report a Bug</h1>
        <p class="subtitle">Help us improve by reporting any issues you encounter</p>
      </div>

      <div class="report-bug-form">
        <form @submit.prevent="submitBug" class="form">
          <div class="form-group">
            <label for="title">Bug Title</label>
            <input
              type="text"
              id="title"
              v-model="bugData.title"
              class="form-control"
              :class="{ 'is-invalid': errors.title }"
              placeholder="Enter a clear and concise title"
              required
            />
            <div class="error-message" v-if="errors.title">{{ errors.title }}</div>
          </div>

          <div class="form-group">
            <label for="description">Description</label>
            <textarea
              id="description"
              v-model="bugData.description"
              class="form-control"
              :class="{ 'is-invalid': errors.description }"
              rows="6"
              placeholder="Provide detailed information about the bug:
• What happened?
• What did you expect to happen?
• Steps to reproduce
• Any relevant error messages"
              required
            ></textarea>
            <div class="error-message" v-if="errors.description">{{ errors.description }}</div>
          </div>

          <div class="form-group">
            <label for="priority">Priority</label>
            <select
              id="priority"
              v-model="bugData.priority"
              class="form-control"
              :class="{ 'is-invalid': errors.priority }"
              required
            >
              <option value="">Select priority level</option>
              <option value="LOW">Low</option>
              <option value="MEDIUM">Medium</option>
              <option value="HIGH">High</option>
              <option value="CRITICAL">Critical</option>
            </select>
            <div class="error-message" v-if="errors.priority">{{ errors.priority }}</div>
          </div>

          <div class="form-group">
            <label for="category">Category</label>
            <select
              id="category"
              v-model="bugData.category"
              class="form-control"
              :class="{ 'is-invalid': errors.category }"
              required
            >
              <option value="">Select a category</option>
              <option value="UI">User Interface</option>
              <option value="FUNCTIONALITY">Functionality</option>
              <option value="PERFORMANCE">Performance</option>
              <option value="SECURITY">Security</option>
              <option value="OTHER">Other</option>
            </select>
            <div class="error-message" v-if="errors.category">{{ errors.category }}</div>
          </div>

          <div class="form-actions">
            <button type="button" @click="resetForm" class="btn btn-secondary">
              Reset Form
            </button>
            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
              <span v-if="isSubmitting" class="spinner"></span>
              {{ isSubmitting ? 'Submitting...' : 'Submit Bug Report' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useBugStore } from '../stores/bug';
import MainLayout from '../layouts/MainLayout.vue';

const router = useRouter();
const bugStore = useBugStore();
const isSubmitting = ref(false);

const bugData = reactive({
  title: '',
  description: '',
  priority: '',
  category: ''
});

const errors = reactive({
  title: '',
  description: '',
  priority: '',
  category: ''
});

const validateForm = () => {
  let isValid = true;
  errors.title = '';
  errors.description = '';
  errors.priority = '';
  errors.category = '';

  if (!bugData.title.trim()) {
    errors.title = 'Title is required';
    isValid = false;
  } else if (bugData.title.length < 5) {
    errors.title = 'Title must be at least 5 characters long';
    isValid = false;
  }

  if (!bugData.description.trim()) {
    errors.description = 'Description is required';
    isValid = false;
  } else if (bugData.description.length < 20) {
    errors.description = 'Description must be at least 20 characters long';
    isValid = false;
  }

  if (!bugData.priority) {
    errors.priority = 'Priority is required';
    isValid = false;
  }

  if (!bugData.category) {
    errors.category = 'Category is required';
    isValid = false;
  }

  return isValid;
};

const resetForm = () => {
  bugData.title = '';
  bugData.description = '';
  bugData.priority = '';
  bugData.category = '';
  Object.keys(errors).forEach(key => errors[key] = '');
};

const submitBug = async () => {
  if (!validateForm()) {
    return;
  }

  isSubmitting.value = true;

  try {
    await bugStore.reportBug(bugData);
    router.push('/bugs');
  } catch (error) {
    console.error('Error reporting bug:', error);
    alert(error.message || 'Failed to report bug. Please try again.');
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped>
.report-bug-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.report-bug-header {
  text-align: center;
  margin-bottom: 2rem;
}

.report-bug-header h1 {
  font-size: 2.5rem;
  color: #2d3748;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #718096;
  font-size: 1.1rem;
}

.report-bug-form {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #4a5568;
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

.form-control.is-invalid {
  border-color: #e53e3e;
}

.error-message {
  color: #e53e3e;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
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
  min-width: 120px;
}

.btn-primary {
  background-color: #4299e1;
  color: white;
  border: none;
}

.btn-primary:hover {
  background-color: #3182ce;
}

.btn-primary:disabled {
  background-color: #a0aec0;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #e2e8f0;
  color: #4a5568;
  border: none;
}

.btn-secondary:hover {
  background-color: #cbd5e0;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid #ffffff;
  border-top: 3px solid transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 8px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

textarea.form-control {
  resize: vertical;
  min-height: 120px;
  font-family: inherit;
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
  