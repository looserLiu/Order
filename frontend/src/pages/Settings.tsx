import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { userApi, authApi, backupApi } from '../services/api'
import { useAuthStore } from '../stores/authStore'
import { useTranslation } from 'react-i18next'
import ThemeSelector from '../components/ThemeSelector'
import { LanguageIcon } from '@heroicons/react/24/outline'

export default function Settings() {
  const { user, logout } = useAuthStore()
  const { t, i18n } = useTranslation()
  
  const languages = [
    { code: 'zh', name: '中文' },
    { code: 'en', name: 'English' },
    { code: 'ja', name: '日本語' },
  ]

  const changeLanguage = (lang: string) => {
    i18n.changeLanguage(lang)
  }
  const queryClient = useQueryClient()
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState({
    nickname: user?.nickname || '',
    currency: user?.currency || 'CNY',
    timezone: 'Asia/Shanghai',
  })

  const { data } = useQuery({
    queryKey: ['user'],
    queryFn: () => userApi.getMe(),
  })

  const updateMutation = useMutation({
    mutationFn: (data: any) => userApi.updateMe(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] })
      setEditing(false)
    },
  })

  const currentUser = data?.data?.data || user

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    updateMutation.mutate(form)
  }

  const handleLogout = () => {
    logout()
    window.location.href = '/login'
  }

  const [showImportModal, setShowImportModal] = useState(false)
  const [importData, setImportData] = useState('')

  const exportMutation = useMutation({
    mutationFn: () => backupApi.exportAll(),
    onSuccess: (data) => {
      const jsonStr = JSON.stringify(data.data.data, null, 2)
      const blob = new Blob([jsonStr], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `backup_${new Date().toISOString().split('T')[0]}.json`
      a.click()
      URL.revokeObjectURL(url)
    },
  })

  const importMutation = useMutation({
    mutationFn: (data: any) => backupApi.importAll(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      setShowImportModal(false)
      setImportData('')
      alert('导入成功！')
    },
    onError: () => {
      alert('导入失败，请检查数据格式')
    },
  })

  const handleImport = () => {
    try {
      const data = JSON.parse(importData)
      importMutation.mutate(data)
    } catch {
      alert('JSON格式错误')
    }
  }

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (event) => {
        setImportData(event.target?.result as string)
      }
      reader.readAsText(file)
    }
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">设置</h2>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">个人资料</h3>
        <div className="flex items-center gap-4 mb-6">
          <div className="w-20 h-20 bg-primary-100 rounded-full flex items-center justify-center">
            <span className="text-3xl text-primary-600 font-bold">
              {currentUser?.nickname?.[0] || currentUser?.email?.[0] || 'U'}
            </span>
          </div>
          <div>
            <p className="font-medium text-lg">{currentUser?.nickname || '未设置昵称'}</p>
            <p className="text-gray-500">{currentUser?.email}</p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4 max-w-md">
          <div>
            <label className="block text-sm font-medium mb-1">昵称</label>
            <input
              type="text"
              value={form.nickname}
              onChange={(e) => setForm({ ...form, nickname: e.target.value })}
              className="input"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">货币</label>
            <select
              value={form.currency}
              onChange={(e) => setForm({ ...form, currency: e.target.value })}
              className="input"
            >
              <option value="CNY">人民币 (CNY)</option>
              <option value="USD">美元 (USD)</option>
              <option value="EUR">欧元 (EUR)</option>
              <option value="JPY">日元 (JPY)</option>
              <option value="GBP">英镑 (GBP)</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">时区</label>
            <select
              value={form.timezone}
              onChange={(e) => setForm({ ...form, timezone: e.target.value })}
              className="input"
            >
              <option value="Asia/Shanghai">中国标准时间 (UTC+8)</option>
              <option value="America/New_York">美国东部时间 (UTC-5)</option>
              <option value="Europe/London">英国时间 (UTC+0)</option>
              <option value="Asia/Tokyo">日本时间 (UTC+9)</option>
            </select>
          </div>
          <button type="submit" className="btn-primary">
            保存修改
          </button>
        </form>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">外观</h3>
        <ThemeSelector />
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">语言设置</h3>
        <div className="flex items-center justify-between py-2">
          <div className="flex items-center gap-2">
            <LanguageIcon className="w-5 h-5 text-gray-500" />
            <div>
              <p className="font-medium">界面语言</p>
              <p className="text-sm text-gray-500">选择您偏好的语言</p>
            </div>
          </div>
          <select
            value={i18n.language}
            onChange={(e) => changeLanguage(e.target.value)}
            className="input w-auto"
          >
            {languages.map((lang) => (
              <option key={lang.code} value={lang.code}>
                {lang.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">安全设置</h3>
        <div className="space-y-3">
          <div className="flex items-center justify-between py-2">
            <div>
              <p className="font-medium">修改密码</p>
              <p className="text-sm text-gray-500">定期更换密码保护账户安全</p>
            </div>
            <button className="btn-secondary text-sm">修改</button>
          </div>
          <div className="flex items-center justify-between py-2">
            <div>
              <p className="font-medium">双因素认证</p>
              <p className="text-sm text-gray-500">增强账户安全保护</p>
            </div>
            <button className="btn-secondary text-sm">设置</button>
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">数据管理</h3>
        <div className="space-y-3">
          <div className="flex items-center justify-between py-2">
            <div>
              <p className="font-medium">导出数据</p>
              <p className="text-sm text-gray-500">导出所有记账数据</p>
            </div>
            <button 
              onClick={() => exportMutation.mutate()} 
              className="btn-secondary text-sm"
              disabled={exportMutation.isPending}
            >
              {exportMutation.isPending ? '导出中...' : '导出'}
            </button>
          </div>
          <div className="flex items-center justify-between py-2">
            <div>
              <p className="font-medium">导入数据</p>
              <p className="text-sm text-gray-500">从其他记账软件导入</p>
            </div>
            <button 
              onClick={() => setShowImportModal(true)}
              className="btn-secondary text-sm"
            >
              导入
            </button>
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">其他</h3>
        <div className="space-y-3">
          <div className="flex items-center justify-between py-2">
            <div>
              <p className="font-medium">帮助与反馈</p>
              <p className="text-sm text-gray-500">获取帮助或提交问题</p>
            </div>
            <button className="btn-secondary text-sm">查看</button>
          </div>
          <div className="flex items-center justify-between py-2">
            <div>
              <p className="font-medium">关于</p>
              <p className="text-sm text-gray-500">版本 v1.0.0</p>
            </div>
          </div>
        </div>
      </div>

      <button
        onClick={handleLogout}
        className="w-full py-3 bg-red-50 text-red-600 rounded-lg font-medium hover:bg-red-100"
      >
        退出登录
      </button>

      {showImportModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl p-6 w-full max-w-lg">
            <h3 className="text-lg font-semibold mb-4">导入数据</h3>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-2">选择文件</label>
              <input
                type="file"
                accept=".json"
                onChange={handleFileUpload}
                className="w-full border rounded-lg p-2"
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-2">或粘贴JSON数据</label>
              <textarea
                value={importData}
                onChange={(e) => setImportData(e.target.value)}
                className="w-full border rounded-lg p-2 h-40 font-mono text-sm"
                placeholder='{"accounts": [...], "transactions": [...]}'
              />
            </div>
            <div className="flex gap-3">
              <button
                onClick={() => {
                  setShowImportModal(false)
                  setImportData('')
                }}
                className="flex-1 btn-secondary"
              >
                取消
              </button>
              <button
                onClick={handleImport}
                disabled={!importData || importMutation.isPending}
                className="flex-1 btn-primary"
              >
                {importMutation.isPending ? '导入中...' : '确认导入'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
