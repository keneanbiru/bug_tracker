import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import BugList from '@/pages/BugList.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBugStore } from '@/stores/bug'

// Mock console methods
console.error = jest.fn()
console.log = jest.fn()

// Mock fetch globally
global.fetch = jest.fn()

// Create a mock router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/bugs/:id', name: 'bug-details', component: { template: '<div>Bug Details</div>' } }
  ]
})

describe('BugList.vue', () => {
  let authStore
  let bugStore
  let pinia

  const mockBugs = [
    {
      id: '1',
      title: 'Test Bug 1',
      description: 'Description 1',
      status: 'OPEN',
      assigned_to: null
    },
    {
      id: '2',
      title: 'Test Bug 2',
      description: 'Description 2',
      status: 'IN_PROGRESS',
      assigned_to: { id: 'dev1', name: 'Developer 1' }
    }
  ]

  const mockDevelopers = [
    { id: 'dev1', name: 'Developer 1' },
    { id: 'dev2', name: 'Developer 2' }
  ]

  beforeEach(() => {
    // Create a fresh Pinia instance
    pinia = createPinia()
    setActivePinia(pinia)
    
    // Get store instances
    authStore = useAuthStore()
    bugStore = useBugStore()
    
    // Mock the auth store state
    authStore.$patch({
      token: 'test-token',
      user: { id: 'user1', role: 'user' },
      role: 'user'
    })

    // Mock the bug store state
    bugStore.$patch({
      bugs: mockBugs,
      loading: false,
      error: null
    })

    // Mock fetchBugs method
    bugStore.fetchBugs = jest.fn().mockResolvedValue(mockBugs)

    // Reset fetch mock
    global.fetch.mockReset()

    // Mock localStorage
    Object.defineProperty(window, 'localStorage', {
      value: {
        getItem: jest.fn(),
        setItem: jest.fn(),
        removeItem: jest.fn()
      },
      writable: true
    })
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  test('renders bug list component', async () => {
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    expect(wrapper.find('.bug-list').exists()).toBe(true)
    expect(wrapper.find('h1').text()).toBe('Bug List')
  })

  test('shows loading state', async () => {
    bugStore.$patch({ loading: true })
    
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    expect(wrapper.find('.loading').exists()).toBe(true)
    expect(wrapper.text()).toContain('Loading bugs...')
  })

  test('shows error state', async () => {
    const errorMessage = 'Failed to load bugs'
    bugStore.$patch({ 
      loading: false,
      error: errorMessage,
      bugs: []
    })
    
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    expect(wrapper.find('.error').exists()).toBe(true)
    expect(wrapper.text()).toContain(errorMessage)
  })

  test('displays bugs correctly', async () => {
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    const bugCards = wrapper.findAll('.bug-card')
    expect(bugCards).toHaveLength(2)
    
    expect(bugCards[0].find('h3').text()).toBe('Test Bug 1')
    expect(bugCards[0].find('.status').text()).toBe('Open')
    expect(bugCards[0].find('.unassigned').exists()).toBe(true)
    
    expect(bugCards[1].find('h3').text()).toBe('Test Bug 2')
    expect(bugCards[1].find('.status').text()).toBe('In Progress')
    expect(bugCards[1].text()).toContain('Developer 1')
  })

  test('filters bugs by search query', async () => {
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    const searchInput = wrapper.find('input[type="text"]')
    await searchInput.setValue('Bug 1')
    await wrapper.vm.$nextTick()
    
    const bugCards = wrapper.findAll('.bug-card')
    expect(bugCards).toHaveLength(1)
    expect(bugCards[0].find('h3').text()).toBe('Test Bug 1')
  })

  test('filters bugs by status', async () => {
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    const statusSelect = wrapper.find('select')
    await statusSelect.setValue('IN_PROGRESS')
    await wrapper.vm.$nextTick()
    
    const bugCards = wrapper.findAll('.bug-card')
    expect(bugCards).toHaveLength(1)
    expect(bugCards[0].find('h3').text()).toBe('Test Bug 2')
  })

  test('shows assign button for managers', async () => {
    // Update auth store to manager role
    authStore.$patch({
      token: 'test-token',
      user: { id: 'manager1', role: 'manager' },
      role: 'manager'
    })

    // Mock localStorage for manager role
    window.localStorage.getItem.mockImplementation((key) => {
      switch (key) {
        case 'role':
          return 'manager'
        case 'token':
          return 'test-token'
        case 'user':
          return JSON.stringify({ id: 'manager1', role: 'manager' })
        default:
          return null
      }
    })
    
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    const assignButtons = wrapper.findAll('[data-test="assign-button"]')
    expect(assignButtons).toHaveLength(1)
  })

  test('shows update status button for assigned developer', async () => {
    // Update auth store to developer role
    authStore.$patch({
      token: 'test-token',
      user: { id: 'dev1', role: 'developer' },
      role: 'developer'
    })

    // Mock localStorage for developer role
    window.localStorage.getItem.mockImplementation((key) => {
      switch (key) {
        case 'role':
          return 'developer'
        case 'token':
          return 'test-token'
        case 'user':
          return JSON.stringify({ id: 'dev1', role: 'developer' })
        default:
          return null
      }
    })
    
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    const updateButtons = wrapper.findAll('[data-test="update-status-button"]')
    expect(updateButtons).toHaveLength(1)
  })

  test('fetches developers when user can assign bugs', async () => {
    // Update auth store to manager role
    authStore.$patch({
      token: 'test-token',
      user: { id: 'manager1', role: 'manager' },
      role: 'manager'
    })

    // Mock localStorage for manager role
    window.localStorage.getItem.mockImplementation((key) => {
      switch (key) {
        case 'role':
          return 'manager'
        case 'token':
          return 'test-token'
        case 'user':
          return JSON.stringify({ id: 'manager1', role: 'manager' })
        default:
          return null
      }
    })
    
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockDevelopers)
    })
    
    const wrapper = mount(BugList, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    await flushPromises()
    
    expect(global.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/auth/developers',
      expect.objectContaining({
        headers: {
          'Authorization': 'Bearer test-token'
        }
      })
    )
  })
}) 