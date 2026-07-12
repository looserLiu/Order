import { FC, ReactNode } from 'react'
import { 
  InboxIcon, 
  DocumentIcon, 
  ChartBarIcon, 
  CreditCardIcon,
  TagIcon,
  CalendarIcon,
  UserGroupIcon,
  TrophyIcon,
  ShieldCheckIcon,
  BanknotesIcon,
} from '@heroicons/react/24/outline'

// Empty state icon types
export type EmptyStateIcon = 'inbox' | 'document' | 'chart' | 'card' | 'tag' | 'calendar' | 'group' | 'goal' | 'insurance' | 'netWorth'

// Icon mapping
const iconMap: Record<EmptyStateIcon, FC<{ className?: string }>> = {
  inbox: InboxIcon,
  document: DocumentIcon,
  chart: ChartBarIcon,
  card: CreditCardIcon,
  tag: TagIcon,
  calendar: CalendarIcon,
  group: UserGroupIcon,
  goal: TrophyIcon,
  insurance: ShieldCheckIcon,
  netWorth: BanknotesIcon,
}

// EmptyState props
export interface EmptyStateProps {
  icon?: EmptyStateIcon
  title: string
  description?: string
  action?: ReactNode
  className?: string
}

/**
 * EmptyState - A reusable empty state component
 */
export const EmptyState: FC<EmptyStateProps> = ({
  icon = 'inbox',
  title,
  description,
  action,
  className,
}) => {
  const Icon = iconMap[icon]

  return (
    <div className={`flex flex-col items-center justify-center py-12 px-4 text-center ${className || ''}`}>
      <div className="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center mb-4">
        <Icon className="w-8 h-8 text-gray-400" />
      </div>
      <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">{title}</h3>
      {description && (
        <p className="text-sm text-gray-500 dark:text-gray-400 max-w-sm mb-4">{description}</p>
      )}
      {action && <div className="mt-2">{action}</div>}
    </div>
  )
}

export default EmptyState