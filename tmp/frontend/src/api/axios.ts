import axios from "axios";

// Базовый URL API из переменных окружения
const BASE_URL = (import.meta as any).env.VITE_API_URL || "/api/v1";

// Создаем экземпляр Axios с настройками
const api = axios.create({
  baseURL: BASE_URL,
  withCredentials: true, // Отправляем cookies
  headers: {
    "Content-Type": "application/json",
  },
});

// Interceptor запросов: не нужен для Kratos, так как он использует cookies
api.interceptors.request.use((config) => {
  return config;
});

// Interceptor ответов: обработка 401 (разлогин)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 401 handle by auth store and ProtectedRoute
    }
    return Promise.reject(error);
  },
);

export default api;
