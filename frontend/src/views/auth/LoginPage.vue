<template>
  <div class="login-page">
    <div class="row justify-content-center">
      <div class="col-md-6 col-lg-5">
        <div class="card">
          <div class="card-header">
            <h4 class="mb-0">Вход в систему</h4>
          </div>
          <div class="card-body">
            <form @submit.prevent="handleLogin">
              <div class="mb-3">
                <label for="email" class="form-label">Email</label>
                <input
                  type="email"
                  class="form-control"
                  id="email"
                  v-model="form.email"
                  required
                  :class="{ 'is-invalid': errors.email }"
                >
                <div v-if="errors.email" class="invalid-feedback">
                  {{ errors.email }}
                </div>
              </div>

              <div class="mb-3">
                <label for="password" class="form-label">Пароль</label>
                <input
                  type="password"
                  class="form-control"
                  id="password"
                  v-model="form.password"
                  required
                  :class="{ 'is-invalid': errors.password }"
                >
                <div v-if="errors.password" class="invalid-feedback">
                  {{ errors.password }}
                </div>
              </div>

              <div v-if="authError" class="alert alert-danger">
                {{ authError }}
              </div>

              <div class="d-grid gap-2">
                <button
                  type="submit"
                  class="btn btn-primary"
                  :disabled="authStore.isLoading"
                >
                  <span v-if="authStore.isLoading" class="spinner-border spinner-border-sm me-2"></span>
                  Войти
                </button>
              </div>
            </form>

            <div class="mt-3 text-center">
              <p class="mb-0">
                Нет аккаунта?
                <router-link to="/auth/register">Зарегистрироваться</router-link>
              </p>
              <p class="mb-0">
                <router-link to="/">Вернуться на главную</router-link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const form = reactive({
  email: '',
  password: ''
});

const errors = reactive({
  email: '',
  password: ''
});

const authError = ref<string>('');

function validateForm() {
  let isValid = true;

  // Очищаем ошибки
  errors.email = '';
  errors.password = '';

  // Валидация email
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!form.email) {
    errors.email = 'Email обязателен';
    isValid = false;
  } else if (!emailRegex.test(form.email)) {
    errors.email = 'Неверный формат email';
    isValid = false;
  }

  // Валидация пароля
  if (!form.password) {
    errors.password = 'Пароль обязателен';
    isValid = false;
  } else if (form.password.length < 5 || form.password.length > 15) {
    errors.password = 'Пароль должен быть от 5 до 15 символов';
    isValid = false;
  }

  return isValid;
}

async function handleLogin() {
  if (!validateForm()) {
    return;
  }

  authError.value = '';

  const result = await authStore.login({
    email: form.email,
    password: form.password
  });

  if (result.success) {
    router.push('/dashboard');
  } else {
    authError.value = result.error || 'Ошибка входа';
  }
}
</script>

<style scoped>
.login-page {
  max-width: 400px;
  margin: 0 auto;
}

.card {
  border: none;
  box-shadow: 0 2px 20px rgba(0,0,0,.1);
}
</style>
