<template>
  <div class="register-container">
    <div class="register-card">
      <h2>Register</h2>
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="name">Full Name</label>
          <input 
            type="text" 
            id="name" 
            v-model="name" 
            required 
            minlength="2"
            placeholder="Enter your full name (min 2 characters)"
            :disabled="isLoading"
          >
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <input 
            type="email" 
            id="email" 
            v-model="email" 
            required 
            placeholder="Enter your email"
            :disabled="isLoading"
          >
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input 
            type="password" 
            id="password" 
            v-model="password" 
            required 
            minlength="6"
            placeholder="Enter your password (min 6 characters)"
            :disabled="isLoading"
          >
        </div>
        <div class="form-group">
          <label for="role">Role</label>
          <select id="role" v-model="role" required :disabled="isLoading">
            <option value="">Select a role</option>
            <option value="developer">Developer</option>
            <option value="manager">Manager</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <button type="submit" class="register-button" :disabled="isLoading">
          {{ isLoading ? 'Registering...' : 'Register' }}
        </button>
        <p class="login-link">
          Already have an account? 
          <router-link to="/login">Login here</router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')
const role = ref('')
const error = ref('')
const isLoading = ref(false)

const validateEmail = (email) => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return re.test(email)
}

const handleRegister = async () => {
  error.value = '' // Clear any previous errors
  isLoading.value = true
  
  try {
    // Validate inputs
    if (!name.value || !email.value || !password.value || !role.value) {
      throw new Error('All fields are required')
    }

    if (name.value.length < 2) {
      throw new Error('Name must be at least 2 characters long')
    }

    if (!validateEmail(email.value)) {
      throw new Error('Please enter a valid email address')
    }

    if (password.value.length < 6) {
      throw new Error('Password must be at least 6 characters long')
    }

    if (!['developer', 'manager', 'admin'].includes(role.value)) {
      throw new Error('Please select a valid role')
    }

    console.log('Sending registration request with data:', {
      name: name.value,
      email: email.value,
      password: password.value,
      role: role.value
    });

    const result = await authStore.register({
      name: name.value,
      email: email.value,
      password: password.value,
      role: role.value
    })

    console.log('Registration result:', result);

    if (result.success) {
      router.push('/login')
    } else {
      error.value = result.error || 'Registration failed. Please try again.'
      console.error('Registration failed:', result.error);
    }
  } catch (err) {
    error.value = err.message || 'An unexpected error occurred. Please try again.'
    console.error('Registration error:', err)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.register-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.register-card h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 500;
  color: #555;
}

.form-group input,
.form-group select {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.register-button {
  background-color: #4CAF50;
  color: white;
  padding: 0.75rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.register-button:hover {
  background-color: #45a049;
}

.login-link {
  text-align: center;
  margin-top: 1rem;
  color: #666;
}

.login-link a {
  color: #4CAF50;
  text-decoration: none;
}

.login-link a:hover {
  text-decoration: underline;
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