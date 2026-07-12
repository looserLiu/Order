import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Dashboard from '../../pages/Dashboard'

// Mock react-query
vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
}))

// Mock recharts
vi.mock('recharts', () => ({
  PieChart: ({ children }: { children: React.ReactNode }) => <div data-testid="pie-chart">{children}</div>,
  Pie: () => <div data-testid="pie" />,
  Cell: () => <div data-testid="cell" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  Tooltip: () => <div data-testid="tooltip" />,
  ResponsiveContainer: ({ children }: { children: React.ReactNode }) => <div data-testid="responsive-container">{children}</div>,
  BarChart: ({ children }: { children: React.ReactNode }) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
}))

// Mock react-router-dom
vi.mock('react-router-dom', () => ({
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to} data-testid="link">
      {children}
    </a>
  ),
}))

// Mock BudgetAlerts
vi.mock('../../components/BudgetAlerts', () => ({
  default: () => <div data-testid="budget-alerts">Budget Alerts</div>,
}))

// Mock API
vi.mock('../../services/api', () => ({
  reportApi: {
    summary: vi.fn(),
    trend: vi.fn(),
    byCategory: vi.fn(),
    monthlyCompare: vi.fn(),
  },
  accountApi: {
    list: vi.fn(),
  },
  transactionApi: {
    list: vi.fn(),
  },
  budgetApi: {
    list: vi.fn(),
  },
}))

import { useQuery } from '@tanstack/react-query'

const mockUseQuery = useQuery as any

describe('Dashboard Component', () => {
  it('renders summary cards with default values', () => {
    mockUseQuery.mockReturnValue({
      data: undefined,
    })

    render(<Dashboard />)

    expect(screen.getByText('本月收入')).toBeInTheDocument()
    expect(screen.getByText('¥0.00')).toBeInTheDocument()
    expect(screen.getByText('本月支出')).toBeInTheDocument()
    expect(screen.getByText('账户余额')).toBeInTheDocument()
    expect(screen.getByText('本月结余')).toBeInTheDocument()
  })

  it('renders with summary data', () => {
    mockUseQuery.mockImplementation(({ queryKey }: { queryKey: string[] }) => {
      if (queryKey[0] === 'summary') {
        return {
          data: {
            data: {
              income: 5000,
              expense: 3000,
            },
          },
        }
      }
      return { data: undefined }
    })

    render(<Dashboard />)

    expect(screen.getByText('¥5000.00')).toBeInTheDocument()
    expect(screen.getByText('¥3000.00')).toBeInTheDocument()
  })

  it('renders budget progress section when budgets exist', () => {
    mockUseQuery.mockImplementation(({ queryKey }: { queryKey: string[] }) => {
      if (queryKey[0] === 'budgets') {
        return {
          data: {
            data: [
              {
                id: '1',
                amount: 5000,
                category: { name: 'Food' },
              },
            ],
          },
        }
      }
      if (queryKey[0] === 'category') {
        return {
          data: {
            data: [
              {
                category_id: '1',
                category_name: 'Food',
                total: 2000,
              },
            ],
          },
        }
      }
      return { data: undefined }
    })

    render(<Dashboard />)

    expect(screen.getByText('预算进度')).toBeInTheDocument()
  })

  it('renders charts', () => {
    mockUseQuery.mockReturnValue({
      data: {
        data: [],
      },
    })

    render(<Dashboard />)

    expect(screen.getByText('收支趋势')).toBeInTheDocument()
    expect(screen.getByText('支出分类')).toBeInTheDocument()
  })

  it('renders recent transactions section', () => {
    mockUseQuery.mockImplementation(({ queryKey }: { queryKey: string[] }) => {
      if (queryKey[0] === 'transactions') {
        return {
          data: {
            data: {
              list: [
                {
                  id: '1',
                  type: 'expense',
                  amount: 100,
                  category: { name: 'Food' },
                  bill_date: '2024-01-01',
                },
              ],
            },
          },
        }
      }
      return { data: undefined }
    })

    render(<Dashboard />)

    expect(screen.getByText('最近交易')).toBeInTheDocument()
  })

  it('renders account overview section', () => {
    mockUseQuery.mockImplementation(({ queryKey }: { queryKey: string[] }) => {
      if (queryKey[0] === 'accounts') {
        return {
          data: {
            data: [
              {
                id: '1',
                name: 'Test Account',
                type: 'bank',
                balance: 1000,
              },
            ],
          },
        }
      }
      return { data: undefined }
    })

    render(<Dashboard />)

    expect(screen.getByText('账户概览')).toBeInTheDocument()
  })

  it('shows empty state when no transactions', () => {
    mockUseQuery.mockImplementation(({ queryKey }: { queryKey: string[] }) => {
      if (queryKey[0] === 'transactions') {
        return {
          data: {
            data: {
              list: [],
            },
          },
        }
      }
      return { data: undefined }
    })

    render(<Dashboard />)

    expect(screen.getByText('暂无交易记录')).toBeInTheDocument()
  })

  it('shows empty state when no accounts', () => {
    mockUseQuery.mockImplementation(({ queryKey }: { queryKey: string[] }) => {
      if (queryKey[0] === 'accounts') {
        return {
          data: {
            data: [],
          },
        }
      }
      return { data: undefined }
    })

    render(<Dashboard />)

    expect(screen.getByText('暂无账户')).toBeInTheDocument()
  })
})