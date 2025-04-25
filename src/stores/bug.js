import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

export const useBugStore = defineStore('bug', {
  state: () => ({
    bugs: [],
    loading: false,
    error: null
  }),

  getters: {
    assignedBugs: (state) => {
      const authStore = useAuthStore()
      return state.bugs.filter(bug => 
        bug.assigned_to?.id === authStore.currentUser?.id
      )
    },
    openBugs: (state) => 
      state.bugs.filter(bug => bug.status === 'open'),
    inProgressBugs: (state) => 
      state.bugs.filter(bug => bug.status === 'in-progress'),
    resolvedBugs: (state) => 
      state.bugs.filter(bug => bug.status === 'resolved')
  },

  actions: {
    async fetchBugs() {
      this.loading = true
      try {
        const response = await fetch('http://localhost:8080/api/bugs', {
          headers: {
            'Authorization': `Bearer ${useAuthStore().token}`
          }
        })
        if (!response.ok) throw new Error('Failed to fetch bugs')
        this.bugs = await response.json()
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async createBug(bugData) {
      try {
        const response = await fetch('http://localhost:8080/api/bugs', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${useAuthStore().token}`
          },
          body: JSON.stringify(bugData)
        })
        if (!response.ok) throw new Error('Failed to create bug')
        const newBug = await response.json()
        this.bugs.push(newBug)
        return newBug
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async updateBugStatus(bugId, status) {
      try {
        const response = await fetch(`http://localhost:8080/api/bugs/${bugId}/status`, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${useAuthStore().token}`
          },
          body: JSON.stringify({ status })
        })
        if (!response.ok) throw new Error('Failed to update bug status')
        const updatedBug = await response.json()
        const index = this.bugs.findIndex(b => b.id === bugId)
        if (index !== -1) {
          this.bugs[index] = updatedBug
        }
        return updatedBug
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async assignBug(bugId, developerId) {
      try {
        const response = await fetch(`http://localhost:8080/api/bugs/${bugId}/assign`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${useAuthStore().token}`
          },
          body: JSON.stringify({ developer_id: developerId })
        })
        if (!response.ok) throw new Error('Failed to assign bug')
        const updatedBug = await response.json()
        const index = this.bugs.findIndex(b => b.id === bugId)
        if (index !== -1) {
          this.bugs[index] = updatedBug
        }
        return updatedBug
      } catch (error) {
        this.error = error.message
        throw error
      }
    }
  }
}) 