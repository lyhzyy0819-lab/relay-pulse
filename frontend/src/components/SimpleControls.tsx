import { RefreshCw, LayoutGrid, List } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { getTimeRanges } from '../constants';
import type { ViewMode } from '../types';

interface SimpleControlsProps {
  timeRange: string;
  viewMode: ViewMode;
  loading: boolean;
  onTimeRangeChange: (range: string) => void;
  onViewModeChange: (mode: ViewMode) => void;
  onRefresh: () => void;
}

export function SimpleControls({
  timeRange,
  viewMode,
  loading,
  onTimeRangeChange,
  onViewModeChange,
  onRefresh,
}: SimpleControlsProps) {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col sm:flex-row gap-3 mb-4">
      {/* 时间范围选择 */}
      <div className="flex-1 bg-slate-900/40 p-2 rounded-2xl border border-slate-800/50 backdrop-blur-md flex items-center gap-1 overflow-x-auto scrollbar-thin scrollbar-thumb-slate-700 scrollbar-track-transparent">
        {getTimeRanges(t).map((range) => (
          <button
            key={range.id}
            onClick={() => onTimeRangeChange(range.id)}
            className={`px-3 py-2 text-xs font-medium rounded-xl transition-all duration-200 whitespace-nowrap flex-shrink-0 ${
              timeRange === range.id
                ? 'bg-gradient-to-br from-cyan-500 to-blue-600 text-white shadow-lg shadow-cyan-500/25'
                : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
            }`}
          >
            {range.label}
          </button>
        ))}
      </div>

      {/* 视图切换和刷新按钮 */}
      <div className="bg-slate-900/40 p-2 rounded-2xl border border-slate-800/50 backdrop-blur-md flex items-center gap-2">
        {/* 视图切换 */}
        <div className="flex bg-slate-800 rounded-lg p-1 border border-slate-700">
          <button
            onClick={() => onViewModeChange('table')}
            className={`p-2.5 rounded min-w-[44px] min-h-[44px] flex items-center justify-center ${
              viewMode === 'table'
                ? 'bg-slate-700 text-cyan-400 shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
            title={t('controls.views.table')}
            aria-label={t('controls.views.switchToTable')}
          >
            <List size={18} />
          </button>
          <button
            onClick={() => onViewModeChange('grid')}
            className={`p-2.5 rounded min-w-[44px] min-h-[44px] flex items-center justify-center ${
              viewMode === 'grid'
                ? 'bg-slate-700 text-cyan-400 shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
            title={t('controls.views.card')}
            aria-label={t('controls.views.switchToCard')}
          >
            <LayoutGrid size={18} />
          </button>
        </div>

        {/* 刷新按钮 */}
        <button
          onClick={onRefresh}
          className="p-2.5 rounded-lg bg-cyan-500/10 text-cyan-400 hover:bg-cyan-500/20 transition-colors border border-cyan-500/20 group min-w-[44px] min-h-[44px] flex items-center justify-center"
          title={t('common.refresh')}
          aria-label={t('common.refresh')}
        >
          <RefreshCw
            size={18}
            className={`transition-transform ${loading ? 'animate-spin' : 'group-hover:rotate-180'}`}
          />
        </button>
      </div>
    </div>
  );
}
