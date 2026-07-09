import { Outlet, NavLink } from 'react-router-dom'
import { useEffect } from 'react'
import { useThemeStore } from '../stores/themeStore'
import { 
  HomeIcon, 
  CreditCardIcon, 
  ArrowsRightLeftIcon, 
  ChartPieIcon, 
  Cog6ToothIcon, 
  TagIcon,
  CurrencyDollarIcon,
  BookmarkIcon,
  UserGroupIcon,
  BellIcon,
  CalendarIcon,
  MagnifyingGlassIcon,
  TrophyIcon,
  ShieldCheckIcon,
  BanknotesIcon,
  UsersIcon,
  SunIcon,
  MoonIcon,
  ArrowPathIcon,
  DocumentArrowDownIcon,
  ChartBarIcon
} from '@heroicons/react/24/outline'
import clsx from 'clsx'
import QuickAddButton from './QuickAddButton'

const navItems = [
  { name: '首页', path: '/', icon: HomeIcon },
  { name: '账户', path: '/accounts', icon: CreditCardIcon },
  { name: '记账', path: '/transactions', icon: ArrowsRightLeftIcon },
  { name: '周期', path: '/recurring', icon: ArrowPathIcon },
  { name: '日历', path: '/calendar', icon: CalendarIcon },
  { name: '预算', path: '/budgets', icon: CurrencyDollarIcon },
  { name: '标签', path: '/tags', icon: BookmarkIcon },
  { name: 'CSV导入', path: '/csv-import', icon: DocumentArrowDownIcon },
  { name: '现金流', path: '/cashflow', icon: ChartBarIcon },
  { name: '家庭', path: '/families', icon: UserGroupIcon },
  { name: 'AA记账', path: '/aa-groups', icon: UsersIcon },
  { name: '报表', path: '/reports', icon: ChartPieIcon },
  { name: '搜索', path: '/search', icon: MagnifyingGlassIcon },
  { name: '目标', path: '/goals', icon: TrophyIcon },
  { name: '保险', path: '/insurances', icon: ShieldCheckIcon },
  { name: '净资产', path: '/net-worth', icon: BanknotesIcon },
  { name: '借出借入', path: '/debts', icon: BanknotesIcon },
  { name: '提醒', path: '/reminders', icon: BellIcon },
  { name: '消息', path: '/notifications', icon: BellIcon },
  { name: '设置', path: '/settings', icon: Cog6ToothIcon },
]

export default function Layout() {
  const { theme, toggleTheme } = useThemeStore()

  useEffect(() => {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [theme])

  return (
    <div className="min-h-screen flex">
      <aside className="w-64 bg-white border-r border-gray-200 hidden lg:block fixed h-full overflow-y-auto dark:bg-gray-800 dark:border-gray-700">
        <div className="p-6">
          <h1 className="text-2xl font-bold text-primary-600">智慧记账</h1>
        </div>
        <nav className="px-3 pb-20 space-y-1">
          {navItems.slice(0, 7).map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              end={item.path === '/'}
              className={({ isActive }) =>
                clsx(
                  'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                  isActive
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-600 hover:bg-gray-50'
                )
              }
            >
              <item.icon className="w-5 h-5" />
              {item.name}
            </NavLink>
          ))}
          
          <div className="pt-4 mt-4 border-t">
            <p className="px-3 text-xs text-gray-400 mb-2">更多</p>
            {navItems.slice(7).map((item) => (
              <NavLink
                key={item.path}
                to={item.path}
                className={({ isActive }) =>
                  clsx(
                    'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-primary-50 text-primary-600'
                      : 'text-gray-600 hover:bg-gray-50'
                  )
                }
              >
                <item.icon className="w-5 h-5" />
                {item.name}
              </NavLink>
            ))}
          </div>
        </nav>
      </aside>
      
      <div className="flex-1 lg:ml-64">
        <header className="bg-white border-b border-gray-200 px-4 md:px-6 py-4 sticky top-0 z-30">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">我的账本</h2>
            <div className="flex items-center gap-2">
              <button
                onClick={toggleTheme}
                className="p-2 hover:bg-gray-100 rounded-lg dark:hover:bg-gray-700"
              >
                {theme === 'light' ? (
                  <MoonIcon className="w-5 h-5 text-gray-600 dark:text-gray-300" />
                ) : (
                  <SunIcon className="w-5 h-5 text-gray-600 dark:text-gray-300" />
                )}
              </button>
              <NavLink
                to="/search"
                className="p-2 hover:bg-gray-100 rounded-lg dark:hover:bg-gray-700"
              >
                <MagnifyingGlassIcon className="w-5 h-5 text-gray-600 dark:text-gray-300" />
              </NavLink>
              <NavLink
                to="/notifications"
                className="p-2 hover:bg-gray-100 rounded-lg relative dark:hover:bg-gray-700"
              >
                <BellIcon className="w-5 h-5 text-gray-600 dark:text-gray-300" />
              </NavLink>
            </div>
          </div>
        </header>
        
        <main className="p-4 md:p-6 bg-gray-50 min-h-[calc(100vh-73px)] pb-24 lg:pb-6 dark:bg-gray-900">
          <Outlet />
        </main>
      </div>
      
      <QuickAddButton />
      
      <nav className="lg:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 px-2 py-2 z-30">
        <div className="flex justify-around">
          {[
            { name: '首页', path: '/', icon: HomeIcon },
            { name: '账户', path: '/accounts', icon: CreditCardIcon },
            { name: '记账', path: '/transactions', icon: ArrowsRightLeftIcon },
            { name: '日历', path: '/calendar', icon: CalendarIcon },
            { name: '更多', path: '/reports', icon: ChartPieIcon },
          ].map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              end={item.path === '/'}
              className={({ isActive }) =>
                clsx(
                  'flex flex-col items-center gap-0.5 text-xs py-1',
                  isActive ? 'text-primary-600' : 'text-gray-500'
                )
              }
            >
              <item.icon className="w-5 h-5" />
              {item.name}
            </NavLink>
          ))}
        </div>
      </nav>
    </div>
  )
}
