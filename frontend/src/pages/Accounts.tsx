import { useState, FormEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { accountApi, Account } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline'
import { LoadingSpinner } from '../components/LoadingSpinner'

// Account type options
const accountTypes = [
  { value: 'cash', label: '现金', color: '#27AE60' },
  { value: 'bank', label: '银行卡', color: '#3498DB' },
  { value: 'credit', label: '信用卡', color: '#E74C3C' },
  { value: 'wechat', label: '微信', color: '#07C160' },
  { value: 'alipay', label: '支付宝', color: '#1677FF' },
  { value: 'investment', label: '投资', color: '#9B59B6' },
]

// Form state type
interface AccountForm {
  name: string
  type: string
  balance: number
  color: string
}

export default function Accounts() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [form, setForm] = useState<AccountForm>({ name: '', type: 'cash', balance: 0, color: '' })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => accountApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: AccountForm) => accountApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<AccountForm> }) => accountApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => accountApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['accounts'] }),
  })

  const accounts = data?.data.data || []

  const openModal = (account?: Account) => {
    if (account) {
      setEditingId(account.id)
      setForm({ name: account.name, type: account.type, balance: account.balance, color: account.color || '' })
    } else {
      setEditingId(null)
      setForm({ name: '', type: 'cash', balance: 0, color: '' })
    }
    setShowModal(true)
  }

  const closeModal = () => {
    setShowModal(false)
    setEditingId(null)
    setForm({ name: '', type: 'cash', balance: 0, color: '' })
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: form })
    } else {
      createMutation.mutate(form)
    }
  }

  const total = accounts.reduce((sum: number, acc: Account) => sum + (acc.balance || 0), 0)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">账户管理</h2>
          <p className="text-gray-500">总余额: ¥{total.toFixed(2)}</p>
        </div>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          新增账户
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {accounts.map((account: Account) => {
          const typeInfo = accountTypes.find(t => t.value === account.type) || accountTypes[0]
          return (
            <div key={account.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ backgroundColor: account.color || typeInfo.color }}>
                    <span className="text-white font-bold">{account.name[0]}</span>
                  </div>
                  <div>
                    <h3 className="font-semibold">{account.name}</h3>
                    <p className="text-sm text-gray-500">{typeInfo.label}</p>
                  </div>
                </div>
                <div className="flex gap-1">
                  <button onClick={() => openModal(account)} className="p-1 hover:bg-gray-100 rounded">
                    <PencilIcon className="w-4 h-4 text-gray-500" />
                  </button>
                  <button onClick={() => deleteMutation.mutate(account.id)} className="p-1 hover:bg-gray-100 rounded">
                    <TrashIcon className="w-4 h-4 text-red-500" />
                  </button>
                </div>
              </div>
              <div className="mt-4 pt-4 border-t border-gray-100">
                <p className="text-2xl font-bold">¥{account.balance?.toFixed(2) || '0.00'}</p>
              </div>
            </div>
          )
        })}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">{editingId ? '编辑账户' : '新增账户'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">账户名称</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">账户类型</label>
                <select
                  value={form.type}
                  onChange={(e) => setForm({ ...form, type: e.target.value })}
                  className="input"
                >
                  {accountTypes.map(t => (
                    <option key={t.value} value={t.value}>{t.label}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">余额</label>
                <input
                  type="number"
                  step="0.01"
                  value={form.balance}
                  onChange={(e) => setForm({ ...form, balance: parseFloat(e.target.value) || 0 })}
                  className="input"
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={closeModal} className="btn-secondary flex-1">取消</button>
                <button type="submit" className="btn-primary flex-1">保存</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
