import apiClient from './api';
import type {
  LinkGroup,
  LinkGroupCreate,
  LinkGroupUpdate,
  SuccessResponse,
  ListResponse,
  PaginationParams
} from '../types';

export const linkGroupService = {
  // Создание группы
  async create(data: LinkGroupCreate): Promise<SuccessResponse<LinkGroup>> {
    try {
      const response = await apiClient.post('/link-groups', data);
      return response.data;
    } catch (error: any) {
      console.error('Error creating link group:', error);
      throw error;
    }
  },

  // Обновление группы
  async update(data: LinkGroupUpdate): Promise<SuccessResponse<LinkGroup>> {
    try {
      const response = await apiClient.put(`/link-groups/${data.id}`, data);
      return response.data;
    } catch (error: any) {
      console.error('Error updating link group:', error);
      throw error;
    }
  },

  // Удаление группы
  async delete(id: number): Promise<SuccessResponse<void>> {
    try {
      const response = await apiClient.delete(`/link-groups/${id}`);
      return response.data;
    } catch (error: any) {
      console.error('Error deleting link group:', error);
      throw error;
    }
  },

  // Получение списка групп - ВАЖНО: всегда возвращаем структуру с массивом
  async list(params: PaginationParams = {}): Promise<SuccessResponse<ListResponse<LinkGroup>>> {
    try {
      const queryParams = new URLSearchParams();

      if (params.page) queryParams.append('page', params.page.toString());
      if (params.page_size) queryParams.append('page_size', params.page_size.toString());
      if (params.name) queryParams.append('name', params.name);

      const queryString = queryParams.toString();
      const url = `/link-groups${queryString ? `?${queryString}` : ''}`;

      const response = await apiClient.get(url);
      return response.data;
    } catch (error: any) {
      console.error('Error fetching link groups:', error);
      // Возвращаем структуру с пустым массивом при ошибке
      return {
        success: true,
        result: {
          data: [],
          total: 0,
          page: 1,
          page_size: params.page_size || 30,
          total_pages: 0
        }
      };
    }
  },

  // Получение группы по ID
  async getById(id: number): Promise<SuccessResponse<LinkGroup>> {
    try {
      const response = await apiClient.get(`/link-groups/${id}`);
      return response.data;
    } catch (error: any) {
      console.error('Error fetching link group by ID:', error);
      throw error;
    }
  }
};
