import { mount } from '@vue/test-utils'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import axios from 'axios'
import { useAuthStore } from '../../src/stores/auth'
import Login from '../../src/pages/Login.vue'
import Register from '../../src/pages/Register.vue'

// Mock axios
jest.mock('axios')

// Mock router
const mockRouter = {
  push: jest.fn()
}

jest.mock('vue-router', () => ({
  useRouter: () => mockRouter
}))

// Mock window object
const mockAlert = jest.fn()
global.window = {
  alert: mockAlert
}

describe('Authentication Integration Tests', () => {
  let app
  let pinia
  let authStore

  beforeEach(() => {
    // Create a fresh app instance for each test
    app = createApp({})
    pinia = createPinia()
    app.use(pinia)
    authStore = useAuthStore()

    // Reset store state
    authStore.$reset()
    mockRouter.push.mockClear()
    mockAlert.mockClear()
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  describe('Login Flow', () => {
    it('should successfully login a user and store token', async () => {
      // Mock successful login response
      const mockResponse = {
        data: {
          token: 'test-token',
          user: {
            id: '1',
            name: 'Test User',
            email: 'test@example.com',
            role: 'developer'
          }
        }
      }
      axios.post.mockResolvedValueOnce(mockResponse)

      // Mount the login component
      const wrapper = mount(Login, {
        global: {
          plugins: [pinia],
          stubs: {
            'router-link': true,
            'router-view': true
          }
        }
      })

      // Fill in login form
      await wrapper.find('input[type="email"]').setValue('test@example.com')
      await wrapper.find('input[type="password"]').setValue('password123')
      
      // Submit form
      await wrapper.find('form').trigger('submit')

      // Wait for API call to complete
      await wrapper.vm.$nextTick()

      // Verify API was called with correct data
      expect(axios.post).toHaveBeenCalledWith(
        expect.stringContaining('/login'),
        {
          email: 'test@example.com',
          password: 'password123'
        }
      )

      // Verify store was updated
      expect(authStore.token).toBe('test-token')
      expect(authStore.user).toEqual(mockResponse.data.user)

      // Verify router navigation
      expect(mockRouter.push).toHaveBeenCalledWith('/')
    })

    it('should handle login errors correctly', async () => {
      // Mock failed login response
      const mockError = {
        response: {
          data: {
            error: 'Invalid credentials'
          }
        }
      }
      axios.post.mockRejectedValueOnce(mockError)

      // Mount the login component
      const wrapper = mount(Login, {
        global: {
          plugins: [pinia],
          stubs: {
            'router-link': true,
            'router-view': true
          }
        }
      })

      // Fill in login form
      await wrapper.find('input[type="email"]').setValue('test@example.com')
      await wrapper.find('input[type="password"]').setValue('wrongpassword')
      
      // Submit form
      await wrapper.find('form').trigger('submit')

      // Wait for API call to complete
      await wrapper.vm.$nextTick()

      // Verify error message is displayed via alert
      expect(window.alert).toHaveBeenCalledWith('Login failed. Please check your credentials.')
    })
  })

  describe('Registration Flow', () => {
    it('should successfully register a new user', async () => {
      // Mock successful registration response
      const mockResponse = {
        data: {
          id: 1,
          name: 'New User',
          email: 'new@example.com',
          role: 'developer'
        }
      }
      axios.post.mockResolvedValueOnce(mockResponse)

      // Mount the register component
      const wrapper = mount(Register, {
        global: {
          plugins: [pinia],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: mockRouter
          }
        }
      })

      // Fill in registration form
      await wrapper.find('input[id="name"]').setValue('New User')
      await wrapper.find('input[id="email"]').setValue('new@example.com')
      await wrapper.find('input[id="password"]').setValue('password123')
      await wrapper.find('select[id="role"]').setValue('developer')

      // Submit form
      await wrapper.find('form').trigger('submit')

      // Wait for async operations
      await wrapper.vm.$nextTick()

      // Verify API was called with correct data
      expect(axios.post).toHaveBeenCalledWith(
        expect.stringContaining('/register'),
        {
          name: 'New User',
          email: 'new@example.com',
          password: 'password123',
          role: 'developer'
        }
      )

      // Verify router navigation
      expect(mockRouter.push).toHaveBeenCalledWith('/login')
    })

    it('should handle registration errors correctly', async () => {
      // Mock failed registration response
      const mockError = {
        response: {
          data: {
            message: 'Email already exists'
          }
        }
      }
      axios.post.mockRejectedValueOnce(mockError)

      // Mount the register component
      const wrapper = mount(Register, {
        global: {
          plugins: [pinia],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: mockRouter
          }
        }
      })

      // Fill in registration form
      await wrapper.find('input[id="name"]').setValue('Test User')
      await wrapper.find('input[id="email"]').setValue('existing@example.com')
      await wrapper.find('input[id="password"]').setValue('password123')
      await wrapper.find('select[id="role"]').setValue('developer')

      // Submit form
      await wrapper.find('form').trigger('submit')

      // Wait for API call to complete and error to be set
      await wrapper.vm.$nextTick()
      await wrapper.vm.$nextTick()

      // Verify error message is displayed
      const errorMessage = wrapper.find('.error-message')
      expect(errorMessage.exists()).toBe(true)
      expect(errorMessage.text()).toBe('Email already exists')
    })
  })
}) 