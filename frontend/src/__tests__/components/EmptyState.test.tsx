import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import EmptyState from '../../components/EmptyState'

describe('EmptyState Component', () => {
  it('renders with default icon and title', () => {
    render(<EmptyState title="No data" />)

    expect(screen.getByText('No data')).toBeInTheDocument()
  })

  it('renders with description', () => {
    render(
      <EmptyState
        title="No data"
        description="There is no data to display"
      />
    )

    expect(screen.getByText('There is no data to display')).toBeInTheDocument()
  })

  it('renders with action button', () => {
    render(
      <EmptyState
        title="No data"
        action={<button>Add Item</button>}
      />
    )

    expect(screen.getByText('Add Item')).toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = render(<EmptyState title="No data" className="custom-class" />)

    const wrapper = container.querySelector('.custom-class')
    expect(wrapper).toBeInTheDocument()
  })

  it('renders different icon types', () => {
    const { container } = render(<EmptyState title="No data" icon="chart" />)

    // Icon should be rendered
    const iconWrapper = container.querySelector('.w-16.h-16')
    expect(iconWrapper).toBeInTheDocument()
  })

  it('renders with card icon', () => {
    const { container } = render(<EmptyState title="No accounts" icon="card" />)

    expect(screen.getByText('No accounts')).toBeInTheDocument()
    const iconWrapper = container.querySelector('.w-16.h-16')
    expect(iconWrapper).toBeInTheDocument()
  })

  it('renders with tag icon', () => {
    render(<EmptyState title="No tags" icon="tag" />)

    expect(screen.getByText('No tags')).toBeInTheDocument()
  })

  it('renders with calendar icon', () => {
    render(<EmptyState title="No events" icon="calendar" />)

    expect(screen.getByText('No events')).toBeInTheDocument()
  })

  it('renders with goal icon', () => {
    render(<EmptyState title="No goals" icon="goal" />)

    expect(screen.getByText('No goals')).toBeInTheDocument()
  })

  it('renders with insurance icon', () => {
    render(<EmptyState title="No insurance" icon="insurance" />)

    expect(screen.getByText('No insurance')).toBeInTheDocument()
  })

  it('renders with netWorth icon', () => {
    render(<EmptyState title="No net worth" icon="netWorth" />)

    expect(screen.getByText('No net worth')).toBeInTheDocument()
  })
})