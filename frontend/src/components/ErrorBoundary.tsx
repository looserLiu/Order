import { Component, ReactNode, ErrorInfo, FC } from 'react'
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline'

// ErrorBoundary props
export interface ErrorBoundaryProps {
  children: ReactNode
  fallback?: ReactNode
  onError?: (error: Error, errorInfo: ErrorInfo) => void
}

// ErrorBoundary state
export interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

/**
 * ErrorBoundary - A reusable error boundary component
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    this.props.onError?.(error, errorInfo)
    console.error('ErrorBoundary caught an error:', error, errorInfo)
  }

  handleReset = (): void => {
    this.setState({ hasError: false, error: null })
  }

  render(): ReactNode {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback
      }

      return (
        <div className="flex flex-col items-center justify-center min-h-[400px] p-6">
          <div className="w-16 h-16 bg-red-100 dark:bg-red-900/20 rounded-full flex items-center justify-center mb-4">
            <ExclamationTriangleIcon className="w-8 h-8 text-red-600" />
          </div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
            出现了一些错误
          </h2>
          <p className="text-gray-500 dark:text-gray-400 text-center max-w-md mb-6">
            {this.state.error?.message || '应用发生了意外错误，请尝试刷新页面或联系支持'}
          </p>
          <button
            onClick={this.handleReset}
            className="btn-primary"
          >
            重试
          </button>
        </div>
      )
    }

    return this.props.children
  }
}

// Hook for error handling
export const useErrorHandler = () => {
  const handleError = (error: Error): void => {
    console.error('Error handled:', error)
  }
  
  return { handleError }
}

export default ErrorBoundary