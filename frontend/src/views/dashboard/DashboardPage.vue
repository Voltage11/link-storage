<template>
  <div class="dashboard-page">
    <div class="row mb-4">
      <div class="col-12">
        <h1 class="h2">Добро пожаловать, {{ user?.name || user?.email }}!</h1>
        <p class="lead text-muted">Управляйте вашими группами ссылок</p>
      </div>
    </div>

    <div class="row">
      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card bg-primary text-white">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title mb-0">Группы ссылок</h6>
                <p class="card-text display-6 mb-0">{{ stats.groups || 0 }}</p>
              </div>
              <i class="bi bi-collection fs-1 opacity-50"></i>
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card bg-success text-white">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title mb-0">Всего ссылок</h6>
                <p class="card-text display-6 mb-0">{{ stats.links || 0 }}</p>
              </div>
              <i class="bi bi-link fs-1 opacity-50"></i>
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card bg-info text-white">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title mb-0">Дата регистрации</h6>
                <p class="card-text mb-0">
                  {{ formatDate(user?.created_at) }}
                </p>
              </div>
              <i class="bi bi-calendar-check fs-1 opacity-50"></i>
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-6 col-lg-3 mb-4">
        <div class="card bg-warning text-dark">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title mb-0">Последний вход</h6>
                <p class="card-text mb-0">
                  {{ formatDate(user?.last_login_at) }}
                </p>
              </div>
              <i class="bi bi-clock-history fs-1 opacity-50"></i>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row mt-5">
      <div class="col-12">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">Быстрые действия</h5>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-md-4 mb-3">
                <router-link to="/dashboard/link-groups/create" class="btn btn-primary w-100 py-3">
                  <i class="bi bi-plus-circle me-2"></i>
                  Создать группу
                </router-link>
              </div>
              <div class="col-md-4 mb-3">
                <router-link to="/dashboard/link-groups" class="btn btn-success w-100 py-3">
                  <i class="bi bi-list-ul me-2"></i>
                  Все группы
                </router-link>
              </div>
              <div class="col-md-4 mb-3">
                <button class="btn btn-outline-secondary w-100 py-3" @click="refreshData">
                  <i class="bi bi-arrow-clockwise me-2"></i>
                  Обновить данные
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useLinkGroupStore } from '@/stores/link-group';
import { storeToRefs } from 'pinia';

const authStore = useAuthStore();
const linkGroupStore = useLinkGroupStore();

const { user } = storeToRefs(authStore);
const { groups } = storeToRefs(linkGroupStore);

const stats = ref({
  groups: 0,
  links: 0
});

function formatDate(dateString?: string): string {
  if (!dateString) return 'Не указано';

  const date = new Date(dateString);
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  });
}

async function loadData() {
  await linkGroupStore.fetchGroups();
  stats.value.groups = groups.value.length;
  // Здесь можно добавить загрузку статистики по ссылкам
}

function refreshData() {
  loadData();
}

onMounted(() => {
  loadData();
});
</script>

<style scoped>
.dashboard-page {
  padding: 20px 0;
}

.card {
  border: none;
  box-shadow: 0 2px 10px rgba(0,0,0,.1);
  transition: transform 0.3s;
}

.card:hover {
  transform: translateY(-5px);
}

.btn {
  transition: all 0.3s;
}

.btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,.2);
}
</style>
