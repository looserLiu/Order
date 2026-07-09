import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { goalApi } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, TrophyIcon } from '@heroicons/react/24/outline'

export default function Goals() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [showAddAmount, setShowAddAmount] = useState<string | null>(null)
  const [form, setForm] = useState({
    name: '',
    target_amount: 0,
    current_amount: 0,
    deadline: '',
    category: 'savings',
    note: '',
  })
  const [amountForm, setAmountForm] = useState({ amount: 0 })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['goals'],
    queryFn: () => goalApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: any) => goalApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => goalApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] })
      closeModal()
    },
  })

  const addAmountMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => goalApi.addAmount(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] })
      setShowAddAmount(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => goalApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['goals'] }),
  })

  const goals = data?.data?.data || []

  const openModal = (goal?: any) => {
    if (goal) {
      setEditingId(goal.id)
      setForm({
        name: goal.name,
        target_amount: goal.target_amount,
        current_amount: goal.current_amount,
        deadline: goal.deadline?.split('T')[0] || '',
        category: goal.category || 'savings',
        note: goal.note || '',
      })
    } else {
      setEditingId(null)
      setForm({
        name: '',
        target_amount: 0,
        current_amount: 0,
        deadline: '',
        category: 'savings',
        note: '',
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

  const getCategoryText = (category: string) => {
    switch (category) {
      case 'savings': return '储蓄'
      case 'investment': return '投资'
      case 'debt': return '还债'
      case 'purchase': return '购物'
      default: return '其他'
    }
  }

  const getCategoryColor = (category: string) => {
    switch (category) {
      case 'savings': return 'bg-green-100 text-green-700'
      case 'investment': return 'bg-blue-100 text-blue-700'
      case 'debt': return 'bg-red-100 text-red-700'
      case 'purchase': return 'bg-yellow-100 text-yellow-700'
      default: return 'bg-gray-100 text-gray-700'
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <TrophyIcon className="w-6 h-6" />
          <h2 className="text-2xl font-bold">财务目标</h2>
        </div>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          新建目标
        </button>
      </div>

      <div className="grid gap-4">
        {goals.map((goal: any) => {
          const progress = goal.target_amount > 0 ? (goal.current_amount / goal.target_amount) * 100 : 0
          const isCompleted = goal.status === 'completed'

          return (
            <div key={goal.id} className={`card ${isCompleted ? 'bg-green-50 border-green-200' : ''}`}>
              <div className="flex items-start justify-between mb-4">
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold text-lg">{goal.name}</h3>
                    <span className={`px-2 py-0.5 rounded text-xs ${getCategoryColor(goal.category)}`}>
                      {getCategoryText(goal.category)}
                    </span>
                    {isCompleted && (
                      <span className="px-2 py-0.5 rounded text-xs bg-green-600 text-white">已完成</span>
                    )}
                  </div>
                  {goal.note && <p className="text-sm text-gray-500 mt-1">{goal.note}</p>}
                </div>
                <div className="flex gap-1">
                  <button onClick={() => openModal(goal)} className="p-2 hover:bg-gray-100 rounded">
                    <PencilIcon className="w-4 h-4 text-gray-500" />
                  </button>
                  <button onClick={() => deleteMutation.mutate(goal.id)} className="p-2 hover:bg-gray-100 rounded">
                    <TrashIcon className="w-4 h-4 text-red-500" />
                  </button>
                </div>
              </div>

              <div className="mb-2">
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-500">当前进度</span>
                  <span className={isCompleted ? 'text-green-600' : ''}>
                    ¥{goal.current_amount?.toFixed(2)} / ¥{goal.target_amount?.toFixed(2)}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div
                    className={`h-3 rounded-full transition-all ${isCompleted ? 'bg-green-600' : 'bg-primary-600'}`}
                    style={{ width: `${Math.min(progress, 100)}%` }}
                  />
                </div>
                <div className="flex justify-between text-xs text-gray-500 mt-1">
                  <span>{progress.toFixed(1)}%</span>
                  {goal.deadline && <span>截止: {goal.deadline.split('T')[0]}</span>}
                </div>
              </div>

              {!isCompleted && (
                <button
                  onClick={() => setShowAddAmount(goal.id)}
                  className="btn-secondary w-full mt-2"
                >
                  存入金额
                </button>
              )}
            </div>
          )
        })}

        {goals.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <TrophyIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
            <p>暂无财务目标</p>
            <p className="text-sm">设定目标，让储蓄更有动力</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">{editingId ? '编辑目标' : '新建目标'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">目标名称</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">目标金额</label>
                  <input
                    type="number"
                    step="0.01"
                    value={form.target_amount}
                    onChange={(e) => setForm({ ...form, target_amount: parseFloat(e.target.value) || 0 })}
                    className="input"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">当前已存</label>
                  <input
                    type="number"
                    step="0.01"
                    value={form.current_amount}
                    onChange={(e) => setForm({ ...form, current_amount: parseFloat(e.target.value) || 0 })}
                    className="input"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">目标类型</label>
                  <select
                    value={form.category}
                    onChange={(e) => setForm({ ...form, category: e.target.value })}
                    className="input"
                  >
                    <option value="savings">储蓄</option>
                    <option value="investment">投资</option>
                    <option value="debt">还债</option>
                    <option value="purchase">购物</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">目标期限</label>
                  <input
                    type="date"
                    value={form.deadline}
                    onChange={(e) => setForm({ ...form, deadline: e.target.value })}
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
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={closeModal} className="btn-secondary flex-1">取消</button>
                <button type="submit" className="btn-primary flex-1">保存</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showAddAmount && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">存入金额</h3>
            <form onSubmit={(e) => {
              e.preventDefault()
              addAmountMutation.mutate({ id: showAddAmount, data: amountForm })
            }} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">存入金额</label>
                <input
                  type="number"
                  step="0.01"
                  value={amountForm.amount}
                  onChange={(e) => setAmountForm({ amount: parseFloat(e.target.value) || 0 })}
                  className="input"
                  required
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setShowAddAmount(null)} className="btn-secondary flex-1">取消</button>
                <button type="submit" className="btn-primary flex-1">确认</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
