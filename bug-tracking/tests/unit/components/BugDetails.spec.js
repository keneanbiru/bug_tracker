import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import BugDetails from '@/pages/BugDetails.vue'
import { useAuthStore } from '@/stores/auth'

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
      path: '/bugs/:id',
      name: 'bug-details',
      component: BugDetails
    }
  ]
})

// Mock bug data
const mockBug = {
  id: '1',
  title: 'Login button not working',
  description: 'Users cannot log in using the login button on the homepage. The button appears to be unresponsive when clicked.',
  status: 'open',
  severity: 'high',
  reporter: 'John Doe',
  assignee: null,
  createdAt: '2024-04-02'
}

// Mock developers
const mockDevelopers = [
  { id: 1, name: 'Mike Johnson' },
  { id: 2, name: 'Sarah Wilson' }
]

describe('BugDetails.vue', () => {
  let wrapper
  let pinia

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks()
    console.error = jest.fn()
    console.log = jest.fn()

    // Setup Pinia
    pinia = createPinia()
    setActivePinia(pinia)

    // Mock route params
    router.push('/bugs/1')

    // Create component wrapper
    wrapper = mount(BugDetails, {
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

  it('renders bug details component', async () => {
    await flushPromises()
    expect(wrapper.find('.bug-details-container').exists()).toBe(true)
  })

  it('displays back button', async () => {
    await flushPromises()
    const backButton = wrapper.find('.back-button')
    expect(backButton.exists()).toBe(true)
    expect(backButton.text()).toBe('Back to List')
  })

  it('displays bug information correctly', async () => {
    await flushPromises()
    
    // Wait for the simulated API call
    await new Promise(resolve => setTimeout(resolve, 500))
    await flushPromises()
    
    const titleSection = wrapper.find('.bug-title-section')
    expect(titleSection.find('h2').text()).toBe(mockBug.title)
    expect(wrapper.find('.status-badge').text()).toBe(mockBug.status)
    expect(wrapper.find('.bug-description p').text()).toBe(mockBug.description)
    expect(wrapper.text()).toContain(mockBug.reporter)
  })

  it('shows update controls for admin', async () => {
    // Mock auth store for admin
    const authStore = useAuthStore()
    authStore.$patch({
      user: { role: 'admin' },
      role: 'admin'
    })
    
    await flushPromises()
    await new Promise(resolve => setTimeout(resolve, 500))
    await flushPromises()
    
    expect(wrapper.find('.bug-actions').exists()).toBe(true)
    expect(wrapper.find('.status-select').exists()).toBe(true)
    expect(wrapper.find('.assignee-select').exists()).toBe(true)
    expect(wrapper.find('.update-button').exists()).toBe(true)
  })

  it('shows update controls for assigned developer', async () => {
    // Mock auth store for developer
    const authStore = useAuthStore()
    authStore.$patch({
      user: { id: 1, name: 'Mike Johnson', role: 'developer' },
      token: 'test-token',
      role: 'developer'
    })
    
    // Set bug assignee to match current user
    const assignedBug = {
      id: '1',
      title: 'Login button not working',
      description: 'Users cannot log in using the login button on the homepage. The button appears to be unresponsive when clicked.',
      status: 'open',
      severity: 'high',
      reporter: 'John Doe',
      assignee: 'Mike Johnson',
      createdAt: '2024-04-02'
    }
    
    // Mock the API call
    global.fetch.mockImplementationOnce(() => Promise.resolve({
      ok: true,
      json: () => Promise.resolve(assignedBug)
    }))
    
    // Wait for component to mount and fetch data
    await flushPromises()
    await new Promise(resolve => setTimeout(resolve, 500))
    await flushPromises()
    
    // Set the bug data manually since we're not using the real API
    wrapper.vm.bug = assignedBug
    wrapper.vm.newStatus = assignedBug.status
    
    await flushPromises()
    
    expect(wrapper.find('.bug-actions').exists()).toBe(true)
    expect(wrapper.find('.status-select').exists()).toBe(true)
    expect(wrapper.find('.assignee-select').exists()).toBe(false)
    expect(wrapper.find('.update-button').exists()).toBe(true)
  })

  it('handles bug update', async () => {
    // Mock auth store for admin
    const authStore = useAuthStore()
    authStore.$patch({
      user: { id: 1, name: 'Admin User', role: 'admin' },
      token: 'test-token',
      role: 'admin'
    })
    
    // Set initial bug data
    wrapper.vm.bug = {
      id: '1',
      title: 'Login button not working',
      description: 'Users cannot log in using the login button on the homepage. The button appears to be unresponsive when clicked.',
      status: 'open',
      severity: 'high',
      reporter: 'John Doe',
      assignee: null,
      createdAt: '2024-04-02'
    }
    wrapper.vm.newStatus = 'open'
    
    await flushPromises()
    
    // Update status
    const statusSelect = wrapper.find('.status-select')
    await statusSelect.setValue('in-progress')
    
    // Update assignee
    const assigneeSelect = wrapper.find('.assignee-select')
    await assigneeSelect.setValue(1)
    
    // Click update button
    const updateButton = wrapper.find('.update-button')
    await updateButton.trigger('click')
    
    await flushPromises()
    await new Promise(resolve => setTimeout(resolve, 1000))
    await flushPromises()
    
    expect(console.log).toHaveBeenCalledWith('Updating bug:', {
      status: 'in-progress',
      assignee: 1
    })
  })

  it('handles error when updating bug', async () => {
    // Mock auth store for admin
    const authStore = useAuthStore()
    authStore.$patch({
      user: { id: 1, name: 'Admin User', role: 'admin' },
      token: 'test-token',
      role: 'admin'
    })
    
    // Set initial bug data
    wrapper.vm.bug = {
      id: '1',
      title: 'Login button not working',
      description: 'Users cannot log in using the login button on the homepage. The button appears to be unresponsive when clicked.',
      status: 'open',
      severity: 'high',
      reporter: 'John Doe',
      assignee: null,
      createdAt: '2024-04-02'
    }
    wrapper.vm.newStatus = 'open'
    
    await flushPromises()
    
    // Mock setTimeout to throw error
    const error = new Error('API Error')
    jest.spyOn(global, 'setTimeout').mockImplementationOnce((callback) => {
      throw error
    })
    
    // Update status and trigger update
    const statusSelect = wrapper.find('.status-select')
    await statusSelect.setValue('in-progress')
    
    const updateButton = wrapper.find('.update-button')
    await updateButton.trigger('click')
    
    await flushPromises()
    
    expect(console.error).toHaveBeenCalledWith('Error updating bug:', error)
  })
}) 