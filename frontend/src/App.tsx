import { useEffect, Suspense, lazy } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import { useThemeStore } from './stores/themeStore'
import Layout from './components/Layout'

// Lazy load pages for better performance
const Login = lazy(() => import('./pages/Login'))
const Dashboard = lazy(() => import('./pages/Dashboard'))
const Accounts = lazy(() => import('./pages/Accounts'))
const Transactions = lazy(() => import('./pages/Transactions'))
const Categories = lazy(() => import('./pages/Categories'))
const Reports = lazy(() => import('./pages/Reports'))
const Budgets = lazy(() => import('./pages/Budgets'))
const Tags = lazy(() => import('./pages/Tags'))
const Families = lazy(() => import('./pages/Families'))
const Settings = lazy(() => import('./pages/Settings'))
const Calendar = lazy(() => import('./pages/Calendar'))
const Search = lazy(() => import('./pages/Search'))
const Notifications = lazy(() => import('./pages/Notifications'))
const Reminders = lazy(() => import('./pages/Reminders'))
const AAGroups = lazy(() => import('./pages/AAGroups'))
const Goals = lazy(() => import('./pages/Goals'))
const Inurances = lazy(() => import('./pages/Inurances'))
const NetWorth = lazy(() => import('./pages/NetWorth'))
const Recurring = lazy(() => import('./pages/Recurring'))
const Debts = lazy(() => import('./pages/Debts'))
const CSVImport = lazy(() => import('./pages/CSVImport'))
const CashFlow = lazy(() => import('./pages/CashFlow'))

// Loading component for Suspense fallback
function LoadingSpinner() {
  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>
  )
}

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { token } = useAuthStore()
  return token ? <>{children}</> : <Navigate to="/login" />
}

export default function App() {
  const { theme, colorScheme } = useThemeStore()

  // Apply theme on mount and change
  useEffect(() => {
    if (theme === 'dark') {
      document.documentElement.setAttribute('data-theme', 'dark')
    } else {
      document.documentElement.removeAttribute('data-theme')
    }
    document.documentElement.setAttribute('data-color-scheme', colorScheme)
  }, [theme, colorScheme])

  return (
    <Suspense fallback={<LoadingSpinner />}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={
          <PrivateRoute>
            <Layout />
          </PrivateRoute>
        }>
          <Route index element={<Dashboard />} />
          <Route path="accounts" element={<Accounts />} />
          <Route path="transactions" element={<Transactions />} />
          <Route path="categories" element={<Categories />} />
          <Route path="budgets" element={<Budgets />} />
          <Route path="tags" element={<Tags />} />
          <Route path="families" element={<Families />} />
          <Route path="reports" element={<Reports />} />
          <Route path="calendar" element={<Calendar />} />
          <Route path="search" element={<Search />} />
          <Route path="notifications" element={<Notifications />} />
          <Route path="reminders" element={<Reminders />} />
          <Route path="aa-groups" element={<AAGroups />} />
          <Route path="goals" element={<Goals />} />
          <Route path="insurances" element={<Inurances />} />
          <Route path="net-worth" element={<NetWorth />} />
          <Route path="recurring" element={<Recurring />} />
          <Route path="debts" element={<Debts />} />
          <Route path="csv-import" element={<CSVImport />} />
          <Route path="cashflow" element={<CashFlow />} />
          <Route path="settings" element={<Settings />} />
        </Route>
      </Routes>
    </Suspense>
  )
}
