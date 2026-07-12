import { useState, ChangeEvent } from 'react'
import { csvApi, ImportResult } from '../services/api'
import { useQueryClient } from '@tanstack/react-query'
import { ArrowUpTrayIcon, DocumentIcon, CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline'

// Parse result type
export interface ParseResult {
  date: string
  type: string
  amount: string
  category: string
  account: string
  merchant: string
  note: string
  valid: boolean
  error?: string
}

// CSV row type for import
interface CSVRow {
  date: string
  type: string
  amount: string
  category: string
  account: string
  merchant: string
  note: string
}

export default function CSVImport() {
  const [file, setFile] = useState<File | null>(null)
  const [parsedData, setParsedData] = useState<ParseResult[]>([])
  const [importing, setImporting] = useState(false)
  const [result, setResult] = useState<ImportResult | null>(null)
  const queryClient = useQueryClient()

  const parseCSV = (content: string): ParseResult[] => {
    const lines = content.trim().split('\n')
    const results: ParseResult[] = []

    // Skip header if present
    const startIndex = lines[0]?.toLowerCase().includes('date') ? 1 : 0

    for (let i = startIndex; i < lines.length; i++) {
      const line = lines[i].trim()
      if (!line) continue

      // Handle both CSV and tab-separated formats
      const parts = line.includes(',') 
        ? line.split(',').map(p => p.trim().replace(/^"|"$/g, ''))
        : line.split('\t').map(p => p.trim())

      const result: ParseResult = {
        date: parts[0] || '',
        type: parts[1] || 'expense',
        amount: parts[2] || '0',
        category: parts[3] || '',
        account: parts[4] || '',
        merchant: parts[5] || '',
        note: parts[6] || '',
        valid: true
      }

      // Validate
      if (!result.date) {
        result.valid = false
        result.error = 'Missing date'
      } else if (!result.amount || isNaN(parseFloat(result.amount))) {
        result.valid = false
        result.error = 'Invalid amount'
      }

      results.push(result)
    }

    return results
  }

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (!selectedFile) return

    setFile(selectedFile)
    setResult(null)

    const reader = new FileReader()
    reader.onload = (event) => {
      const content = event.target?.result as string
      const parsed = parseCSV(content)
      setParsedData(parsed)
    }
    reader.readAsText(selectedFile)
  }

  const handleImport = async () => {
    if (parsedData.length === 0) return

    setImporting(true)
    try {
      const validData: CSVRow[] = parsedData
        .filter(d => d.valid)
        .map(d => ({
          date: d.date,
          type: d.type,
          amount: d.amount,
          category: d.category,
          account: d.account,
          merchant: d.merchant,
          note: d.note
        }))

      const { data } = await csvApi.importCSV({ transactions: validData })
      setResult(data.data)
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
    } catch (error) {
      console.error('Import failed:', error)
    } finally {
      setImporting(false)
    }
  }

  const sampleCSV = `日期,类型,金额,分类,账户,商家,备注
2024-01-15,支出,50.00,餐饮,现金,麦当劳,午餐
2024-01-16,收入,5000.00,工资,银行卡,公司,月薪
2024-01-17,支出,120.00,购物,支付宝,淘宝,日用品`

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2">
        <ArrowUpTrayIcon className="w-6 h-6" />
        <h2 className="text-2xl font-bold">CSV 导入</h2>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">导入说明</h3>
        <div className="text-sm text-gray-600 space-y-2">
          <p>1. 准备 CSV 文件，支持以下格式：</p>
          <div className="bg-gray-50 p-3 rounded-lg overflow-x-auto">
            <code className="text-xs">日期,类型,金额,分类,账户,商家,备注</code>
          </div>
          <p>2. 类型支持：支出/收入/转账 (或 expense/income/transfer)</p>
          <p>3. 日期格式：YYYY-MM-DD 或 YYYY/MM/DD</p>
        </div>

        <div className="mt-4">
          <label className="block text-sm font-medium mb-2">下载示例模板</label>
          <button
            onClick={() => {
              const blob = new Blob([sampleCSV], { type: 'text/csv' })
              const url = URL.createObjectURL(blob)
              const a = document.createElement('a')
              a.href = url
              a.download = 'template.csv'
              a.click()
            }}
            className="btn-secondary text-sm"
          >
            下载 CSV 模板
          </button>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-semibold mb-4">选择文件</h3>
        <input
          type="file"
          accept=".csv,.txt"
          onChange={handleFileChange}
          className="w-full border border-gray-300 rounded-lg p-2"
        />
      </div>

      {parsedData.length > 0 && (
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">
              预览 ({parsedData.length} 条记录)
            </h3>
            <button
              onClick={handleImport}
              disabled={importing || parsedData.filter(d => d.valid).length === 0}
              className="btn-primary"
            >
              {importing ? '导入中...' : '开始导入'}
            </button>
          </div>

          <div className="overflow-x-auto max-h-96">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 sticky top-0">
                <tr>
                  <th className="px-3 py-2 text-left">状态</th>
                  <th className="px-3 py-2 text-left">日期</th>
                  <th className="px-3 py-2 text-left">类型</th>
                  <th className="px-3 py-2 text-right">金额</th>
                  <th className="px-3 py-2 text-left">分类</th>
                  <th className="px-3 py-2 text-left">账户</th>
                  <th className="px-3 py-2 text-left">备注</th>
                </tr>
              </thead>
              <tbody>
                {parsedData.slice(0, 20).map((row, i) => (
                  <tr key={i} className="border-b">
                    <td className="px-3 py-2">
                      {row.valid ? (
                        <CheckCircleIcon className="w-5 h-5 text-green-500" />
                      ) : (
                        <XCircleIcon className="w-5 h-5 text-red-500" title={row.error} />
                      )}
                    </td>
                    <td className="px-3 py-2">{row.date}</td>
                    <td className="px-3 py-2">
                      <span className={`px-2 py-0.5 rounded text-xs ${
                        row.type === '收入' || row.type === 'income' ? 'bg-green-100 text-green-700' :
                        row.type === '支出' || row.type === 'expense' ? 'bg-red-100 text-red-700' : 'bg-blue-100 text-blue-700'
                      }`}>
                        {row.type}
                      </span>
                    </td>
                    <td className="px-3 py-2 text-right font-medium">{row.amount}</td>
                    <td className="px-3 py-2">{row.category}</td>
                    <td className="px-3 py-2">{row.account}</td>
                    <td className="px-3 py-2 text-gray-500">{row.note}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            {parsedData.length > 20 && (
              <p className="text-center py-2 text-gray-500">
                ... 还有 {parsedData.length - 20} 条记录
              </p>
            )}
          </div>

          <div className="mt-4 flex gap-4 text-sm">
            <span className="flex items-center gap-1">
              <CheckCircleIcon className="w-4 h-4 text-green-500" />
              有效: {parsedData.filter(d => d.valid).length}
            </span>
            <span className="flex items-center gap-1">
              <XCircleIcon className="w-4 h-4 text-red-500" />
              无效: {parsedData.filter(d => !d.valid).length}
            </span>
          </div>
        </div>
      )}

      {result && (
        <div className="card bg-green-50 border-green-200">
          <h3 className="text-lg font-semibold text-green-700 mb-2">导入完成</h3>
          <p className="text-green-600">
            成功导入: {result.imported} 条记录
            {result.failed > 0 && `，失败: ${result.failed} 条`}
          </p>
        </div>
      )}
    </div>
  )
}
