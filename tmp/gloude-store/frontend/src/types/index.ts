// Типы данных приложения Gloude Store

export interface User {
  id: number;
  email: string;
  username: string;
  quota_max: number;
  created_at: string;
}

export interface FileItem {
  id: number;
  user_id: number;
  original_name: string;
  storage_name: string;
  extension: string;
  mime_type: string;
  size: number;
  path: string;
  is_favorite: boolean;
  created_at: string;
  updated_at: string;
}

export interface QuotaInfo {
  total: number;
  used: number;
  free: number;
  percentage: number;
}

export interface ActivityByDay {
  date: string;  // YYYY-MM-DD
  count: number;
}

export interface FilesResponse {
  files: FileItem[];
  total: number;
  page: number;
  limit: number;
}

export interface FileFilter {
  extension?: string;
  search?: string;
  min_size?: number;
  max_size?: number;
  page?: number;
  limit?: number;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiError {
  error: string;
}

export type ViewMode = 'grid' | 'list';

export type FileCategory = 'all' | 'images' | 'videos' | 'documents' | 'other' | 'favorites';
