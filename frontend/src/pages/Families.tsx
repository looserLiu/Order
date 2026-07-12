import { useState, FormEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { familyApi, accountApi, categoryApi, Account, Category, Transaction } from '../services/api'
import { PlusIcon, PencilIcon, TrashIcon, UserGroupIcon, ArrowRightIcon } from '@heroicons/react/24/outline'

// Form state types
interface FamilyForm {
  name: string
}

interface MemberForm {
  email: string
  role: string
}

interface TransactionForm {
  account_id: string
  category_id: string
  type: string
  amount: number
  note: string
  bill_date: string
}

// Family type for API response
interface Family {
  id: string
  name: string
  members?: FamilyMember[]
}

interface FamilyMember {
  id: string
  user?: {
    email?: { nickname?: string }
    name?: string
  }
  role: string
}

export default function Families() {
  const [showModal, setShowModal] = useState(false)
  const [showAddMember, setShowAddMember] = useState<string | null>(null)
  const [showTxModal, setShowTxModal] = useState<string | null>(null)
  const [familyForm, setFamilyForm] = useState<FamilyForm>({ name: '' })
  const [memberForm, setMemberForm] = useState<MemberForm>({ email: '', role: 'member' })
  const [txForm, setTxForm] = useState<TransactionForm>({
    account_id: '',
    category_id: '',
    type: 'expense',
    amount: 0,
    note: '',
    bill_date: new Date().toISOString().split('T')[0],
  })
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['families'],
    queryFn: () => familyApi.list(),
  })

  const { data: accountsData } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => accountApi.list(),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(),
  })

  const { data: txData } = useQuery({
    queryKey: ['family-tx', showTxModal],
    queryFn: () => showTxModal ? familyApi.getTransactions(showTxModal) : null,
    enabled: !!showTxModal,
  })

  const createFamilyMutation = useMutation({
    mutationFn: (data: FamilyForm) => familyApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['families'] })
      closeModal()
    },
  })

  const deleteFamilyMutation = useMutation({
    mutationFn: (id: string) => familyApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['families'] }),
  })

  const addMemberMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: MemberForm }) => familyApi.addMember(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['families'] })
      setShowAddMember(null)
    },
  })

  const removeMemberMutation = useMutation({
    mutationFn: ({ id, memberId }: { id: string; memberId: string }) => familyApi.removeMember(id, memberId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['families'] }),
  })

  const createTxMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Record<string, unknown> }) => familyApi.createTransaction(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['family-tx'] })
      setShowTxModal(null)
    },
  })

  const families = (data?.data.data || []) as Family[]
  const accounts = (accountsData?.data.data || []) as Account[]
  const categories = (categoriesData?.data.data || []) as Category[]
  const transactions = (txData?.data.data || []) as Transaction[]

  const closeModal = () => {
    setShowModal(false)
    setFamilyForm({ name: '' })
  }

  const handleCreateFamily = (e: FormEvent) => {
    e.preventDefault()
    createFamilyMutation.mutate(familyForm)
  }

  const handleAddMember = (e: FormEvent, familyId: string) => {
    e.preventDefault()
    addMemberMutation.mutate({ id: familyId, data: memberForm })
  }

  const handleCreateTx = (e: FormEvent, familyId: string) => {
    e.preventDefault()
    createTxMutation.mutate({
      id: familyId,
      data: {
        ...txForm,
        account_id: txForm.account_id || undefined,
        category_id: txForm.category_id || undefined,
        amount: Number(txForm.amount),
      },
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">家庭/团队记账</h2>
        <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
          <PlusIcon className="w-5 h-5" />
          创建账本
        </button>
      </div>

      <div className="grid gap-4">
        {families.map((family: Family) => (
          <div key={family.id} className="card">
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 bg-primary-100 rounded-xl flex items-center justify-center">
                  <UserGroupIcon className="w-6 h-6 text-primary-600" />
                </div>
                <div>
                  <h3 className="font-semibold text-lg">{family.name}</h3>
                  <p className="text-sm text-gray-500">
                    成员 {family.members?.length || 0} 人
                  </p>
                </div>
              </div>
              <button onClick={() => deleteFamilyMutation.mutate(family.id)} className="p-2 hover:bg-gray-100 rounded-lg">
                <TrashIcon className="w-4 h-4 text-red-500" />
              </button>
            </div>

            <div className="border-t border-gray-100 pt-4 mb-4">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium">成员列表</span>
                <button
                  onClick={() => setShowAddMember(family.id)}
                  className="text-sm text-primary-600 hover:underline"
                >
                  添加成员
                </button>
              </div>
              <div className="flex flex-wrap gap-2">
                {family.members?.map((member: FamilyMember) => (
                  <div
                    key={member.id}
                    className="flex items-center gap-2 px-3 py-1 bg-gray-100 rounded-full text-sm"
                  >
                    <span>{member.user?.email?.nickname || member.user?.name || 'Unknown'}</span>
                    <span className="text-gray-500 text-xs">({member.role})</span>
                    {member.role !== 'owner' && (
                      <button
                        onClick={() => removeMemberMutation.mutate({ id: family.id, memberId: member.id })}
                        className="text-red-500 hover:text-red-700"
                      >
                        ×
                      </button>
                    )}
                  </div>
                ))}
              </div>
            </div>

            <button
              onClick={() => setShowTxModal(family.id)}
              className="w-full btn-secondary flex items-center justify-center gap-2"
            >
              <ArrowRightIcon className="w-4 h-4" />
              查看/添加记录
            </button>
          </div>
        ))}

        {families.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <p>暂无家庭/团队账本</p>
            <p className="text-sm">点击上方按钮创建</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">创建账本</h3>
            <form onSubmit={handleCreateFamily} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">账本名称</label>
                <input
                  type="text"
                  value={familyForm.name}
                  onChange={(e) => setFamilyForm({ ...familyForm, name: e.target.value })}
                  className="input"
                  placeholder="例如：家庭共同账本"
                  required
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={closeModal} className="btn-secondary flex-1">取消</button>
                <button type="submit" className="btn-primary flex-1">创建</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showAddMember && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">添加成员</h3>
            <form onSubmit={(e) => handleAddMember(e, showAddMember)} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">成员邮箱</label>
                <input
                  type="email"
                  value={memberForm.email}
                  onChange={(e) => setMemberForm({ ...memberForm, email: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">角色</label>
                <select
                  value={memberForm.role}
                  onChange={(e) => setMemberForm({ ...memberForm, role: e.target.value })}
                  className="input"
                >
                  <option value="member">成员</option>
                  <option value="admin">管理员</option>
                </select>
              </div>
              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setShowAddMember(null)} className="btn-secondary flex-1">取消</button>
                <button type="submit" className="btn-primary flex-1">添加</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showTxModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-lg max-h-[80vh] overflow-y-auto">
            <h3 className="text-lg font-semibold mb-4">账本记录</h3>
            
            <form onSubmit={(e) => handleCreateTx(e, showTxModal)} className="space-y-3 mb-6 pb-6 border-b">
              <p className="font-medium text-sm">添加新记录</p>
              <div className="flex gap-2">
                {['expense', 'income'].map(t => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => setTxForm({ ...txForm, type: t })}
                    className={`flex-1 py-2 rounded-lg text-sm font-medium ${
                      txForm.type === t
                        ? t === 'expense' ? 'bg-red-500 text-white' : 'bg-green-600 text-white'
                        : 'bg-gray-100 text-gray-600'
                    }`}
                  >
                    {t === 'expense' ? '支出' : '收入'}
                  </button>
                ))}
              </div>
              <div className="grid grid-cols-2 gap-2">
                <input
                  type="number"
                  placeholder="金额"
                  value={txForm.amount}
                  onChange={(e) => setTxForm({ ...txForm, amount: parseFloat(e.target.value) || 0 })}
                  className="input"
                  required
                />
                <input
                  type="date"
                  value={txForm.bill_date}
                  onChange={(e) => setTxForm({ ...txForm, bill_date: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <select
                value={txForm.account_id}
                onChange={(e) => setTxForm({ ...txForm, account_id: e.target.value })}
                className="input"
              >
                <option value="">选择账户</option>
                {accounts.map((acc: Account) => (
                  <option key={acc.id} value={acc.id}>{acc.name}</option>
                ))}
              </select>
              <input
                type="text"
                placeholder="备注"
                value={txForm.note}
                onChange={(e) => setTxForm({ ...txForm, note: e.target.value })}
                className="input"
              />
              <button type="submit" className="btn-primary w-full">添加记录</button>
            </form>

            <div className="space-y-2">
              {transactions.map((tx: Transaction) => (
                <div key={tx.id} className="flex items-center justify-between py-2 border-b border-gray-100">
                  <div>
                    <p className="font-medium">{tx.category?.name || '未分类'}</p>
                    <p className="text-xs text-gray-500">{tx.bill_date}</p>
                  </div>
                  <p className={`font-semibold ${tx.type === 'income' ? 'text-green-600' : 'text-red-500'}`}>
                    {tx.type === 'income' ? '+' : '-'}¥{tx.amount?.toFixed(2)}
                  </p>
                </div>
              ))}
              {transactions.length === 0 && (
                <p className="text-center text-gray-500 py-4">暂无记录</p>
              )}
            </div>

            <button
              onClick={() => setShowTxModal(null)}
              className="btn-secondary w-full mt-4"
            >
              关闭
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
