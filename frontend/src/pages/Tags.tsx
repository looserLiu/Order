import { useState, FormEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { tagApi, Tag } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline'

// Form state type
interface TagForm {
  name: string
  color: string
}

const defaultColors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8', '#F7DC6F', '#BDC3C7', '#E74C3C']

export default function Tags() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [form, setForm] = useState<TagForm>({ name: '', color: defaultColors[0] })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['tags'],
    queryFn: () => tagApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: TagForm) => tagApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tags'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<TagForm> }) => tagApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tags'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => tagApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tags'] }),
  })

  const tags = data?.data.data || []

  const openModal = (tag?: Tag) => {
    if (tag) {
      setEditingId(tag.id)
      setForm({ name: tag.name, color: tag.color || defaultColors[0] })
    } else {
      setEditingId(null)
      setForm({ name: '', color: defaultColors[Math.floor(Math.random() * defaultColors.length)] })
    }
    setShowModal(true)
  }

  const closeModal = () => {
    setShowModal(false)
    setEditingId(null)
    setForm({ name: '', color: defaultColors[0] })
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
        <h2 className="text-2xl font-bold">标签管理</h2>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          新增标签
        </button>
      </div>

      <div className="flex flex-wrap gap-3">
        {tags.map((tag: Tag) => (
          <div
            key={tag.id}
            className="flex items-center gap-2 px-3 py-2 rounded-lg bg-white border border-gray-200"
          >
            <div
              className="w-3 h-3 rounded-full"
              style={{ backgroundColor: tag.color || '#999' }}
            />
            <span className="font-medium">{tag.name}</span>
            <button onClick={() => openModal(tag)} className="p-1 hover:bg-gray-100 rounded">
              <PencilIcon className="w-3 h-3 text-gray-500" />
            </button>
            <button onClick={() => deleteMutation.mutate(tag.id)} className="p-1 hover:bg-gray-100 rounded">
              <TrashIcon className="w-3 h-3 text-red-500" />
            </button>
          </div>
        ))}

        {tags.length === 0 && (
          <div className="text-center py-12 text-gray-500 w-full">
            <p>暂无标签</p>
            <p className="text-sm">点击上方按钮创建标签</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">{editingId ? '编辑标签' : '新增标签'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">标签名称</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="input"
                  required
                />
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
