<template>
    <div class="not-found-container">
      <div class="not-found-content">
        <h1 class="not-found-title">404</h1>
        <p class="not-found-text">Oops! Page not found.</p>
        <router-link to="/" class="btn btn-primary">
          Go to Dashboard
        </router-link>
      </div>
    </div>
  </template>
  
  <style scoped>
  .not-found-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background-color: #f5f5f5;
  }
  
  .not-found-content {
    text-align: center;
  }
  
  .not-found-title {
    font-size: 6rem;
    font-weight: 700;
    color: #e74c3c;
    line-height: 1;
  }
  
  .not-found-text {
    font-size: 1.25rem;
    color: #4a5568;
    margin: 1rem 0;
  }
  </style>
  