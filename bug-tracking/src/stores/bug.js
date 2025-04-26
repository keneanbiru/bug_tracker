import { defineStore } from 'pinia';
import api from '../services/api';

export const useBugStore = defineStore('bug', {
    state: () => ({
        bugs: [],
        loading: false,
        error: null,
        currentBug: null
    }),

    getters: {
        inProgressBugs: (state) => state.bugs.filter(bug => bug.status === 'in_progress'),
        resolvedBugs: (state) => state.bugs.filter(bug => bug.status === 'resolved'),
        openBugs: (state) => state.bugs.filter(bug => bug.status === 'open'),
        closedBugs: (state) => state.bugs.filter(bug => bug.status === 'closed')
    },

    actions: {
        async fetchBugs() {
            this.loading = true;
            try {
                console.log('Fetching bugs...');
                console.log('Auth token:', localStorage.getItem('token'));
                const response = await api.get('/bugs');
                console.log('Bugs response:', response.data);
                this.bugs = response.data;
            } catch (error) {
                console.error('Error fetching bugs:', error);
                console.error('Error response:', error.response);
                this.error = error.response?.data?.error || 'Failed to fetch bugs';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async fetchBugById(id) {
            this.loading = true;
            try {
                const response = await api.get(`/bugs/${id}`);
                this.currentBug = response.data;
                return response.data;
            } catch (error) {
                this.error = error.response?.data?.error || 'Failed to fetch bug';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async createBug(bugData) {
            try {
                const response = await api.post('/bugs', bugData);
                this.bugs.push(response.data);
                return response.data;
            } catch (error) {
                this.error = error.response?.data?.error || 'Failed to create bug';
                throw error;
            }
        },

        async updateBugStatus(bugId, status) {
            try {
                console.log('Updating bug status:', { bugId, status });
                
                // Try with a different format for the status value
                const statusMap = {
                    'OPEN': 'open',
                    'IN_PROGRESS': 'in_progress',
                    'RESOLVED': 'resolved',
                    'CLOSED': 'closed'
                };
                
                const mappedStatus = statusMap[status] || status;
                console.log('Mapped status:', mappedStatus);
                
                // Get the current bug
                const currentBug = this.bugs.find(b => b.id === bugId);
                if (!currentBug) {
                    throw new Error('Bug not found');
                }
                
                // Create a complete updated bug object
                const updatedBug = {
                    ...currentBug,
                    status: mappedStatus
                };
                
                const response = await api.put(`/bugs/${bugId}`, updatedBug);
                
                // Update the bug in the store
                const index = this.bugs.findIndex(b => b.id === bugId);
                if (index !== -1) {
                    this.bugs[index] = response.data;
                }
                
                return response.data;
            } catch (error) {
                this.error = error.response?.data?.error || 'Failed to update bug status';
                throw error;
            }
        },

        async assignBug(bugId, developerId) {
            try {
                const response = await api.post(`/bugs/${bugId}/assign`, { developer_id: developerId });
                const index = this.bugs.findIndex(b => b.id === bugId);
                if (index !== -1) {
                    this.bugs[index] = response.data;
                }
                return response.data;
            } catch (error) {
                this.error = error.response?.data?.error || 'Failed to assign bug';
                throw error;
            }
        }
    }
}); 