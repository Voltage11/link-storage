<template>
  <router-view />
</template>

<script setup lang="ts">
// Инициализация приложения
import { onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();

onMounted(() => {
  authStore.init();
});
</script>

<style>
/* Глобальные стили */
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f8f9fa;
  min-height: 100vh;
}

/* Кастомные стили для Bootstrap */
.btn {
  border-radius: 0.375rem;
}

.card {
  border-radius: 0.5rem;
  overflow: hidden;
}

.form-control:focus {
  box-shadow: 0 0 0 0.25rem rgba(13, 110, 253, 0.25);
}

/* Утилиты */
.cursor-pointer {
  cursor: pointer;
}

/* Анимации */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
