import { useState, FormEvent, ChangeEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { budgetApi, categoryApi, Budget, Category } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline'

// Form state type
interface BudgetForm {
  category_id: string
  amount: number
  period: string
  start_date: string
  alert_threshold: number
}

export default function Budgets() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [form, setForm] = useState<BudgetForm>({
    category_id: '',
    amount: 0,
    period: 'monthly',
    start_date: new Date().toISOString().split('T')[0],
    alert_threshold: 0.8,
  })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['budgets'],
    queryFn: () => budgetApi.list(),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(),
  })

  const { data: progressData } = useQuery({
    queryKey: ['budgets-progress'],
    queryFn: async () => {
      const budgets = data?.data.data || []
      const progressPromises = budgets.map((b: Budget) => budgetApi.getProgress(b.id))
      return Promise.all(progressPromises.map(p => p.catch(() => null)))
    },
    enabled: !!data?.data.data?.length,
  })

  const createMutation = useMutation({
    mutationFn: (data: any) => budgetApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budgets'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => budgetApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budgets'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => budgetApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['budgets'] }),
  })

  const budgets = data?.data.data || []
  const categories = categoriesData?.data.data || []
  const progressList = (progressData || []).map((p: any) => p?.data).filter(Boolean)

  const openModal = (budget?: Budget) => {
    if (budget) {
      setEditingId(budget.id)
      setForm({
        category_id: budget.category_id || '',
        amount: budget.amount,
        period: budget.period,
        start_date: budget.start_date?.split('T')[0] || new Date().toISOString().split('T')[0],
        alert_threshold: budget.alert_threshold || 0.8,
      })
    } else {
      setEditingId(null)
      setForm({
        category_id: '',
        amount: 0,
        period: 'monthly',
        start_date: new Date().toISOString().split('T')[0],
        alert_threshold: 0.8,
      })
    }
    setShowModal(true)
  }

  const closeModal = () => {
    setShowModal(false)
    setEditingId(null)
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const submitData = {
      ...form,
      category_id: form.category_id || null,
      amount: Number(form.amount),
    }
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: submitData })
    } else {
      createMutation.mutate(submitData)
    }
  }

  const getCategoryName = (categoryId: string) => {
    const cat = categories.find((c: Category) => c.id === categoryId)
    return cat?.name || '全部分类'
  }

  const getProgress = (budgetId: string) => {
    return progressList.find((p: any) => p?.budget?.id === budgetId)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">预算管理</h2>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          新增预算
        </button>
      </div>

      <div className="grid gap-4">
        {budgets.map((budget: Budget) => {
          const progress = getProgress(budget.id)
          const spent = progress?.spent || 0
          const remaining = progress?.remaining || 0
          const progressPercent = Math.min(progress?.progress || 0, 100)
          const isAlert = progress?.alert

          return (
            <div key={budget.id} className="card">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="font-semibold text-lg">{getCategoryName(budget.category_id || '')}</h3>
                  <p className="text-sm text-gray-500">
                    {budget.period === 'monthly' ? '每月' : '每年'} · ¥{budget.amount?.toFixed(2)}
                  </p>
                </div>
                <div className="flex gap-1">
                  <button onClick={() => openModal(budget)} className="p-2 hover:bg-gray-100 rounded-lg">
                    <PencilIcon className="w-4 h-4 text-gray-500" />
                  </button>
                  <button onClick={() => deleteMutation.mutate(budget.id)} className="p-2 hover:bg-gray-100 rounded-lg">
                    <TrashIcon className="w-4 h-4 text-red-500" />
                  </button>
                </div>
              </div>

              <div className="mb-2 flex justify-between text-sm">
                <span className={isAlert ? 'text-red-500' : 'text-gray-600'}>
                  已花费 ¥{spent?.toFixed(2)}
                </span>
                <span className={remaining < 0 ? 'text-red-500' : 'text-green-600'}>
                  剩余 ¥{remaining?.toFixed(2)}
                </span>
              </div>

              <div className="w-full bg-gray-200 rounded-full h-3 mb-2">
                <div
                  className={`h-3 rounded-full transition-all ${isAlert ? 'bg-red-500' : 'bg-primary-600'}`}
                  style={{ width: `${progressPercent}%` }}
                />
              </div>

              {isAlert && (
                <div className="flex items-center gap-2 text-red-500 text-sm">
                  <ExclamationTriangleIcon className="w-4 h-4" />
                  <span>预算已超过{Math.floor(budget.alert_threshold * 100)}%</span>
                </div>
              )}
            </div>
          )
        })}

        {budgets.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <p>暂无预算</p>
            <p className="text-sm">点击上方按钮创建预算</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">{editingId ? '编辑预算' : '新增预算'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">分类</label>
                <select
                  value={form.category_id}
                  onChange={(e) => setForm({ ...form, category_id: e.target.value })}
                  className="input"
                >
                  <option value="">全部分类</option>
                  {categories.filter((c: Category) => c.type === 'expense').map((cat: Category) => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">预算金额</label>
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
                <label className="block text-sm font-medium mb-1">周期</label>
                <select
                  value={form.period}
                  onChange={(e) => setForm({ ...form, period: e.target.value })}
                  className="input"
                >
                  <option value="monthly">每月</option>
                  <option value="yearly">每年</option>
                </select>
              </div>
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
                <label className="block text-sm font-medium mb-1">提醒阈值 ({Math.round(form.alert_threshold * 100)}%)</label>
                <input
                  type="range"
                  min="0.5"
                  max="1"
                  step="0.05"
                  value={form.alert_threshold}
                  onChange={(e) => setForm({ ...form, alert_threshold: parseFloat(e.target.value) })}
                  className="w-full"
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
