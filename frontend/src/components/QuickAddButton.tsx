import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { transactionApi, accountApi, categoryApi, tagApi } from '../services/api'
import { PlusIcon, XMarkIcon } from '@heroicons/react/24/solid'

export default function QuickAddButton() {
  const [isOpen, setIsOpen] = useState(false)
  const [form, setForm] = useState({
    type: 'expense',
    account_id: '',
    category_id: '',
    amount: 0,
    merchant: '',
    note: '',
    bill_date: new Date().toISOString().split('T')[0],
    currency: 'CNY',
  })
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const queryClient = useQueryClient()

  const { data: accountsData } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => accountApi.list(),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(),
  })

  const { data: tagsData } = useQuery({
    queryKey: ['tags'],
    queryFn: () => tagApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: (data: any) => transactionApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['summary'] })
      setIsOpen(false)
      setForm({
        type: 'expense',
        account_id: accountsData?.data?.data[0]?.id || '',
        category_id: '',
        amount: 0,
        merchant: '',
        note: '',
        bill_date: new Date().toISOString().split('T')[0],
      })
    },
  })

  const accounts = accountsData?.data?.data || []
  const categories = categoriesData?.data?.data || []
  const expenseCategories = categories.filter((c: any) => c.type === 'expense')
  const incomeCategories = categories.filter((c: any) => c.type === 'income')
  const tags = tagsData?.data?.data || []

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const submitData = {
      ...form,
      account_id: form.account_id || undefined,
      category_id: form.category_id || undefined,
      tags: selectedTags,
    }
    createMutation.mutate(submitData)
  }

  const toggleTag = (tagId: string) => {
    setSelectedTags(prev => 
      prev.includes(tagId) 
        ? prev.filter(t => t !== tagId)
        : [...prev, tagId]
    )
  }

  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="fixed bottom-20 right-6 md:bottom-6 w-14 h-14 bg-primary-600 text-white rounded-full shadow-lg hover:bg-primary-700 transition-all flex items-center justify-center z-40"
      >
        <PlusIcon className="w-8 h-8" />
      </button>

      {isOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-end md:items-center justify-center z-50">
          <div className="bg-white rounded-t-2xl md:rounded-xl w-full md:max-w-md p-4 md:p-6 animate-slide-up md:animate-none">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">快速记账</h3>
              <button onClick={() => setIsOpen(false)} className="p-1">
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="flex gap-2">
                {['expense', 'income', 'transfer'].map(t => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => setForm({ ...form, type: t })}
                    className={`flex-1 py-2 rounded-lg text-sm font-medium ${
                      form.type === t
                        ? t === 'expense' ? 'bg-red-500 text-white' : 
                          t === 'income' ? 'bg-green-600 text-white' : 'bg-blue-500 text-white'
                        : 'bg-gray-100 text-gray-600'
                    }`}
                  >
                    {t === 'expense' ? '支出' : t === 'income' ? '收入' : '转账'}
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
                  className="w-full text-3xl font-bold text-center py-2 border-b-2 border-gray-200 focus:border-primary-500 outline-none"
                  autoFocus
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
                <input
                  type="date"
                  value={form.bill_date}
                  onChange={(e) => setForm({ ...form, bill_date: e.target.value })}
                  className="input"
                  required
                />
                <input
                  type="text"
                  placeholder="商家(可选)"
                  value={form.merchant}
                  onChange={(e) => setForm({ ...form, merchant: e.target.value })}
                  className="input"
                />
              </div>

              <input
                type="text"
                placeholder="备注(可选)"
                value={form.note}
                onChange={(e) => setForm({ ...form, note: e.target.value })}
                className="input"
              />

              {tags.length > 0 && (
                <div>
                  <p className="text-sm text-gray-500 mb-2">标签</p>
                  <div className="flex flex-wrap gap-2">
                    {tags.map((tag: any) => (
                      <button
                        key={tag.id}
                        type="button"
                        onClick={() => toggleTag(tag.id)}
                        className={`px-3 py-1 rounded-full text-xs transition-colors ${
                          selectedTags.includes(tag.id)
                            ? 'bg-primary-600 text-white'
                            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                        }`}
                      >
                        {tag.name}
                      </button>
                    ))}
                  </div>
                </div>
              )}

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
    </>
  )
}
