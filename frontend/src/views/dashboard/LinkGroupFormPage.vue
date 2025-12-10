<template>
  <div class="link-group-form-page">
    <div class="row mb-4">
      <div class="col-12">
        <h1 class="h2">
          {{ isEditMode ? 'Редактирование группы' : 'Создание группы' }}
        </h1>
        <p class="text-muted">
          {{ isEditMode ? 'Измените данные группы ссылок' : 'Создайте новую группу для организации ваших ссылок' }}
        </p>
      </div>
    </div>

    <div class="row justify-content-center">
      <div class="col-md-8 col-lg-6">
        <div class="card">
          <div class="card-body">
            <form @submit.prevent="handleSubmit">
              <!-- Название группы -->
              <div class="mb-3">
                <label for="name" class="form-label">
                  Название группы <span class="text-danger">*</span>
                </label>
                <input
                  type="text"
                  class="form-control"
                  id="name"
                  v-model="form.name"
                  required
                  :class="{ 'is-invalid': errors.name }"
                  placeholder="Введите название группы"
                >
                <div v-if="errors.name" class="invalid-feedback">
                  {{ errors.name }}
                </div>
                <div class="form-text">
                  Название должно быть от 3 до 50 символов
                </div>
              </div>

              <!-- Описание -->
              <div class="mb-3">
                <label for="description" class="form-label">Описание</label>
                <textarea
                  class="form-control"
                  id="description"
                  v-model="form.description"
                  rows="3"
                  :class="{ 'is-invalid': errors.description }"
                  placeholder="Добавьте описание группы (необязательно)"
                ></textarea>
                <div v-if="errors.description" class="invalid-feedback">
                  {{ errors.description }}
                </div>
              </div>

              <!-- Цвет -->
              <div class="mb-4">
                <label class="form-label">Цвет группы</label>
                <div class="color-picker">
                  <div class="row g-2">
                    <div
                      v-for="color in colorOptions"
                      :key="color.value"
                      class="col-3 col-sm-2"
                    >
                      <div
                        class="color-option"
                        :style="{ backgroundColor: color.value }"
                        :class="{ active: form.color === color.value }"
                        @click="form.color = color.value"
                      >
                        <span v-if="form.color === color.value" class="checkmark">
                          ✓
                        </span>
                      </div>
                      <small class="d-block text-center mt-1">{{ color.label }}</small>
                    </div>
                  </div>
                </div>
                <div class="mt-2">
                  <input
                    type="color"
                    class="form-control form-control-color"
                    v-model="form.color"
                    title="Выберите цвет"
                  >
                  <div class="form-text">Или выберите любой цвет</div>
                </div>
              </div>

              <!-- Ошибки -->
              <div v-if="submitError" class="alert alert-danger">
                {{ submitError }}
              </div>

              <!-- Успех -->
              <div v-if="successMessage" class="alert alert-success">
                {{ successMessage }}
              </div>

              <!-- Кнопки -->
              <div class="d-flex gap-2">
                <button
                  type="submit"
                  class="btn btn-primary flex-grow-1"
                  :disabled="linkGroupStore.isLoading"
                >
                  <span v-if="linkGroupStore.isLoading" class="spinner-border spinner-border-sm me-2"></span>
                  {{ isEditMode ? 'Сохранить изменения' : 'Создать группу' }}
                </button>

                <router-link
                  to="/dashboard/link-groups"
                  class="btn btn-outline-secondary"
                >
                  Отмена
                </router-link>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useLinkGroupStore } from '@/stores/link-group';
import type { LinkGroupCreate, LinkGroupUpdate } from '@/types';

const route = useRoute();
const router = useRouter();
const linkGroupStore = useLinkGroupStore();

const isEditMode = computed(() => !!route.params.id);

// Форма
const form = reactive({
  id: 0,
  name: '',
  description: '',
  color: '#007bff' // Bootstrap primary color по умолчанию
});

const errors = reactive({
  name: '',
  description: ''
});

const submitError = ref<string>('');
const successMessage = ref<string>('');

// Цвета по умолчанию
const colorOptions = [
  { value: '#007bff', label: 'Синий' },
  { value: '#28a745', label: 'Зеленый' },
  { value: '#dc3545', label: 'Красный' },
  { value: '#ffc107', label: 'Желтый' },
  { value: '#6c757d', label: 'Серый' },
  { value: '#17a2b8', label: 'Голубой' },
  { value: '#6610f2', label: 'Фиолет.' },
  { value: '#e83e8c', label: 'Розовый' }
];

// Валидация
function validateForm(): boolean {
  let isValid = true;

  errors.name = '';
  errors.description = '';

  if (!form.name.trim()) {
    errors.name = 'Название обязательно';
    isValid = false;
  } else if (form.name.length < 3 || form.name.length > 50) {
    errors.name = 'Название должно быть от 3 до 50 символов';
    isValid = false;
  }

  return isValid;
}

// Загрузка данных для редактирования
async function loadGroupData() {
  const id = parseInt(route.params.id as string);

  if (!id) return;

  const result = await linkGroupStore.fetchGroupById(id);

  if (result.success && linkGroupStore.currentGroup) {
    form.id = linkGroupStore.currentGroup.id;
    form.name = linkGroupStore.currentGroup.name;
    form.description = linkGroupStore.currentGroup.description || '';
    form.color = linkGroupStore.currentGroup.color || '#007bff';
  }
}

// Отправка формы
async function handleSubmit() {
  if (!validateForm()) {
    return;
  }

  submitError.value = '';
  successMessage.value = '';

  let result;

  if (isEditMode.value) {
    const updateData: LinkGroupUpdate = {
      id: form.id,
      name: form.name,
      description: form.description,
      color: form.color
    };

    result = await linkGroupStore.updateGroup(updateData);
  } else {
    const createData: LinkGroupCreate = {
      name: form.name,
      description: form.description,
      color: form.color
    };

    result = await linkGroupStore.createGroup(createData);
  }

  if (result.success) {
    successMessage.value = isEditMode.value
      ? 'Группа успешно обновлена!'
      : 'Группа успешно создана!';

    // Редирект через 2 секунды
    setTimeout(() => {
      router.push('/dashboard/link-groups');
    }, 2000);
  } else {
    submitError.value = result.error || 'Ошибка сохранения';
  }
}

onMounted(() => {
  if (isEditMode.value) {
    loadGroupData();
  }
});
</script>

<style scoped>
.link-group-form-page {
  padding: 20px 0;
}

.card {
  border: none;
  box-shadow: 0 2px 20px rgba(0,0,0,.1);
}

.color-picker {
  padding: 10px;
  border: 1px solid #dee2e6;
  border-radius: 0.375rem;
  background-color: #f8f9fa;
}

.color-option {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.color-option:hover {
  transform: scale(1.1);
  border-color: #6c757d;
}

.color-option.active {
  border-color: #212529;
  transform: scale(1.1);
}

.checkmark {
  color: white;
  font-weight: bold;
  text-shadow: 0 0 2px rgba(0,0,0,0.5);
}

.form-control-color {
  height: 40px;
  padding: 0;
  border: none;
  background: transparent;
}
</style>
