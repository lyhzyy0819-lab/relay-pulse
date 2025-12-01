import { Activity, CheckCircle, AlertTriangle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { TypewriterText } from './TypewriterText';

interface HeaderProps {
  stats: {
    total: number;
    healthy: number;
    issues: number;
  };
}

export function Header({ stats }: HeaderProps) {
  const { t } = useTranslation();

  return (
    <header className="flex flex-col md:flex-row justify-between items-start md:items-center mb-4 gap-4 border-b border-slate-800/50 pb-3">
      {/* 左侧：Logo 和标语 */}
      <div>
        <div className="flex items-center gap-2 sm:gap-3 mb-1 sm:mb-2">
          <div className="p-1.5 sm:p-2 bg-amber-500/10 rounded-lg border border-amber-500/20">
            <Activity className="w-5 h-5 sm:w-6 sm:h-6 text-amber-400" />
          </div>
          <h1 className="text-2xl sm:text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-amber-400 via-amber-500 to-orange-500">
            <TypewriterText text="Code-CLI" className="inline-block min-w-[120px] sm:min-w-[160px]" />
          </h1>
        </div>
        <p className="text-slate-400 text-xs sm:text-sm flex items-center gap-2">
          <span className="inline-block w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          {t('header.tagline')}
        </p>
      </div>

      {/* 右侧：统计卡片 */}
      <div className="grid grid-cols-2 sm:flex gap-2 sm:gap-4">
        <div className="px-3 sm:px-4 py-2 rounded-xl bg-slate-900/50 border border-slate-800 backdrop-blur-sm flex items-center gap-2 sm:gap-3 shadow-lg">
          <div className="p-1 sm:p-1.5 rounded-full bg-emerald-500/10 text-emerald-400">
            <CheckCircle size={14} className="sm:w-4 sm:h-4" />
          </div>
          <div>
            <div className="text-slate-400 text-[10px] sm:text-xs">{t('header.stats.healthy')}</div>
            <div className="font-mono font-bold text-emerald-400 text-sm sm:text-base">
              {stats.healthy}
            </div>
          </div>
        </div>
        <div className="px-3 sm:px-4 py-2 rounded-xl bg-slate-900/50 border border-slate-800 backdrop-blur-sm flex items-center gap-2 sm:gap-3 shadow-lg">
          <div className="p-1 sm:p-1.5 rounded-full bg-rose-500/10 text-rose-400">
            <AlertTriangle size={14} className="sm:w-4 sm:h-4" />
          </div>
          <div>
            <div className="text-slate-400 text-[10px] sm:text-xs">{t('header.stats.issues')}</div>
            <div className="font-mono font-bold text-rose-400 text-sm sm:text-base">
              {stats.issues}
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
