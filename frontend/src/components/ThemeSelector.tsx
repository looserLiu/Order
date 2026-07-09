import { useThemeStore, getColorScheme } from '../stores/themeStore'

const colorSchemes = [
  { id: 'blue', name: '蓝色', colors: ['#3B82F6', '#60A5FA', '#93C5FD'] },
  { id: 'green', name: '绿色', colors: ['#10B981', '#34D399', '#6EE7B7'] },
  { id: 'purple', name: '紫色', colors: ['#8B5CF6', '#A78BFA', '#C4B5FD'] },
  { id: 'orange', name: '橙色', colors: ['#F97316', '#FB923C', '#FDBA74'] },
  { id: 'pink', name: '粉色', colors: ['#EC4899', '#F472B6', '#F9A8D4'] },
  { id: 'cyan', name: '青色', colors: ['#06B6D4', '#22D3EE', '#67E8F9'] },
] as const

export default function ThemeSelector() {
  const { colorScheme, setColorScheme, theme, setTheme } = useThemeStore()

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-2">主题模式</label>
        <div className="flex gap-2">
          <button
            onClick={() => {
              setTheme('light')
              document.documentElement.removeAttribute('data-theme')
            }}
            className={`flex-1 py-2 px-4 rounded-lg border-2 transition-all ${
              theme === 'light'
                ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20'
                : 'border-gray-200 dark:border-gray-700'
            }`}
          >
            ☀️ 浅色
          </button>
          <button
            onClick={() => {
              setTheme('dark')
              document.documentElement.setAttribute('data-theme', 'dark')
            }}
            className={`flex-1 py-2 px-4 rounded-lg border-2 transition-all ${
              theme === 'dark'
                ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20'
                : 'border-gray-200 dark:border-gray-700'
            }`}
          >
            🌙 深色
          </button>
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium mb-2">主题颜色</label>
        <div className="grid grid-cols-3 gap-2">
          {colorSchemes.map((scheme) => (
            <button
              key={scheme.id}
              onClick={() => {
                setColorScheme(scheme.id)
                document.documentElement.setAttribute('data-color-scheme', scheme.id)
              }}
              className={`p-3 rounded-lg border-2 transition-all ${
                colorScheme === scheme.id
                  ? 'border-primary-600'
                  : 'border-gray-200 dark:border-gray-700 hover:border-gray-300'
              }`}
            >
              <div className="flex gap-1 justify-center mb-1">
                {scheme.colors.map((color, i) => (
                  <div
                    key={i}
                    className="w-4 h-4 rounded-full"
                    style={{ backgroundColor: color }}
                  />
                ))}
              </div>
              <span className="text-xs">{scheme.name}</span>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}

// Hook to apply theme on mount
export function useApplyTheme() {
  const { theme, colorScheme } = useThemeStore()

  // Apply theme on mount
  if (typeof window !== 'undefined') {
    if (theme === 'dark') {
      document.documentElement.setAttribute('data-theme', 'dark')
    }
    document.documentElement.setAttribute('data-color-scheme', colorScheme)
  }
}
