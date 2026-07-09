import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import { useThemeStore } from './stores/themeStore'
import Layout from './components/Layout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import Accounts from './pages/Accounts'
import Transactions from './pages/Transactions'
import Categories from './pages/Categories'
import Reports from './pages/Reports'
import Budgets from './pages/Budgets'
import Tags from './pages/Tags'
import Families from './pages/Families'
import Settings from './pages/Settings'
import Calendar from './pages/Calendar'
import Search from './pages/Search'
import Notifications from './pages/Notifications'
import Reminders from './pages/Reminders'
import AAGroups from './pages/AAGroups'
import Goals from './pages/Goals'
import Inurances from './pages/Inurances'
import NetWorth from './pages/NetWorth'
import Recurring from './pages/Recurring'
import Debts from './pages/Debts'
import CSVImport from './pages/CSVImport'
import CashFlow from './pages/CashFlow'

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
  )
}
