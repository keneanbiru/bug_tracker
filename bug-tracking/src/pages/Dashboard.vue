<template>
    <MainLayout>
      <div class="dashboard-grid">
        <div class="card">
          <h3 class="card-title">Total Bugs</h3>
          <p class="card-number text-red">12</p>
        </div>
        <div class="card">
          <h3 class="card-title">In Progress</h3>
          <p class="card-number text-yellow">5</p>
        </div>
        <div class="card">
          <h3 class="card-title">Resolved</h3>
          <p class="card-number text-green">7</p>
        </div>
      </div>
    </MainLayout>
  </template>
  
  <script>
  import MainLayout from '../layouts/MainLayout.vue';
  
  export default {
    components: { MainLayout }
  };
  </script>
  
  <style scoped>
  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
  }
  
  .card-title {
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }
  
  .card-number {
    font-size: 1.875rem;
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
  </style>
  