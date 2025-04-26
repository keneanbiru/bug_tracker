import axios from 'axios';

// Get the base URL from environment or use default
const getBaseURL = () => {
    if (process.env.NODE_ENV === 'test') {
        return 'http://localhost:8080/api';
    }
    try {
        // @ts-ignore
        return window.__VITE_API_URL__ || 'http://localhost:8080/api';
    } catch {
        return 'http://localhost:8080/api';
    }
};

const api = axios.create({
    baseURL: getBaseURL(),
});

// Add request interceptor to add auth token
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export default api; 