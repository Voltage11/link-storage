<template>
  <div class="link-groups-page">
    <div class="row mb-4">
      <div class="col-md-8">
        <h1 class="h2">Мои группы ссылок</h1>
        <p class="text-muted">Управляйте вашими группами закладок</p>
      </div>
      <div class="col-md-4 text-md-end">
        <router-link to="/dashboard/link-groups/create" class="btn btn-primary">
          <i class="bi bi-plus-circle me-2"></i>
          Новая группа
        </router-link>
      </div>
    </div>

    <!-- Поиск и фильтры -->
    <div class="row mb-4">
      <div class="col-md-6">
        <div class="input-group">
          <input
            type="text"
            class="form-control"
            placeholder="Поиск по названию..."
            v-model="searchQuery"
            @input="handleSearch"
          >
          <button class="btn btn-outline-secondary" type="button" @click="handleSearch">
            <i class="bi bi-search"></i>
          </button>
        </div>
      </div>
      <div class="col-md-6 text-md-end">
        <div class="btn-group" role="group">
          <button
            type="button"
            class="btn btn-outline-secondary"
            :class="{ active: sortBy === 'name' }"
            @click="setSort('name')"
          >
            По названию
          </button>
          <button
            type="button"
            class="btn btn-outline-secondary"
            :class="{ active: sortBy === 'date' }"
            @click="setSort('date')"
          >
            По дате
          </button>
        </div>
      </div>
    </div>

    <!-- Сообщения об ошибках -->
    <div v-if="linkGroupStore.error" class="alert alert-danger alert-dismissible fade show">
      {{ linkGroupStore.error }}
      <button type="button" class="btn-close" @click="clearError"></button>
    </div>

    <!-- Загрузка -->
    <div v-if="linkGroupStore.isLoading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Загрузка...</span>
      </div>
      <p class="mt-2">Загрузка групп...</p>
    </div>

    <!-- Список групп -->
    <div v-else-if="linkGroupStore.groups && linkGroupStore.groups.length > 0" class="row">
      <div
        v-for="group in sortedGroups"
        :key="group.id"
        class="col-md-6 col-lg-4 mb-4"
      >
        <div class="card h-100" :style="{ borderLeft: `4px solid ${group.color || '#007bff'}` }">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-3">
              <h5 class="card-title mb-0">{{ group.name }}</h5>
              <span class="badge" :style="{ backgroundColor: group.color || '#007bff', color: 'white' }">
                {{ group.position || 0 }}
              </span>
            </div>

            <p class="card-text text-muted" v-if="group.description">
              {{ group.description }}
            </p>
            <p class="card-text text-muted" v-else>
              Нет описания
            </p>

            <div class="small text-muted mb-3">
              Создано: {{ formatDate(group.created_at) }}
            </div>
          </div>

          <div class="card-footer bg-transparent">
            <div class="btn-group w-100" role="group">
              <router-link
                :to="`/dashboard/link-groups/edit/${group.id}`"
                class="btn btn-outline-primary btn-sm"
              >
                <i class="bi bi-pencil me-1"></i> Редактировать
              </router-link>
              <button
                type="button"
                class="btn btn-outline-danger btn-sm"
                @click="confirmDelete(group)"
              >
                <i class="bi bi-trash me-1"></i> Удалить
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Нет групп -->
    <div v-else class="text-center my-5">
      <div class="card border-0">
        <div class="card-body py-5">
          <i class="bi bi-collection text-muted fs-1 mb-3"></i>
          <h5 class="card-title">Группы не найдены</h5>
          <p class="card-text text-muted">Создайте свою первую группу ссылок</p>
          <router-link to="/dashboard/link-groups/create" class="btn btn-primary">
            Создать группу
          </router-link>
        </div>
      </div>
    </div>

    <!-- Пагинация -->
    <div v-if="linkGroupStore.pagination.total_pages > 1" class="row mt-4">
      <div class="col-12">
        <nav aria-label="Пагинация">
          <ul class="pagination justify-content-center">
            <li class="page-item" :class="{ disabled: currentPage === 1 }">
              <button
                class="page-link"
                @click="changePage(currentPage - 1)"
              >
                Назад
              </button>
            </li>

            <li
              v-for="page in pages"
              :key="page"
              class="page-item"
              :class="{ active: page === currentPage }"
            >
              <button
                class="page-link"
                @click="changePage(page)"
              >
                {{ page }}
              </button>
            </li>

            <li class="page-item" :class="{ disabled: currentPage === linkGroupStore.pagination.total_pages }">
              <button
                class="page-link"
                @click="changePage(currentPage + 1)"
              >
                Вперед
              </button>
            </li>
          </ul>
        </nav>
        <p class="text-center text-muted small">
          Показано {{ linkGroupStore.groups.length }} из {{ linkGroupStore.pagination.total }} групп
        </p>
      </div>
    </div>

    <!-- Модальное окно подтверждения удаления -->
    <div v-if="showDeleteModal" class="modal fade show d-block" tabindex="-1" style="background: rgba(0,0,0,0.5)">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Подтверждение удаления</h5>
            <button type="button" class="btn-close" @click="cancelDelete"></button>
          </div>
          <div class="modal-body">
            <p>Вы уверены, что хотите удалить группу <strong>{{ groupToDelete?.name }}</strong>?</p>
            <p class="text-danger small">Это действие нельзя отменить.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="cancelDelete">
              Отмена
            </button>
            <button
              type="button"
              class="btn btn-danger"
              @click="performDelete"
              :disabled="linkGroupStore.isLoading"
            >
              <span v-if="linkGroupStore.isLoading" class="spinner-border spinner-border-sm me-2"></span>
              Удалить
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useLinkGroupStore } from '@/stores/link-group';
import type { LinkGroup } from '@/types';

const route = useRoute();
const router = useRouter();
const linkGroupStore = useLinkGroupStore();

// Состояние
const searchQuery = ref<string>('');
const sortBy = ref<'name' | 'date'>('name');
const currentPage = ref<number>(1);
const showDeleteModal = ref<boolean>(false);
const groupToDelete = ref<LinkGroup | null>(null);

// Вычисляемые свойства
const sortedGroups = computed(() => {
  if (!linkGroupStore.groups) return [];

  const groups = [...linkGroupStore.groups];

  if (sortBy.value === 'name') {
    return groups.sort((a, b) => a.name.localeCompare(b.name));
  } else {
    return groups.sort((a, b) =>
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    );
  }
});

const pages = computed(() => {
  const totalPages = linkGroupStore.pagination.total_pages;
  const pageCount = 5; // Максимальное количество отображаемых страниц
  const current = currentPage.value;

  if (totalPages <= pageCount) {
    return Array.from({ length: totalPages }, (_, i) => i + 1);
  }

  let start = Math.max(1, current - Math.floor(pageCount / 2));
  let end = Math.min(totalPages, start + pageCount - 1);

  if (end - start + 1 < pageCount) {
    start = Math.max(1, end - pageCount + 1);
  }

  return Array.from({ length: end - start + 1 }, (_, i) => start + i);
});

// Методы
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  });
}

async function loadGroups() {
  await linkGroupStore.fetchGroups({
    page: currentPage.value,
    page_size: linkGroupStore.pagination.page_size,
    name: searchQuery.value || undefined
  });
}

let searchTimeout: NodeJS.Timeout;

function handleSearch() {
  clearTimeout(searchTimeout);

  // Дебаунс поиска - ждем 300мс после последнего ввода
  searchTimeout = setTimeout(() => {
    currentPage.value = 1;
    loadGroups();

    // Обновляем URL
    router.push({
      query: {
        search: searchQuery.value || undefined
      }
    });
  }, 300);
}

function clearSearch() {
  searchQuery.value = '';
  handleSearch();
}

function setSort(type: 'name' | 'date') {
  sortBy.value = type;
}

function changePage(page: number) {
  if (page < 1 || page > linkGroupStore.pagination.total_pages) {
    return;
  }

  currentPage.value = page;
  loadGroups();

  // Обновляем URL с query параметром
  router.push({
    query: {
      page: page.toString(),
      search: searchQuery.value || undefined
    }
  });
}

function confirmDelete(group: LinkGroup) {
  groupToDelete.value = group;
  showDeleteModal.value = true;
}

function cancelDelete() {
  showDeleteModal.value = false;
  groupToDelete.value = null;
}

async function performDelete() {
  if (!groupToDelete.value) return;

  const result = await linkGroupStore.deleteGroup(groupToDelete.value.id);

  if (result.success) {
    showDeleteModal.value = false;
    groupToDelete.value = null;
    loadGroups();
  }
}

function clearError() {
  linkGroupStore.error = null;
}

// Обработка query параметров из URL
watch(
  () => route.query,
  (query) => {
    if (query.page) {
      const page = parseInt(query.page as string);
      if (!isNaN(page) && page !== currentPage.value) {
        currentPage.value = page;
      }
    }

    if (query.search && query.search !== searchQuery.value) {
      searchQuery.value = query.search as string;
    }

    loadGroups();
  },
  { immediate: true }
);

onMounted(() => {
  // Если нет query параметров, загружаем первую страницу
  if (!route.query.page) {
    loadGroups();
  }
});
</script>

<style scoped>
.link-groups-page {
  padding: 20px 0;
}

.card {
  border: none;
  box-shadow: 0 2px 10px rgba(0,0,0,.1);
  transition: transform 0.3s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0,0,0,.15);
}

.page-link {
  cursor: pointer;
}

.modal {
  background: rgba(0,0,0,0.5);
}
</style>
