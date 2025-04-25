import { setActivePinia, createPinia } from 'pinia'
import { useBugStore } from '@/stores/bug'
import api from '@/services/api'

// Mock the API
jest.mock('@/services/api')

describe('Bug Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test('fetches bugs successfully', async () => {
    const store = useBugStore()
    const mockBugs = [
      { id: '1', title: 'Bug 1' },
      { id: '2', title: 'Bug 2' }
    ]

    api.get.mockResolvedValueOnce({ data: mockBugs })

    await store.fetchBugs()

    expect(store.bugs).toEqual(mockBugs)
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  test('handles fetch bugs error', async () => {
    const store = useBugStore()
    const error = new Error('Failed to fetch')
    api.get.mockRejectedValueOnce(error)

    await expect(store.fetchBugs()).rejects.toThrow()
    expect(store.loading).toBe(false)
    expect(store.error).toBeTruthy()
  })

  test('creates bug successfully', async () => {
    const store = useBugStore()
    const newBug = {
      title: 'New Bug',
      description: 'Description',
      priority: 'high'
    }
    const createdBug = { ...newBug, id: '1' }

    api.post.mockResolvedValueOnce({ data: createdBug })

    const result = await store.createBug(newBug)

    expect(result).toEqual(createdBug)
    expect(store.bugs).toContainEqual(createdBug)
  })

  test('updates bug status successfully', async () => {
    const store = useBugStore()
    const bug = { id: '1', status: 'open' }
    const updatedBug = { ...bug, status: 'in-progress' }

    store.bugs = [bug]
    api.patch.mockResolvedValueOnce({ data: updatedBug })

    await store.updateBugStatus('1', 'in-progress')

    expect(store.bugs[0].status).toBe('in-progress')
  })

  test('assigns bug successfully', async () => {
    const store = useBugStore()
    const bug = { id: '1', assigned_to: null }
    const updatedBug = { ...bug, assigned_to: { id: 'dev1' } }

    store.bugs = [bug]
    api.post.mockResolvedValueOnce({ data: updatedBug })

    await store.assignBug('1', 'dev1')

    expect(store.bugs[0].assigned_to).toEqual({ id: 'dev1' })
  })
}) 