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

// Interceptor запросов: добавляем JWT токен если он есть
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("auth_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor ответов: обработка 401 (разлогин)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Очищаем токен и редиректим на логин
      localStorage.removeItem("auth_token");
      localStorage.removeItem("auth_user");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

export default api;
