import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { transactionApi, accountApi, categoryApi } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, FunnelIcon } from '@heroicons/react/24/outline'

export default function Transactions() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [filters, setFilters] = useState({ type: '', category_id: '', account_id: '' })
  const [form, setForm] = useState({
    type: 'expense',
    account_id: '',
    category_id: '',
    amount: 0,
    merchant: '',
    note: '',
    bill_date: new Date().toISOString().split('T')[0],
  })
  const queryClient = useQueryClient()

  const { data: txData } = useQuery({
    queryKey: ['transactions', filters],
    queryFn: () => transactionApi.list(filters),
  })

  const { data: accountsData } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => accountApi.list(),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: any) => transactionApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['summary'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => transactionApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['summary'] })
    },
  })

  const transactions = txData?.data?.data?.list || []
  const accounts = accountsData?.data?.data || []
  const categories = categoriesData?.data?.data || []

  const expenseCategories = categories.filter((c: any) => c.type === 'expense')
  const incomeCategories = categories.filter((c: any) => c.type === 'income')

  const dateRanges = [
    { label: '今天', days: 0 },
    { label: '本周', days: 7 },
    { label: '本月', days: 30 },
    { label: '上月', days: 60 },
  ]

  const setDateFilter = (days: number) => {
    const end = new Date()
    const start = new Date()
    start.setDate(start.getDate() - days)
    setFilters({
      ...filters,
      start_date: start.toISOString().split('T')[0],
      end_date: end.toISOString().split('T')[0],
    })
  }

  const openModal = (tx?: any) => {
    if (tx) {
      setEditingId(tx.id)
      setForm({
        type: tx.type,
        account_id: tx.account_id,
        category_id: tx.category_id || '',
        amount: tx.amount,
        merchant: tx.merchant || '',
        note: tx.note || '',
        bill_date: tx.bill_date?.split('T')[0] || new Date().toISOString().split('T')[0],
      })
    } else {
      setEditingId(null)
      setForm({
        type: 'expense',
        account_id: accounts[0]?.id || '',
        category_id: '',
        amount: 0,
        merchant: '',
        note: '',
        bill_date: new Date().toISOString().split('T')[0],
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
    const submitData = {
      ...form,
      account_id: form.account_id || undefined,
      category_id: form.category_id || undefined,
    }
    createMutation.mutate(submitData)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <h2 className="text-2xl font-bold">记账流水</h2>
        <div className="flex gap-2">
          {dateRanges.map(range => (
            <button
              key={range.days}
              onClick={() => setDateFilter(range.days)}
              className="px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
            >
              {range.label}
            </button>
          ))}
        </div>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          记一笔
        </button>
      </div>

      <div className="card">
        <div className="flex flex-wrap gap-3 mb-4">
          <select
            value={filters.type}
            onChange={(e) => setFilters({ ...filters, type: e.target.value })}
            className="input w-auto"
          >
            <option value="">全部类型</option>
            <option value="income">收入</option>
            <option value="expense">支出</option>
            <option value="transfer">转账</option>
          </select>
          <select
            value={filters.account_id}
            onChange={(e) => setFilters({ ...filters, account_id: e.target.value })}
            className="input w-auto"
          >
            <option value="">全部账户</option>
            {accounts.map((acc: any) => (
              <option key={acc.id} value={acc.id}>{acc.name}</option>
            ))}
          </select>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left text-sm text-gray-500 border-b">
                <th className="pb-3">日期</th>
                <th className="pb-3">类型</th>
                <th className="pb-3">分类</th>
                <th className="pb-3">账户</th>
                <th className="pb-3">商家</th>
                <th className="pb-3">金额</th>
                <th className="pb-3"></th>
              </tr>
            </thead>
            <tbody>
              {transactions.map((tx: any) => (
                <tr key={tx.id} className="border-b border-gray-100">
                  <td className="py-3">{tx.bill_date?.split('T')[0]}</td>
                  <td className="py-3">
                    <span className={`px-2 py-1 rounded text-xs ${
                      tx.type === 'income' ? 'bg-green-100 text-green-700' :
                      tx.type === 'expense' ? 'bg-red-100 text-red-700' : 'bg-blue-100 text-blue-700'
                    }`}>
                      {tx.type === 'income' ? '收入' : tx.type === 'expense' ? '支出' : '转账'}
                    </span>
                  </td>
                  <td className="py-3">{tx.category?.name || '-'}</td>
                  <td className="py-3">{tx.account?.name || '-'}</td>
                  <td className="py-3">{tx.merchant || '-'}</td>
                  <td className="py-3 font-medium">
                    <span className={tx.type === 'income' ? 'text-green-600' : 'text-red-500'}>
                      {tx.type === 'income' ? '+' : '-'}¥{tx.amount?.toFixed(2)}
                    </span>
                  </td>
                  <td className="py-3">
                    <button onClick={() => deleteMutation.mutate(tx.id)} className="p-1 hover:bg-gray-100 rounded">
                      <TrashIcon className="w-4 h-4 text-red-500" />
                    </button>
                  </td>
                </tr>
              ))}
              {transactions.length === 0 && (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-gray-500">暂无记录</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">记一笔</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="flex gap-2">
                {['expense', 'income', 'transfer'].map(t => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => setForm({ ...form, type: t })}
                    className={`flex-1 py-2 rounded-lg text-sm font-medium ${
                      form.type === t
                        ? t === 'expense' ? 'bg-red-500 text-white' : t === 'income' ? 'bg-green-600 text-white' : 'bg-blue-500 text-white'
                        : 'bg-gray-100 text-gray-600'
                    }`}
                  >
                    {t === 'expense' ? '支出' : t === 'income' ? '收入' : '转账'}
                  </button>
                ))}
              </div>
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
                <label className="block text-sm font-medium mb-1">账户</label>
                <select
                  value={form.account_id}
                  onChange={(e) => setForm({ ...form, account_id: e.target.value })}
                  className="input"
                  required
                >
                  <option value="">选择账户</option>
                  {accounts.map((acc: any) => (
                    <option key={acc.id} value={acc.id}>{acc.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">分类</label>
                <select
                  value={form.category_id}
                  onChange={(e) => setForm({ ...form, category_id: e.target.value })}
                  className="input"
                >
                  <option value="">选择分类</option>
                  {(form.type === 'expense' ? expenseCategories : incomeCategories).map((cat: any) => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">日期</label>
                <input
                  type="date"
                  value={form.bill_date}
                  onChange={(e) => setForm({ ...form, bill_date: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">商家</label>
                <input
                  type="text"
                  value={form.merchant}
                  onChange={(e) => setForm({ ...form, merchant: e.target.value })}
                  className="input"
                />
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
