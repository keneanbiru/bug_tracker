import { mount } from '@vue/test-utils'
import { createApp } from 'vue'
import { createPinia, setActivePinia } from 'pinia'
import axios from 'axios'
import { useAuthStore } from '../../src/stores/auth'
import { useBugStore } from '../../src/stores/bug'
import BugList from '../../src/pages/BugList.vue'
import NewBug from '../../src/pages/NewBug.vue'
import { flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createRouter, createWebHistory } from 'vue-router'

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  clear: jest.fn()
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock axios
jest.mock('axios')

// Mock router
const router = createRouter({
  history: createWebHistory(),
  routes: []
})

// Mock window object
const mockAlert = jest.fn()
global.window = {
  alert: mockAlert
}

describe('Bug Management Integration Tests', () => {
  let wrapper = null
  let authStore;
  let bugStore;
  let router
  
  beforeEach(() => {
    // Reset all mocks
    jest.clearAllMocks()
    
    // Mock successful API responses
    axios.get.mockResolvedValue({ 
      data: [
        { 
          id: 1, 
          title: 'Test Bug 1', 
          status: 'OPEN', 
          assigned_to: {
            id: 1,
            name: 'testuser'
          }
        }
      ]
    })
    
    axios.post.mockResolvedValue({ 
      data: { 
        id: 1, 
        title: 'New Bug',
        status: 'OPEN'
      } 
    })
    
    axios.put.mockResolvedValue({ 
      data: { 
        id: 1,
        status: 'IN_PROGRESS'
      } 
    })

    // Set up router mock
    router = {
      push: jest.fn(),
      replace: jest.fn()
    }

    // Set up fetch mock
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ status: 'IN_PROGRESS' })
      })
    )

    // Mount the bug list component with properly initialized stores
    wrapper = mount(BugList, {
      global: {
        plugins: [createTestingPinia({
          initialState: {
            auth: {
              isAuthenticated: true,
              currentUser: { id: 1, name: 'Test User' },
              role: 'developer',
              canEditBugStatus: true,
              isManager: true,
              canAssignBugs: true
            },
            bug: {
              bugs: [
                {
                  id: 1,
                  title: 'Test Bug 1',
                  description: 'Description 1',
                  status: 'OPEN',
                  priority: 'high',
                  assigned_to: {
                    id: 1,
                    name: 'testuser'
                  }
                }
              ],
              loading: false,
              error: null
            }
          },
          stubActions: false
        })],
        stubs: {
          RouterLink: true,
          'router-link': true,
          'router-view': true
        },
        mocks: {
          $router: router
        }
      }
    })

    // Get store instances after mounting
    authStore = useAuthStore()
    bugStore = useBugStore()
  })

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount()
    }
    jest.clearAllMocks()
  })

  const mountComponent = (initialState = {}) => {
    return mount(BugList, {
      global: {
        plugins: [createTestingPinia({
          initialState: {
            auth: {
              user: {
                id: 1,
                name: 'testuser',
                role: 'user'
              },
              role: 'user',
              token: 'test-token',
              isManager: false,
              isAdmin: false,
              canEditBugStatus: true,
              currentUser: {
                id: 1,
                name: 'testuser',
                role: 'user'
              },
              ...initialState.auth
            },
            bug: {
              bugs: [],
              loading: false,
              error: null,
              ...initialState.bug
            }
          },
          stubActions: false
        })],
        stubs: {
          'router-link': true,
          'router-view': true
        },
        mocks: {
          $router: router
        }
      }
    })
  }

  describe('Bug List Flow', () => {
    it('should fetch and display bugs for a developer', async () => {
      // Wait for component to mount and API calls to complete
      await flushPromises()

      // Verify API was called
      expect(axios.get).toHaveBeenCalledWith('/api/bugs')

      // Verify bugs are displayed
      expect(wrapper.text()).toContain('Test Bug 1')
    })

    it('should handle error when fetching bugs', async () => {
      // Mock failed response
      axios.get.mockRejectedValueOnce(new Error('Failed to fetch bugs'))

      // Remount component to trigger error
      wrapper = mount(BugList, {
        global: {
          plugins: [createTestingPinia({
            initialState: {
              auth: {
                isAuthenticated: true,
                currentUser: { id: 1, name: 'Test User' },
                role: 'developer'
              },
              bug: {
                bugs: [],
                loading: false,
                error: 'Failed to fetch bugs'
              }
            }
          })],
          stubs: {
            RouterLink: true,
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: router
          }
        }
      })

      // Wait for error to be displayed
      await flushPromises()

      // Verify error message is displayed
      expect(wrapper.find('.error').exists()).toBe(true)
      expect(wrapper.find('.error').text()).toBe('Failed to fetch bugs')
    })
  })

  describe('Bug Creation Flow', () => {
    it('should successfully create a new bug', async () => {
      // Mount the new bug component
      wrapper = mount(NewBug, {
        global: {
          plugins: [createTestingPinia({
            initialState: {
              auth: {
                user: { id: 1, name: 'Test User', role: 'developer' },
                role: 'developer',
                token: 'test-token'
              }
            }
          })],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: router
          }
        }
      })

      // Fill in bug form
      await wrapper.find('input[id="title"]').setValue('New Bug')
      await wrapper.find('textarea[id="description"]').setValue('New Description')
      await wrapper.find('select[id="severity"]').setValue('high')
      
      // Submit form
      await wrapper.find('form').trigger('submit')

      // Wait for API call to complete
      await flushPromises()

      // Verify API was called with correct data
      expect(axios.post).toHaveBeenCalledWith('/api/bugs', {
        title: 'New Bug',
        description: 'New Description',
        severity: 'high',
        assignee: '',
        reporter: 'Test User',
        status: 'open'
      })

      // Verify navigation
      expect(router.push).toHaveBeenCalledWith('/bugs')
    })

    it('should handle error when creating bug', async () => {
      // Mock failed response
      axios.post.mockRejectedValueOnce({
        response: {
          data: {
            message: 'Failed to create bug'
          }
        }
      })

      // Mount component
      wrapper = mount(NewBug, {
        global: {
          plugins: [createTestingPinia({
            initialState: {
              auth: {
                user: { id: 1, name: 'Test User', role: 'developer' },
                role: 'developer'
              }
            }
          })],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: router
          }
        }
      })

      // Fill in and submit form
      await wrapper.find('input[id="title"]').setValue('New Bug')
      await wrapper.find('textarea[id="description"]').setValue('New Description')
      await wrapper.find('select[id="severity"]').setValue('high')
      await wrapper.find('form').trigger('submit')

      // Wait for error to be displayed
      await flushPromises()

      // Verify error message
      const errorMessage = wrapper.find('.error-message')
      expect(errorMessage.exists()).toBe(true)
      expect(errorMessage.text()).toBe('Failed to create bug')
    })
  })

  describe('Bug status update', () => {
    it('should update bug status when user has permission', async () => {
      // Mount component with initialized stores
      wrapper = mount(BugList, {
        global: {
          plugins: [createTestingPinia({
            initialState: {
              auth: {
                isAuthenticated: true,
                currentUser: { id: 1, name: 'Test User' },
                role: 'manager',
                isManager: true,
                canEditBugStatus: true
              },
              bug: {
                bugs: [{
                  id: 1,
                  title: 'Test Bug',
                  status: 'OPEN',
                  assigned_to: {
                    id: 1,
                    name: 'Test User'
                  }
                }]
              }
            }
          })],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: router
          }
        }
      })

      // Wait for component to render
      await flushPromises()

      // Find and click update button
      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(true)
      await updateButton.trigger('click')

      // Select new status
      const statusSelect = wrapper.find('[data-testid="status-select"]')
      await statusSelect.setValue('IN_PROGRESS')
      
      // Confirm update
      const confirmButton = wrapper.find('[data-testid="confirm-status-update"]')
      await confirmButton.trigger('click')

      // Wait for API call
      await flushPromises()

      // Verify API call
      expect(axios.put).toHaveBeenCalledWith(
        '/api/bugs/1/status',
        { status: 'IN_PROGRESS' }
      )
    })

    it('should not show update button without permission', async () => {
      const bug = {
        id: 1,
        title: 'Test Bug',
        description: 'Test Description',
        status: 'OPEN',
        priority: 'high',
        assigned_to: {
          id: 2, // Different from current user
          name: 'other user'
        }
      }

      // Mount with no special permissions
      wrapper = mountComponent({
        auth: {
          isManager: false,
          isAdmin: false,
          canEditBugStatus: false
        },
        bug: {
          bugs: [bug]
        }
      })

      await flushPromises()

      // Verify update button is not visible
      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(false)
    })

    it('should update bug status successfully', async () => {
      // Mock successful status update response
      const bug = {
        id: 1,
        title: 'Test Bug',
        description: 'Test Description',
        status: 'OPEN',
        assigned_to: {
          id: 1,
          name: 'testuser'
        }
      }

      const updatedBug = { ...bug, status: 'IN_PROGRESS' }
      axios.put.mockResolvedValueOnce({ data: updatedBug })

      // Set up auth store with manager role (can always update status)
      authStore.$patch({
        user: {
          id: 1,
          name: 'testuser',
          role: 'manager'
        },
        role: 'manager',
        token: 'test-token'
      })

      // Set up bug store with initial bug
      bugStore.$patch({
        bugs: [bug],
        loading: false,
        error: null
      })

      // Mount component with initialized stores
      wrapper = mount(BugList, {
        global: {
          plugins: [createTestingPinia({
            initialState: {
              auth: {
                user: {
                  id: 1,
                  name: 'testuser',
                  role: 'manager'
                },
                role: 'manager',
                token: 'test-token'
              },
              bug: {
                bugs: [{
                  id: 1,
                  title: 'Test Bug',
                  description: 'Test Description',
                  status: 'OPEN',
                  assigned_to: {
                    id: 1,
                    name: 'testuser'
                  }
                }],
                loading: false,
                error: null
              }
            }
          })],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: router
          }
        }
      })

      // Wait for component to render
      await flushPromises()

      // Debug output
      console.log('Auth store state:', {
        role: authStore.role,
        isManager: authStore.isManager,
        canEditBugStatus: authStore.canEditBugStatus,
        currentUser: authStore.currentUser
      })
      console.log('Bug store state:', {
        bugs: bugStore.bugs
      })
      console.log('Button exists:', wrapper.find('[data-testid="update-status-button"]').exists())

      // Find and click the update button
      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(true)
      await updateButton.trigger('click')

      // Select new status and confirm
      const statusSelect = wrapper.find('[data-testid="status-select"]')
      await statusSelect.setValue('IN_PROGRESS')
      
      const confirmButton = wrapper.find('[data-testid="confirm-status-update"]')
      await confirmButton.trigger('click')

      // Wait for API call to complete
      await flushPromises()

      // Verify API call and store update
      expect(axios.put).toHaveBeenCalledWith(
        '/api/bugs/1/status',
        { status: 'IN_PROGRESS' }
      )
      expect(bugStore.bugs[0].status).toBe('IN_PROGRESS')
    })

    it('should cancel status update', async () => {
      // Set up auth store with manager role (can always update status)
      authStore.$patch({
        user: {
          id: 1,
          name: 'testuser',
          role: 'manager'
        },
        role: 'manager',
        token: 'test-token'
      })

      // Set up initial bug data
      const bug = {
        id: 1,
        title: 'Test Bug',
        description: 'Test Description',
        status: 'OPEN',
        priority: 'high',
        assigned_to: {
          id: 1,
          name: 'testuser'
        }
      }

      // Set up bug store with initial bug
      bugStore.$patch({
        bugs: [bug],
        loading: false,
        error: null
      })

      // Render the component
      wrapper = mount(BugList, {
        global: {
          plugins: [createTestingPinia()],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: router
          }
        }
      })

      // Wait for the component to update
      await flushPromises()

      // Debug output
      console.log('Auth store state:', {
        role: authStore.role,
        isManager: authStore.isManager,
        canEditBugStatus: authStore.canEditBugStatus,
        currentUser: authStore.currentUser
      })
      console.log('Bug store state:', {
        bugs: bugStore.bugs
      })
      console.log('Button exists:', wrapper.find('[data-testid="update-status-button"]').exists())

      // Find and click the update status button
      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(true)
      await updateButton.trigger('click')

      // Wait for the modal to appear
      await flushPromises()

      // Click cancel button
      const cancelButton = wrapper.find('[data-testid="cancel-status-update"]')
      await cancelButton.trigger('click')

      // Wait for the modal to close
      await flushPromises()

      // Verify the modal is closed
      expect(wrapper.find('[data-testid="status-select"]').exists()).toBe(false)

      // Verify the bug status wasn't changed
      expect(bugStore.bugs[0].status).toBe('OPEN')
    })
  })

  describe('Update status button visibility', () => {
    it('should show update status button for manager regardless of assignment', async () => {
      // Mount with manager permissions
      wrapper = mount(BugList, {
        global: {
          plugins: [createTestingPinia({
            initialState: {
              auth: {
                isAuthenticated: true,
                currentUser: { id: 1, name: 'Manager' },
                role: 'manager',
                isManager: true,
                canEditBugStatus: true
              },
              bug: {
                bugs: [{
                  id: 1,
                  title: 'Test Bug',
                  status: 'OPEN',
                  assigned_to: {
                    id: 2, // Different from current user
                    name: 'Other User'
                  }
                }]
              }
            }
          })],
          stubs: {
            'router-link': true,
            'router-view': true
          },
          mocks: {
            $router: router
          }
        }
      })

      await flushPromises()

      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(true)
    })

    it('should show update status button for assigned developer with permission', async () => {
      // Set up developer assigned to bug
      authStore.$patch({
        isManager: false,
        canEditBugStatus: true,
        currentUser: {
          username: 'testuser'
        }
      })

      bugStore.$patch({
        bugs: [{
          id: 1,
          assignee: 'testuser'
        }]
      })

      await wrapper.vm.$nextTick()

      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(true)
    })

    it('should hide update status button for non-assigned developer even with permission', async () => {
      // Set up developer not assigned to bug
      authStore.$patch({
        isManager: false,
        canEditBugStatus: true,
        currentUser: {
          username: 'otheruser'
        }
      })

      await wrapper.vm.$nextTick()

      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(false)
    })

    it('should hide update status button when user lacks canEditBugStatus permission', async () => {
      // Remove edit permission
      authStore.$patch({
        canEditBugStatus: false
      })

      await wrapper.vm.$nextTick()

      const updateButton = wrapper.find('[data-testid="update-status-button"]')
      expect(updateButton.exists()).toBe(false)
    })
  })
}) 