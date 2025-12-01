import { Github, Tag, Bug } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useVersionInfo } from '../hooks/useVersionInfo';
import { FEEDBACK_URLS } from '../constants';

export function Footer() {
  const { t } = useTranslation();
  const { versionInfo } = useVersionInfo();

  return (
    <footer className="mt-4 bg-slate-900/60 border border-slate-800 rounded-2xl p-4 sm:p-5 text-slate-400">
      {/* 简化的免责声明 */}
      <div className="rounded-lg border border-slate-700 bg-slate-800/50 p-4 mb-4">
        <p className="text-sm text-slate-400 text-center">
          {t('footer.disclaimer.text')}
        </p>
      </div>

      {/* GitHub 链接与版本信息 */}
      <div className="border-t border-slate-800/50 pt-4 flex flex-col sm:flex-row items-center justify-center gap-2 text-xs">
        <div className="flex items-center gap-2 flex-wrap justify-center">
          <a
            href="https://github.com/prehisle/relay-pulse"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800/50 border border-slate-700/50 text-slate-300 hover:text-cyan-300 hover:border-cyan-500/30 transition min-h-[36px]"
          >
            <Github size={14} />
            <span>GitHub</span>
          </a>
          <span className="hidden sm:inline text-slate-600">·</span>
          <a
            href={FEEDBACK_URLS.BUG_REPORT}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800/50 border border-slate-700/50 text-slate-300 hover:text-rose-300 hover:border-rose-500/30 transition min-h-[36px]"
          >
            <Bug size={14} />
            <span>{t('footer.issuesBtn')}</span>
          </a>
          <span className="hidden sm:inline text-slate-600">·</span>
          <span className="text-slate-500 text-[11px] sm:text-xs">{t('footer.openSourceLabel')}</span>
        </div>
        {versionInfo && (
          <>
            <span className="hidden sm:inline text-slate-600">·</span>
            <div
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800/50 border border-slate-700/50 text-slate-400"
              title={`Commit: ${versionInfo.git_commit} | Built: ${versionInfo.build_time}`}
            >
              <Tag size={14} className="text-slate-500" />
              <span className="text-slate-400">{versionInfo.version}</span>
            </div>
          </>
        )}
      </div>

      {/* ICP 备案信息 */}
      <div className="mt-4 pt-4 border-t border-slate-700">
        <p className="text-xs text-slate-500 text-center">
          {t('footer.icp')}
        </p>
      </div>
    </footer>
  );
}
