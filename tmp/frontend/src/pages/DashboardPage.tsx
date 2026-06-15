import { useEffect, useState, useRef } from "react";
import { LayoutGrid, List, RefreshCw } from "lucide-react";
import toast from "react-hot-toast";
import { useStorage } from "../hooks/useStorage";
import FileCard from "../components/files/FileCard";
import FilterBar from "../components/files/FilterBar";
import Uploader from "../components/files/Uploader";
import ConfirmDeleteModal from "../components/ui/ConfirmDeleteModal";
import ActivityChart from "../components/analytics/ActivityChart";
import { FileGridSkeleton, FileListSkeleton } from "../components/ui/Skeleton";
import type { FileItem, FileFilter, FileCategory, ViewMode } from "../types";

interface DashboardContentProps {
  activeCategory: FileCategory;
}

const CATEGORY_TITLES: Record<FileCategory, string> = {
  all: "Все файлы",
  images: "Изображения",
  videos: "Видео",
  documents: "Документы",
  favorites: "Избранное",
  other: "Другие файлы",
};

export default function DashboardContent({
  activeCategory,
}: DashboardContentProps) {
  const [viewMode, setViewMode] = useState<ViewMode>("grid");
  const [fileToDelete, setFileToDelete] = useState<FileItem | null>(null);

  // Храним фильтр в ref чтобы не вызывать лишних ре-рендеров
  const currentFilterRef = useRef<FileFilter>({});

  const {
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
  } = useStorage();

  // Загружаем квоту и активность только один раз при монтировании
  useEffect(() => {
    fetchQuota();
    fetchActivity();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Перезагружаем файлы при смене категории
  useEffect(() => {
    fetchFiles(currentFilterRef.current);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeCategory]);

  const handleFilterChange = (filter: FileFilter) => {
    currentFilterRef.current = filter;
    fetchFiles(filter);
  };

  const handleRefresh = () => {
    fetchFiles(currentFilterRef.current);
  };

  const handleUploadComplete = () => {
    fetchFiles(currentFilterRef.current);
    fetchQuota();
    fetchActivity();
  };

  const handleDownload = async (file: FileItem) => {
    try {
      await downloadFile(file.id, file.name);
    } catch {
      toast.error("Ошибка при скачивании файла");
    }
  };

  const handleDeleteConfirm = async () => {
    if (!fileToDelete) return;
    try {
      await deleteNode(fileToDelete.id);
      toast.success("Файл удален");
    } catch {
      toast.error("Ошибка при удалении файла");
    }
    setFileToDelete(null);
  };

  const handleToggleFavorite = async (file: FileItem) => {
    try {
      await toggleFavorite(file.id, !file.is_favorite);
    } catch {
      toast.error("Ошибка при обновлении");
    }
  };

  const displayedFiles =
    activeCategory === "favorites" ? files.filter((f) => f.is_favorite) : files;

  return (
    <div className="flex-1 overflow-auto bg-surface-950">
      <div className="max-w-7xl mx-auto p-6 space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-bold text-white">
              {CATEGORY_TITLES[activeCategory]}
            </h2>
            <p className="text-sm text-slate-500 mt-0.5">
              {total} файл{total === 1 ? "" : total < 5 ? "а" : "ов"}
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={handleRefresh}
              className="p-2.5 rounded-xl bg-surface-900 border border-surface-800 text-slate-400 hover:text-white hover:border-surface-700 transition-all"
            >
              <RefreshCw className="w-4 h-4" />
            </button>

            <div className="flex bg-surface-900 border border-surface-800 rounded-xl p-1">
              <button
                onClick={() => setViewMode("grid")}
                className={`p-2 rounded-lg transition-all ${viewMode === "grid" ? "bg-primary-600 text-white" : "text-slate-400 hover:text-white"}`}
              >
                <LayoutGrid className="w-4 h-4" />
              </button>
              <button
                onClick={() => setViewMode("list")}
                className={`p-2 rounded-lg transition-all ${viewMode === "list" ? "bg-primary-600 text-white" : "text-slate-400 hover:text-white"}`}
              >
                <List className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        {/* Uploader */}
        <Uploader quota={quota} onUploadComplete={handleUploadComplete} />

        {/* Filters */}
        <FilterBar onFilterChange={handleFilterChange} />

        {/* Activity chart */}
        {activeCategory === "all" && activity.length > 0 && (
          <ActivityChart activity={activity} />
        )}

        {/* Files */}
        {isLoading ? (
          viewMode === "grid" ? (
            <FileGridSkeleton />
          ) : (
            <FileListSkeleton />
          )
        ) : displayedFiles.length === 0 ? (
          <div className="text-center py-24">
            <div className="w-16 h-16 rounded-2xl bg-surface-900 border border-surface-800 flex items-center justify-center mx-auto mb-4">
              <span className="text-2xl">📭</span>
            </div>
            <h3 className="text-lg font-semibold text-white mb-2">
              Файлы не найдены
            </h3>
            <p className="text-slate-500 text-sm">
              Загрузите первый файл или измените фильтры
            </p>
          </div>
        ) : viewMode === "grid" ? (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {displayedFiles.map((file) => (
              <FileCard
                key={file.id}
                file={file}
                viewMode="grid"
                onDownload={handleDownload}
                onDelete={(f) => setFileToDelete(f)}
                onToggleFavorite={handleToggleFavorite}
              />
            ))}
          </div>
        ) : (
          <div className="bg-surface-900 border border-surface-800 rounded-2xl overflow-hidden">
            <div className="flex items-center gap-4 px-4 py-2 border-b border-surface-800">
              <div className="w-9" />
              <span className="flex-1 text-xs font-medium text-slate-500 uppercase tracking-wider">
                Имя
              </span>
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider w-20 text-right">
                Размер
              </span>
              <div className="w-24" />
            </div>
            {displayedFiles.map((file) => (
              <FileCard
                key={file.id}
                file={file}
                viewMode="list"
                onDownload={handleDownload}
                onDelete={(f) => setFileToDelete(f)}
                onToggleFavorite={handleToggleFavorite}
              />
            ))}
          </div>
        )}
      </div>

      <ConfirmDeleteModal
        file={fileToDelete}
        onConfirm={handleDeleteConfirm}
        onCancel={() => setFileToDelete(null)}
      />
    </div>
  );
}
