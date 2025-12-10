import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { LinkGroup, LinkGroupCreate, LinkGroupUpdate, PaginationParams } from '@/types';
import { linkGroupService } from '@/services';

export const useLinkGroupStore = defineStore('linkGroup', () => {
  // Состояние
  const groups = ref<LinkGroup[]>([]);
  const currentGroup = ref<LinkGroup | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const pagination = ref({
    page: 1,
    page_size: 30,
    total: 0,
    total_pages: 0
  });

  // Действия
  async function fetchGroups(params: PaginationParams = {}) {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await linkGroupService.list({
        page: params.page || pagination.value.page,
        page_size: params.page_size || pagination.value.page_size,
        name: params.name
      });

      // Убедимся, что всегда есть массив
      groups.value = response.result.data || [];
      pagination.value = {
        page: response.result.page,
        page_size: response.result.page_size,
        total: response.result.total,
        total_pages: response.result.total_pages
      };

      return { success: true };
    } catch (err: any) {
      // В случае ошибки сбрасываем группы в пустой массив
      groups.value = [];
      error.value = err.response?.data?.error?.message || 'Ошибка загрузки групп';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  async function createGroup(data: LinkGroupCreate) {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await linkGroupService.create(data);
      // Добавляем в начало массива
      groups.value.unshift(response.result);
      return { success: true, group: response.result };
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Ошибка создания группы';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  async function updateGroup(data: LinkGroupUpdate) {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await linkGroupService.update(data);

      // Обновляем группу в списке
      const index = groups.value.findIndex(g => g.id === data.id);
      if (index !== -1) {
        groups.value[index] = response.result;
      }

      return { success: true, group: response.result };
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Ошибка обновления группы';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteGroup(id: number) {
    isLoading.value = true;
    error.value = null;

    try {
      await linkGroupService.delete(id);

      // Удаляем группу из списка
      groups.value = groups.value.filter(g => g.id !== id);

      return { success: true };
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Ошибка удаления группы';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchGroupById(id: number) {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await linkGroupService.getById(id);
      currentGroup.value = response.result;
      return { success: true };
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Ошибка загрузки группы';
      return { success: false, error: error.value };
    } finally {
      isLoading.value = false;
    }
  }

  function setCurrentGroup(group: LinkGroup | null) {
    currentGroup.value = group;
  }

  return {
    // Состояние
    groups,
    currentGroup,
    isLoading,
    error,
    pagination,

    // Действия
    fetchGroups,
    createGroup,
    updateGroup,
    deleteGroup,
    fetchGroupById,
    setCurrentGroup
  };
});
