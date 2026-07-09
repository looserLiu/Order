import { useQuery } from '@tanstack/react-query'
import { reportApi, transactionApi, accountApi, budgetApi } from '../services/api'
import { ArrowTrendingUpIcon, ArrowTrendingDownIcon, WalletIcon, ChartBarIcon } from '@heroicons/react/24/solid'
import { PieChart, Pie, Cell, XAxis, YAxis, Tooltip, ResponsiveContainer, BarChart, Bar, LineChart, Line } from 'recharts'
import { Link } from 'react-router-dom'
import BudgetAlerts from '../components/BudgetAlerts'

const COLORS = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8', '#F7DC6F']

export default function Dashboard() {
  const { data: summary } = useQuery({
    queryKey: ['summary'],
    queryFn: () => reportApi.summary(),
  })

  const { data: trend } = useQuery({
    queryKey: ['trend'],
    queryFn: () => reportApi.trend(),
  })

  const { data: categoryData } = useQuery({
    queryKey: ['category'],
    queryFn: () => reportApi.byCategory({ type: 'expense' }),
  })

  const { data: accounts } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => accountApi.list(),
  })

  const { data: transactions } = useQuery({
    queryKey: ['transactions', 'recent'],
    queryFn: () => transactionApi.list({ page_size: 5 }),
  })

  const { data: budgets } = useQuery({
    queryKey: ['budgets'],
    queryFn: () => budgetApi.list(),
  })

  const { data: monthlyData } = useQuery({
    queryKey: ['monthly'],
    queryFn: () => reportApi.monthlyCompare(),
  })

  const summaryData = summary?.data?.data
  const trendData = trend?.data?.data || []
  const categoryChartData = categoryData?.data?.data || []
  const accountList = accounts?.data?.data || []
  const recentTransactions = transactions?.data?.data?.list || []
  const budgetList = budgets?.data?.data || []
  const monthlyCompareData = monthlyData?.data?.data || []

  const totalBalance = accountList.reduce((sum: number, acc: any) => sum + (acc.balance || 0), 0)
  const currentMonth = new Date().toISOString().slice(0, 7)
  
  const budgetProgress = budgetList.map((b: any) => {
    const spent = categoryChartData.find((c: any) => c.category_id === b.category_id)?.total || 0
    return { ...b, spent, progress: b.amount > 0 ? (spent / b.amount) * 100 : 0 }
  }).slice(0, 3)

  return (
    <div className="space-y-6">
      <BudgetAlerts />
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">本月收入</p>
              <p className="text-2xl font-bold text-green-600">
                ¥{summaryData?.income?.toFixed(2) || '0.00'}
              </p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
              <ArrowTrendingUpIcon className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>
        
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">本月支出</p>
              <p className="text-2xl font-bold text-red-500">
                ¥{summaryData?.expense?.toFixed(2) || '0.00'}
              </p>
            </div>
            <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
              <ArrowTrendingDownIcon className="w-6 h-6 text-red-500" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">账户余额</p>
              <p className="text-2xl font-bold text-primary-600">
                ¥{totalBalance.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-primary-100 rounded-full flex items-center justify-center">
              <WalletIcon className="w-6 h-6 text-primary-600" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">本月结余</p>
              <p className={`text-2xl font-bold ${(summaryData?.income - summaryData?.expense) >= 0 ? 'text-green-600' : 'text-red-500'}`}>
                ¥{((summaryData?.income || 0) - (summaryData?.expense || 0)).toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <ChartBarIcon className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </div>
      </div>

      {budgetProgress.length > 0 && (
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">预算进度</h3>
            <Link to="/budgets" className="text-sm text-primary-600 hover:underline">查看全部</Link>
          </div>
          <div className="space-y-4">
            {budgetProgress.map((b: any) => (
              <div key={b.id}>
                <div className="flex justify-between text-sm mb-1">
                  <span>{b.category?.name || '未分类'}</span>
                  <span className={b.progress > 100 ? 'text-red-500' : 'text-gray-500'}>
                    ¥{b.spent.toFixed(0)} / ¥{b.amount}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className={`h-2 rounded-full ${b.progress > 100 ? 'bg-red-500' : b.progress > 80 ? 'bg-yellow-500' : 'bg-green-500'}`}
                    style={{ width: `${Math.min(b.progress, 100)}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">收支趋势</h3>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={trendData.slice(-7)}>
              <XAxis dataKey="date" fontSize={12} />
              <YAxis fontSize={12} />
              <Tooltip />
              <Bar dataKey="income" fill="#27AE60" name="收入" radius={[4, 4, 0, 0]} />
              <Bar dataKey="expense" fill="#E74C3C" name="支出" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
        
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">支出分类</h3>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie
                data={categoryChartData}
                cx="50%"
                cy="50%"
                innerRadius={60}
                outerRadius={80}
                paddingAngle={5}
                dataKey="total"
              >
                {categoryChartData.map((_: any, index: number) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
          <div className="flex flex-wrap gap-2 mt-4">
            {categoryChartData.slice(0, 5).map((item: any, index: number) => (
              <div key={item.category_name} className="flex items-center gap-1 text-xs">
                <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[index % COLORS.length] }} />
                {item.category_name}
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">最近交易</h3>
          <div className="space-y-3">
            {recentTransactions.map((tx: any) => (
              <div key={tx.id} className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <div className="flex items-center gap-3">
                  <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                    tx.type === 'income' ? 'bg-green-100' : 'bg-red-100'
                  }`}>
                    <span className={tx.type === 'income' ? 'text-green-600' : 'text-red-500'}>
                      {tx.type === 'income' ? '+' : '-'}
                    </span>
                  </div>
                  <div>
                    <p className="font-medium">{tx.category?.name || '未分类'}</p>
                    <p className="text-sm text-gray-500">{tx.note || tx.merchant || '-'}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className={`font-semibold ${tx.type === 'income' ? 'text-green-600' : 'text-red-500'}`}>
                    {tx.type === 'income' ? '+' : '-'}¥{tx.amount?.toFixed(2)}
                  </p>
                  <p className="text-xs text-gray-400">{new Date(tx.bill_date).toLocaleDateString()}</p>
                </div>
              </div>
            ))}
            {recentTransactions.length === 0 && (
              <p className="text-center text-gray-500 py-4">暂无交易记录</p>
            )}
          </div>
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold mb-4">账户概览</h3>
          <div className="space-y-3">
            {accountList.slice(0, 5).map((acc: any) => (
              <div key={acc.id} className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full flex items-center justify-center bg-gray-100">
                    <WalletIcon className="w-5 h-5 text-gray-600" />
                  </div>
                  <div>
                    <p className="font-medium">{acc.name}</p>
                    <p className="text-xs text-gray-500">{acc.type}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-semibold">¥{acc.balance?.toFixed(2) || '0.00'}</p>
                </div>
              </div>
            ))}
            {accountList.length === 0 && (
              <p className="text-center text-gray-500 py-4">暂无账户</p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
