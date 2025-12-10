import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User, LoginRequest, RegisterRequest } from '@/types';
import { authService } from '@/services';

export const useAuthStore = defineStore('auth', () => {
  // Состояние
  const user = ref<User | null>(authService.getCurrentUser());
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Геттеры
  const isAuthenticated = computed(() => !!user.value);
  const isAdmin = computed(() => user.value?.is_admin || false);
  const currentUser = computed(() => user.value);

  // Действия
  async function register(data: RegisterRequest) {
    isLoading.value = true;
    error.value = null;

    try {
      await authService.register(data);
      return { success: true };
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Ошибка регистрации';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  async function login(data: LoginRequest) {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await authService.login(data);
      user.value = response.result.user;
      return { success: true };
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Ошибка входа';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  async function logout() {
    authService.logout();
    user.value = null;
  }

  async function loadProfile() {
    try {
      const response = await authService.profile();
      user.value = response.result;
      localStorage.setItem('user', JSON.stringify(user.value));
    } catch (err) {
      console.error('Failed to load profile:', err);
    }
  }

  async function activateAccount(token: string, code: string) {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await authService.activate(token, code);
      return { success: true, user: response.result };
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Ошибка активации';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  // Инициализация при загрузке
  function init() {
    const storedUser = authService.getCurrentUser();
    if (storedUser) {
      user.value = storedUser;
    }
  }

  return {
    // Состояние
    user,
    isLoading,
    error,

    // Геттеры
    isAuthenticated,
    isAdmin,
    currentUser,

    // Действия
    register,
    login,
    logout,
    loadProfile,
    activateAccount,
    init
  };
});
