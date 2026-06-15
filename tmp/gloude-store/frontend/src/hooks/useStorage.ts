import { useState, useCallback } from 'react';
import { storageApi } from '../api';
import type { FileItem, QuotaInfo, FileFilter, ActivityByDay } from '../types';

// Хук для управления файлами в хранилище
export function useStorage() {
  const [files, setFiles] = useState<FileItem[]>([]);
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [quota, setQuota] = useState<QuotaInfo | null>(null);
  const [activity, setActivity] = useState<ActivityByDay[]>([]);

  // Загрузить список файлов
  const fetchFiles = useCallback(async (filter: FileFilter = {}) => {
    setIsLoading(true);
    try {
      const response = await storageApi.listFiles(filter);
      setFiles(response.files || []);
      setTotal(response.total);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Загрузить информацию о квоте
  const fetchQuota = useCallback(async () => {
    const q = await storageApi.getQuota();
    setQuota(q);
  }, []);

  // Загрузить данные активности
  const fetchActivity = useCallback(async () => {
    const { activity } = await storageApi.getActivity();
    setActivity(activity || []);
  }, []);

  // Удалить файл
  const deleteFile = useCallback(async (fileId: number) => {
    await storageApi.deleteFile(fileId);
    setFiles((prev) => prev.filter((f) => f.id !== fileId));
    setTotal((prev) => prev - 1);
    // Обновляем квоту
    fetchQuota();
  }, [fetchQuota]);

  // Переключить избранное
  const toggleFavorite = useCallback(async (fileId: number, isFavorite: boolean) => {
    await storageApi.toggleFavorite(fileId, isFavorite);
    setFiles((prev) =>
      prev.map((f) => (f.id === fileId ? { ...f, is_favorite: isFavorite } : f))
    );
  }, []);

  // Скачать файл
  const downloadFile = useCallback(async (fileId: number, fileName: string) => {
    await storageApi.download(fileId, fileName);
  }, []);

  return {
    files,
    total,
    isLoading,
    quota,
    activity,
    fetchFiles,
    fetchQuota,
    fetchActivity,
    deleteFile,
    toggleFavorite,
    downloadFile,
  };
}
