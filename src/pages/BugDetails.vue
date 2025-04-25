<template>
  <div class="bug-details">
    <div v-if="loading" class="loading">
      Loading bug details...
    </div>
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    <div v-else-if="bug" class="bug-content">
      <div class="header">
        <h1>{{ bug.title }}</h1>
        <span :class="['status', bug.status]">{{ bug.status }}</span>
      </div>

      <div class="meta">
        <div class="meta-item">
          <strong>Priority:</strong> {{ bug.priority }}
        </div>
        <div class="meta-item">
          <strong>Reported by:</strong> {{ bug.reported_by.name }}
        </div>
        <div class="meta-item">
          <strong>Assigned to:</strong> 
          {{ bug.assigned_to ? bug.assigned_to.name : 'Unassigned' }}
        </div>
        <div class="meta-item">
          <strong>Created:</strong> 
          {{ new Date(bug.created_at).toLocaleDateString() }}
        </div>
      </div>

      <div class="description">
        <h2>Description</h2>
        <p>{{ bug.description }}</p>
      </div>

      <div class="actions">
        <button 
          v-if="canUpdateStatus"
          @click="openStatusModal"
          class="btn btn-primary"
        >
          Update Status
        </button>
        <button 
          v-if="authStore.canAssignBugs && bug.status !== 'resolved'"
          @click="openAssignModal"
          class="btn btn-secondary"
        >
          Assign Bug
        </button>
      </div>
    </div>

    <!-- Status Modal -->
    <div v-if="showStatusModal" class="modal">
      <div class="modal-content">
        <h3>Update Status</h3>
        <select v-model="selectedStatus" class="form-control">
          <option value="open">Open</option>
          <option value="in-progress">In Progress</option>
          <option value="resolved">Resolved</option>
        </select>
        <div class="modal-actions">
          <button @click="closeStatusModal" class="btn btn-secondary">Cancel</button>
          <button @click="updateStatus" class="btn btn-primary">Update</button>
        </div>
      </div>
    </div>

    <!-- Assignment Modal -->
    <div v-if="showAssignModal" class="modal">
      <div class="modal-content">
        <h3>Assign Bug</h3>
        <select v-model="selectedDeveloper" class="form-control">
          <option value="">Select Developer</option>
          <option 
            v-for="dev in developers" 
            :key="dev.id" 
            :value="dev.id"
          >
            {{ dev.name }}
          </option>
        </select>
        <div class="modal-actions">
          <button @click="closeAssignModal" class="btn btn-secondary">Cancel</button>
          <button @click="assignBug" class="btn btn-primary">Assign</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-
</script> 