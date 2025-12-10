<template>
  <div class="home-page">
    <div class="hero bg-primary text-white py-5 mb-5">
      <div class="container">
        <h1 class="display-4">Добро пожаловать в Link Storage</h1>
        <p class="lead">Управляйте своими закладками и ссылками в одном месте</p>
        <div class="mt-4">
          <router-link v-if="!isAuthenticated" to="/auth/register" class="btn btn-light btn-lg me-3">
            Начать бесплатно
          </router-link>
          <router-link v-if="!isAuthenticated" to="/auth/login" class="btn btn-outline-light btn-lg">
            Войти
          </router-link>
          <router-link v-else to="/dashboard" class="btn btn-light btn-lg">
            Перейти в приложение
          </router-link>
        </div>
      </div>
    </div>

    <div class="container">
      <div class="row">
        <div class="col-md-4 mb-4">
          <div class="card h-100">
            <div class="card-body text-center">
              <i class="bi bi-collection fs-1 text-primary mb-3"></i>
              <h5 class="card-title">Группируйте ссылки</h5>
              <p class="card-text">Создавайте группы для удобной организации ваших закладок</p>
            </div>
          </div>
        </div>
        <div class="col-md-4 mb-4">
          <div class="card h-100">
            <div class="card-body text-center">
              <i class="bi bi-tags fs-1 text-success mb-3"></i>
              <h5 class="card-title">Цветовые метки</h5>
              <p class="card-text">Используйте цвета для быстрой идентификации групп</p>
            </div>
          </div>
        </div>
        <div class="col-md-4 mb-4">
          <div class="card h-100">
            <div class="card-body text-center">
              <i class="bi bi-shield-lock fs-1 text-danger mb-3"></i>
              <h5 class="card-title">Безопасность</h5>
              <p class="card-text">Ваши данные защищены и доступны только вам</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { storeToRefs } from 'pinia';

const authStore = useAuthStore();
const { isAuthenticated } = storeToRefs(authStore);
</script>

<style scoped>
.hero {
  border-radius: 0 0 20px 20px;
}

.card {
  border: none;
  box-shadow: 0 2px 10px rgba(0,0,0,.1);
  transition: transform 0.3s;
}

.card:hover {
  transform: translateY(-5px);
}
</style>
