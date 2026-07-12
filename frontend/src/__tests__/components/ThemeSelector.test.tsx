import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import ThemeSelector from '../../components/ThemeSelector'
import { useThemeStore } from '../../stores/themeStore'

// Mock the theme store
vi.mock('../../stores/themeStore', () => ({
  useThemeStore: vi.fn(),
}))

describe('ThemeSelector Component', () => {
  const mockSetTheme = vi.fn()
  const mockSetColorScheme = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useThemeStore as any).mockReturnValue({
      theme: 'light',
      colorScheme: 'blue',
      setTheme: mockSetTheme,
      setColorScheme: mockSetColorScheme,
    })
  })

  it('renders theme mode label', () => {
    render(<ThemeSelector />)

    expect(screen.getByText('主题模式')).toBeInTheDocument()
  })

  it('renders color scheme label', () => {
    render(<ThemeSelector />)

    expect(screen.getByText('主题颜色')).toBeInTheDocument()
  })

  it('renders light and dark theme buttons', () => {
    render(<ThemeSelector />)

    expect(screen.getByText('☀️ 浅色')).toBeInTheDocument()
    expect(screen.getByText('🌙 深色')).toBeInTheDocument()
  })

  it('renders all color scheme options', () => {
    render(<ThemeSelector />)

    expect(screen.getByText('蓝色')).toBeInTheDocument()
    expect(screen.getByText('绿色')).toBeInTheDocument()
    expect(screen.getByText('紫色')).toBeInTheDocument()
    expect(screen.getByText('橙色')).toBeInTheDocument()
    expect(screen.getByText('粉色')).toBeInTheDocument()
    expect(screen.getByText('青色')).toBeInTheDocument()
  })

  it('calls setTheme with light when light button is clicked', () => {
    render(<ThemeSelector />)

    const lightButton = screen.getByText('☀️ 浅色')
    fireEvent.click(lightButton)

    expect(mockSetTheme).toHaveBeenCalledWith('light')
  })

  it('calls setTheme with dark when dark button is clicked', () => {
    render(<ThemeSelector />)

    const darkButton = screen.getByText('🌙 深色')
    fireEvent.click(darkButton)

    expect(mockSetTheme).toHaveBeenCalledWith('dark')
  })

  it('calls setColorScheme when color scheme button is clicked', () => {
    render(<ThemeSelector />)

    const greenButton = screen.getByText('绿色')
    fireEvent.click(greenButton)

    expect(mockSetColorScheme).toHaveBeenCalledWith('green')
  })

  it('highlights the selected color scheme', () => {
    ;(useThemeStore as any).mockReturnValue({
      theme: 'light',
      colorScheme: 'green',
      setTheme: mockSetTheme,
      setColorScheme: mockSetColorScheme,
    })

    render(<ThemeSelector />)

    // The green button should have the selected class
    const greenButton = screen.getByText('绿色').closest('button')
    expect(greenButton).toHaveClass('border-primary-600')
  })
})