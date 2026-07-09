import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { transactionApi, reportApi } from '../services/api'
import { 
  format, 
  startOfMonth, 
  endOfMonth, 
  eachDayOfInterval, 
  isSameDay, 
  isSameMonth,
  addMonths,
  subMonths,
  startOfWeek,
  endOfWeek
} from 'date-fns'
import { zhCN } from 'date-fns/locale'
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/outline'

export default function Calendar() {
  const [currentDate, setCurrentDate] = useState(new Date())
  const [selectedDate, setSelectedDate] = useState<string | null>(null)

  const start = startOfWeek(startOfMonth(currentDate), { weekStartsOn: 0 })
  const end = endOfWeek(endOfMonth(currentDate), { weekStartsOn: 0 })
  const days = eachDayOfInterval({ start, end })

  const monthStart = startOfMonth(currentDate).toISOString().split('T')[0]
  const monthEnd = endOfMonth(currentDate).toISOString().split('T')[0]

  const { data: txData } = useQuery({
    queryKey: ['transactions', monthStart, monthEnd],
    queryFn: () => transactionApi.list({ start_date: monthStart, end_date: monthEnd }),
  })

  const transactions = txData?.data?.data?.list || []

  const getDayData = (day: Date) => {
    const dateStr = format(day, 'yyyy-MM-dd')
    const dayTx = transactions.filter((t: any) => t.bill_date === dateStr)
    const income = dayTx.filter((t: any) => t.type === 'income').reduce((sum: number, t: any) => sum + t.amount, 0)
    const expense = dayTx.filter((t: any) => t.type === 'expense').reduce((sum: number, t: any) => sum + t.amount, 0)
    return { income, expense, count: dayTx.length }
  }

  const selectedDateTx = selectedDate 
    ? transactions.filter((t: any) => t.bill_date === selectedDate)
    : []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">日历视图</h2>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setCurrentDate(subMonths(currentDate, 1))}
            className="p-2 hover:bg-gray-100 rounded-lg"
          >
            <ChevronLeftIcon className="w-5 h-5" />
          </button>
          <span className="font-medium min-w-[120px] text-center">
            {format(currentDate, 'yyyy年MM月', { locale: zhCN })}
          </span>
          <button
            onClick={() => setCurrentDate(addMonths(currentDate, 1))}
            className="p-2 hover:bg-gray-100 rounded-lg"
          >
            <ChevronRightIcon className="w-5 h-5" />
          </button>
        </div>
      </div>

      <div className="grid grid-cols-7 gap-1">
        {['日', '一', '二', '三', '四', '五', '六'].map((day) => (
          <div key={day} className="text-center text-sm font-medium text-gray-500 py-2">
            {day}
          </div>
        ))}
        
        {days.map((day) => {
          const dayData = getDayData(day)
          const dateStr = format(day, 'yyyy-MM-dd')
          const isCurrentMonth = isSameMonth(day, currentDate)
          const isToday = isSameDay(day, new Date())
          const isSelected = selectedDate === dateStr

          return (
            <button
              key={dateStr}
              onClick={() => setSelectedDate(dateStr)}
              className={`
                min-h-[80px] p-2 rounded-lg text-left transition-colors
                ${isCurrentMonth ? 'bg-white' : 'bg-gray-50'}
                ${isSelected ? 'ring-2 ring-primary-500' : ''}
                ${isToday ? 'bg-primary-50' : ''}
                hover:bg-gray-50
              `}
            >
              <div className={`text-sm ${isCurrentMonth ? 'text-gray-900' : 'text-gray-400'} ${isToday ? 'font-bold text-primary-600' : ''}`}>
                {format(day, 'd')}
              </div>
              {dayData.count > 0 && (
                <div className="mt-1 space-y-1">
                  {dayData.income > 0 && (
                    <div className="text-xs text-green-600">+¥{dayData.income.toFixed(0)}</div>
                  )}
                  {dayData.expense > 0 && (
                    <div className="text-xs text-red-500">-¥{dayData.expense.toFixed(0)}</div>
                  )}
                </div>
              )}
            </button>
          )
        })}
      </div>

      {selectedDate && (
        <div className="card">
          <h3 className="font-semibold mb-4">{selectedDate} 记账明细</h3>
          <div className="space-y-2">
            {selectedDateTx.map((tx: any) => (
              <div key={tx.id} className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <div>
                  <p className="font-medium">{tx.category?.name || '未分类'}</p>
                  <p className="text-sm text-gray-500">{tx.note || tx.merchant || '-'}</p>
                </div>
                <p className={`font-semibold ${tx.type === 'income' ? 'text-green-600' : 'text-red-500'}`}>
                  {tx.type === 'income' ? '+' : '-'}¥{tx.amount?.toFixed(2)}
                </p>
              </div>
            ))}
            {selectedDateTx.length === 0 && (
              <p className="text-center text-gray-500 py-4">当天无记录</p>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
