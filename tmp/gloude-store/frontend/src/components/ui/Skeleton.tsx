interface SkeletonProps {
  className?: string;
}

export function Skeleton({ className = '' }: SkeletonProps) {
  return (
    <div className={`animate-pulse bg-surface-800 rounded-lg ${className}`} />
  );
}

export function FileGridSkeleton() {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      {Array.from({ length: 12 }).map((_, i) => (
        <div key={i} className="bg-surface-900 border border-surface-800 rounded-xl p-4">
          <Skeleton className="w-12 h-12 mb-3" />
          <Skeleton className="h-3 w-full mb-2" />
          <Skeleton className="h-2 w-2/3" />
        </div>
      ))}
    </div>
  );
}

export function FileListSkeleton() {
  return (
    <div className="space-y-2">
      {Array.from({ length: 8 }).map((_, i) => (
        <div key={i} className="flex items-center gap-4 px-4 py-3">
          <Skeleton className="w-9 h-9 rounded-lg flex-shrink-0" />
          <div className="flex-1">
            <Skeleton className="h-3 w-2/3 mb-2" />
            <Skeleton className="h-2 w-1/3" />
          </div>
          <Skeleton className="h-2 w-16" />
        </div>
      ))}
    </div>
  );
}
