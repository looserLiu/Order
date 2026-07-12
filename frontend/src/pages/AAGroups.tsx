import { useState, FormEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { aaGroupApi, AAGroup, AAMember } from '../services/api'
import { PlusIcon, TrashIcon, CalculatorIcon } from '@heroicons/react/24/outline'

// Form state types
interface GroupForm {
  name: string
  description: string
  members: { name: string }[]
}

interface ExpenseForm {
  member_id: string
  amount: number
  note: string
}

// Settlement type
interface Settlement {
  from_id: string
  to_id: string
  amount: number
}

export default function AAGroups() {
  const [showModal, setShowModal] = useState(false)
  const [showExpenseModal, setShowExpenseModal] = useState<string | null>(null)
  const [form, setForm] = useState<GroupForm>({ name: '', description: '', members: [{ name: '' }] })
  const [expenseForm, setExpenseForm] = useState<ExpenseForm>({ member_id: '', amount: 0, note: '' })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['aa-groups'],
    queryFn: () => aaGroupApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: { name: string; description?: string; members: { name: string }[] }) => aaGroupApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['aa-groups'] })
      setShowModal(false)
      setForm({ name: '', description: '', members: [{ name: '' }] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => aaGroupApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['aa-groups'] }),
  })

  const addExpenseMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: ExpenseForm }) => aaGroupApi.addExpense(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['aa-groups'] })
      setShowExpenseModal(null)
    },
  })

  const { data: settlementsData } = useQuery({
    queryKey: ['aa-settlements', showExpenseModal],
    queryFn: () => showExpenseModal ? aaGroupApi.getSettlements(showExpenseModal) : null,
    enabled: !!showExpenseModal,
  })

  const groups = (data?.data?.data || []) as AAGroup[]
  const settlements = (settlementsData?.data?.data || []) as Settlement[]

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(form)
  }

  const handleAddMember = () => {
    setForm({ ...form, members: [...form.members, { name: '' }] })
  }

  const handleAddExpense = (e: React.FormEvent) => {
    e.preventDefault()
    if (showExpenseModal) {
      addExpenseMutation.mutate({ id: showExpenseModal, data: expenseForm })
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <CalculatorIcon className="w-6 h-6" />
          <h2 className="text-2xl font-bold">AA记账</h2>
        </div>
        <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          新建AA账本
        </button>
      </div>

      <div className="grid gap-4">
        {groups.map((group: AAGroup) => (
          <div key={group.id} className="card">
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="font-semibold text-lg">{group.name}</h3>
                {group.description && <p className="text-sm text-gray-500">{group.description}</p>}
              </div>
              <button onClick={() => deleteMutation.mutate(group.id)} className="p-2 hover:bg-gray-100 rounded">
                <TrashIcon className="w-4 h-4 text-red-500" />
              </button>
            </div>

            <div className="mb-4">
              <p className="text-sm text-gray-500 mb-2">总费用: ¥{group.total_amount?.toFixed(2)}</p>
              <div className="space-y-2">
                {group.members?.map((member: AAMember) => (
                  <div key={member.id} className="flex justify-between text-sm">
                    <span>{member.name}</span>
                    <div className="flex gap-4">
                      <span className="text-green-600">已付: ¥{member.paid?.toFixed(2)}</span>
                      <span className={member.owe > 0 ? 'text-red-500' : 'text-gray-400'}>
                        {member.owe > 0 ? `应摊: ¥${member.owe?.toFixed(2)}` : '已结清'}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex gap-2">
              <button
                onClick={() => setShowExpenseModal(group.id)}
                className="btn-secondary flex-1"
              >
                添加费用
              </button>
              <button
                onClick={() => aaGroupApi.getSettlements(group.id).then(r => {
                  const settlements = (r.data.data || []) as Settlement[]
                  if (settlements.length > 0) {
                    alert('结算方案: ' + settlements.map((s) =>
                      `${s.from_id} -> ${s.to_id}: ¥${s.amount}`
                    ).join('\n'))
                  }
                })}
                className="btn-secondary flex-1"
              >
                查看结算
              </button>
            </div>
          </div>
        ))}

        {groups.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <CalculatorIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
            <p>暂无AA账本</p>
            <p className="text-sm">适用于旅游、聚餐等多人均摊场景</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">新建AA账本</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">账本名称</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="input"
                  placeholder="例如：国庆旅游"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">描述</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  className="input"
                  rows={2}
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">成员</label>
                {form.members.map((m, i) => (
                  <input
                    key={i}
                    type="text"
                    value={m.name}
                    onChange={(e) => {
                      const newMembers = [...form.members]
                      newMembers[i].name = e.target.value
                      setForm({ ...form, members: newMembers })
                    }}
                    className="input mb-2"
                    placeholder="成员名称"
                  />
                ))}
                <button type="button" onClick={handleAddMember} className="text-sm text-primary-600">
                  + 添加成员
                </button>
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setShowModal(false)} className="btn-secondary flex-1">取消</button>
                <button type="submit" className="btn-primary flex-1">创建</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showExpenseModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">添加费用</h3>
            <form onSubmit={handleAddExpense} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">付款人</label>
                <select
                  value={expenseForm.member_id}
                  onChange={(e) => setExpenseForm({ ...expenseForm, member_id: e.target.value })}
                  className="input"
                  required
                >
                  <option value="">选择成员</option>
                  {groups.find((g) => g.id === showExpenseModal)?.members?.map((m) => (
                    <option key={m.id} value={m.id}>{m.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">金额</label>
                <input
                  type="number"
                  step="0.01"
                  value={expenseForm.amount}
                  onChange={(e) => setExpenseForm({ ...expenseForm, amount: parseFloat(e.target.value) || 0 })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">备注</label>
                <input
                  type="text"
                  value={expenseForm.note}
                  onChange={(e) => setExpenseForm({ ...expenseForm, note: e.target.value })}
                  className="input"
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setShowExpenseModal(null)} className="btn-secondary flex-1">取消</button>
                <button type="submit" className="btn-primary flex-1">添加</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
