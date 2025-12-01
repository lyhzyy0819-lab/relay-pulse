export function Footer() {
  return (
    <footer className="fixed bottom-0 left-0 right-0 w-full bg-slate-900/95 border-t border-slate-700/50 backdrop-blur-sm z-20">
      <p className="text-xs text-slate-500 text-center py-2 px-4">
        <a
          href="https://beian.miit.gov.cn/"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-slate-400 transition-colors"
        >
          沪ICP备2025148260号
        </a>
      </p>
    </footer>
  );
}
