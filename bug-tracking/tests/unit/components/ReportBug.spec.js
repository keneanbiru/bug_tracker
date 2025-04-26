import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ReportBug from '../../../src/pages/ReportBug.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../../../src/stores/auth'

// Mock fetch globally
global.fetch = jest.fn()

// Mock window.alert
window.alert = jest.fn()

// Create a mock router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/bugs', name: 'bugs', component: { template: '<div>Bugs List</div>' } }
  ]
})

describe('ReportBug.vue', () => {
  let authStore
  let pinia

  beforeEach(() => {
    // Create a fresh Pinia instance with initial state
    pinia = createPinia()
    setActivePinia(pinia)
    
    // Get store instance
    authStore = useAuthStore()
    
    // Mock the auth store state
    authStore.$patch({
      token: 'test-token',
      user: {
        role: 'user',
        id: 'test-user-id'
      },
      role: 'user' // Make sure to set the role state
    })

    // Reset fetch mock
    global.fetch.mockReset()
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  test('renders report bug form', async () => {
    const wrapper = mount(ReportBug, {
      global: {
        plugins: [router, pinia]
      }
    })
    
    // Wait for component to update
    await wrapper.vm.$nextTick()
    
    // Check if the form renders
    expect(wrapper.find('form.bug-form').exists()).toBe(true)
    
    // Check for required form fields
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('textarea').exists()).toBe(true)
    expect(wrapper.find('select').exists()).toBe(true)
    
    // Check form labels
    expect(wrapper.text()).toContain('Title')
    expect(wrapper.text()).toContain('Description')
    expect(wrapper.text()).toContain('Priority')
  })

  test('validates required fields', async () => {
    const wrapper = mount(ReportBug, {
      global: {
        plugins: [router, pinia]
      }
    })

    // Try to submit empty form
    await wrapper.find('form').trigger('submit')
    
    // Get form controls
    const titleInput = wrapper.find('input[type="text"]')
    const descriptionTextarea = wrapper.find('textarea')
    
    // Check validity states
    expect(titleInput.element.validity.valid).toBe(false)
    expect(descriptionTextarea.element.validity.valid).toBe(false)
  })

  test('submits bug report successfully', async () => {
    // Mock successful API response
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ id: 'new-bug-id' })
    })

    const wrapper = mount(ReportBug, {
      global: {
        plugins: [router, pinia]
      }
    })

    // Fill in the form
    await wrapper.find('input[type="text"]').setValue('Test Bug')
    await wrapper.find('textarea').setValue('Test Description')
    await wrapper.find('select').setValue('High')

    // Submit the form
    await wrapper.find('form').trigger('submit')

    // Check if API was called with correct data
    expect(global.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/bugs',
      expect.objectContaining({
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token'
        },
        body: JSON.stringify({
          title: 'Test Bug',
          description: 'Test Description',
          priority: 'high'
        })
      })
    )

    // Verify that success message is shown (via alert)
    expect(window.alert).toHaveBeenCalledWith('Bug submitted successfully!')
  })

  test('handles API errors', async () => {
    // Mock API error response
    global.fetch.mockResolvedValueOnce({
      ok: false,
      json: () => Promise.resolve({ error: 'Server error' })
    })

    const wrapper = mount(ReportBug, {
      global: {
        plugins: [router, pinia]
      }
    })

    // Fill in the form
    await wrapper.find('input[type="text"]').setValue('Test Bug')
    await wrapper.find('textarea').setValue('Test Description')
    await wrapper.find('select').setValue('High')

    // Submit the form
    await wrapper.find('form').trigger('submit')

    // Verify that error is logged
    expect(console.error).toHaveBeenCalled()
  })

  test('shows developer selection for managers', async () => {
    // Mock the auth store with manager permissions
    authStore.$patch({
      token: 'test-token',
      user: {
        role: 'manager',
        id: 'test-manager-id'
      },
      role: 'manager' // This is crucial for the canAssignBugs getter
    })

    // Mock developers API response
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([
        { id: 'dev1', name: 'Developer 1' },
        { id: 'dev2', name: 'Developer 2' }
      ])
    })

    const wrapper = mount(ReportBug, {
      global: {
        plugins: [router, pinia]
      }
    })

    // Wait for all promises to resolve
    await flushPromises()
    
    // Check if developer selection is shown
    const assignToLabel = wrapper.findAll('label.form-label').find(label => label.text() === 'Assign To')
    expect(assignToLabel).toBeTruthy()
    
    // Find all selects
    const selects = wrapper.findAll('select')
    expect(selects.length).toBe(2) // One for priority, one for developer assignment
    
    // Get the developer select (second select element)
    const developerSelect = selects[1]
    const options = developerSelect.findAll('option')
    
    // Check options
    expect(options.length).toBe(3) // Unassigned + 2 developers
    expect(options[0].text()).toBe('Unassigned')
    expect(options[1].text()).toBe('Developer 1')
    expect(options[2].text()).toBe('Developer 2')
  })
}) 