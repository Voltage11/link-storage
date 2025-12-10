import apiClient from './api';
import type {
  LoginRequest,
  RegisterRequest,
  ActivateRequest,
  LoginResponse,
  SuccessResponse,
  User,
  SessionResponse
} from '@/types';

export const authService = {
  // Регистрация
  async register(data: RegisterRequest): Promise<SuccessResponse<void>> {
    const response = await apiClient.post('/auth/register', data);
    return response.data;
  },

  // Активация аккаунта
  async activate(token: string, code: string): Promise<SuccessResponse<User>> {
    const response = await apiClient.post(`/auth/activate/${token}`, { code });
    return response.data;
  },

  // Вход
  async login(data: LoginRequest): Promise<SuccessResponse<SessionResponse>> {
    const response = await apiClient.post('/auth/login', data);
    const result = response.data.result;

    // Сохраняем токены и пользователя
    if (result.access_token && result.refresh_token) {
      localStorage.setItem('access_token', result.access_token);
      localStorage.setItem('refresh_token', result.refresh_token);
      localStorage.setItem('user', JSON.stringify(result.user));
    }

    return response.data;
  },

  // Профиль пользователя
  async profile(): Promise<SuccessResponse<User>> {
    const response = await apiClient.get('/auth/profile');
    return response.data;
  },

  // Обновление токена
  async refreshToken(refreshToken: string): Promise<SuccessResponse<SessionResponse>> {
    const response = await apiClient.post('/auth/refresh-token', {
      refresh_token: refreshToken
    });
    const result = response.data.result;

    if (result.access_token && result.refresh_token) {
      localStorage.setItem('access_token', result.access_token);
      localStorage.setItem('refresh_token', result.refresh_token);
      localStorage.setItem('user', JSON.stringify(result.user));
    }

    return response.data;
  },

  // Выход
  logout(): void {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
  },

  // Проверка авторизации
  isAuthenticated(): boolean {
    return !!localStorage.getItem('access_token');
  },

  // Получение текущего пользователя
  getCurrentUser(): User | null {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  }
};
