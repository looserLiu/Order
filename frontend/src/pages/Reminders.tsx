import { useState, FormEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { reminderApi, categoryApi, Reminder, Category } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, BellIcon } from '@heroicons/react/24/outline'

// Form state type
interface ReminderForm {
  title: string
  content: string
  remind_time: string
  repeat_type: string
  category_id: string
  is_active: boolean
}

export default function Reminders() {
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [form, setForm] = useState<ReminderForm>({
    title: '',
    content: '',
    remind_time: '',
    repeat_type: 'none',
    category_id: '',
    is_active: true,
  })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['reminders'],
    queryFn: () => reminderApi.list(),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: Partial<Reminder>) => reminderApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reminders'] })
      closeModal()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Reminder> }) => reminderApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reminders'] })
      closeModal()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => reminderApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['reminders'] }),
  })

  const reminders = (data?.data?.data || []) as Reminder[]
  const categories = (categoriesData?.data?.data || []) as Category[]

  const openModal = (reminder?: Reminder) => {
    if (reminder) {
      setEditingId(reminder.id)
      setForm({
        title: reminder.title,
        content: reminder.content || '',
        remind_time: reminder.remind_time?.slice(0, 16) || '',
        repeat_type: reminder.repeat_type || 'none',
        category_id: reminder.category_id || '',
        is_active: reminder.is_active,
      })
    } else {
      setEditingId(null)
      setForm({
        title: '',
        content: '',
        remind_time: '',
        repeat_type: 'none',
        category_id: '',
        is_active: true,
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
    const submitData: Partial<Reminder> = {
      ...form,
      category_id: form.category_id || undefined,
    }
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: submitData })
    } else {
      createMutation.mutate(submitData)
    }
  }

  const getRepeatText = (type: string) => {
    switch (type) {
      case 'daily': return '每天'
      case 'weekly': return '每周'
      case 'monthly': return '每月'
      case 'yearly': return '每年'
      default: return '不重复'
    }
  }

  const formatTime = (time: string) => {
    if (!time) return ''
    const date = new Date(time)
    return date.toLocaleString()
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <BellIcon className="w-6 h-6" />
          <h2 className="text-2xl font-bold">提醒</h2>
        </div>
        <button onClick={() => openModal()} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          新增提醒
        </button>
      </div>

      <div className="space-y-3">
        {reminders.map((reminder: Reminder) => (
          <div key={reminder.id} className={`card ${!reminder.is_active ? 'opacity-50' : ''}`}>
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <h3 className="font-semibold">{reminder.title}</h3>
                  {!reminder.is_active && (
                    <span className="text-xs bg-gray-200 px-2 py-0.5 rounded">已暂停</span>
                  )}
                </div>
                {reminder.content && (
                  <p className="text-sm text-gray-500 mt-1">{reminder.content}</p>
                )}
                <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
                  <span>{formatTime(reminder.remind_time)}</span>
                  <span>{getRepeatText(reminder.repeat_type || 'none')}</span>
                  {reminder.category_id && (
                    <span className="text-primary-600">
                      {categories.find((c) => c.id === reminder.category_id)?.name}
                    </span>
                  )}
                </div>
              </div>
              <div className="flex gap-1">
                <button onClick={() => openModal(reminder)} className="p-2 hover:bg-gray-100 rounded">
                  <PencilIcon className="w-4 h-4 text-gray-500" />
                </button>
                <button onClick={() => deleteMutation.mutate(reminder.id)} className="p-2 hover:bg-gray-100 rounded">
                  <TrashIcon className="w-4 h-4 text-red-500" />
                </button>
              </div>
            </div>
          </div>
        ))}

        {reminders.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <BellIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
            <p>暂无提醒</p>
            <p className="text-sm">点击上方按钮创建提醒</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">
              {editingId ? '编辑提醒' : '新增提醒'}
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">标题</label>
                <input
                  type="text"
                  value={form.title}
                  onChange={(e) => setForm({ ...form, title: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">内容</label>
                <textarea
                  value={form.content}
                  onChange={(e) => setForm({ ...form, content: e.target.value })}
                  className="input"
                  rows={2}
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">提醒时间</label>
                <input
                  type="datetime-local"
                  value={form.remind_time}
                  onChange={(e) => setForm({ ...form, remind_time: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">重复</label>
                <select
                  value={form.repeat_type}
                  onChange={(e) => setForm({ ...form, repeat_type: e.target.value })}
                  className="input"
                >
                  <option value="none">不重复</option>
                  <option value="daily">每天</option>
                  <option value="weekly">每周</option>
                  <option value="monthly">每月</option>
                  <option value="yearly">每年</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">关联分类</label>
                <select
                  value={form.category_id}
                  onChange={(e) => setForm({ ...form, category_id: e.target.value })}
                  className="input"
                >
                  <option value="">无</option>
                  {categories.filter((c) => c.type === 'expense').map((cat) => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
              </div>
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={form.is_active}
                  onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                />
                <span className="text-sm">启用提醒</span>
              </label>
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
