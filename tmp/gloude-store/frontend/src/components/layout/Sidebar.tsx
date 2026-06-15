import { Link, useLocation } from "react-router-dom";
import {
  Cloud,
  HardDrive,
  Activity,
  Settings,
  LogOut,
  Star,
  Image,
  Film,
  FileText,
  FolderOpen,
} from "lucide-react";
import toast from "react-hot-toast";
import { useAuthStore } from "../../hooks/useAuthStore";
import { useStorage } from "../../hooks/useStorage";
import {
  formatFileSize,
  getQuotaColor,
  getQuotaTextColor,
} from "../../utils/files";
import type { QuotaInfo, FileCategory } from "../../types";

interface SidebarProps {
  activeCategory: FileCategory;
  onCategoryChange: (category: FileCategory) => void;
}

const NAV_ITEMS: Array<{
  id: FileCategory;
  label: string;
  icon: React.ElementType;
}> = [
  { id: "all", label: "Все файлы", icon: FolderOpen },
  { id: "images", label: "Изображения", icon: Image },
  { id: "videos", label: "Видео", icon: Film },
  { id: "documents", label: "Документы", icon: FileText },
  { id: "favorites", label: "Избранное", icon: Star },
];

export default function Sidebar({
  activeCategory,
  onCategoryChange,
}: SidebarProps) {
  const location = useLocation();
  const { user, logout } = useAuthStore();
  const { quota } = useStorage();

  const handleLogout = async () => {
    await logout();
    toast.success("Вы вышли из системы");
  };

  const pct = quota?.percentage ?? 0;

  return (
    <aside className="w-64 bg-surface-900 border-r border-surface-800 flex flex-col h-full">
      {/* Logo */}
      <div className="p-6 border-b border-surface-800">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl bg-primary-600 flex items-center justify-center shadow-lg shadow-primary-600/30">
            <Cloud className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="font-bold text-white text-sm">Gloude Store</h1>
            <p className="text-slate-500 text-xs truncate">{user?.email}</p>
          </div>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-4 space-y-1">
        <p className="text-xs font-medium text-slate-500 uppercase tracking-wider px-3 mb-3">
          Хранилище
        </p>
        {NAV_ITEMS.map(({ id, label, icon: Icon }) => (
          <button
            key={id}
            onClick={() => onCategoryChange(id)}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all ${
              activeCategory === id
                ? "bg-primary-600/20 text-primary-400 border border-primary-600/30"
                : "text-slate-400 hover:bg-surface-800 hover:text-white"
            }`}
          >
            <Icon className="w-4 h-4 flex-shrink-0" />
            {label}
          </button>
        ))}

        <div className="pt-4">
          <p className="text-xs font-medium text-slate-500 uppercase tracking-wider px-3 mb-3">
            Система
          </p>
          <Link
            to="/settings"
            className={`flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all ${
              location.pathname === "/settings"
                ? "bg-primary-600/20 text-primary-400"
                : "text-slate-400 hover:bg-surface-800 hover:text-white"
            }`}
          >
            <Settings className="w-4 h-4" />
            Настройки
          </Link>
          <Link
            to="/activity"
            className={`flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all ${
              location.pathname === "/activity"
                ? "bg-primary-600/20 text-primary-400"
                : "text-slate-400 hover:bg-surface-800 hover:text-white"
            }`}
          >
            <Activity className="w-4 h-4" />
            Активность
          </Link>
        </div>
      </nav>

      {/* Quota Widget */}
      {quota && (
        <div className="p-4 border-t border-surface-800">
          <div className="bg-surface-800 rounded-xl p-4">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <HardDrive className="w-4 h-4 text-slate-400" />
                <span className="text-sm font-medium text-slate-300">
                  Хранилище
                </span>
              </div>
              <span className={`text-xs font-bold ${getQuotaTextColor(pct)}`}>
                {pct.toFixed(1)}%
              </span>
            </div>

            {/* Progress bar */}
            <div className="h-2 bg-surface-700 rounded-full overflow-hidden mb-2">
              <div
                className={`h-full rounded-full transition-all duration-500 ${getQuotaColor(pct)}`}
                style={{ width: `${Math.min(pct, 100)}%` }}
              />
            </div>

            <div className="flex justify-between text-xs text-slate-500">
              <span>{formatFileSize(quota.used)} использовано</span>
              <span>{formatFileSize(quota.total)}</span>
            </div>
          </div>
        </div>
      )}

      {/* User & Logout */}
      <div className="p-4 border-t border-surface-800">
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-slate-400 hover:bg-red-500/10 hover:text-red-400 transition-all"
        >
          <LogOut className="w-4 h-4" />
          Выйти
        </button>
      </div>
    </aside>
  );
}
