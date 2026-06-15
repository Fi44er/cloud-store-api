import { useMemo, useState } from 'react';
import { format, eachDayOfInterval, subYears, startOfDay } from 'date-fns';
import { ru } from 'date-fns/locale';
import type { ActivityByDay } from '../../types';

interface ActivityChartProps {
  activity: ActivityByDay[];
}

const MONTHS = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'];
const WEEKDAYS = ['Пн', '', 'Ср', '', 'Пт', '', ''];

function getIntensityClass(count: number): string {
  if (count === 0) return 'bg-surface-800 border-surface-700';
  if (count <= 2) return 'bg-primary-900 border-primary-800';
  if (count <= 5) return 'bg-primary-700 border-primary-600';
  if (count <= 10) return 'bg-primary-500 border-primary-400';
  return 'bg-primary-400 border-primary-300';
}

export default function ActivityChart({ activity }: ActivityChartProps) {
  const [tooltip, setTooltip] = useState<{ date: string; count: number; x: number; y: number } | null>(null);

  // Строим карту дата -> количество
  const activityMap = useMemo(() => {
    const map: Record<string, number> = {};
    activity.forEach(({ date, count }) => {
      map[date] = count;
    });
    return map;
  }, [activity]);

  // Генерируем массив дней за последний год
  const days = useMemo(() => {
    const end = startOfDay(new Date());
    const start = subYears(end, 1);
    return eachDayOfInterval({ start, end });
  }, []);

  // Группируем по неделям для отображения в сетке
  const weeks = useMemo(() => {
    const result: Date[][] = [];
    let currentWeek: Date[] = [];

    // Добавляем пустые дни в начало если первый день не понедельник
    const firstDay = days[0];
    const dayOfWeek = (firstDay.getDay() + 6) % 7; // Понедельник = 0
    for (let i = 0; i < dayOfWeek; i++) {
      currentWeek.push(null as unknown as Date);
    }

    days.forEach((day) => {
      currentWeek.push(day);
      if (currentWeek.length === 7) {
        result.push(currentWeek);
        currentWeek = [];
      }
    });

    if (currentWeek.length > 0) {
      while (currentWeek.length < 7) currentWeek.push(null as unknown as Date);
      result.push(currentWeek);
    }

    return result;
  }, [days]);

  // Метки месяцев для заголовка
  const monthLabels = useMemo(() => {
    const labels: Array<{ month: string; weekIndex: number }> = [];
    let lastMonth = -1;

    weeks.forEach((week, weekIndex) => {
      const firstValidDay = week.find((d) => d !== null);
      if (firstValidDay) {
        const month = firstValidDay.getMonth();
        if (month !== lastMonth) {
          labels.push({ month: MONTHS[month], weekIndex });
          lastMonth = month;
        }
      }
    });

    return labels;
  }, [weeks]);

  // Итого за год
  const totalUploads = useMemo(() => activity.reduce((sum, d) => sum + d.count, 0), [activity]);

  return (
    <div className="bg-surface-900 border border-surface-800 rounded-2xl p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-base font-semibold text-white">Активность загрузок</h3>
          <p className="text-sm text-slate-500 mt-0.5">{totalUploads} файлов за последний год</p>
        </div>

        {/* Legend */}
        <div className="flex items-center gap-2 text-xs text-slate-500">
          <span>Меньше</span>
          {['bg-surface-800', 'bg-primary-900', 'bg-primary-700', 'bg-primary-500', 'bg-primary-400'].map((cls) => (
            <div key={cls} className={`w-3 h-3 rounded-sm ${cls}`} />
          ))}
          <span>Больше</span>
        </div>
      </div>

      {/* Chart */}
      <div className="overflow-x-auto">
        <div className="min-w-max">
          {/* Month labels */}
          <div className="flex mb-1 ml-8">
            {weeks.map((_, weekIndex) => {
              const label = monthLabels.find((l) => l.weekIndex === weekIndex);
              return (
                <div key={weekIndex} className="w-3 mr-[2px] text-xs text-slate-600 text-center">
                  {label ? label.month : ''}
                </div>
              );
            })}
          </div>

          <div className="flex gap-[2px]">
            {/* Weekday labels */}
            <div className="flex flex-col gap-[2px] mr-1">
              {WEEKDAYS.map((day, i) => (
                <div key={i} className="h-3 text-[10px] text-slate-600 flex items-center justify-end pr-1 w-7">
                  {day}
                </div>
              ))}
            </div>

            {/* Weeks grid */}
            {weeks.map((week, weekIndex) => (
              <div key={weekIndex} className="flex flex-col gap-[2px]">
                {week.map((day, dayIndex) => {
                  if (!day) {
                    return <div key={dayIndex} className="w-3 h-3" />;
                  }

                  const dateStr = format(day, 'yyyy-MM-dd');
                  const count = activityMap[dateStr] || 0;

                  return (
                    <div
                      key={dayIndex}
                      className={`w-3 h-3 rounded-sm border cursor-pointer transition-all hover:opacity-80 ${getIntensityClass(count)}`}
                      onMouseEnter={(e) => {
                        const rect = (e.target as HTMLElement).getBoundingClientRect();
                        setTooltip({
                          date: format(day, 'd MMMM yyyy', { locale: ru }),
                          count,
                          x: rect.left,
                          y: rect.top,
                        });
                      }}
                      onMouseLeave={() => setTooltip(null)}
                    />
                  );
                })}
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Tooltip */}
      {tooltip && (
        <div
          className="fixed z-50 px-3 py-2 rounded-lg bg-surface-800 border border-surface-700 text-xs text-white shadow-xl pointer-events-none"
          style={{
            left: tooltip.x + 12,
            top: tooltip.y - 40,
          }}
        >
          <div className="font-medium">{tooltip.date}</div>
          <div className="text-slate-400">
            {tooltip.count === 0 ? 'Нет загрузок' : `${tooltip.count} файл${tooltip.count === 1 ? '' : tooltip.count < 5 ? 'а' : 'ов'}`}
          </div>
        </div>
      )}
    </div>
  );
}
