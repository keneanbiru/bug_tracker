import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import Login from '../../../src/pages/Login.vue'
import { useAuthStore } from '../../../src/stores/auth'

// Mock console methods
const originalConsoleError = console.error
const originalConsoleLog = console.log

// Mock fetch
global.fetch = jest.fn()

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  clear: jest.fn()
}
global.localStorage = localStorageMock

// Mock window.alert
window.alert = jest.fn()

// Create router instance
const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: { template: '<div>Home</div>' }
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/register',
      name: 'register',
      component: { template: '<div>Register</div>' }
    }
  ]
})

describe('Login.vue', () => {
  let wrapper
  let pinia
  let authStore

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks()
    console.error = jest.fn()
    console.log = jest.fn()
    window.alert.mockClear()
    localStorageMock.getItem.mockClear()
    localStorageMock.setItem.mockClear()
    localStorageMock.clear.mockClear()
    global.fetch.mockReset()

    // Setup Pinia
    pinia = createPinia()
    setActivePinia(pinia)

    // Get auth store
    authStore = useAuthStore()

    // Mock router push
    router.push = jest.fn()

    // Create component wrapper
    wrapper = mount(Login, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': true
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  afterAll(() => {
    console.error = originalConsoleError
    console.log = originalConsoleLog
  })

  it('renders login form correctly', () => {
    expect(wrapper.find('.login-container').exists()).toBe(true)
    expect(wrapper.find('.login-card').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('Login')
    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
    expect(wrapper.find('.register-link').exists()).toBe(true)
  })

  it('updates form data when inputs change', async () => {
    const emailInput = wrapper.find('input[type="email"]')
    const passwordInput = wrapper.find('input[type="password"]')

    await emailInput.setValue('test@example.com')
    await passwordInput.setValue('password123')

    expect(wrapper.vm.email).toBe('test@example.com')
    expect(wrapper.vm.password).toBe('password123')
  })

  it('calls login method when form is submitted', async () => {
    // Mock auth store login method
    authStore.login = jest.fn().mockResolvedValue(true)

    // Fill in form
    await wrapper.find('input[type="email"]').setValue('test@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit')

    // Check if login was called with correct credentials
    expect(authStore.login).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password123'
    })
  })

  it('redirects to home page on successful login', async () => {
    // Mock auth store login method to return success
    authStore.login = jest.fn().mockResolvedValue(true)

    // Fill in form
    await wrapper.find('input[type="email"]').setValue('test@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if router.push was called with correct path
    expect(router.push).toHaveBeenCalledWith('/')
  })

  it('shows error message on failed login', async () => {
    // Mock auth store login method to return failure
    authStore.login = jest.fn().mockResolvedValue(false)

    // Fill in form
    await wrapper.find('input[type="email"]').setValue('test@example.com')
    await wrapper.find('input[type="password"]').setValue('wrongpassword')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if alert was called with error message
    expect(window.alert).toHaveBeenCalledWith('Login failed. Please check your credentials.')
  })

  it('validates required fields', async () => {
    // Mock the login method
    authStore.login = jest.fn()

    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia],
        stubs: {
          'router-link': true
        }
      }
    })

    // Get form controls
    const emailInput = wrapper.find('input[type="email"]')
    const passwordInput = wrapper.find('input[type="password"]')
    
    // Check validity states before submission
    expect(emailInput.element.validity.valid).toBe(false)
    expect(passwordInput.element.validity.valid).toBe(false)
    
    // Try to submit empty form
    await wrapper.find('form').trigger('submit')

    // Check if form validation prevents submission
    expect(authStore.login).not.toHaveBeenCalled()
  })
}) 