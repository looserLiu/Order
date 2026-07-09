import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { notificationApi } from '../services/api'
import { BellIcon, CheckIcon, TrashIcon } from '@heroicons/react/24/outline'

export default function Notifications() {
  const [showUnreadOnly, setShowUnreadOnly] = useState(false)
  const queryClient = useQueryClient()

  const { data } = useQuery({
    queryKey: ['notifications', showUnreadOnly],
    queryFn: () => notificationApi.list({ unread: showUnreadOnly ? 'true' : 'false' }),
  })

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) => notificationApi.markAsRead(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['notifications'] }),
  })

  const markAllAsReadMutation = useMutation({
    mutationFn: () => notificationApi.markAllAsRead(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['notifications'] }),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => notificationApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['notifications'] }),
  })

  const notifications = data?.data?.data?.list || []
  const unreadCount = data?.data?.data?.unread_count || 0

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'budget_alert': return 'bg-yellow-100 text-yellow-700'
      case 'reminder': return 'bg-blue-100 text-blue-700'
      default: return 'bg-gray-100 text-gray-700'
    }
  }

  const formatTime = (time: string) => {
    const date = new Date(time)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (minutes < 1) return '刚刚'
    if (minutes < 60) return `${minutes}分钟前`
    if (hours < 24) return `${hours}小时前`
    if (days < 7) return `${days}天前`
    return date.toLocaleDateString()
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <BellIcon className="w-6 h-6" />
          <h2 className="text-2xl font-bold">消息通知</h2>
          {unreadCount > 0 && (
            <span className="bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">
              {unreadCount}
            </span>
          )}
        </div>
        <div className="flex gap-2">
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={showUnreadOnly}
              onChange={(e) => setShowUnreadOnly(e.target.checked)}
              className="rounded"
            />
            只看未读
          </label>
          {unreadCount > 0 && (
            <button
              onClick={() => markAllAsReadMutation.mutate()}
              className="btn-secondary text-sm"
            >
              全部已读
            </button>
          )}
        </div>
      </div>

      <div className="space-y-2">
        {notifications.map((notification: any) => (
          <div
            key={notification.id}
            className={`card flex items-start gap-3 ${!notification.is_read ? 'bg-primary-50' : ''}`}
          >
            <div className={`px-2 py-1 rounded text-xs ${getTypeColor(notification.type)}`}>
              {notification.type === 'budget_alert' ? '预算提醒' : 
               notification.type === 'reminder' ? '定时提醒' : '系统通知'}
            </div>
            <div className="flex-1 min-w-0">
              <p className="font-medium">{notification.title}</p>
              {notification.content && (
                <p className="text-sm text-gray-500 mt-1">{notification.content}</p>
              )}
              <p className="text-xs text-gray-400 mt-2">
                {formatTime(notification.created_at)}
              </p>
            </div>
            <div className="flex gap-1">
              {!notification.is_read && (
                <button
                  onClick={() => markAsReadMutation.mutate(notification.id)}
                  className="p-2 hover:bg-gray-100 rounded"
                  title="标记已读"
                >
                  <CheckIcon className="w-4 h-4 text-gray-500" />
                </button>
              )}
              <button
                onClick={() => deleteMutation.mutate(notification.id)}
                className="p-2 hover:bg-gray-100 rounded"
                title="删除"
              >
                <TrashIcon className="w-4 h-4 text-red-500" />
              </button>
            </div>
          </div>
        ))}

        {notifications.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <BellIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
            <p>暂无消息通知</p>
          </div>
        )}
      </div>
    </div>
  )
}
