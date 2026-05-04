import { NavLink, Outlet } from 'react-router';
import { Shield, Settings as SettingsIcon, Activity, ScanSearch, X, Minus, Square } from 'lucide-react';
import { Toaster } from 'sonner';
import { useScanner } from '../context/ScannerContext';
import { translations } from '../i18n';

export function Layout() {
  const { language } = useScanner();
  const t = translations[language];

  return (
    <div className="flex h-screen w-full bg-slate-950 text-slate-200 font-sans overflow-hidden selection:bg-cyan-500/30">
      {/* Sidebar Navigation */}
      <aside className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col z-10">
        <div className="h-16 flex items-center px-6 border-b border-slate-800">
          <div className="flex items-center gap-3 text-cyan-400">
            <Shield className="w-6 h-6" />
            <h1 className="text-xl font-bold tracking-wider">{t.app_name}</h1>
          </div>
        </div>

        <nav className="flex-1 px-4 py-6 space-y-2">
          <NavLink
            to="/"
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                isActive
                  ? 'bg-cyan-500/10 text-cyan-400'
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
              }`
            }
          >
            <ScanSearch className="w-5 h-5" />
            <span className="font-medium">{t.scanner}</span>
          </NavLink>

          <NavLink
            to="/history"
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                isActive
                  ? 'bg-cyan-500/10 text-cyan-400'
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
              }`
            }
          >
            <Activity className="w-5 h-5" />
            <span className="font-medium">{t.history}</span>
          </NavLink>
        </nav>

        <div className="px-4 py-6 border-t border-slate-800">
          <NavLink
            to="/settings"
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                isActive
                  ? 'bg-cyan-500/10 text-cyan-400'
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
              }`
            }
          >
            <SettingsIcon className="w-5 h-5" />
            <span className="font-medium">{t.settings}</span>
          </NavLink>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col relative">
        {/* Mock Window Title Bar (macOS/Win hybrid look) */}
        <header className="h-12 border-b border-slate-800/50 flex items-center justify-between px-4 bg-slate-900/50 backdrop-blur-sm select-none">
          <div className="text-xs text-slate-500 font-medium tracking-wide flex-1 text-center">
            {t.window_title}
          </div>
          <div className="flex items-center gap-4 text-slate-500">
            <Minus className="w-4 h-4 cursor-pointer hover:text-slate-300" />
            <Square className="w-3.5 h-3.5 cursor-pointer hover:text-slate-300" />
            <X className="w-4 h-4 cursor-pointer hover:text-red-400" />
          </div>
        </header>
        
        <div className="flex-1 overflow-auto">
          <Outlet />
        </div>
      </main>
      
      {/* Toast notifications */}
      <Toaster theme="dark" position="bottom-right" className="!bg-slate-900 !border-slate-800 !text-slate-200" />
    </div>
  );
}
