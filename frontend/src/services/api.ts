import axios from 'axios';
import type { ErrorResponse } from '../types';

const apiClient = axios.create({
  baseURL: '/api/v1', // Теперь соответствует настройке прокси
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
});

// Интерцептор для добавления токена
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Логируем запросы для отладки
    console.log(`[API Request] ${config.method?.toUpperCase()} ${config.baseURL}${config.url}`);
    return config;
  },
  (error) => {
    console.error('Request error:', error);
    return Promise.reject(error);
  }
);

// Интерцептор для обработки ответов
apiClient.interceptors.response.use(
  (response) => {
    console.log(`[API Response] ${response.status} ${response.config.url}`);
    return response;
  },
  async (error) => {
    console.error(`[API Error] ${error.code} ${error.config?.url}:`, error.message);

    // Если сервер вернул HTML вместо JSON
    const contentType = error.response?.headers?.['content-type'];
    if (contentType && contentType.includes('text/html')) {
      console.error('Server returned HTML instead of JSON. Check proxy configuration.');

      error.response.data = {
        success: false,
        error: {
          type: 'INVALID_RESPONSE',
          message: 'Сервер вернул некорректный ответ (HTML вместо JSON)'
        }
      };
    }

    return Promise.reject(error);
  }
);

export default apiClient;
