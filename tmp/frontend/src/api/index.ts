import api from "./axios";
import type {
  FileFilter,
  FilesResponse,
  QuotaInfo,
  ActivityByDay,
  KratosSession,
} from "../types";

// ===== AUTH (Kratos) =====

export const authApi = {
  // Инициализация флоу
  initLoginFlow: async (): Promise<{ flow_id: string }> => {
    const { data } = await api.get<{ flow_id: string }>("/auth/login/flow");
    return data;
  },
  
  initRegistrationFlow: async (): Promise<{ flow_id: string }> => {
    const { data } = await api.get<{ flow_id: string }>("/auth/registration/flow");
    return data;
  },

  // Сабмит флоу
  login: async (flowID: string, identifier: string, password: string) => {
    const { data } = await api.post("/auth/login", {
      flow_id: flowID,
      identifier,
      password,
    });
    return data;
  },

  register: async (flowID: string, email: string, username: string, password: string) => {
    const { data } = await api.post("/auth/registration", {
      flow_id: flowID,
      email,
      username,
      password,
    });
    return data;
  },

  // Получение текущей сессии
  whoAmI: async (): Promise<KratosSession> => {
    const { data } = await api.get<KratosSession>("/account/me");
    return data;
  },
  
  // Выход
  logout: async () => {
    const { data } = await api.post("/account/logout");
    return data;
  }
};

// ===== STORAGE =====

export const storageApi = {
  // Загрузка файла с прогрессом
  upload: async (file: File, parentId: string | null = null, onProgress?: (progress: number) => void) => {
    const formData = new FormData();
    formData.append("file", file);
    if (parentId) formData.append("parent_id", parentId);

    const { data } = await api.post("/files/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
      onUploadProgress: (event: any) => {
        if (event.total && onProgress) {
          const progress = Math.round((event.loaded * 100) / event.total);
          onProgress(progress);
        }
      },
    });
    return data;
  },

  // Создание папки
  createFolder: async (name: string, parentId: string | null = null) => {
    const { data } = await api.post("/files/folder", { name, parent_id: parentId });
    return data;
  },

  // Список файлов с фильтрами
  listFiles: async (filter: FileFilter = {}): Promise<FilesResponse> => {
    const params = new URLSearchParams();
    if (filter.extension) params.set("extension", filter.extension);
    if (filter.search) params.set("search", filter.search);
    if (filter.page) params.set("page", String(filter.page));
    if (filter.limit) params.set("limit", String(filter.limit));

    const { data } = await api.get<FilesResponse>(`/files`, { params });
    return data;
  },

  // Удаление файла
  deleteFile: async (fileId: string) => {
    const { data } = await api.delete(`/files/${fileId}`);
    return data;
  },

  // Скачивание файла
  getDownloadUrl: (fileId: string): string => {
    const base = (import.meta as any).env.VITE_API_URL || "/api/v1";
    return `${base}/files/download/${fileId}`;
  },

  download: async (fileId: string, fileName: string) => {
    const response = await api.get(`/files/download/${fileId}`, {
      responseType: "blob",
    });

    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", fileName);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  },

  // Квота
  getQuota: async (): Promise<QuotaInfo> => {
    const { data } = await api.get<QuotaInfo>("/files/quota");
    return data;
  },

  // Активность
  getActivity: async (): Promise<ActivityByDay[]> => {
    const { data } = await api.get<ActivityByDay[]>("/files/activity");
    return data;
  },

  // Избранное
  toggleFavorite: async (fileId: string, isFavorite: boolean) => {
    const { data } = await api.post(`/files/favorite/${fileId}`, {
      is_favorite: isFavorite,
    });
    return data;
  },
};
