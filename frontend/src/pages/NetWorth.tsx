import { useQuery } from '@tanstack/react-query'
import { netWorthApi, assetApi, accountApi, Account, AssetChange } from '../services/api'
import {
  BanknotesIcon,
  CreditCardIcon,
  ChartBarIcon,
  ArrowTrendingUpIcon,
  WalletIcon
} from '@heroicons/react/24/solid'
import { PieChart, Pie, Cell, ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip } from 'recharts'

const COLORS = ['#27AE60', '#E74C3C', '#3498DB', '#F39C12', '#9B59B6']

// NetWorth response type
interface NetWorthResponse {
  net_worth: number
  total_assets: number
  total_investment: number
  total_debt_owed: number
  total_debt_owing: number
}

export default function NetWorth() {
  const { data: netWorthData } = useQuery({
    queryKey: ['net-worth'],
    queryFn: () => netWorthApi.get(),
  })

  const { data: accountsData } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => accountApi.list(),
  })

  const { data: assetsData } = useQuery({
    queryKey: ['assets'],
    queryFn: () => assetApi.list({ type: 'investment' }),
  })

  const netWorth = (netWorthData?.data?.data || {}) as NetWorthResponse
  const accounts = (accountsData?.data?.data || []) as Account[]
  const investments = (assetsData?.data?.data || []) as AssetChange[]

  const assetData = [
    { name: '账户余额', value: netWorth.total_assets || 0 },
    { name: '借出款', value: netWorth.total_debt_owed || 0 },
    { name: '投资资产', value: netWorth.total_investment || 0 },
  ].filter(d => d.value > 0)

  const liabilityData = [
    { name: '借入款', value: netWorth.total_debt_owing || 0 },
  ].filter(d => d.value > 0)

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2">
        <ChartBarIcon className="w-6 h-6" />
        <h2 className="text-2xl font-bold">净资产</h2>
      </div>

      <div className="card bg-gradient-to-r from-primary-600 to-primary-700 text-white">
        <p className="text-sm opacity-80">净资产总额</p>
        <p className="text-4xl font-bold mt-2">
          ¥{(netWorth.net_worth || 0).toFixed(2)}
        </p>
        <p className="text-sm opacity-80 mt-2">
          计算公式: 资产 + 借出 - 借入
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
              <WalletIcon className="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">总资产</p>
              <p className="text-xl font-bold text-green-600">
                ¥{((netWorth.total_assets || 0) + (netWorth.total_investment || 0)).toFixed(2)}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
              <BanknotesIcon className="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">账户余额</p>
              <p className="text-xl font-bold text-blue-600">
                ¥{(netWorth.total_assets || 0).toFixed(2)}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center">
              <ArrowTrendingUpIcon className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">投资资产</p>
              <p className="text-xl font-bold text-purple-600">
                ¥{(netWorth.total_investment || 0).toFixed(2)}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center">
              <CreditCardIcon className="w-6 h-6 text-red-500" />
            </div>
            <div>
              <p className="text-sm text-gray-500">总负债</p>
              <p className="text-xl font-bold text-red-500">
                ¥{(netWorth.total_debt_owing || 0).toFixed(2)}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="font-semibold mb-4">资产构成</h3>
          {assetData.length > 0 ? (
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={assetData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={80}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {assetData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => `¥${Number(value).toFixed(2)}`} />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <div className="text-center py-8 text-gray-500">暂无资产数据</div>
          )}
          <div className="flex flex-wrap gap-2 mt-4">
            {assetData.map((item, index) => (
              <div key={item.name} className="flex items-center gap-1 text-xs">
                <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[index % COLORS.length] }} />
                {item.name}
              </div>
            ))}
          </div>
        </div>

        <div className="card">
          <h3 className="font-semibold mb-4">账户明细</h3>
          <div className="space-y-3">
            {accounts.slice(0, 6).map((acc: Account) => (
              <div key={acc.id} className="flex justify-between items-center">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full" style={{ backgroundColor: acc.color || '#999' }} />
                  <span>{acc.name}</span>
                </div>
                <span className="font-medium">¥{acc.balance?.toFixed(2)}</span>
              </div>
            ))}
            {accounts.length === 0 && (
              <p className="text-center py-4 text-gray-500">暂无账户</p>
            )}
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="font-semibold mb-4">净资产趋势</h3>
        <div className="grid grid-cols-2 gap-4">
          <div className="p-4 bg-green-50 rounded-lg">
            <p className="text-sm text-green-600">资产总计</p>
            <p className="text-xl font-bold text-green-600">
              ¥{((netWorth.total_assets || 0) + (netWorth.total_investment || 0) + (netWorth.total_debt_owed || 0)).toFixed(2)}
            </p>
          </div>
          <div className="p-4 bg-red-50 rounded-lg">
            <p className="text-sm text-red-600">负债总计</p>
            <p className="text-xl font-bold text-red-600">
              ¥{(netWorth.total_debt_owing || 0).toFixed(2)}
            </p>
          </div>
        </div>
        <div className="mt-4 p-4 bg-primary-50 rounded-lg">
          <p className="text-sm text-primary-600">净资产 = 资产 - 负债</p>
          <p className="text-xl font-bold text-primary-600">
            ¥{(netWorth.net_worth || 0).toFixed(2)}
          </p>
        </div>
      </div>

      {investments.length > 0 && (
        <div className="card">
          <h3 className="font-semibold mb-4">投资账户</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {investments.map((inv: AssetChange) => (
              <div key={inv.id} className="p-4 border rounded-lg">
                <p className="font-medium">{inv.name}</p>
                <p className="text-xl font-bold text-purple-600 mt-1">¥{inv.amount?.toFixed(2)}</p>
                {inv.interest_rate && inv.interest_rate > 0 && (
                  <p className="text-xs text-gray-500 mt-1">利率: {inv.interest_rate}%</p>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
