import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { cashFlowApi, accountApi } from '../services/api'
import { 
  AreaChart, 
  Area, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  BarChart,
  Bar
} from 'recharts'
import { ArrowTrendingUpIcon, ArrowTrendingDownIcon, CurrencyDollarIcon, CalendarIcon } from '@heroicons/react/24/outline'

export default function CashFlow() {
  const [days, setDays] = useState(30)

  const { data: balanceData } = useQuery({
    queryKey: ['totalBalance'],
    queryFn: async () => {
      const res = await accountApi.getTotalBalance()
      return res.data.data
    },
  })

  const { data: projectionData, isLoading } = useQuery({
    queryKey: ['cashflow', days],
    queryFn: async () => {
      const res = await cashFlowApi.getProjection(days)
      return res.data.data
    },
  })

  const currentBalance = balanceData?.total || 0
  const projections = projectionData?.projections || []
  const endBalance = projections.length > 0 ? projections[projections.length - 1].projected_bal : currentBalance

  const totalIncome = projections.reduce((sum: number, p: any) => sum + p.income, 0)
  const totalExpense = projections.reduce((sum: number, p: any) => sum + p.expense, 0)

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('zh-CN', {
      style: 'currency',
      currency: 'CNY',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">现金流预测</h1>
        <div className="flex items-center gap-2">
          <CalendarIcon className="w-5 h-5 text-gray-500" />
          <select
            value={days}
            onChange={(e) => setDays(Number(e.target.value))}
            className="px-3 py-2 border border-gray-300 rounded-lg dark:border-gray-600 dark:bg-gray-700"
          >
            <option value={7}>7天</option>
            <option value={14}>14天</option>
            <option value={30}>30天</option>
            <option value={60}>60天</option>
            <option value={90}>90天</option>
          </select>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
              <CurrencyDollarIcon className="w-6 h-6 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">当前余额</p>
              <p className="text-xl font-bold text-gray-900 dark:text-white">{formatCurrency(currentBalance)}</p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
              <ArrowTrendingUpIcon className="w-6 h-6 text-green-600 dark:text-green-400" />
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">预计收入</p>
              <p className="text-xl font-bold text-green-600 dark:text-green-400">+{formatCurrency(totalIncome)}</p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
              <ArrowTrendingDownIcon className="w-6 h-6 text-red-600 dark:text-red-400" />
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">预计支出</p>
              <p className="text-xl font-bold text-red-600 dark:text-red-400">-{formatCurrency(totalExpense)}</p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3">
            <div className={`p-3 rounded-lg ${endBalance >= currentBalance ? 'bg-green-100 dark:bg-green-900' : 'bg-red-100 dark:bg-red-900'}`}>
              <ArrowTrendingUpIcon className={`w-6 h-6 ${endBalance >= currentBalance ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`} />
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">预测余额</p>
              <p className={`text-xl font-bold ${endBalance >= currentBalance ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                {formatCurrency(endBalance)}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Balance Projection Chart */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">余额趋势预测</h2>
        {isLoading ? (
          <div className="h-80 flex items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          </div>
        ) : (
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart data={projections}>
              <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
              <XAxis 
                dataKey="date" 
                tick={{ fill: '#6B7280', fontSize: 12 }}
                tickFormatter={(value) => value.slice(5)}
              />
              <YAxis 
                tick={{ fill: '#6B7280', fontSize: 12 }}
                tickFormatter={(value) => `${(value / 1000).toFixed(0)}k`}
              />
              <Tooltip 
                contentStyle={{ 
                  backgroundColor: '#1F2937', 
                  border: 'none', 
                  borderRadius: '8px',
                  color: '#fff'
                }}
                formatter={(value: number) => [formatCurrency(value), '余额']}
                labelFormatter={(label) => `日期: ${label}`}
              />
              <Area 
                type="monotone" 
                dataKey="projected_bal" 
                stroke="#3B82F6" 
                fill="#3B82F6" 
                fillOpacity={0.3}
                name="预测余额"
              />
            </AreaChart>
          </ResponsiveContainer>
        )}
      </div>

      {/* Income/Expense Chart */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">每日收支</h2>
        {isLoading ? (
          <div className="h-80 flex items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          </div>
        ) : (
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={projections}>
              <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
              <XAxis 
                dataKey="date" 
                tick={{ fill: '#6B7280', fontSize: 12 }}
                tickFormatter={(value) => value.slice(5)}
              />
              <YAxis tick={{ fill: '#6B7280', fontSize: 12 }} />
              <Tooltip 
                contentStyle={{ 
                  backgroundColor: '#1F2937', 
                  border: 'none', 
                  borderRadius: '8px',
                  color: '#fff'
                }}
                formatter={(value: number, name: string) => [
                  formatCurrency(value),
                  name === 'income' ? '收入' : '支出'
                ]}
              />
              <Bar dataKey="income" fill="#22C55E" name="收入" />
              <Bar dataKey="expense" fill="#EF4444" name="支出" />
            </BarChart>
          </ResponsiveContainer>
        )}
      </div>

      {/* Projections Table */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">详细预测</h2>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-500 dark:text-gray-400">日期</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500 dark:text-gray-400">收入</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500 dark:text-gray-400">支出</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500 dark:text-gray-400">预测余额</th>
              </tr>
            </thead>
            <tbody>
              {projections.slice(0, 14).map((day: any, index: number) => (
                <tr key={index} className="border-b border-gray-100 dark:border-gray-700">
                  <td className="py-3 px-4 text-sm text-gray-900 dark:text-white">{day.date}</td>
                  <td className="py-3 px-4 text-sm text-right text-green-600 dark:text-green-400">
                    {day.income > 0 ? formatCurrency(day.income) : '-'}
                  </td>
                  <td className="py-3 px-4 text-sm text-right text-red-600 dark:text-red-400">
                    {day.expense > 0 ? formatCurrency(day.expense) : '-'}
                  </td>
                  <td className="py-3 px-4 text-sm text-right font-medium text-gray-900 dark:text-white">
                    {formatCurrency(day.projected_bal)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
