<template>
  <MainLayout>
    <div class="dashboard-grid">
      <div class="card">
        <h3 class="card-title">Total Bugs</h3>
        <p class="card-number text-red">{{ bugStore.bugs.length }}</p>
      </div>
      <div class="card">
        <h3 class="card-title">In Progress</h3>
        <p class="card-number text-yellow">{{ inProgressCount }}</p>
      </div>
      <div class="card">
        <h3 class="card-title">Resolved</h3>
        <p class="card-number text-green">{{ resolvedCount }}</p>
      </div>
    </div>

    <!-- Recent Bugs Section -->
    <div class="recent-bugs">
      <h2>Recent Bugs</h2>
      <div v-if="bugStore.loading" class="loading">
        Loading recent bugs...
      </div>
      <div v-else-if="bugStore.error" class="error">
        {{ bugStore.error }}
      </div>
      <div v-else class="bugs-list">
        <div v-for="bug in recentBugs" :key="bug.id" class="bug-card">
          <div class="bug-header">
            <h3>{{ bug.title }}</h3>
            <span :class="['status', bug.status]">{{ bug.status }}</span>
          </div>
          <p class="description">{{ bug.description }}</p>
          <div class="bug-meta">
            <span>Reported by: {{ bug.reported_by?.name }}</span>
            <span v-if="bug.assigned_to">Assigned to: {{ bug.assigned_to.name }}</span>
            <span v-else class="unassigned">Unassigned</span>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import MainLayout from '../layouts/MainLayout.vue';
import { useBugStore } from '../stores/bug';

const bugStore = useBugStore();

// Fetch bugs when component mounts
onMounted(async () => {
  await bugStore.fetchBugs();
});

// Get the 5 most recent bugs
const recentBugs = computed(() => {
  return [...bugStore.bugs]
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    .slice(0, 5);
});

const inProgressCount = computed(() => {
  return bugStore.bugs.filter(bug => bug.status === 'in-progress').length;
});

const resolvedCount = computed(() => {
  return bugStore.bugs.filter(bug => bug.status === 'resolved').length;
});
</script>

<style scoped>
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #4a5568;
}

.card-number {
  font-size: 2rem;
  font-weight: 700;
}

.text-red {
  color: #e74c3c;
}

.text-yellow {
  color: #f1c40f;
}

.text-green {
  color: #2ecc71;
}

.recent-bugs {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.recent-bugs h2 {
  margin-bottom: 1.5rem;
  color: #2d3748;
}

.bugs-list {
  display: grid;
  gap: 1rem;
}

.bug-card {
  padding: 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
}

.bug-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 0.5rem;
}

.bug-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #2d3748;
}

.status {
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.875rem;
  text-transform: capitalize;
}

.status.open {
  background-color: #fff3cd;
  color: #856404;
}

.status.in-progress {
  background-color: #cce5ff;
  color: #004085;
}

.status.resolved {
  background-color: #d4edda;
  color: #155724;
}

.description {
  color: #4a5568;
  margin-bottom: 0.75rem;
  font-size: 0.875rem;
}

.bug-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
  color: #718096;
}

.unassigned {
  color: #e53e3e;
  font-style: italic;
}

.loading {
  text-align: center;
  color: #718096;
  padding: 2rem;
}

.error {
  text-align: center;
  color: #e53e3e;
  padding: 2rem;
}
</style>
  