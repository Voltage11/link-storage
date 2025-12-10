<template>
  <div class="activate-page">
    <div class="row justify-content-center">
      <div class="col-md-6 col-lg-5">
        <div class="card">
          <div class="card-header">
            <h4 class="mb-0">Активация аккаунта</h4>
          </div>
          <div class="card-body">
            <div v-if="!isActivated">
              <div class="alert alert-info">
                <p class="mb-2">На вашу почту отправлен 5-значный код подтверждения.</p>
                <p class="mb-0">Введите его ниже для активации аккаунта.</p>
              </div>

              <form @submit.prevent="handleActivate">
                <div class="mb-3">
                  <label for="code" class="form-label">Код подтверждения</label>
                  <input
                    type="text"
                    class="form-control"
                    id="code"
                    v-model="code"
                    required
                    maxlength="5"
                    :class="{ 'is-invalid': error }"
                    placeholder="Введите 5 цифр"
                  >
                  <div v-if="error" class="invalid-feedback">
                    {{ error }}
                  </div>
                </div>

                <div class="d-grid gap-2">
                  <button
                    type="submit"
                    class="btn btn-primary"
                    :disabled="authStore.isLoading"
                  >
                    <span v-if="authStore.isLoading" class="spinner-border spinner-border-sm me-2"></span>
                    Активировать
                  </button>
                </div>
              </form>
            </div>

            <div v-else>
              <div class="alert alert-success">
                <h5 class="alert-heading">Аккаунт успешно активирован!</h5>
                <p class="mb-0">Теперь вы можете войти в систему.</p>
              </div>

              <div class="text-center mt-4">
                <router-link to="/auth/login" class="btn btn-success">
                  Перейти к входу
                </router-link>
              </div>
            </div>

            <div class="mt-3 text-center">
              <router-link to="/">Вернуться на главную</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const token = ref<string>('');
const code = ref<string>('');
const error = ref<string>('');
const isActivated = ref<boolean>(false);

onMounted(() => {
  token.value = route.params.token as string;
});

async function handleActivate() {
  if (!code.value || code.value.length !== 5) {
    error.value = 'Введите 5-значный код';
    return;
  }

  error.value = '';

  const result = await authStore.activateAccount(token.value, code.value);

  if (result.success) {
    isActivated.value = true;
  } else {
    error.value = result.error || 'Ошибка активации';
  }
}
</script>

<style scoped>
.activate-page {
  max-width: 400px;
  margin: 0 auto;
}

.card {
  border: none;
  box-shadow: 0 2px 20px rgba(0,0,0,.1);
}
</style>
