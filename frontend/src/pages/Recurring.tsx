import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { transactionApi, accountApi, categoryApi } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, ClockIcon, ArrowPathIcon } from '@heroicons/react/24/outline'

const recurringOptions = [
  { value: 'daily', label: '每天' },
  { value: 'weekly', label: '每周' },
  { value: 'monthly', label: '每月' },
  { value: 'yearly', label: '每年' },
]

export default function RecurringTransactions() {
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({
    type: 'expense',
    account_id: '',
    category_id: '',
    amount: 0,
    merchant: '',
    note: '',
    bill_date: new Date().toISOString().split('T')[0],
    is_recurring: true,
    recurring_rule: 'monthly',
    next_date: new Date().toISOString().split('T')[0],
  })
  const queryClient = useQueryClient()

  const { data: accountsData } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => accountApi.list(),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(),
  })

  const { data: recurringData } = useQuery({
    queryKey: ['transactions', 'recurring'],
    queryFn: () => transactionApi.list({ is_recurring: true }),
  })

  const createMutation = useMutation({
    mutationFn: (data: any) => transactionApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      setShowModal(false)
      resetForm()
    },
  })

  const accounts = accountsData?.data?.data || []
  const categories = categoriesData?.data?.data || []
  const recurringTransactions = recurringData?.data?.data?.list || []

  const expenseCategories = categories.filter((c: any) => c.type === 'expense')
  const incomeCategories = categories.filter((c: any) => c.type === 'income')

  const resetForm = () => {
    setForm({
      type: 'expense',
      account_id: accounts[0]?.id || '',
      category_id: '',
      amount: 0,
      merchant: '',
      note: '',
      bill_date: new Date().toISOString().split('T')[0],
      is_recurring: true,
      recurring_rule: 'monthly',
      next_date: new Date().toISOString().split('T')[0],
    })
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
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">周期记账</h2>
        <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          添加周期记账
        </button>
      </div>

      <div className="card">
        <div className="flex items-center gap-2 mb-4 text-blue-600">
          <ClockIcon className="w-5 h-5" />
          <span className="text-sm">周期记账将自动在指定日期创建交易记录</span>
        </div>

        <div className="space-y-3">
          {recurringTransactions.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <ArrowPathIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
              <p>暂无周期记账</p>
              <p className="text-sm">点击上方按钮添加周期记账</p>
            </div>
          ) : (
            recurringTransactions.map((tx: any) => (
              <div key={tx.id} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50">
                <div className="flex items-center gap-4">
                  <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                    tx.type === 'income' ? 'bg-green-100' : 'bg-red-100'
                  }`}>
                    <span className={`text-lg font-bold ${tx.type === 'income' ? 'text-green-600' : 'text-red-500'}`}>
                      {tx.type === 'income' ? '+' : '-'}
                    </span>
                  </div>
                  <div>
                    <p className="font-medium">{tx.merchant || tx.category?.name || '未分类'}</p>
                    <p className="text-sm text-gray-500">
                      {tx.recurring_rule === 'daily' ? '每天' : 
                       tx.recurring_rule === 'weekly' ? '每周' : 
                       tx.recurring_rule === 'monthly' ? '每月' : '每年'}
                      {' · '}下次: {tx.next_date ? new Date(tx.next_date).toLocaleDateString() : '-'}
                    </p>
                  </div>
                </div>
                <div className="text-right">
                  <p className={`font-semibold text-lg ${tx.type === 'income' ? 'text-green-600' : 'text-red-500'}`}>
                    {tx.type === 'income' ? '+' : '-'}¥{tx.amount?.toFixed(2)}
                  </p>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-end md:items-center justify-center z-50 p-4">
          <div className="bg-white rounded-t-2xl md:rounded-xl w-full md:max-w-md p-4 md:p-6 animate-slide-up md:animate-none max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">添加周期记账</h3>
              <button onClick={() => setShowModal(false)} className="p-1">
                <TrashIcon className="w-6 h-6" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="flex gap-2">
                {['expense', 'income'].map(t => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => setForm({ ...form, type: t })}
                    className={`flex-1 py-2 rounded-lg text-sm font-medium ${
                      form.type === t
                        ? t === 'expense' ? 'bg-red-500 text-white' : 'bg-green-600 text-white'
                        : 'bg-gray-100 text-gray-600'
                    }`}
                  >
                    {t === 'expense' ? '支出' : '收入'}
                  </button>
                ))}
              </div>

              <div>
                <input
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  value={form.amount || ''}
                  onChange={(e) => setForm({ ...form, amount: parseFloat(e.target.value) || 0 })}
                  className="w-full text-2xl font-bold text-center py-2 border-b-2 border-gray-200 focus:border-primary-500 outline-none"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-2">
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

              <div className="grid grid-cols-2 gap-2">
                <select
                  value={form.recurring_rule}
                  onChange={(e) => setForm({ ...form, recurring_rule: e.target.value })}
                  className="input"
                >
                  {recurringOptions.map(opt => (
                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                  ))}
                </select>

                <input
                  type="date"
                  value={form.next_date}
                  onChange={(e) => setForm({ ...form, next_date: e.target.value })}
                  className="input"
                  required
                />
              </div>

              <input
                type="text"
                placeholder="商家(可选)"
                value={form.merchant}
                onChange={(e) => setForm({ ...form, merchant: e.target.value })}
                className="input"
              />

              <input
                type="text"
                placeholder="备注(可选)"
                value={form.note}
                onChange={(e) => setForm({ ...form, note: e.target.value })}
                className="input"
              />

              <button
                type="submit"
                disabled={createMutation.isPending}
                className="btn-primary w-full py-3"
              >
                {createMutation.isPending ? '保存中...' : '保存'}
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
