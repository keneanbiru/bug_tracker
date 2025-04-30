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

        async createBug(bugData) {
            this.loading = true;
            this.error = null;
            try {
                console.log('Creating bug:', bugData);
                const response = await api.post('/bugs', bugData);
                console.log('Create bug response:', response.data);
                this.bugs.push(response.data);
                return response.data;
            } catch (error) {
                console.error('Error creating bug:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status
                });
                this.error = error.response?.data?.message || 'Failed to create bug';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async updateBugStatus(bugId, status) {
            this.loading = true;
            this.error = null;
            try {
                // Log the current bug data
                const currentBug = this.bugs.find(bug => bug.id === bugId);
                console.log('Current bug data:', currentBug);
                console.log('Updating bug status:', { 
                    bugId, 
                    status,
                    currentStatus: currentBug?.status,
                    bugExists: !!currentBug
                });
                
                // Validate status
                const validStatuses = ['OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED'];
                if (!validStatuses.includes(status)) {
                    throw new Error('Invalid status value');
                }

                // Map status to backend format
                const statusMap = {
                    'OPEN': 'open',
                    'IN_PROGRESS': 'in-progress',
                    'RESOLVED': 'resolved'
                };

                const backendStatus = statusMap[status];
                if (!backendStatus) {
                    throw new Error('Invalid status value');
                }

                // Make the API call with the correct request format
                const requestData = {
                    status: backendStatus
                };

                console.log('Sending request:', {
                    url: `/bugs/${bugId}/status`,
                    method: 'PATCH',
                    data: requestData
                });

                const response = await api.patch(`/bugs/${bugId}/status`, requestData);
                
                console.log('Update status response:', response.data);
                
                if (!response.data) {
                    throw new Error('No response data received');
                }

                // Update the bug in the store
                const index = this.bugs.findIndex(bug => bug.id === bugId);
                if (index !== -1) {
                    const updatedBug = { 
                        ...this.bugs[index], 
                        status: status // Keep the frontend status format
                    };
                    this.bugs[index] = updatedBug;
                    console.log('Updated bug in store:', updatedBug);
                } else {
                    console.warn('Bug not found in store:', bugId);
                }
                
                return response.data;
            } catch (error) {
                console.error('Error updating bug status:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status,
                    config: {
                        url: error.config?.url,
                        method: error.config?.method,
                        data: error.config?.data,
                        headers: error.config?.headers
                    }
                });

                // Handle specific error cases
                if (error.response?.status === 404) {
                    this.error = 'Bug not found';
                } else if (error.response?.status === 403) {
                    this.error = 'You do not have permission to update this bug';
                } else if (error.response?.status === 400) {
                    this.error = error.response.data.message || 'Invalid status value';
                } else {
                    this.error = error.response?.data?.message || 'Failed to update bug status';
                }
                throw error;
            } finally {
                this.loading = false;
            }
        },

        async assignBug(bugId, developerId) {
            this.loading = true;
            this.error = null;
            try {
                console.log('Assigning bug in store:', { 
                    bugId, 
                    developerId,
                    bugIdType: typeof bugId,
                    developerIdType: typeof developerId,
                    bugExists: this.bugs.some(bug => bug.id === bugId)
                });
                const requestData = {
                    developer_id: developerId
                };
                console.log('Sending request:', {
                    url: `/bugs/${bugId}/assign`,
                    method: 'POST',
                    data: requestData
                });
                const response = await api.post(`/bugs/${bugId}/assign`, requestData);
                console.log('Assign bug response:', response.data);
                const index = this.bugs.findIndex(bug => bug.id === bugId);
                if (index !== -1) {
                    this.bugs[index] = { ...this.bugs[index], ...response.data };
                }
                return response.data;
            } catch (error) {
                console.error('Error assigning bug:', {
                    message: error.message,
                    response: error.response?.data,
                    status: error.response?.status,
                    config: {
                        url: error.config?.url,
                        method: error.config?.method,
                        data: error.config?.data
                    }
                });
                this.error = error.response?.data?.message || 'Failed to assign bug';
                throw error;
            } finally {
                this.loading = false;
            }
        }
    }
}); 