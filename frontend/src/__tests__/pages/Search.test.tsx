import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Search from '../../pages/Search'

// Mock react-query
vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
}))

// Mock react-router-dom
vi.mock('react-router-dom', () => ({
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to} data-testid="link">
      {children}
    </a>
  ),
}))

// Mock API
vi.mock('../../services/api', () => ({
  searchApi: {
    search: vi.fn(),
  },
}))

import { useQuery } from '@tanstack/react-query'

const mockUseQuery = useQuery as any

describe('Search Component', () => {
  it('renders search title', () => {
    mockUseQuery.mockReturnValue({
      data: undefined,
    })

    render(<Search />)

    expect(screen.getByText('搜索')).toBeInTheDocument()
  })

  it('renders search input with placeholder', () => {
    mockUseQuery.mockReturnValue({
      data: undefined,
    })

    render(<Search />)

    const input = screen.getByPlaceholderText('搜索交易、账户、分类...')
    expect(input).toBeInTheDocument()
  })

  it('renders search type selector', () => {
    mockUseQuery.mockReturnValue({
      data: undefined,
    })

    render(<Search />)

    const select = screen.getByRole('combobox')
    expect(select).toBeInTheDocument()
  })

  it('shows empty state when no keyword', () => {
    mockUseQuery.mockReturnValue({
      data: undefined,
    })

    render(<Search />)

    expect(screen.getByText('输入关键词开始搜索')).toBeInTheDocument()
  })

  it('updates keyword on input change', () => {
    mockUseQuery.mockReturnValue({
      data: undefined,
    })

    render(<Search />)

    const input = screen.getByPlaceholderText('搜索交易、账户、分类...')
    fireEvent.change(input, { target: { value: 'test' } })

    expect(input).toHaveValue('test')
  })

  it('updates search type on select change', () => {
    mockUseQuery.mockReturnValue({
      data: undefined,
    })

    render(<Search />)

    const select = screen.getByRole('combobox')
    fireEvent.change(select, { target: { value: 'transactions' } })

    expect(select).toHaveValue('transactions')
  })

  it('renders search results when keyword is provided', () => {
    mockUseQuery.mockReturnValue({
      data: {
        data: {
          data: {
            transactions: [
              {
                id: '1',
                amount: 100,
                type: 'expense',
                bill_date: '2024-01-01',
                account: { name: 'Test Account' },
                category: { name: 'Test Category' },
              },
            ],
            accounts: [],
            categories: [],
          },
        },
      },
    })

    render(<Search />)

    // Need to set keyword to trigger search
    const input = screen.getByPlaceholderText('搜索交易、账户、分类...')
    fireEvent.change(input, { target: { value: 'test' } })

    expect(screen.getByText('交易记录 (1)')).toBeInTheDocument()
  })

  it('renders accounts in search results', () => {
    mockUseQuery.mockReturnValue({
      data: {
        data: {
          data: {
            transactions: [],
            accounts: [
              {
                id: '1',
                name: 'Test Account',
                balance: 1000,
              },
            ],
            categories: [],
          },
        },
      },
    })

    render(<Search />)

    const input = screen.getByPlaceholderText('搜索交易、账户、分类...')
    fireEvent.change(input, { target: { value: 'test' } })

    expect(screen.getByText('账户 (1)')).toBeInTheDocument()
    expect(screen.getByText('Test Account')).toBeInTheDocument()
  })

  it('renders categories in search results', () => {
    mockUseQuery.mockReturnValue({
      data: {
        data: {
          data: {
            transactions: [],
            accounts: [],
            categories: [
              {
                id: '1',
                name: 'Test Category',
                color: '#FF0000',
              },
            ],
          },
        },
      },
    })

    render(<Search />)

    const input = screen.getByPlaceholderText('搜索交易、账户、分类...')
    fireEvent.change(input, { target: { value: 'test' } })

    expect(screen.getByText('分类 (1)')).toBeInTheDocument()
    expect(screen.getByText('Test Category')).toBeInTheDocument()
  })
})