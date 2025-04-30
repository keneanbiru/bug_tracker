import axios from 'axios';
import { useRouter } from 'vue-router';

// Get the base URL from environment or use default
const getBaseURL = () => {
    if (process.env.NODE_ENV === 'test') {
        return 'http://localhost:8080/api';
    }
    try {
        // @ts-ignore
        return window.__VITE_API_URL__ || 'https://bug-tracker-3.onrender.com/api';
    } catch {
        return 'https://bug-tracker-3.onrender.com/api';
    }
};

const api = axios.create({
    baseURL: getBaseURL(),
    headers: {
        'Content-Type': 'application/json'
    }
});

// Add request interceptor to add auth token
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, error => {
    console.error('Request error:', error);
    return Promise.reject(error);
});

// Add response interceptor for better error handling
api.interceptors.response.use(
    response => response,
    error => {
        console.error('API Error:', {
            url: error.config?.url,
            method: error.config?.method,
            status: error.response?.status,
            data: error.response?.data,
            message: error.message,
            requestData: error.config?.data
        });

        // Handle authentication errors
        if (error.response?.status === 401) {
            // Clear auth data
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            localStorage.removeItem('user');
            
            // Redirect to login
            const router = useRouter();
            router.push('/login');
        }

        return Promise.reject(error);
    }
);

export default api; 