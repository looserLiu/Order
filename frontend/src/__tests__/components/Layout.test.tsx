import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Layout from '../../components/Layout'
import { useThemeStore } from '../../stores/themeStore'

// Mock the theme store
vi.mock('../../stores/themeStore', () => ({
  useThemeStore: vi.fn(),
}))

// Mock Outlet
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    Outlet: () => <div data-testid="outlet">Outlet Content</div>,
    NavLink: ({ children, to, end }: { children: React.ReactNode; to: string; end?: boolean }) => (
      <a href={to} data-testid={`navlink-${to}`}>
        {children}
      </a>
    ),
  }
})

describe('Layout Component', () => {
  it('renders the application title', () => {
    ;(useThemeStore as any).mockReturnValue({
      theme: 'light',
      toggleTheme: vi.fn(),
    })

    render(<Layout />)

    expect(screen.getByText('智慧记账')).toBeInTheDocument()
  })

  it('renders navigation items', () => {
    ;(useThemeStore as any).mockReturnValue({
      theme: 'light',
      toggleTheme: vi.fn(),
    })

    render(<Layout />)

    // Check for main navigation items
    expect(screen.getByText('首页')).toBeInTheDocument()
    expect(screen.getByText('账户')).toBeInTheDocument()
    expect(screen.getByText('记账')).toBeInTheDocument()
    expect(screen.getByText('预算')).toBeInTheDocument()
  })

  it('renders theme toggle button', () => {
    ;(useThemeStore as any).mockReturnValue({
      theme: 'light',
      toggleTheme: vi.fn(),
    })

    render(<Layout />)

    // Theme toggle button should be present
    const themeButton = screen.getByRole('button')
    expect(themeButton).toBeInTheDocument()
  })

  it('applies dark theme class when theme is dark', () => {
    const toggleTheme = vi.fn()
    ;(useThemeStore as any).mockReturnValue({
      theme: 'dark',
      toggleTheme,
    })

    render(<Layout />)

    expect(document.documentElement.classList.contains('dark')).toBe(true)
  })

  it('renders Outlet for page content', () => {
    ;(useThemeStore as any).mockReturnValue({
      theme: 'light',
      toggleTheme: vi.fn(),
    })

    render(<Layout />)

    expect(screen.getByTestId('outlet')).toBeInTheDocument()
  })
})