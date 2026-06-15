// Утилиты для работы с файлами

// Форматирование размера файла в человекочитаемый формат
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
}

// Определение иконки по расширению файла
export function getFileCategory(extension: string, mimeType: string): 'image' | 'video' | 'audio' | 'document' | 'archive' | 'code' | 'other' {
  const ext = extension.toLowerCase().replace('.', '');
  const mime = mimeType.toLowerCase();

  if (mime.startsWith('image/') || ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico'].includes(ext)) {
    return 'image';
  }
  if (mime.startsWith('video/') || ['mp4', 'avi', 'mkv', 'mov', 'webm', 'flv'].includes(ext)) {
    return 'video';
  }
  if (mime.startsWith('audio/') || ['mp3', 'wav', 'ogg', 'flac', 'aac'].includes(ext)) {
    return 'audio';
  }
  if (['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt', 'odt', 'ods'].includes(ext)) {
    return 'document';
  }
  if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2'].includes(ext)) {
    return 'archive';
  }
  if (['js', 'ts', 'py', 'go', 'java', 'cpp', 'c', 'h', 'html', 'css', 'json', 'xml', 'sh', 'rb', 'php'].includes(ext)) {
    return 'code';
  }
  return 'other';
}

// Цвет иконки по категории файла
export const FILE_CATEGORY_COLORS: Record<string, string> = {
  image: 'text-emerald-400',
  video: 'text-purple-400',
  audio: 'text-pink-400',
  document: 'text-blue-400',
  archive: 'text-amber-400',
  code: 'text-cyan-400',
  other: 'text-slate-400',
};

// Расширения по категориям (для фильтра)
export const CATEGORY_EXTENSIONS: Record<string, string[]> = {
  images: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp'],
  videos: ['mp4', 'avi', 'mkv', 'mov', 'webm'],
  documents: ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt'],
  other: [],
};

// Форматирование даты
export function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

// Получение цвета для прогресс-бара квоты
export function getQuotaColor(percentage: number): string {
  if (percentage < 60) return 'bg-emerald-500';
  if (percentage < 80) return 'bg-amber-500';
  return 'bg-red-500';
}

export function getQuotaTextColor(percentage: number): string {
  if (percentage < 60) return 'text-emerald-400';
  if (percentage < 80) return 'text-amber-400';
  return 'text-red-400';
}
