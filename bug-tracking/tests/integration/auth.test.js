import { mount } from '@vue/test-utils'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import axios from 'axios'
import { useAuthStore } from '../../src/stores/auth'
import Login from '../../src/pages/Login.vue'

// Mock axios
jest.mock('axios', () => {
  const mockAxios = {
    get: jest.fn(),
    post: jest.fn(),
    put: jest.fn(),
    delete: jest.fn(),
    create: jest.fn(() => ({
      get: jest.fn(),
      post: jest.fn(),
      put: jest.fn(),
      delete: jest.fn(),
      interceptors: {
        request: {
          use: jest.fn(),
          eject: jest.fn()
        },
        response: {
          use: jest.fn(),
          eject: jest.fn()
        }
      }
    })),
    interceptors: {
      request: {
        use: jest.fn(),
        eject: jest.fn()
      },
      response: {
        use: jest.fn(),
        eject: jest.fn()
      }
    }
  }
  return mockAxios
})

// Mock router
const mockRouter = {
  push: jest.fn()
}

jest.mock('vue-router', () => ({
  useRouter: () => mockRouter
}))

// Mock window object
const mockAlert = jest.fn()
Object.defineProperty(window, 'alert', {
  value: mockAlert,
  writable: true
})

describe('Authentication Tests', () => {
  let wrapper
  let store

  beforeEach(() => {
    // Create a fresh Pinia instance for each test
    const pinia = createPinia()
    const app = createApp({})
    app.use(pinia)
    store = useAuthStore(pinia)

    // Reset store state
    store.$reset()

    // Reset all mocks
    jest.clearAllMocks()
    mockRouter.push.mockClear()
    mockAlert.mockClear()

    // Clear localStorage
    localStorage.clear()

    // Mock localStorage methods
    const localStorageMock = {
      getItem: jest.fn(),
      setItem: jest.fn(),
      removeItem: jest.fn(),
      clear: jest.fn()
    }
    Object.defineProperty(window, 'localStorage', {
      value: localStorageMock
    })

    // Mount the component with router
    wrapper = mount(Login, {
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
  })

  afterEach(() => {
    // Clean up after each test
    wrapper.unmount()
  })

  it('should login successfully with valid credentials', async () => {
    // Mock successful login response
    const mockResponse = {
      data: {
        token: 'test-token',
        user: {
          id: 1,
          username: 'testuser',
          email: 'test@example.com',
          role: 'developer'
        }
      }
    }
    axios.post.mockResolvedValueOnce(mockResponse)

    // Fill in login form
    await wrapper.find('input[type="email"]').setValue('test@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit')

    // Wait for async operations
    await wrapper.vm.$nextTick()

    // Verify API was called with correct data
    expect(axios.post).toHaveBeenCalledWith(
      '/api/auth/login',
      {
        email: 'test@example.com',
        password: 'password123'
      }
    )

    // Verify store was updated
    expect(store.token).toBe('test-token')
    expect(store.user).toEqual(mockResponse.data.user)
    expect(store.role).toBe('developer')
    expect(store.isAuthenticated).toBe(true)

    // Verify router navigation
    expect(mockRouter.push).toHaveBeenCalledWith('/')
  })

  it('should show error message with invalid credentials', async () => {
    // Mock failed login response
    axios.post.mockRejectedValueOnce({
      response: {
        data: {
          message: 'Invalid credentials'
        }
      }
    })

    // Fill in login form
    await wrapper.find('input[type="email"]').setValue('test@example.com')
    await wrapper.find('input[type="password"]').setValue('wrongpassword')

    // Submit form
    await wrapper.find('form').trigger('submit')

    // Wait for async operations
    await wrapper.vm.$nextTick()

    // Verify error message
    expect(mockAlert).toHaveBeenCalledWith('Login failed. Please check your credentials.')

    // Verify store was not updated
    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.role).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })

  it('should logout successfully', async () => {
    // Set initial authenticated state
    store.$patch({
      token: 'test-token',
      user: {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        role: 'developer'
      },
      role: 'developer'
    })

    // Call logout and wait for it to complete
    await store.logout()

    // Verify store was cleared
    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.role).toBeNull()
    expect(store.isAuthenticated).toBe(false)

    // Verify router navigation
    expect(mockRouter.push).toHaveBeenCalledWith('/login')
  })
}) 