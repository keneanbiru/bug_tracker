import axios from 'axios';

// Get the base URL from environment or use default
const getBaseURL = () => {
    if (import.meta.env.MODE === 'test') {
        return 'http://localhost:8080/api';
    }
    return 'https://bug-tracker-8.onrender.com/api';
};

const api = axios.create({
    baseURL: getBaseURL(),
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    },
    withCredentials: false // Disable credentials for cross-origin requests since we're using token-based auth
});

// Add request interceptor to add auth token
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    // Log request for debugging
    console.log('Making request to:', `${config.baseURL}${config.url}`, {
        method: config.method,
        headers: config.headers,
        data: config.data
    });
    return config;
}, error => {
    console.error('Request error:', error);
    return Promise.reject(error);
});

// Add response interceptor for better error handling
api.interceptors.response.use(
    response => {
        // Log successful response for debugging
        console.log('Response received:', {
            url: response.config.url,
            status: response.status,
            data: response.data,
            headers: response.headers
        });
        return response;
    },
    error => {
        // Log detailed error information
        console.error('API Error:', {
            url: error.config?.url,
            method: error.config?.method,
            status: error.response?.status,
            data: error.response?.data,
            message: error.message,
            requestData: error.config?.data,
            headers: error.config?.headers,
            responseHeaders: error.response?.headers
        });

        // Handle authentication errors
        if (error.response?.status === 401) {
            // Clear auth data
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            localStorage.removeItem('user');
            
            // Redirect to login
            window.location.href = '/login';
        }

        return Promise.reject(error);
    }
);

export default api; 