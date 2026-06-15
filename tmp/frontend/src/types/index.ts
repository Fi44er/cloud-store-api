// Типы данных приложения Gloude Store

// Типы данных приложения Gloude Store

export interface User {
  id: string;
  email: string;
  username: string;
  traits: {
    email?: string;
    username?: string;
  };
}

export interface FileItem {
  id: string;
  user_id: string;
  name: string;
  is_dir: boolean;
  mime_type: string;
  size: number;
  extension: string;
  is_favorite: boolean;
  parent_id: string | null;
  created_at: string;
  updated_at: string;
}

export interface QuotaInfo {
  total: number;
  used: number;
  pct: number;
}

export interface ActivityByDay {
  date: string;  // YYYY-MM-DD
  count: number;
}

export interface FilesResponse {
  files: FileItem[];
  total: number;
}

export interface FileFilter {
  extension?: string;
  search?: string;
  min_size?: number;
  max_size?: number;
  page?: number;
  limit?: number;
}

export interface ApiError {
  error: string;
  message?: string;
}

export type ViewMode = 'grid' | 'list';

export type FileCategory = 'all' | 'images' | 'videos' | 'documents' | 'other' | 'favorites';

// Kratos specific types
export interface KratosSession {
  active: boolean;
  identity: {
    id: string;
    traits: {
      email?: string;
      username?: string;
    };
  };
}
