import {
  Download,
  Trash2,
  Star,
  Image,
  Film,
  Music,
  FileText,
  Archive,
  Code,
  File,
} from 'lucide-react';
import { formatFileSize, formatDate, getFileCategory, FILE_CATEGORY_COLORS } from '../../utils/files';
import type { FileItem, ViewMode } from '../../types';

interface FileCardProps {
  file: FileItem;
  viewMode: ViewMode;
  onDownload: (file: FileItem) => void;
  onDelete: (file: FileItem) => void;
  onToggleFavorite: (file: FileItem) => void;
}

const CATEGORY_ICONS: Record<string, React.ElementType> = {
  image: Image,
  video: Film,
  audio: Music,
  document: FileText,
  archive: Archive,
  code: Code,
  other: File,
};

export default function FileCard({ file, viewMode, onDownload, onDelete, onToggleFavorite }: FileCardProps) {
  const category = getFileCategory(file.extension, file.mime_type);
  const IconComponent = CATEGORY_ICONS[category] || File;
  const iconColor = FILE_CATEGORY_COLORS[category] || 'text-slate-400';

  if (viewMode === 'grid') {
    return (
      <div className="group bg-surface-900 border border-surface-800 rounded-xl p-4 hover:border-surface-700 transition-all hover:shadow-lg hover:shadow-black/20 relative">
        {/* File icon */}
        <div className={`w-12 h-12 rounded-xl bg-surface-800 flex items-center justify-center mb-3 ${iconColor}`}>
          <IconComponent className="w-6 h-6" />
        </div>

        {/* File name */}
        <p className="text-sm font-medium text-white truncate mb-1" title={file.original_name}>
          {file.original_name}
        </p>
        <p className="text-xs text-slate-500">{formatFileSize(file.size)}</p>
        <p className="text-xs text-slate-600 mt-0.5">{formatDate(file.created_at)}</p>

        {/* Favorite */}
        <button
          onClick={() => onToggleFavorite(file)}
          className={`absolute top-3 right-3 transition-all ${
            file.is_favorite
              ? 'text-amber-400'
              : 'text-slate-600 opacity-0 group-hover:opacity-100 hover:text-amber-400'
          }`}
        >
          <Star className="w-4 h-4" fill={file.is_favorite ? 'currentColor' : 'none'} />
        </button>

        {/* Actions */}
        <div className="flex gap-2 mt-3 pt-3 border-t border-surface-800 opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            onClick={() => onDownload(file)}
            className="flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg bg-surface-800 hover:bg-primary-600/20 text-slate-400 hover:text-primary-400 text-xs font-medium transition-all"
          >
            <Download className="w-3.5 h-3.5" />
            Скачать
          </button>
          <button
            onClick={() => onDelete(file)}
            className="flex items-center justify-center px-2.5 py-1.5 rounded-lg bg-surface-800 hover:bg-red-500/20 text-slate-400 hover:text-red-400 transition-all"
          >
            <Trash2 className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    );
  }

  // List view
  return (
    <div className="group flex items-center gap-4 px-4 py-3 hover:bg-surface-800/50 rounded-xl transition-all border border-transparent hover:border-surface-700">
      {/* Icon */}
      <div className={`w-9 h-9 rounded-lg bg-surface-800 flex items-center justify-center flex-shrink-0 ${iconColor}`}>
        <IconComponent className="w-4.5 h-4.5" />
      </div>

      {/* Name & info */}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-white truncate">{file.original_name}</p>
        <p className="text-xs text-slate-500 mt-0.5">{formatDate(file.created_at)}</p>
      </div>

      {/* Size */}
      <span className="text-xs text-slate-400 w-20 text-right flex-shrink-0">{formatFileSize(file.size)}</span>

      {/* Actions */}
      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          onClick={() => onToggleFavorite(file)}
          className={`p-2 rounded-lg hover:bg-surface-700 transition-all ${
            file.is_favorite ? 'text-amber-400' : 'text-slate-500 hover:text-amber-400'
          }`}
        >
          <Star className="w-4 h-4" fill={file.is_favorite ? 'currentColor' : 'none'} />
        </button>
        <button
          onClick={() => onDownload(file)}
          className="p-2 rounded-lg hover:bg-primary-600/20 text-slate-500 hover:text-primary-400 transition-all"
        >
          <Download className="w-4 h-4" />
        </button>
        <button
          onClick={() => onDelete(file)}
          className="p-2 rounded-lg hover:bg-red-500/20 text-slate-500 hover:text-red-400 transition-all"
        >
          <Trash2 className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
