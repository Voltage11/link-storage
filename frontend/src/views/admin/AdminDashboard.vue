<template>
  <div class="admin-dashboard">
    <div class="row mb-4">
      <div class="col-12">
        <h1 class="h2">Админ-панель</h1>
        <p class="lead text-muted">Панель управления для администраторов</p>
      </div>
    </div>

    <div class="row">
      <div class="col-md-6 mb-4">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Статистика</h5>
          </div>
          <div class="card-body">
            <div class="list-group list-group-flush">
              <div class="list-group-item d-flex justify-content-between align-items-center">
                Всего пользователей
                <span class="badge bg-primary rounded-pill">{{ stats.totalUsers || 0 }}</span>
              </div>
              <div class="list-group-item d-flex justify-content-between align-items-center">
                Активных пользователей
                <span class="badge bg-success rounded-pill">{{ stats.activeUsers || 0 }}</span>
              </div>
              <div class="list-group-item d-flex justify-content-between align-items-center">
                Всего групп ссылок
                <span class="badge bg-info rounded-pill">{{ stats.totalGroups || 0 }}</span>
              </div>
              <div class="list-group-item d-flex justify-content-between align-items-center">
                Администраторов
                <span class="badge bg-warning rounded-pill">{{ stats.adminUsers || 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-6 mb-4">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Быстрые действия</h5>
          </div>
          <div class="card-body">
            <div class="d-grid gap-2">
              <button class="btn btn-outline-primary" @click="refreshStats">
                <i class="bi bi-arrow-clockwise me-2"></i>
                Обновить статистику
              </button>
              <router-link to="/dashboard" class="btn btn-outline-success">
                <i class="bi bi-box-arrow-right me-2"></i>
                Перейти в приложение
              </router-link>
              <button class="btn btn-outline-secondary" @click="viewLogs">
                <i class="bi bi-journal-text me-2"></i>
                Просмотреть логи
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Сообщения -->
    <div v-if="message" class="alert alert-info alert-dismissible fade show mt-4">
      {{ message }}
      <button type="button" class="btn-close" @click="message = ''"></button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

const stats = ref({
  totalUsers: 0,
  activeUsers: 0,
  totalGroups: 0,
  adminUsers: 0
});

const message = ref<string>('');

// Загрузка статистики
async function loadStats() {
  try {
    // Здесь можно добавить API вызовы для получения статистики
    // Для примера используем фиктивные данные
    stats.value = {
      totalUsers: 150,
      activeUsers: 120,
      totalGroups: 450,
      adminUsers: 3
    };
  } catch (error) {
    console.error('Ошибка загрузки статистики:', error);
  }
}

function refreshStats() {
  loadStats();
  message.value = 'Статистика обновлена';
}

function viewLogs() {
  message.value = 'Функция просмотра логов в разработке';
}

onMounted(() => {
  loadStats();
});
</script>

<style scoped>
.admin-dashboard {
  padding: 20px 0;
}

.card {
  border: none;
  box-shadow: 0 2px 10px rgba(0,0,0,.1);
}

.list-group-item {
  border-color: rgba(0,0,0,.125);
}
</style>
