<template>
  <div class="register-page">
    <div class="row justify-content-center">
      <div class="col-md-6 col-lg-5">
        <div class="card">
          <div class="card-header">
            <h4 class="mb-0">Регистрация</h4>
          </div>
          <div class="card-body">
            <form @submit.prevent="handleRegister">
              <div class="mb-3">
                <label for="name" class="form-label">Имя</label>
                <input
                  type="text"
                  class="form-control"
                  id="name"
                  v-model="form.name"
                  :class="{ 'is-invalid': errors.name }"
                >
                <div v-if="errors.name" class="invalid-feedback">
                  {{ errors.name }}
                </div>
                <div class="form-text">Если оставить пустым, будет использован email</div>
              </div>

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
                <div class="form-text">Пароль должен быть от 5 до 15 символов</div>
              </div>

              <div class="mb-3">
                <label for="confirmPassword" class="form-label">Подтверждение пароля</label>
                <input
                  type="password"
                  class="form-control"
                  id="confirmPassword"
                  v-model="form.confirmPassword"
                  required
                  :class="{ 'is-invalid': errors.confirmPassword }"
                >
                <div v-if="errors.confirmPassword" class="invalid-feedback">
                  {{ errors.confirmPassword }}
                </div>
              </div>

              <div v-if="successMessage" class="alert alert-success">
                {{ successMessage }}
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
                  Зарегистрироваться
                </button>
              </div>
            </form>

            <div class="mt-3 text-center">
              <p class="mb-0">
                Уже есть аккаунт?
                <router-link to="/auth/login">Войти</router-link>
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
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
});

const errors = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
});

const authError = ref<string>('');
const successMessage = ref<string>('');

function validateForm() {
  let isValid = true;

  // Очищаем ошибки
  errors.name = '';
  errors.email = '';
  errors.password = '';
  errors.confirmPassword = '';

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

  // Подтверждение пароля
  if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Пароли не совпадают';
    isValid = false;
  }

  return isValid;
}

async function handleRegister() {
  if (!validateForm()) {
    return;
  }

  authError.value = '';
  successMessage.value = '';

  const result = await authStore.register({
    name: form.name || form.email,
    email: form.email,
    password: form.password
  });

  if (result.success) {
    successMessage.value = 'Регистрация успешна! Проверьте вашу почту для подтверждения аккаунта.';

    // Очищаем форму
    form.name = '';
    form.email = '';
    form.password = '';
    form.confirmPassword = '';

    // Автоматический редирект на страницу входа через 3 секунды
    setTimeout(() => {
      router.push('/auth/login');
    }, 3000);
  } else {
    authError.value = result.error || 'Ошибка регистрации';
  }
}
</script>

<style scoped>
.register-page {
  max-width: 400px;
  margin: 0 auto;
}

.card {
  border: none;
  box-shadow: 0 2px 20px rgba(0,0,0,.1);
}
</style>
