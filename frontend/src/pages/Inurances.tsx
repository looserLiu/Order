import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { insuranceApi } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, ShieldCheckIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline'

export default function Insurances() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [form, setForm] = useState({
    name: '',
    insurance_type: 'health',
    company: '',
    premium: 0,
    payment_type: 'yearly',
    start_date: '',
    end_date: '',
    coverage: 0,
    beneficiary: '',
    note: '',
  })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['insurances'],
    queryFn: () => insuranceApi.list(),
  })

  const { data: summaryData } = useQuery({
    queryKey: ['insurances-summary'],
    queryFn: () => insuranceApi.getSummary(),
  })

  const createMutation = useMutation({
    mutationFn: (data: any) => insuranceApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['insurances'] })
      queryClient.invalidateQueries({ queryKey: ['insurances-summary'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => insuranceApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['insurances'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => insuranceApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['insurances'] })
      queryClient.invalidateQueries({ queryKey: ['insurances-summary'] })
    },
  })

  const insurances = data?.data?.data || []
  const summary = summaryData?.data?.data || {}

  const openModal = (ins?: any) => {
    if (ins) {
      setEditingId(ins.id)
      setForm({
        name: ins.name,
        insurance_type: ins.insurance_type || 'health',
        company: ins.company || '',
        premium: ins.premium || 0,
        payment_type: ins.payment_type || 'yearly',
        start_date: ins.start_date?.split('T')[0] || '',
        end_date: ins.end_date?.split('T')[0] || '',
        coverage: ins.coverage || 0,
        beneficiary: ins.beneficiary || '',
        note: ins.note || '',
      })
    } else {
      setEditingId(null)
      setForm({
        name: '',
        insurance_type: 'health',
        company: '',
        premium: 0,
        payment_type: 'yearly',
        start_date: '',
        end_date: '',
        coverage: 0,
        beneficiary: '',
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

  const getTypeText = (type: string) => {
    switch (type) {
      case 'health': return '医疗险'
      case 'life': return '寿险'
      case 'car': return '车险'
      case 'property': return '财产险'
      case 'travel': return '旅行险'
      default: return '其他'
    }
  }

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'health': return 'bg-blue-100 text-blue-700'
      case 'life': return 'bg-purple-100 text-purple-700'
      case 'car': return 'bg-orange-100 text-orange-700'
      case 'property': return 'bg-green-100 text-green-700'
      case 'travel': return 'bg-cyan-100 text-cyan-700'
      default: return 'bg-gray-100 text-gray-700'
    }
  }

  const isExpiringSoon = (nextPaymentDate: string) => {
    if (!nextPaymentDate) return false
    const date = new Date(nextPaymentDate)
    const now = new Date()
    const days = (date.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
    return days <= 30 && days > 0
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <ShieldCheckIcon className="w-6 h-6" />
          <h2 className="text-2xl font-bold">保险管理</h2>
        </div>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          添加保单
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="card text-center">
          <p className="text-sm text-gray-500">年度保费</p>
          <p className="text-2xl font-bold text-primary-600">¥{summary.total_premium?.toFixed(2) || '0'}</p>
        </div>
        <div className="card text-center">
          <p className="text-sm text-gray-500">有效保单</p>
          <p className="text-2xl font-bold">{summary.active_count || 0}</p>
        </div>
        <div className="card text-center">
          <p className="text-sm text-gray-500">即将到期</p>
          <p className="text-2xl font-bold text-orange-500">{summary.expiring_count || 0}</p>
        </div>
        <div className="card text-center">
          <p className="text-sm text-gray-500">总保额</p>
          <p className="text-2xl font-bold text-green-600">¥{(summary.total_coverage || 0).toFixed(0)}</p>
        </div>
      </div>

      <div className="grid gap-4">
        {insurances.map((ins: any) => (
          <div key={ins.id} className="card">
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-2">
                  <h3 className="font-semibold text-lg">{ins.name}</h3>
                  <span className={`px-2 py-0.5 rounded text-xs ${getTypeColor(ins.insurance_type)}`}>
                    {getTypeText(ins.insurance_type)}
                  </span>
                  {ins.status === 'expired' && (
                    <span className="px-2 py-0.5 rounded text-xs bg-red-100 text-red-700">已过期</span>
                  )}
                </div>
                {ins.company && <p className="text-sm text-gray-500">保险公司: {ins.company}</p>}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-2 mt-3 text-sm">
                  <div>
                    <span className="text-gray-500">保费: </span>
                    <span className="font-medium">¥{ins.premium?.toFixed(2)}/{ins.payment_type === 'yearly' ? '年' : ins.payment_type === 'quarterly' ? '季' : '月'}</span>
                  </div>
                  <div>
                    <span className="text-gray-500">保额: </span>
                    <span className="font-medium">¥{ins.coverage?.toFixed(0)}</span>
                  </div>
                  <div>
                    <span className="text-gray-500">生效: </span>
                    <span className="">{ins.start_date?.split('T')[0]}</span>
                  </div>
                  <div>
                    <span className="text-gray-500">到期: </span>
                    <span className="">{ins.end_date?.split('T')[0] || '-'}</span>
                  </div>
                </div>
                {ins.next_payment_date && isExpiringSoon(ins.next_payment_date) && (
                  <div className="flex items-center gap-2 mt-2 text-orange-500 text-sm">
                    <ExclamationTriangleIcon className="w-4 h-4" />
                    <span>即将到期: {ins.next_payment_date.split('T')[0]}</span>
                  </div>
                )}
              </div>
              <div className="flex gap-1">
                <button onClick={() => openModal(ins)} className="p-2 hover:bg-gray-100 rounded">
                  <PencilIcon className="w-4 h-4 text-gray-500" />
                </button>
                <button onClick={() => deleteMutation.mutate(ins.id)} className="p-2 hover:bg-gray-100 rounded">
                  <TrashIcon className="w-4 h-4 text-red-500" />
                </button>
              </div>
            </div>
          </div>
        ))}

        {insurances.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <ShieldCheckIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
            <p>暂无保单</p>
            <p className="text-sm">管理您的保险资产</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md max-h-[80vh] overflow-y-auto">
            <h3 className="text-lg font-semibold mb-4">{editingId ? '编辑保单' : '添加保单'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">保单名称</label>
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
                  <label className="block text-sm font-medium mb-1">险种</label>
                  <select
                    value={form.insurance_type}
                    onChange={(e) => setForm({ ...form, insurance_type: e.target.value })}
                    className="input"
                  >
                    <option value="health">医疗险</option>
                    <option value="life">寿险</option>
                    <option value="car">车险</option>
                    <option value="property">财产险</option>
                    <option value="travel">旅行险</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">保险公司</label>
                  <input
                    type="text"
                    value={form.company}
                    onChange={(e) => setForm({ ...form, company: e.target.value })}
                    className="input"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">保费</label>
                  <input
                    type="number"
                    step="0.01"
                    value={form.premium}
                    onChange={(e) => setForm({ ...form, premium: parseFloat(e.target.value) || 0 })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">缴费方式</label>
                  <select
                    value={form.payment_type}
                    onChange={(e) => setForm({ ...form, payment_type: e.target.value })}
                    className="input"
                  >
                    <option value="yearly">年缴</option>
                    <option value="quarterly">季缴</option>
                    <option value="monthly">月缴</option>
                  </select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">生效日期</label>
                  <input
                    type="date"
                    value={form.start_date}
                    onChange={(e) => setForm({ ...form, start_date: e.target.value })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">到期日期</label>
                  <input
                    type="date"
                    value={form.end_date}
                    onChange={(e) => setForm({ ...form, end_date: e.target.value })}
                    className="input"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">保额</label>
                  <input
                    type="number"
                    value={form.coverage}
                    onChange={(e) => setForm({ ...form, coverage: parseFloat(e.target.value) || 0 })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">受益人</label>
                  <input
                    type="text"
                    value={form.beneficiary}
                    onChange={(e) => setForm({ ...form, beneficiary: e.target.value })}
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
    </div>
  )
}
