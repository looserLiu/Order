import { useQuery } from '@tanstack/react-query'
import { budgetAlertApi, BudgetAlert } from '../services/api'
import { ExclamationTriangleIcon, ExclamationCircleIcon } from '@heroicons/react/24/outline'
import { LoadingSpinner } from './LoadingSpinner'

export default function BudgetAlerts() {
  const { data: alerts, isLoading } = useQuery({
    queryKey: ['budgetAlerts'],
    queryFn: async () => {
      const res = await budgetAlertApi.getAlerts()
      return res.data.data
    },
  })

  if (isLoading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <div className="animate-pulse space-y-3">
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/3"></div>
          <div className="h-16 bg-gray-200 dark:bg-gray-700 rounded"></div>
        </div>
      </div>
    )
  }

  if (!alerts || alerts.length === 0) {
    return null
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
      <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
        <ExclamationTriangleIcon className="w-5 h-5 text-yellow-500" />
        预算提醒
      </h2>
      <div className="space-y-3">
        {alerts.map((alert: BudgetAlert) => (
          <div
            key={alert.budget_id}
            className={`p-4 rounded-lg border ${
              alert.alert_type === 'exceeded'
                ? 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800'
                : 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800'
            }`}
          >
            <div className="flex items-start gap-3">
              {alert.alert_type === 'exceeded' ? (
                <ExclamationCircleIcon className="w-5 h-5 text-red-500 mt-0.5" />
              ) : (
                <ExclamationTriangleIcon className="w-5 h-5 text-yellow-500 mt-0.5" />
              )}
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <span className="font-medium text-gray-900 dark:text-white">
                    {alert.category_name}
                  </span>
                  <span
                    className={`text-sm font-medium ${
                      alert.alert_type === 'exceeded'
                        ? 'text-red-600 dark:text-red-400'
                        : 'text-yellow-600 dark:text-yellow-400'
                    }`}
                  >
                    {alert.percentage.toFixed(0)}%
                  </span>
                </div>
                <div className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  预算: ¥{alert.budget_amount.toFixed(2)} | 已花费: ¥{alert.spent_amount.toFixed(2)}
                </div>
                <div className="mt-2 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full ${
                      alert.alert_type === 'exceeded' ? 'bg-red-500' : 'bg-yellow-500'
                    }`}
                    style={{ width: `${Math.min(alert.percentage, 100)}%` }}
                  ></div>
                </div>
                <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  {alert.alert_type === 'exceeded'
                    ? `已超出预算 ¥${Math.abs(alert.remaining).toFixed(2)}`
                    : `剩余 ¥${alert.remaining.toFixed(2)}`}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
