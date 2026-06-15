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
    } catch (error) {
      console.error("Failed to fetch files:", error);
      setFiles([]);
      setTotal(0);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Загрузить информацию о квоте
  const fetchQuota = useCallback(async () => {
    try {
      const q = await storageApi.getQuota();
      setQuota(q);
    } catch (error) {
      console.error("Failed to fetch quota:", error);
    }
  }, []);

  // Загрузить данные активности
  const fetchActivity = useCallback(async () => {
    try {
      const data = await storageApi.getActivity();
      setActivity(data || []);
    } catch (error) {
      console.error("Failed to fetch activity:", error);
      setActivity([]);
    }
  }, []);

  // Удалить файл
  const deleteNode = useCallback(async (fileId: string) => {
    try {
      await storageApi.deleteFile(fileId);
      setFiles((prev: FileItem[]) => prev.filter((f: FileItem) => f.id !== fileId));
      setTotal((prev: number) => prev - 1);
      // Обновляем квоту
      fetchQuota();
    } catch (error) {
      console.error("Failed to delete node:", error);
      throw error;
    }
  }, [fetchQuota]);

  // Переключить избранное
  const toggleFavorite = useCallback(async (fileId: string, isFavorite: boolean) => {
    try {
      await storageApi.toggleFavorite(fileId, isFavorite);
      setFiles((prev: FileItem[]) =>
        prev.map((f: FileItem) => (f.id === fileId ? { ...f, is_favorite: isFavorite } : f))
      );
    } catch (error) {
      console.error("Failed to toggle favorite:", error);
      throw error;
    }
  }, []);

  // Скачать файл
  const downloadFile = useCallback(async (fileId: string, fileName: string) => {
    try {
      await storageApi.download(fileId, fileName);
    } catch (error) {
      console.error("Failed to download file:", error);
      throw error;
    }
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
    deleteNode,
    toggleFavorite,
    downloadFile,
  };
}
