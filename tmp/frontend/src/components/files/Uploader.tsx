import { useState, useRef, useCallback } from 'react';
import { Upload, X, CheckCircle, AlertCircle, CloudUpload } from 'lucide-react';
import toast from 'react-hot-toast';
import { formatFileSize } from '../../utils/files';
import { storageApi } from '../../api';
import type { QuotaInfo } from '../../types';

interface UploadItem {
  file: File;
  progress: number;
  status: 'pending' | 'uploading' | 'done' | 'error';
  error?: string;
}

interface UploaderProps {
  quota: QuotaInfo | null;
  onUploadComplete: () => void;
}

export default function Uploader({ quota, onUploadComplete }: UploaderProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [uploads, setUploads] = useState<UploadItem[]>([]);
  const [isVisible, setIsVisible] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const processFiles = useCallback(async (files: FileList | File[]) => {
    const fileArray = Array.from(files);

    // Валидация: проверяем размер против квоты
    if (quota) {
      const free = quota.total - quota.used;
      const totalUploadSize = fileArray.reduce((acc, f) => acc + f.size, 0);
      if (totalUploadSize > free) {
        toast.error(`Недостаточно места! Нужно ${formatFileSize(totalUploadSize)}, доступно ${formatFileSize(free)}`);
        return;
      }
    }

    const newUploads: UploadItem[] = fileArray.map((file) => ({
      file,
      progress: 0,
      status: 'pending',
    }));

    setUploads((prev) => [...prev, ...newUploads]);
    setIsVisible(true);

    // Загружаем файлы последовательно
    for (let i = 0; i < fileArray.length; i++) {
      const file = fileArray[i];
      const uploadIndex = uploads.length + i;

      setUploads((prev) =>
        prev.map((u, idx) =>
          idx === uploadIndex ? { ...u, status: 'uploading' } : u
        )
      );

      try {
        await storageApi.upload(file, null, (progress: number) => {
          setUploads((prev) =>
            prev.map((u, idx) =>
              idx === uploadIndex ? { ...u, progress } : u
            )
          );
        });

        setUploads((prev) =>
          prev.map((u, idx) =>
            idx === uploadIndex ? { ...u, status: 'done', progress: 100 } : u
          )
        );

        toast.success(`${file.name} загружен!`);
        onUploadComplete();
      } catch (err: unknown) {
        const error = err as { response?: { data?: { error?: string } } };
        const errMsg = error?.response?.data?.error || 'Ошибка загрузки';
        setUploads((prev) =>
          prev.map((u, idx) =>
            idx === uploadIndex ? { ...u, status: 'error', error: errMsg } : u
          )
        );
        toast.error(`${file.name}: ${errMsg}`);
      }
    }
  }, [quota, uploads.length, onUploadComplete]);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files.length > 0) {
      processFiles(e.dataTransfer.files);
    }
  }, [processFiles]);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => setIsDragging(false);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      processFiles(e.target.files);
      e.target.value = '';
    }
  };

  const clearDone = () => {
    setUploads((prev) => prev.filter((u) => u.status !== 'done'));
    if (uploads.every((u) => u.status === 'done')) {
      setIsVisible(false);
    }
  };

  return (
    <div className="space-y-4">
      {/* Drop Zone */}
      <div
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onClick={() => fileInputRef.current?.click()}
        className={`relative border-2 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all ${
          isDragging
            ? 'border-primary-500 bg-primary-600/10 scale-[1.02]'
            : 'border-surface-700 hover:border-surface-600 hover:bg-surface-800/50'
        }`}
      >
        <input
          ref={fileInputRef}
          type="file"
          multiple
          className="hidden"
          onChange={handleFileSelect}
        />

        <div className={`transition-transform ${isDragging ? 'scale-110' : ''}`}>
          <CloudUpload className={`w-12 h-12 mx-auto mb-3 transition-colors ${isDragging ? 'text-primary-400' : 'text-slate-600'}`} />
        </div>

        <p className="text-sm font-medium text-slate-300 mb-1">
          {isDragging ? 'Отпустите файлы' : 'Перетащите файлы сюда'}
        </p>
        <p className="text-xs text-slate-500">
          или <span className="text-primary-400 hover:underline">выберите файлы</span>
        </p>

        {quota && (
          <p className="text-xs text-slate-600 mt-3">
            Доступно: {formatFileSize(quota.total - quota.used)}
          </p>
        )}
      </div>

      {/* Upload queue */}
      {isVisible && uploads.length > 0 && (
        <div className="bg-surface-900 border border-surface-800 rounded-xl overflow-hidden">
          <div className="flex items-center justify-between px-4 py-3 border-b border-surface-800">
            <span className="text-sm font-medium text-slate-300">Очередь загрузки</span>
            <button
              onClick={clearDone}
              className="text-xs text-slate-500 hover:text-slate-300 transition-colors"
            >
              Очистить
            </button>
          </div>

          <div className="divide-y divide-surface-800 max-h-64 overflow-y-auto">
            {uploads.map((upload, i) => (
              <div key={i} className="px-4 py-3">
                <div className="flex items-center gap-3 mb-2">
                  {/* Status icon */}
                  {upload.status === 'done' && <CheckCircle className="w-4 h-4 text-emerald-400 flex-shrink-0" />}
                  {upload.status === 'error' && <AlertCircle className="w-4 h-4 text-red-400 flex-shrink-0" />}
                  {(upload.status === 'uploading' || upload.status === 'pending') && (
                    <Upload className="w-4 h-4 text-primary-400 flex-shrink-0 animate-pulse" />
                  )}

                  <p className="text-sm text-slate-300 truncate flex-1">{upload.file.name}</p>
                  <span className="text-xs text-slate-500 flex-shrink-0">{formatFileSize(upload.file.size)}</span>
                </div>

                {/* Progress bar */}
                {upload.status !== 'error' && (
                  <div className="h-1 bg-surface-700 rounded-full overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all duration-300 ${
                        upload.status === 'done' ? 'bg-emerald-500' : 'bg-primary-500'
                      }`}
                      style={{ width: `${upload.progress}%` }}
                    />
                  </div>
                )}

                {upload.error && (
                  <p className="text-xs text-red-400 mt-1">{upload.error}</p>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
