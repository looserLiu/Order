import { useState, FormEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { categoryApi, Category } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline'

// Form state type
interface CategoryForm {
  name: string
  type: string
  color: string
}

const defaultColors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8', '#F7DC6F']

export default function Categories() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [form, setForm] = useState<CategoryForm>({ name: '', type: 'expense', color: defaultColors[0] })
  const [filterType, setFilterType] = useState('expense')
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: CategoryForm) => categoryApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CategoryForm> }) => categoryApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => categoryApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['categories'] }),
  })

  const categories = data?.data.data || []
  const filteredCategories = categories.filter((c: Category) => c.type === filterType)

  const openModal = (cat?: Category) => {
    if (cat) {
      setEditingId(cat.id)
      setForm({ name: cat.name, type: cat.type, color: cat.color || defaultColors[0] })
    } else {
      setEditingId(null)
      setForm({ name: '', type: filterType, color: defaultColors[Math.floor(Math.random() * defaultColors.length)] })
    }
    setShowModal(true)
  }

  const closeModal = () => {
    setShowModal(false)
    setEditingId(null)
    setForm({ name: '', type: 'expense', color: defaultColors[0] })
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: form })
    } else {
      createMutation.mutate(form)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">分类管理</h2>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          新增分类
        </button>
      </div>

      <div className="flex gap-2">
        {['expense', 'income', 'transfer'].map(t => (
          <button
            key={t}
            onClick={() => setFilterType(t)}
            className={`px-4 py-2 rounded-lg text-sm font-medium ${
              filterType === t
                ? 'bg-primary-600 text-white'
                : 'bg-white text-gray-600 hover:bg-gray-50'
            }`}
          >
            {t === 'expense' ? '支出' : t === 'income' ? '收入' : '转账'}
          </button>
        ))}
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
        {filteredCategories.map((cat: Category) => (
          <div key={cat.id} className="card">
            <div className="flex items-center gap-3">
              <div
                className="w-10 h-10 rounded-lg flex items-center justify-center"
                style={{ backgroundColor: cat.color || '#E0E0E0' }}
              >
                <span className="text-white text-lg">📁</span>
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="font-medium truncate">{cat.name}</h3>
                <p className="text-xs text-gray-500">{cat.is_system ? '系统' : '自定义'}</p>
              </div>
            </div>
            {!cat.is_system && (
              <div className="flex gap-1 mt-3 pt-3 border-t border-gray-100">
                <button onClick={() => openModal(cat)} className="flex-1 p-1 hover:bg-gray-100 rounded text-xs">
                  <PencilIcon className="w-4 h-4 text-gray-500 mx-auto" />
                </button>
                <button onClick={() => deleteMutation.mutate(cat.id)} className="flex-1 p-1 hover:bg-gray-100 rounded text-xs">
                  <TrashIcon className="w-4 h-4 text-red-500 mx-auto" />
                </button>
              </div>
            )}
          </div>
        ))}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">{editingId ? '编辑分类' : '新增分类'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">分类名称</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">类型</label>
                <select
                  value={form.type}
                  onChange={(e) => setForm({ ...form, type: e.target.value })}
                  className="input"
                >
                  <option value="expense">支出</option>
                  <option value="income">收入</option>
                  <option value="transfer">转账</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">颜色</label>
                <div className="flex gap-2 flex-wrap">
                  {defaultColors.map(color => (
                    <button
                      key={color}
                      type="button"
                      onClick={() => setForm({ ...form, color })}
                      className={`w-8 h-8 rounded-lg ${form.color === color ? 'ring-2 ring-offset-2 ring-primary-500' : ''}`}
                      style={{ backgroundColor: color }}
                    />
                  ))}
                </div>
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
