import { defineStore } from "pinia";
import api from '../services/api';
import { useRouter } from 'vue-router';

export const useAuthStore = defineStore("auth", {
  state: () => {
    // Initialize user from localStorage if available
    const storedUser = localStorage.getItem("user");
    const user = storedUser ? JSON.parse(storedUser) : null;
    
    return {
      user: user,
      token: localStorage.getItem("token") || null,
      role: localStorage.getItem("role") || null,
    };
  },

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.role === "admin",
    isDeveloper: (state) => state.role === "developer",
    isManager: (state) => state.role === "manager",
    currentUser: (state) => state.user,
    canReportBugs: (state) =>
      ["developer", "manager", "admin"].includes(state.role),
    canAssignBugs: (state) => ["manager", "admin"].includes(state.role),
    canEditBugStatus: (state) =>
      ["developer", "manager", "admin"].includes(state.role),
    canViewAllBugs: (state) => ["manager", "admin"].includes(state.role),
  },

  actions: {
    async login(credentials) {
      try {
        console.log('Attempting login with credentials:', credentials);
        console.log('API URL:', api.defaults.baseURL);
        
        const response = await api.post("/auth/login", credentials);
        console.log('Login response:', response.data);
        
        const data = response.data;
        
        this.token = data.token;
        this.user = data.user;
        this.role = data.user.role;

        // Store user data in localStorage
        localStorage.setItem("token", data.token);
        localStorage.setItem("role", data.user.role);
        localStorage.setItem("user", JSON.stringify(data.user));

        return true;
      } catch (error) {
        console.error("Login error details:", {
          message: error.message,
          response: error.response?.data,
          status: error.response?.status,
          config: {
            url: error.config?.url,
            method: error.config?.method,
            headers: error.config?.headers,
            data: error.config?.data
          }
        });
        
        if (error.response?.status === 401) {
          window.alert('Invalid email or password');
        } else if (error.response?.status === 0) {
          window.alert('Unable to connect to the server. Please check if the server is running.');
        } else {
          window.alert(error.response?.data?.message || 'Login failed. Please try again.');
        }
        return false;
      }
    },

    async register(userData) {
      try {
        console.log('Attempting registration with data:', userData);
        console.log('API URL:', api.defaults.baseURL);
        
        const response = await api.post("/auth/register", userData);
        console.log('Registration response:', response.data);
        
        const data = response.data;
        return { success: true, data };
      } catch (error) {
        console.error("Registration error details:", {
          message: error.message,
          response: error.response?.data,
          status: error.response?.status,
          config: {
            url: error.config?.url,
            method: error.config?.method,
            headers: error.config?.headers,
            data: error.config?.data
          }
        });
        
        if (error.response?.status === 409) {
          return { success: false, error: 'Email already exists' };
        } else if (error.response?.status === 0) {
          return { success: false, error: 'Unable to connect to the server. Please check if the server is running.' };
        } else if (error.response?.status === 400) {
          return { success: false, error: error.response.data.message || 'Invalid registration data' };
        } else {
          return { success: false, error: error.response?.data?.message || 'Registration failed. Please try again.' };
        }
      }
    },

    async logout() {
      this.user = null;
      this.token = null;
      this.role = null;
      localStorage.removeItem("token");
      localStorage.removeItem("role");
      localStorage.removeItem("user");
      
      // Use router for navigation
      const router = useRouter();
      await router.push('/login');
    }
  },
});
