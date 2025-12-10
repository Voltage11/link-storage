// Типы для аутентификации
export interface User {
  id: number;
  name: string;
  email: string;
  is_active: boolean;
  is_admin: boolean;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CurrentUser {
  id: number;
  name: string;
  email: string;
  is_active: boolean;
  is_admin: boolean;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface SessionResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

export interface LoginResponse {
  success: boolean;
  result: SessionResponse;
}

export interface ActivateRequest {
  code: string;
}

// Типы для групп ссылок
export interface LinkGroup {
  id: number;
  user_id: number;
  name: string;
  description: string;
  position: number;
  color: string;
  created_at: string;
  updated_at: string;
}

export interface LinkGroupCreate {
  name: string;
  description: string;
  color: string;
}

export interface LinkGroupUpdate {
  id: number;
  name: string;
  description: string;
  color: string;
}

export interface PaginationParams {
  page?: number;
  page_size?: number;
  name?: string;
}

export interface ListResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// Общие типы ответов
export interface SuccessResponse<T = any> {
  success: boolean;
  result: T;
}

export interface ErrorResponse {
  success: boolean;
  error: {
    type: string;
    message: string;
    code?: string;
    details?: any;
  };
}
