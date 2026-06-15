import api from "./axios";
import type {
  AuthResponse,
  FileFilter,
  FilesResponse,
  QuotaInfo,
  ActivityByDay,
  User,
} from "../types";

// ===== AUTH =====

export const authApi = {
  register: async (email: string, username: string, password: string) => {
    const { data } = await api.post("/auth/register", {
      email,
      username,
      password,
    });
    return data;
  },

  login: async (email: string, password: string): Promise<AuthResponse> => {
    const { data } = await api.post<AuthResponse>("/auth/login", {
      email,
      password,
    });
    return data;
  },

  logout: async () => {
    const { data } = await api.post("/auth/logout");
    return data;
  },

  me: async (): Promise<User> => {
    const { data } = await api.get<User>("/account/me");
    return data;
  },
};

// ===== STORAGE =====

export const storageApi = {
  // Загрузка файла с прогрессом
  upload: async (file: File, onProgress?: (progress: number) => void) => {
    const formData = new FormData();
    formData.append("file", file);

    const { data } = await api.post("/storage/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
      onUploadProgress: (event) => {
        if (event.total && onProgress) {
          const progress = Math.round((event.loaded * 100) / event.total);
          onProgress(progress);
        }
      },
    });
    return data;
  },

  // Список файлов с фильтрами
  listFiles: async (filter: FileFilter = {}): Promise<FilesResponse> => {
    const params = new URLSearchParams();
    if (filter.extension) params.set("extension", filter.extension);
    if (filter.search) params.set("search", filter.search);
    if (filter.min_size) params.set("min_size", String(filter.min_size));
    if (filter.max_size) params.set("max_size", String(filter.max_size));
    if (filter.page) params.set("page", String(filter.page));
    if (filter.limit) params.set("limit", String(filter.limit));

    const { data } = await api.get<FilesResponse>(`/storage/files?${params}`);
    return data;
  },

  // Скачивание файла
  getDownloadUrl: (fileId: number): string => {
    const token = localStorage.getItem("auth_token");
    const base = (import.meta as any).env.VITE_API_URL || "/api/v1";

    return `${base}/storage/download/${fileId}?token=${token}`;
  },

  download: async (fileId: number, fileName: string) => {
    const token = localStorage.getItem("auth_token");
    const response = await api.get(`/storage/download/${fileId}`, {
      responseType: "blob",
      headers: { Authorization: `Bearer ${token}` },
    });

    // Создаем ссылку для скачивания
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", fileName);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  },

  // Удаление файла
  deleteFile: async (fileId: number) => {
    const { data } = await api.delete(`/storage/${fileId}`);
    return data;
  },

  // Квота
  getQuota: async (): Promise<QuotaInfo> => {
    const { data } = await api.get<QuotaInfo>("/storage/quota");
    return data;
  },

  // Активность
  getActivity: async (): Promise<{ activity: ActivityByDay[] }> => {
    const { data } = await api.get("/storage/activity");
    return data;
  },

  // Избранное
  toggleFavorite: async (fileId: number, isFavorite: boolean) => {
    const { data } = await api.patch(`/storage/${fileId}/favorite`, {
      is_favorite: isFavorite,
    });
    return data;
  },
};
