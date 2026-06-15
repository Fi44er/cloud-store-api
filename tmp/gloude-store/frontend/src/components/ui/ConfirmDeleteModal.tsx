import { Trash2, X } from 'lucide-react';
import type { FileItem } from '../../types';

interface ConfirmDeleteModalProps {
  file: FileItem | null;
  onConfirm: () => void;
  onCancel: () => void;
}

export default function ConfirmDeleteModal({ file, onConfirm, onCancel }: ConfirmDeleteModalProps) {
  if (!file) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={onCancel}
      />

      {/* Modal */}
      <div className="relative bg-surface-900 border border-surface-800 rounded-2xl p-6 w-full max-w-sm shadow-2xl animate-slide-up">
        <button
          onClick={onCancel}
          className="absolute top-4 right-4 text-slate-500 hover:text-white transition-colors"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="w-12 h-12 rounded-2xl bg-red-500/20 flex items-center justify-center mb-4">
          <Trash2 className="w-6 h-6 text-red-400" />
        </div>

        <h3 className="text-lg font-semibold text-white mb-2">Удалить файл?</h3>
        <p className="text-sm text-slate-400 mb-6">
          <span className="text-white font-medium">«{file.original_name}»</span> будет удален без возможности восстановления.
        </p>

        <div className="flex gap-3">
          <button
            onClick={onCancel}
            className="flex-1 py-2.5 rounded-xl border border-surface-700 text-sm font-medium text-slate-300 hover:bg-surface-800 hover:text-white transition-all"
          >
            Отмена
          </button>
          <button
            onClick={onConfirm}
            className="flex-1 py-2.5 rounded-xl bg-red-600 hover:bg-red-700 text-sm font-medium text-white transition-all shadow-lg shadow-red-600/20"
          >
            Удалить
          </button>
        </div>
      </div>
    </div>
  );
}
