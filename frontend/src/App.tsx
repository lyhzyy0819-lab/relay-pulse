import { useState } from 'react';
import { Server } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Helmet } from 'react-helmet-async';
import { Header } from './components/Header';
import { SimpleControls } from './components/SimpleControls';
import { StatusTable } from './components/StatusTable';
import { StatusCard } from './components/StatusCard';
import { Tooltip } from './components/Tooltip';
import { Footer } from './components/Footer';
import { useMonitorData } from './hooks/useMonitorData';
import { useUrlState } from './hooks/useUrlState';
import type { TooltipState, ProcessedMonitorData } from './types';

function App() {
  const { t, i18n } = useTranslation();

  // 使用 URL 状态同步
  const [state, actions] = useUrlState();
  const { timeRange, viewMode, sortConfig } = state;
  const { setTimeRange, setViewMode, setSortConfig } = actions;

  // 保留硬编码的筛选器状态（不使用，但 useMonitorData 需要）
  const filterProvider = 'all';
  const filterService = 'all';
  const filterChannel = 'all';
  const filterCategory = 'all';

  // 刷新令牌状态（暂未使用，预留给未来手动刷新功能）
  // const [reloadToken, setReloadToken] = useState(0);

  const [tooltip, setTooltip] = useState<TooltipState>({
    show: false,
    x: 0,
    y: 0,
    data: null,
  });

  const { loading, error, data, stats, slowLatencyMs } = useMonitorData({
    timeRange,
    filterService,
    filterProvider,
    filterChannel,
    filterCategory,
    sortConfig,
  });

  const handleSort = (key: string) => {
    let direction: 'asc' | 'desc' = 'desc';
    if (sortConfig.key === key && sortConfig.direction === 'desc') {
      direction = 'asc';
    }
    setSortConfig({ key, direction });
  };

  const handleRefresh = () => {
    // TODO: 实现手动刷新逻辑
    window.location.reload();
  };

  const handleBlockHover = (
    e: React.MouseEvent<HTMLDivElement>,
    point: ProcessedMonitorData['history'][number]
  ) => {
    const rect = e.currentTarget.getBoundingClientRect();
    setTooltip({
      show: true,
      x: rect.left + rect.width / 2,
      y: rect.top - 10,
      data: point,
    });
  };

  const handleBlockLeave = () => {
    setTooltip((prev) => ({ ...prev, show: false }));
  };

  return (
    <>
      {/* 动态更新 HTML meta 标签 */}
      <Helmet>
        <html lang={i18n.language} />
        <title>{t('meta.title')}</title>
        <meta name="description" content={t('meta.description')} />
      </Helmet>

      <div className="min-h-screen bg-slate-950 text-slate-200 font-sans selection:bg-cyan-500 selection:text-white overflow-x-hidden">
        {/* 全局 Tooltip */}
        <Tooltip tooltip={tooltip} onClose={handleBlockLeave} slowLatencyMs={slowLatencyMs} />

        {/* 背景装饰 */}
        <div className="fixed top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
          <div className="absolute top-[-10%] right-[-10%] w-[600px] h-[600px] bg-blue-600/10 rounded-full blur-[120px]" />
          <div className="absolute bottom-[-10%] left-[-10%] w-[600px] h-[600px] bg-cyan-600/10 rounded-full blur-[120px]" />
        </div>

        <div className="relative z-10 max-w-7xl mx-auto px-4 py-4 sm:py-6 sm:px-6 lg:px-8">
          {/* 头部 */}
          <Header stats={stats} />

          {/* 精简控制栏 */}
          <SimpleControls
            timeRange={timeRange}
            viewMode={viewMode}
            loading={loading}
            onTimeRangeChange={setTimeRange}
            onViewModeChange={setViewMode}
            onRefresh={handleRefresh}
          />

          {/* 内容区域 */}
          {error ? (
            <div className="flex flex-col items-center justify-center py-20 text-rose-400">
              <Server size={64} className="mb-4 opacity-20" />
              <p className="text-lg">{t('common.error', { message: error })}</p>
            </div>
          ) : loading && data.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-64 text-slate-500 gap-4">
              <div className="w-12 h-12 border-4 border-cyan-500/20 border-t-cyan-500 rounded-full animate-spin" />
              <p className="animate-pulse">{t('common.loading')}</p>
            </div>
          ) : data.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-20 text-slate-600">
              <Server size={64} className="mb-4 opacity-20" />
              <p className="text-lg">{t('common.noData')}</p>
            </div>
          ) : viewMode === 'table' ? (
            <StatusTable
              data={data}
              sortConfig={sortConfig}
              timeRange={timeRange}
              slowLatencyMs={slowLatencyMs}
              onSort={handleSort}
              onBlockHover={handleBlockHover}
              onBlockLeave={handleBlockLeave}
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {data.map((item) => (
                <StatusCard
                  key={`${item.providerId}-${item.serviceType}`}
                  item={item}
                  timeRange={timeRange}
                  slowLatencyMs={slowLatencyMs}
                  onBlockHover={handleBlockHover}
                  onBlockLeave={handleBlockLeave}
                />
              ))}
            </div>
          )}

          {/* 免责声明 */}
          <Footer />
        </div>
      </div>
    </>
  );
}

export default App;
