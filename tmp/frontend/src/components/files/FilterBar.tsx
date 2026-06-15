import { useState, useCallback } from 'react';
import { Search, Filter, X } from 'lucide-react';
import type { FileFilter } from '../../types';

interface FilterBarProps {
  onFilterChange: (filter: FileFilter) => void;
}

const EXTENSIONS = [
  { label: 'Все', value: '' },
  { label: 'JPG', value: '.jpg' },
  { label: 'PNG', value: '.png' },
  { label: 'PDF', value: '.pdf' },
  { label: 'MP4', value: '.mp4' },
  { label: 'ZIP', value: '.zip' },
  { label: 'DOCX', value: '.docx' },
];

export default function FilterBar({ onFilterChange }: FilterBarProps) {
  const [search, setSearch] = useState('');
  const [extension, setExtension] = useState('');
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [minSize, setMinSize] = useState('');
  const [maxSize, setMaxSize] = useState('');

  const applyFilters = useCallback((
    s: string,
    ext: string,
    min: string,
    max: string
  ) => {
    const filter: FileFilter = {
      search: s || undefined,
      extension: ext || undefined,
      min_size: min ? Number(min) * 1024 : undefined,
      max_size: max ? Number(max) * 1024 * 1024 : undefined,
      page: 1,
    };
    onFilterChange(filter);
  }, [onFilterChange]);

  const handleSearchChange = (val: string) => {
    setSearch(val);
    applyFilters(val, extension, minSize, maxSize);
  };

  const handleExtensionChange = (val: string) => {
    setExtension(val);
    applyFilters(search, val, minSize, maxSize);
  };

  const handleMinSizeChange = (val: string) => {
    setMinSize(val);
    applyFilters(search, extension, val, maxSize);
  };

  const handleMaxSizeChange = (val: string) => {
    setMaxSize(val);
    applyFilters(search, extension, minSize, val);
  };

  const handleReset = () => {
    setSearch('');
    setExtension('');
    setMinSize('');
    setMaxSize('');
    onFilterChange({ page: 1 });
  };

  const hasActiveFilters = search || extension || minSize || maxSize;

  return (
    <div className="space-y-3">
      <div className="flex gap-2">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
          <input
            type="text"
            value={search}
            onChange={(e) => handleSearchChange(e.target.value)}
            placeholder="Поиск по имени файла..."
            className="w-full bg-surface-900 border border-surface-800 rounded-xl pl-10 pr-4 py-2.5 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-primary-500 focus:ring-1 focus:ring-primary-500 transition-colors"
          />
          {search && (
            <button
              onClick={() => handleSearchChange('')}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300"
            >
              <X className="w-4 h-4" />
            </button>
          )}
        </div>

        <button
          onClick={() => setShowAdvanced(!showAdvanced)}
          className={`flex items-center gap-2 px-4 py-2.5 rounded-xl border text-sm font-medium transition-all ${
            showAdvanced || hasActiveFilters
              ? 'bg-primary-600/20 border-primary-600/30 text-primary-400'
              : 'bg-surface-900 border-surface-800 text-slate-400 hover:border-surface-700 hover:text-white'
          }`}
        >
          <Filter className="w-4 h-4" />
          Фильтры
          {hasActiveFilters && <span className="w-2 h-2 rounded-full bg-primary-400" />}
        </button>

        {hasActiveFilters && (
          <button
            onClick={handleReset}
            className="px-4 py-2.5 rounded-xl border border-surface-800 bg-surface-900 text-sm text-slate-400 hover:text-red-400 hover:border-red-500/30 hover:bg-red-500/10 transition-all"
          >
            <X className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* Extension quick filters */}
      <div className="flex gap-2 flex-wrap">
        {EXTENSIONS.map(({ label, value }) => (
          <button
            key={value}
            onClick={() => handleExtensionChange(value)}
            className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all ${
              extension === value
                ? 'bg-primary-600 text-white'
                : 'bg-surface-800 text-slate-400 hover:bg-surface-700 hover:text-white'
            }`}
          >
            {label}
          </button>
        ))}
      </div>

      {/* Advanced filters */}
      {showAdvanced && (
        <div className="grid grid-cols-2 gap-3 p-4 bg-surface-900 border border-surface-800 rounded-xl">
          <div>
            <label className="block text-xs text-slate-500 mb-1.5">Мин. размер (KB)</label>
            <input
              type="number"
              value={minSize}
              onChange={(e) => handleMinSizeChange(e.target.value)}
              placeholder="0"
              className="w-full bg-surface-800 border border-surface-700 rounded-lg px-3 py-2 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-primary-500 transition-colors"
            />
          </div>
          <div>
            <label className="block text-xs text-slate-500 mb-1.5">Макс. размер (MB)</label>
            <input
              type="number"
              value={maxSize}
              onChange={(e) => handleMaxSizeChange(e.target.value)}
              placeholder="∞"
              className="w-full bg-surface-800 border border-surface-700 rounded-lg px-3 py-2 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-primary-500 transition-colors"
            />
          </div>
        </div>
      )}
    </div>
  );
}
