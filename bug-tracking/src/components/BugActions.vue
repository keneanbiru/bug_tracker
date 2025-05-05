<template>
    <div class="bug-actions" v-if="isManager">
        <!-- Delete Button -->
        <button 
            @click="confirmDelete" 
            class="btn btn-danger btn-sm"
            :disabled="loading"
        >
            <i class="fas fa-trash"></i> Remove
        </button>

        <!-- Reassign Button - Only show for assigned bugs -->
        <button 
            v-if="hasAssignee"
            @click="handleReassignClick" 
            class="btn btn-info btn-sm"
            :disabled="loading"
        >
            <i class="fas fa-user-edit"></i> Reassign
        </button>

        <!-- Reassign Modal -->
        <div v-if="showReassignModal" class="modal-overlay">
            <div class="modal-content">
                <h3>Reassign Bug to Developer</h3>
                <div class="form-group">
                    <label for="developer">Select Developer:</label>
                    <select 
                        v-model="selectedDeveloper" 
                        class="form-control" 
                        id="developer"
                    >
                        <option value="">Choose a developer...</option>
                        <option 
                            v-for="dev in developers" 
                            :key="dev.id" 
                            :value="dev.id"
                        >
                            {{ dev.name }}
                        </option>
                    </select>
                </div>
                <div class="modal-actions">
                    <button 
                        @click="showReassignModal = false" 
                        class="btn btn-secondary"
                    >
                        Cancel
                    </button>
                    <button 
                        @click="handleReassign"
                        class="btn btn-primary"
                        :disabled="!selectedDeveloper || loading"
                    >
                        Reassign
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useBugStore } from '../stores/bug';
import { storeToRefs } from 'pinia';

const props = defineProps({
    bugId: {
        type: String,
        required: true
    },
    hasAssignee: {
        type: Boolean,
        required: true
    }
});

const emit = defineEmits(['bugDeleted', 'bugReassigned']);

const authStore = useAuthStore();
const bugStore = useBugStore();
const { loading } = storeToRefs(bugStore);

const showReassignModal = ref(false);
const selectedDeveloper = ref('');
const developers = ref([]);

const isManager = computed(() => {
    return ['manager', 'admin'].includes(authStore.role);
});

// Fetch developers when component is mounted
const fetchDevelopers = async () => {
    try {
        const response = await fetch('https://bug-tracker-8.onrender.com/api/auth/developers', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        if (!response.ok) throw new Error('Failed to fetch developers');
        developers.value = await response.json();
    } catch (error) {
        console.error('Error fetching developers:', error);
    }
};

// Call fetchDevelopers when the reassign modal is opened
const handleReassignClick = () => {
    showReassignModal.value = true;
    fetchDevelopers();
};

const confirmDelete = async () => {
    if (confirm('Are you sure you want to delete this bug?')) {
        try {
            await bugStore.deleteBug(props.bugId);
            emit('bugDeleted');
        } catch (error) {
            alert(error.response?.data?.error || 'Failed to delete bug');
        }
    }
};

const handleReassign = async () => {
    if (!selectedDeveloper.value) return;
    
    try {
        await bugStore.reassignBug(props.bugId, selectedDeveloper.value);
        showReassignModal.value = false;
        selectedDeveloper.value = '';
        emit('bugReassigned');
    } catch (error) {
        alert(error.response?.data?.error || 'Failed to reassign bug');
    }
};
</script>

<style scoped>
.bug-actions {
    display: inline-flex;
    gap: 0.5rem;
    margin-left: 0.5rem;
}

.btn {
    white-space: nowrap;
}

.btn-danger {
    background-color: #dc3545;
    color: white;
}

.btn-danger:hover {
    background-color: #c82333;
}

.btn-info {
    background-color: #17a2b8;
    color: white;
}

.btn-info:hover {
    background-color: #138496;
}

.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.modal-content {
    background: white;
    padding: 2rem;
    border-radius: 12px;
    width: 100%;
    max-width: 500px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-content h3 {
    margin-top: 0;
    margin-bottom: 1.5rem;
    color: #2d3748;
}

.form-group {
    margin-bottom: 1.5rem;
}

.form-group label {
    display: block;
    margin-bottom: 0.5rem;
    color: #4a5568;
    font-weight: 500;
}

.form-control {
    width: 100%;
    padding: 0.75rem;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    font-size: 1rem;
    transition: all 0.3s ease;
}

.form-control:focus {
    outline: none;
    border-color: #4299e1;
    box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.2);
}

.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
}

.btn i {
    margin-right: 0.25rem;
}

/* Make buttons more compact */
.btn-sm {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
}
</style> 