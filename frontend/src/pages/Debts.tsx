import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { assetApi } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, CurrencyDollarIcon, ArrowUpIcon, ArrowDownIcon } from '@heroicons/react/24/outline'

export default function Debts() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [filter, setFilter] = useState('all')
  const [form, setForm] = useState({
    asset_type: 'debt_owed',
    name: '',
    related_user: '',
    amount: 0,
    interest_rate: 0,
    start_date: new Date().toISOString().split('T')[0],
    end_date: '',
    note: '',
    status: 'active',
  })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['assets', filter],
    queryFn: () => assetApi.list({ asset_type: filter === 'all' ? undefined : filter }),
  })

  const { data: summary } = useQuery({
    queryKey: ['assets', 'summary'],
    queryFn: () => assetApi.getSummary(),
  })

  const createMutation = useMutation({
    mutationFn: (data: any) => assetApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assets'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => assetApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assets'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => assetApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assets'] })
    },
  })

  const assets = data?.data?.data || []
  const summaryData = summary?.data?.data || {}

  const openModal = (asset?: any) => {
    if (asset) {
      setEditingId(asset.id)
      setForm({
        asset_type: asset.asset_type,
        name: asset.name,
        related_user: asset.related_user || '',
        amount: asset.amount,
        interest_rate: asset.interest_rate || 0,
        start_date: asset.start_date?.split('T')[0] || new Date().toISOString().split('T')[0],
        end_date: asset.end_date?.split('T')[0] || '',
        note: asset.note || '',
        status: asset.status || 'active',
      })
    } else {
      setEditingId(null)
      setForm({
        asset_type: 'debt_owed',
        name: '',
        related_user: '',
        amount: 0,
        interest_rate: 0,
        start_date: new Date().toISOString().split('T')[0],
        end_date: '',
        note: '',
        status: 'active',
      })
    }
    setShowModal(true)
  }

  const closeModal = () => {
    setShowModal(false)
    setEditingId(null)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: form })
    } else {
      createMutation.mutate(form)
    }
  }

  const settleMutation = useMutation({
    mutationFn: ({ id, amount }: { id: string; amount: number }) => 
      assetApi.update(id, { status: 'settled', note: `已还清 ¥${amount}` }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assets'] })
    },
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">借出借入</h2>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          添加
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
              <ArrowUpIcon className="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">借出(别人欠我)</p>
              <p className="text-xl font-bold text-green-600">¥{summaryData.total_debt_owed?.toFixed(2) || '0.00'}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
              <ArrowDownIcon className="w-6 h-6 text-red-500" />
            </div>
            <div>
              <p className="text-sm text-gray-500">借入(我欠别人)</p>
              <p className="text-xl font-bold text-red-500">¥{summaryData.total_debt_owing?.toFixed(2) || '0.00'}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <CurrencyDollarIcon className="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">投资</p>
              <p className="text-xl font-bold text-blue-600">¥{summaryData.total_investment?.toFixed(2) || '0.00'}</p>
            </div>
          </div>
        </div>
      </div>

      <div className="flex gap-2 mb-4">
        {['all', 'debt_owed', 'debt_owing', 'investment'].map(f => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-4 py-2 rounded-lg text-sm ${
              filter === f 
                ? 'bg-primary-600 text-white' 
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            {f === 'all' ? '全部' : f === 'debt_owed' ? '借出' : f === 'debt_owing' ? '借入' : '投资'}
          </button>
        ))}
      </div>

      <div className="card">
        <div className="space-y-3">
          {assets.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <CurrencyDollarIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
              <p>暂无记录</p>
            </div>
          ) : (
            assets.map((asset: any) => (
              <div key={asset.id} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50">
                <div className="flex items-center gap-4">
                  <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                    asset.asset_type === 'debt_owed' ? 'bg-green-100' : 
                    asset.asset_type === 'debt_owing' ? 'bg-red-100' : 'bg-blue-100'
                  }`}>
                    {asset.asset_type === 'debt_owed' ? (
                      <ArrowUpIcon className="w-6 h-6 text-green-600" />
                    ) : asset.asset_type === 'debt_owing' ? (
                      <ArrowDownIcon className="w-6 h-6 text-red-500" />
                    ) : (
                      <CurrencyDollarIcon className="w-6 h-6 text-blue-600" />
                    )}
                  </div>
                  <div>
                    <p className="font-medium">{asset.name}</p>
                    <p className="text-sm text-gray-500">
                      {asset.related_user || '-'} · {asset.start_date ? new Date(asset.start_date).toLocaleDateString() : '-'}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-4">
                  <div className="text-right">
                    <p className={`font-semibold text-lg ${
                      asset.asset_type === 'debt_owed' ? 'text-green-600' : 
                      asset.asset_type === 'debt_owing' ? 'text-red-500' : 'text-blue-600'
                    }`}>
                      ¥{asset.amount?.toFixed(2)}
                    </p>
                    <p className={`text-xs ${asset.status === 'active' ? 'text-green-500' : 'text-gray-400'}`}>
                      {asset.status === 'active' ? '进行中' : '已还清'}
                    </p>
                  </div>
                  {asset.status === 'active' && (
                    <button
                      onClick={() => settleMutation.mutate({ id: asset.id, amount: asset.amount })}
                      className="p-2 text-green-600 hover:bg-green-50 rounded-lg"
                      title="标记为已还清"
                    >
                      ✓
                    </button>
                  )}
                  <button
                    onClick={() => openModal(asset)}
                    className="p-2 text-gray-400 hover:text-gray-600"
                  >
                    <PencilIcon className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => deleteMutation.mutate(asset.id)}
                    className="p-2 text-gray-400 hover:text-red-500"
                  >
                    <TrashIcon className="w-5 h-5" />
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-end md:items-center justify-center z-50 p-4">
          <div className="bg-white rounded-t-2xl md:rounded-xl w-full md:max-w-md p-4 md:p-6 animate-slide-up md:animate-none">
            <h3 className="text-lg font-semibold mb-4">
              {editingId ? '编辑记录' : '添加记录'}
            </h3>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">类型</label>
                <select
                  value={form.asset_type}
                  onChange={(e) => setForm({ ...form, asset_type: e.target.value })}
                  className="input"
                >
                  <option value="debt_owed">借出(别人欠我)</option>
                  <option value="debt_owing">借入(我欠别人)</option>
                  <option value="investment">投资</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">名称</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="input"
                  placeholder="如：借给张三"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">对方</label>
                <input
                  type="text"
                  value={form.related_user}
                  onChange={(e) => setForm({ ...form, related_user: e.target.value })}
                  className="input"
                  placeholder="对方姓名(可选)"
                />
              </div>

              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-sm font-medium mb-1">金额</label>
                  <input
                    type="number"
                    step="0.01"
                    value={form.amount}
                    onChange={(e) => setForm({ ...form, amount: parseFloat(e.target.value) || 0 })}
                    className="input"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">利率(%)</label>
                  <input
                    type="number"
                    step="0.01"
                    value={form.interest_rate}
                    onChange={(e) => setForm({ ...form, interest_rate: parseFloat(e.target.value) || 0 })}
                    className="input"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-sm font-medium mb-1">开始日期</label>
                  <input
                    type="date"
                    value={form.start_date}
                    onChange={(e) => setForm({ ...form, start_date: e.target.value })}
                    className="input"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">结束日期</label>
                  <input
                    type="date"
                    value={form.end_date}
                    onChange={(e) => setForm({ ...form, end_date: e.target.value })}
                    className="input"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">备注</label>
                <textarea
                  value={form.note}
                  onChange={(e) => setForm({ ...form, note: e.target.value })}
                  className="input"
                  rows={2}
                />
              </div>

              <div className="flex gap-3">
                <button type="button" onClick={closeModal} className="flex-1 btn-secondary">
                  取消
                </button>
                <button type="submit" disabled={createMutation.isPending || updateMutation.isPending} className="flex-1 btn-primary">
                  {createMutation.isPending || updateMutation.isPending ? '保存中...' : '保存'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
