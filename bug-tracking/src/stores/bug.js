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
        inProgressBugs: (state) => state.bugs.filter(bug => bug.status === 'IN_PROGRESS'),
        resolvedBugs: (state) => state.bugs.filter(bug => bug.status === 'RESOLVED'),
        openBugs: (state) => state.bugs.filter(bug => bug.status === 'OPEN'),
        closedBugs: (state) => state.bugs.filter(bug => bug.status === 'CLOSED')
    },

    actions: {
        async fetchBugs() {
            this.loading = true;
            this.error = null;
            try {
                console.log('Fetching bugs...');
                const response = await api.get('/bugs');
                console.log('Bugs response:', response.data);
                this.bugs = response.data;
            } catch (error) {
                console.error('Error fetching bugs:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status
                });
                this.error = error.response?.data?.message || 'Failed to fetch bugs';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async fetchBugById(id) {
            this.loading = true;
            try {
                console.log('Fetching bug by ID:', id);
                const response = await api.get(`/bugs/${id}`);
                console.log('Bug response:', response.data);
                this.currentBug = response.data;
                return response.data;
            } catch (error) {
                console.error('Error fetching bug:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status
                });
                this.error = error.response?.data?.message || 'Failed to fetch bug';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async reportBug(bugData) {
            this.loading = true;
            this.error = null;
            try {
                console.log('Reporting bug:', bugData);
                const response = await api.post('/bugs', bugData);
                console.log('Bug report response:', response.data);
                // Add the new bug to the list
                this.bugs.push(response.data);
                return response.data;
            } catch (error) {
                console.error('Error reporting bug:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status
                });
                this.error = error.response?.data?.message || 'Failed to report bug';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async updateBugStatus(bugId, status) {
            this.loading = true;
            this.error = null;
            try {
                console.log('Updating bug status:', { bugId, status });
                const response = await api.patch(`/bugs/${bugId}/status`, { status });
                console.log('Status update response:', response.data);
                
                // Update the bug in the list
                const index = this.bugs.findIndex(bug => bug.id === bugId);
                if (index !== -1) {
                    this.bugs[index] = { ...this.bugs[index], status: response.data.status };
                }
                
                return response.data;
            } catch (error) {
                console.error('Error updating bug status:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status
                });
                this.error = error.response?.data?.message || 'Failed to update bug status';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async assignBug(bugId, developerId) {
            this.loading = true;
            this.error = null;
            try {
                console.log('Assigning bug:', { bugId, developerId });
                const response = await api.post(`/bugs/${bugId}/assign`, { developerId });
                console.log('Bug assignment response:', response.data);
                
                // Update the bug in the list
                const index = this.bugs.findIndex(bug => bug.id === bugId);
                if (index !== -1) {
                    this.bugs[index] = { ...this.bugs[index], assigned_to: response.data.assigned_to };
                }
                
                return response.data;
            } catch (error) {
                console.error('Error assigning bug:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status
                });
                this.error = error.response?.data?.message || 'Failed to assign bug';
                throw error;
            } finally {
                this.loading = false;
            }
        }
    }
}); 