import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import Register from '../../../src/pages/Register.vue'
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
      component: { template: '<div>Login</div>' }
    },
    {
      path: '/register',
      name: 'register',
      component: Register
    }
  ]
})

describe('Register.vue', () => {
  let wrapper
  let pinia
  let authStore

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks()
    console.error = jest.fn()
    console.log = jest.fn()
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
    wrapper = mount(Register, {
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

  it('renders register form correctly', () => {
    expect(wrapper.find('.register-container').exists()).toBe(true)
    expect(wrapper.find('.register-card').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('Register')
    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('select').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
    expect(wrapper.find('.login-link').exists()).toBe(true)
  })

  it('updates form data when inputs change', async () => {
    const nameInput = wrapper.find('input[type="text"]')
    const emailInput = wrapper.find('input[type="email"]')
    const passwordInput = wrapper.find('input[type="password"]')
    const roleSelect = wrapper.find('select')

    await nameInput.setValue('John Doe')
    await emailInput.setValue('john@example.com')
    await passwordInput.setValue('password123')
    await roleSelect.setValue('developer')

    expect(wrapper.vm.name).toBe('John Doe')
    expect(wrapper.vm.email).toBe('john@example.com')
    expect(wrapper.vm.password).toBe('password123')
    expect(wrapper.vm.role).toBe('developer')
  })

  it('validates email format', async () => {
    // Fill in form with invalid email
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('invalid-email')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if error message is displayed
    expect(wrapper.vm.error).toBe('Please enter a valid email address')
  })

  it('validates password length', async () => {
    // Fill in form with short password
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('12345')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if error message is displayed
    expect(wrapper.vm.error).toBe('Password must be at least 6 characters long')
  })

  it('validates name length', async () => {
    // Fill in form with short name
    await wrapper.find('input[type="text"]').setValue('J')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if error message is displayed
    expect(wrapper.vm.error).toBe('Name must be at least 2 characters long')
  })

  it('validates role selection', async () => {
    // Fill in form without selecting role
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if error message is displayed
    expect(wrapper.vm.error).toBe('All fields are required')
  })

  it('calls register method when form is submitted with valid data', async () => {
    // Mock auth store register method
    authStore.register = jest.fn().mockResolvedValue({ success: true })

    // Fill in form with valid data
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if register was called with correct data
    expect(authStore.register).toHaveBeenCalledWith({
      name: 'John Doe',
      email: 'john@example.com',
      password: 'password123',
      role: 'developer'
    })
  })

  it('redirects to login page on successful registration', async () => {
    // Mock auth store register method to return success
    authStore.register = jest.fn().mockResolvedValue({ success: true })

    // Fill in form with valid data
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if router.push was called with correct path
    expect(router.push).toHaveBeenCalledWith('/login')
  })

  it('displays error message on failed registration', async () => {
    // Mock auth store register method to return failure
    authStore.register = jest.fn().mockResolvedValue({ 
      success: false, 
      error: 'Email already exists' 
    })

    // Fill in form with valid data
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if error message is displayed
    expect(wrapper.vm.error).toBe('Email already exists')
  })

  it('handles API errors during registration', async () => {
    // Mock auth store register method to throw error
    authStore.register = jest.fn().mockRejectedValue(new Error('API Error'))

    // Fill in form with valid data
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    // Check if error message is displayed
    expect(wrapper.vm.error).toBe('API Error')
    expect(console.error).toHaveBeenCalledWith('Registration error:', expect.any(Error))
  })

  it('disables form inputs during submission', async () => {
    // Mock auth store register method to delay response
    authStore.register = jest.fn().mockImplementation(() => {
      return new Promise(resolve => {
        setTimeout(() => {
          resolve({ success: true })
        }, 100)
      })
    })

    // Fill in form with valid data
    await wrapper.find('input[type="text"]').setValue('John Doe')
    await wrapper.find('input[type="email"]').setValue('john@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('select').setValue('developer')

    // Submit form
    await wrapper.find('form').trigger('submit')

    // Check if inputs are disabled
    expect(wrapper.find('input[type="text"]').element.disabled).toBe(true)
    expect(wrapper.find('input[type="email"]').element.disabled).toBe(true)
    expect(wrapper.find('input[type="password"]').element.disabled).toBe(true)
    expect(wrapper.find('select').element.disabled).toBe(true)
    expect(wrapper.find('button[type="submit"]').element.disabled).toBe(true)
    expect(wrapper.find('button[type="submit"]').text()).toBe('Registering...')

    // Wait for registration to complete
    await flushPromises()
    await new Promise(resolve => setTimeout(resolve, 100))
    await flushPromises()

    // Check if inputs are enabled again
    expect(wrapper.find('input[type="text"]').element.disabled).toBe(false)
    expect(wrapper.find('input[type="email"]').element.disabled).toBe(false)
    expect(wrapper.find('input[type="password"]').element.disabled).toBe(false)
    expect(wrapper.find('select').element.disabled).toBe(false)
    expect(wrapper.find('button[type="submit"]').element.disabled).toBe(false)
    expect(wrapper.find('button[type="submit"]').text()).toBe('Register')
  })
}) 