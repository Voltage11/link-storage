<template>
  <div class="admin-layout">
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
      <div class="container">
        <router-link class="navbar-brand" to="/admin">Админ-панель</router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav">
            <li class="nav-item">
              <router-link class="nav-link" to="/admin">Статистика</router-link>
            </li>
            <li class="nav-item">
              <router-link class="nav-link" to="/dashboard">Вернуться в приложение</router-link>
            </li>
          </ul>
          <ul class="navbar-nav ms-auto">
            <li class="nav-item">
              <a class="nav-link" href="#" @click="logout">Выход</a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <main class="container-fluid mt-3">
      <div class="row">
        <div class="col-12">
          <router-view />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const router = useRouter();
const authStore = useAuthStore();

const logout = async () => {
  await authStore.logout();
  router.push('/');
};
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background-color: #f0f0f0;
}
</style>
