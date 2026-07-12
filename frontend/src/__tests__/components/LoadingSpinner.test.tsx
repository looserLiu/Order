import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import LoadingSpinner from '../../components/LoadingSpinner'

describe('LoadingSpinner Component', () => {
  it('renders with default size (md)', () => {
    const { container } = render(<LoadingSpinner />)

    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()
    expect(spinner).toHaveClass('h-8', 'w-8', 'border-2')
  })

  it('renders with small size', () => {
    const { container } = render(<LoadingSpinner size="sm" />)

    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()
    expect(spinner).toHaveClass('h-4', 'w-4', 'border-2')
  })

  it('renders with large size', () => {
    const { container } = render(<LoadingSpinner size="lg" />)

    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()
    expect(spinner).toHaveClass('h-12', 'w-12', 'border-b-2')
  })

  it('applies custom className', () => {
    const { container } = render(<LoadingSpinner className="custom-class" />)

    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toHaveClass('custom-class')
  })

  it('renders full screen version when fullScreen is true', () => {
    const { container } = render(<LoadingSpinner fullScreen />)

    const wrapper = container.querySelector('.min-h-screen')
    expect(wrapper).toBeInTheDocument()
    expect(wrapper).toHaveClass('flex', 'items-center', 'justify-center')
  })

  it('does not render full screen wrapper when fullScreen is false', () => {
    const { container } = render(<LoadingSpinner fullScreen={false} />)

    const wrapper = container.querySelector('.min-h-screen')
    expect(wrapper).not.toBeInTheDocument()
  })

  it('has primary color border', () => {
    const { container } = render(<LoadingSpinner />)

    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toHaveClass('border-primary-600')
  })
})