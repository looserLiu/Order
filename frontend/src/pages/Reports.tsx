import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { reportApi, accountApi } from '../services/api'
import { PieChart, Pie, Cell, XAxis, YAxis, Tooltip, ResponsiveContainer, BarChart, Bar, LineChart, Line, Legend } from 'recharts'

const COLORS = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8', '#F7DC6F', '#F8B500', '#00D4AA']

export default function Reports() {
  const [period, setPeriod] = useState('month')
  const [activeTab, setActiveTab] = useState('overview')

  const getDateRange = () => {
    const end = new Date()
    const start = new Date()
    if (period === 'week') {
      start.setDate(end.getDate() - 7)
    } else if (period === 'month') {
      start.setMonth(end.getMonth() - 1)
    } else {
      start.setFullYear(end.getFullYear() - 1)
    }
    return {
      start_date: start.toISOString().split('T')[0],
      end_date: end.toISOString().split('T')[0],
    }
  }

  const params = getDateRange()

  const { data: summary } = useQuery({
    queryKey: ['summary', params],
    queryFn: () => reportApi.summary(params),
  })

  const { data: trend } = useQuery({
    queryKey: ['trend', params],
    queryFn: () => reportApi.trend(params),
  })

  const { data: expenseCategory } = useQuery({
    queryKey: ['category-expense', params],
    queryFn: () => reportApi.byCategory({ ...params, type: 'expense' }),
  })

  const { data: incomeCategory } = useQuery({
    queryKey: ['category-income', params],
    queryFn: () => reportApi.byCategory({ ...params, type: 'income' }),
  })

  const { data: accountData } = useQuery({
    queryKey: ['report-account', params],
    queryFn: () => reportApi.byAccount(params),
  })

  const { data: merchantData } = useQuery({
    queryKey: ['report-merchant', params],
    queryFn: () => reportApi.byMerchant(params),
  })

  const { data: monthly } = useQuery({
    queryKey: ['monthly', params],
    queryFn: () => reportApi.monthlyCompare(),
  })

  const summaryData = summary?.data?.data
  const trendData = trend?.data?.data || []
  const expenseData = expenseCategory?.data?.data || []
  const incomeData = incomeCategory?.data?.data || []
  const accountReportData = accountData?.data?.data || []
  const merchantReportData = merchantData?.data?.data || []
  const monthlyData = monthly?.data?.data || []

  const totalExpense = expenseData.reduce((sum: number, item: any) => sum + item.total, 0)
  const totalIncome = incomeData.reduce((sum: number, item: any) => sum + item.total, 0)

  const renderOverview = () => (
    <>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="card text-center">
          <p className="text-gray-500 mb-2">收入</p>
          <p className="text-3xl font-bold text-green-600">¥{summaryData?.income?.toFixed(2) || '0.00'}</p>
        </div>
        <div className="card text-center">
          <p className="text-gray-500 mb-2">支出</p>
          <p className="text-3xl font-bold text-red-500">¥{summaryData?.expense?.toFixed(2) || '0.00'}</p>
        </div>
        <div className="card text-center">
          <p className="text-gray-500 mb-2">结余</p>
          <p className="text-3xl font-bold text-primary-600">¥{((summaryData?.income || 0) - (summaryData?.expense || 0)).toFixed(2)}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">收支趋势</h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={trendData}>
              <XAxis dataKey="date" fontSize={12} />
              <YAxis fontSize={12} />
              <Tooltip />
              <Legend />
              <Line type="monotone" dataKey="income" stroke="#27AE60" strokeWidth={2} name="收入" dot={false} />
              <Line type="monotone" dataKey="expense" stroke="#E74C3C" strokeWidth={2} name="支出" dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold mb-4">支出分类</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={expenseData}
                cx="50%"
                cy="50%"
                outerRadius={100}
                dataKey="total"
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
              >
                {expenseData.map((_: any, index: number) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">支出明细</h3>
        <table className="w-full">
          <thead>
            <tr className="text-left text-sm text-gray-500 border-b">
              <th className="pb-3">分类</th>
              <th className="pb-3 text-right">金额</th>
              <th className="pb-3 text-right">占比</th>
              <th className="pb-3 w-24">占比</th>
            </tr>
          </thead>
          <tbody>
            {expenseData.map((item: any, index: number) => (
              <tr key={item.category_name} className="border-b border-gray-100">
                <td className="py-3">
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[index % COLORS.length] }} />
                    {item.category_name}
                  </div>
                </td>
                <td className="py-3 text-right">¥{item.total?.toFixed(2)}</td>
                <td className="py-3 text-right">{((item.total / totalExpense) * 100).toFixed(1)}%</td>
                <td className="py-3">
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div className="h-2 rounded-full" style={{ width: `${(item.total / totalExpense) * 100}%`, backgroundColor: COLORS[index % COLORS.length] }} />
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  )

  const renderIncome = () => (
    <div className="space-y-6">
      <div className="card">
        <h3 className="text-lg font-semibold mb-4">收入分类</h3>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={incomeData}
              cx="50%"
              cy="50%"
              outerRadius={100}
              dataKey="total"
              label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
            >
              {incomeData.map((_: any, index: number) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip />
          </PieChart>
        </ResponsiveContainer>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">收入明细</h3>
        <table className="w-full">
          <thead>
            <tr className="text-left text-sm text-gray-500 border-b">
              <th className="pb-3">分类</th>
              <th className="pb-3 text-right">金额</th>
              <th className="pb-3 text-right">占比</th>
            </tr>
          </thead>
          <tbody>
            {incomeData.map((item: any, index: number) => (
              <tr key={item.category_name} className="border-b border-gray-100">
                <td className="py-3">
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[index % COLORS.length] }} />
                    {item.category_name}
                  </div>
                </td>
                <td className="py-3 text-right text-green-600">¥{item.total?.toFixed(2)}</td>
                <td className="py-3 text-right">{totalIncome > 0 ? ((item.total / totalIncome) * 100).toFixed(1) : 0}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )

  const renderAccounts = () => (
    <div className="card">
      <h3 className="text-lg font-semibold mb-4">账户收支</h3>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={accountReportData} layout="vertical">
          <XAxis type="number" fontSize={12} />
          <YAxis dataKey="account_name" type="category" fontSize={12} width={80} />
          <Tooltip />
          <Legend />
          <Bar dataKey="income" fill="#27AE60" name="收入" />
          <Bar dataKey="expense" fill="#E74C3C" name="支出" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )

  const renderMerchant = () => (
    <div className="card">
      <h3 className="text-lg font-semibold mb-4">商家消费排名</h3>
      <table className="w-full">
        <thead>
          <tr className="text-left text-sm text-gray-500 border-b">
            <th className="pb-3">排名</th>
            <th className="pb-3">商家</th>
            <th className="pb-3 text-right">消费金额</th>
            <th className="pb-3 text-right">占比</th>
          </tr>
        </thead>
        <tbody>
          {merchantReportData.slice(0, 10).map((item: any, index: number) => (
            <tr key={item.merchant} className="border-b border-gray-100">
              <td className="py-3">
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                  index === 0 ? 'bg-yellow-100 text-yellow-700' :
                  index === 1 ? 'bg-gray-100 text-gray-700' :
                  index === 2 ? 'bg-orange-100 text-orange-700' :
                  'bg-gray-50 text-gray-500'
                }`}>
                  #{index + 1}
                </span>
              </td>
              <td className="py-3">{item.merchant || '-'}</td>
              <td className="py-3 text-right text-red-500">¥{item.total?.toFixed(2)}</td>
              <td className="py-3 text-right">{item.percentage?.toFixed(1)}%</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )

  const renderMonthly = () => (
    <div className="space-y-6">
      <div className="card">
        <h3 className="text-lg font-semibold mb-4">月度对比</h3>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={monthlyData}>
            <XAxis dataKey="month" fontSize={12} />
            <YAxis fontSize={12} />
            <Tooltip />
            <Legend />
            <Bar dataKey="income" fill="#27AE60" name="收入" radius={[4, 4, 0, 0]} />
            <Bar dataKey="expense" fill="#E74C3C" name="支出" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">月度明细</h3>
        <table className="w-full">
          <thead>
            <tr className="text-left text-sm text-gray-500 border-b">
              <th className="pb-3">月份</th>
              <th className="pb-3 text-right">收入</th>
              <th className="pb-3 text-right">支出</th>
              <th className="pb-3 text-right">结余</th>
            </tr>
          </thead>
          <tbody>
            {monthlyData.map((item: any) => (
              <tr key={item.month} className="border-b border-gray-100">
                <td className="py-3">{item.month}</td>
                <td className="py-3 text-right text-green-600">¥{item.income?.toFixed(2)}</td>
                <td className="py-3 text-right text-red-500">¥{item.expense?.toFixed(2)}</td>
                <td className={`py-3 text-right font-medium ${item.income - item.expense >= 0 ? 'text-green-600' : 'text-red-500'}`}>
                  ¥{(item.income - item.expense).toFixed(2)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <h2 className="text-2xl font-bold">财务报表</h2>
        <div className="flex gap-2">
          {['week', 'month', 'year'].map(p => (
            <button
              key={p}
              onClick={() => setPeriod(p)}
              className={`px-4 py-2 rounded-lg text-sm font-medium ${
                period === p ? 'bg-primary-600 text-white' : 'bg-white text-gray-600 hover:bg-gray-50'
              }`}
            >
              {p === 'week' ? '本周' : p === 'month' ? '本月' : '本年'}
            </button>
          ))}
        </div>
      </div>

      <div className="flex gap-2 overflow-x-auto pb-2">
        {[
          { id: 'overview', label: '概览' },
          { id: 'income', label: '收入' },
          { id: 'accounts', label: '账户' },
          { id: 'merchant', label: '商家' },
          { id: 'monthly', label: '月度' },
        ].map(tab => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap ${
              activeTab === tab.id 
                ? 'bg-primary-600 text-white' 
                : 'bg-white text-gray-600 hover:bg-gray-50'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'overview' && renderOverview()}
      {activeTab === 'income' && renderIncome()}
      {activeTab === 'accounts' && renderAccounts()}
      {activeTab === 'merchant' && renderMerchant()}
      {activeTab === 'monthly' && renderMonthly()}
    </div>
  )
}
