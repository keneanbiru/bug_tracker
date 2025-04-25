import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import BugList from '@/pages/BugList.vue'
import { useBugStore } from '@/stores/bug'
import { useAuthStore } from '@/stores/auth'

describe('BugList.vue', () => {
  const createWrapper = (initialState = {}) => {
    return mount(BugList, {
      global: {
        plugins: [
          createTestingPinia({
            initialState: {
              bug: {
                bugs: [],
                loading: false,
                error: null,
                ...initialState.bug
              },
              auth: {
                user: null,
                role: null,
                ...initialState.auth
              }
            }
          })
        ]
      }
    })
  }

  test('displays loading state', () => {
    const wrapper = createWrapper({
      bug: { loading: true }
    })
    expect(wrapper.find('.loading').exists()).toBe(true)
  })

  test('displays error message', () => {
    const wrapper = createWrapper({
      bug: { error: 'Failed to load bugs' }
    })
    expect(wrapper.find('.error').text()).toBe('Failed to load bugs')
  })

  test('displays bug list', () => {
    const bugs = [
      {
        id: '1',
        title: 'Test Bug',
        description: 'Test Description',
        status: 'open',
        priority: 'high',
        reported_by: { name: 'John' },
        assigned_to: null
      }
    ]

    const wrapper = createWrapper({
      bug: { bugs }
    })

    expect(wrapper.findAll('.bug-card')).toHaveLength(1)
    expect(wrapper.find('.bug-header h3').text()).toBe('Test Bug')
  })

  test('filters bugs by status', async () => {
    const bugs = [
      { id: '1', status: 'open', title: 'Bug 1' },
      { id: '2', status: 'in-progress', title: 'Bug 2' }
    ]

    const wrapper = createWrapper({
      bug: { bugs }
    })

    await wrapper.find('.status-filter select').setValue('open')
    expect(wrapper.findAll('.bug-card')).toHaveLength(1)
    expect(wrapper.find('.bug-header h3').text()).toBe('Bug 1')
  })

  test('shows assign button only for managers', () => {
    const wrapper = createWrapper({
      auth: { role: 'manager' },
      bug: {
        bugs: [{
          id: '1',
          status: 'open',
          title: 'Bug 1'
        }]
      }
    })

    expect(wrapper.find('button:contains("Assign")').exists()).toBe(true)
  })

  test('shows update status button only for assigned developer', () => {
    const wrapper = createWrapper({
      auth: {
        role: 'developer',
        currentUser: { id: 'dev1' }
      },
      bug: {
        bugs: [{
          id: '1',
          status: 'open',
          title: 'Bug 1',
          assigned_to: { id: 'dev1' }
        }]
      }
    })

    expect(wrapper.find('button:contains("Update Status")').exists()).toBe(true)
  })

  test('opens status update modal when update status button is clicked', async () => {
    const wrapper = createWrapper({
      auth: {
        role: 'developer',
        currentUser: { id: 'dev1' }
      },
      bug: {
        bugs: [{
          id: '1',
          status: 'open',
          title: 'Bug 1',
          assigned_to: { id: 'dev1' }
        }]
      }
    })

    await wrapper.find('button:contains("Update Status")').trigger('click')
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
    expect(wrapper.find('.modal-content h3').text()).toBe('Update Bug Status')
  })

  test('updates bug status when status is selected and update button is clicked', async () => {
    // Mock fetch function
    global.fetch = jest.fn().mockImplementation(() => 
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({
          id: '1',
          status: 'in-progress',
          title: 'Bug 1',
          assigned_to: { id: 'dev1' }
        })
      })
    )

    const wrapper = createWrapper({
      auth: {
        role: 'developer',
        currentUser: { id: 'dev1' },
        token: 'test-token'
      },
      bug: {
        bugs: [{
          id: '1',
          status: 'open',
          title: 'Bug 1',
          assigned_to: { id: 'dev1' }
        }]
      }
    })

    // Open the modal
    await wrapper.find('button:contains("Update Status")').trigger('click')
    
    // Select a status
    await wrapper.find('.modal-content select').setValue('IN_PROGRESS')
    
    // Click update button
    await wrapper.find('.modal-actions button:contains("Update")').trigger('click')
    
    // Check if fetch was called with the correct parameters
    expect(global.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/bugs/1/status',
      expect.objectContaining({
        method: 'PATCH',
        headers: expect.objectContaining({
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token'
        }),
        body: JSON.stringify({ status: 'in-progress' })
      })
    )
    
    // Check if the modal is closed
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })
}) 