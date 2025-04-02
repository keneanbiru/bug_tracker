<template>
  <div class="app-container">
    <!-- Navigation Bar (Simple) -->
    <nav class="nav">
      <div class="container">
        <router-link to="/" class="nav-brand">Bug Tracker</router-link>
        <div class="nav-links">
          <router-link to="/" class="nav-link">Dashboard</router-link>
          <router-link to="/bugs" class="nav-link">Bug List</router-link>
          <router-link to="/report" class="nav-link">Report Bug</router-link>
          <router-link to="/login" class="nav-link">Login</router-link>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: 'App'
};
</script>

<style>
/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #f5f5f5;
}

/* App container */
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Navigation */
.nav {
  background-color: #3498db;
  padding: 1rem 0;
}

.nav .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-brand {
  color: white;
  font-size: 1.25rem;
  font-weight: 700;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 1rem;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
}

.nav-link:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

/* Main content */
.main-content {
  flex: 1;
  padding: 1.5rem 0;
}

/* Container */
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

/* Buttons */
.btn {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover {
  background-color: #2980b9;
}

.btn-danger {
  background-color: #e74c3c;
  color: white;
}

.btn-danger:hover {
  background-color: #c0392b;
}

/* Forms */
.form-group {
  margin-bottom: 1rem;
}

.form-control {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 0.25rem;
  font-size: 1rem;
}

.form-label {
  display: block;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

/* Cards */
.card {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1rem;
  margin-bottom: 1rem;
}
</style>
