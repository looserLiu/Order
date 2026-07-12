import { create } from "zustand";
import { persist } from "zustand/middleware";

// Theme types
export type Theme = "light" | "dark";
export type ColorScheme =
  | "blue"
  | "green"
  | "purple"
  | "orange"
  | "pink"
  | "cyan";

// Color scheme configuration
export interface ColorSchemeConfig {
  primary: string;
  primaryHover: string;
  primaryLight: string;
}

// Theme state interface
export interface ThemeState {
  theme: Theme;
  colorScheme: ColorScheme;
  setTheme: (theme: Theme) => void;
  setColorScheme: (scheme: ColorScheme) => void;
  toggleTheme: () => void;
}

// Color schemes configuration
const colorSchemes: Record<ColorScheme, ColorSchemeConfig> = {
  blue: {
    primary: "#3B82F6",
    primaryHover: "#2563EB",
    primaryLight: "#DBEAFE",
  },
  green: {
    primary: "#10B981",
    primaryHover: "#059669",
    primaryLight: "#D1FAE5",
  },
  purple: {
    primary: "#8B5CF6",
    primaryHover: "#7C3AED",
    primaryLight: "#EDE9FE",
  },
  orange: {
    primary: "#F97316",
    primaryHover: "#EA580C",
    primaryLight: "#FFEDD5",
  },
  pink: {
    primary: "#EC4899",
    primaryHover: "#DB2777",
    primaryLight: "#FCE7F3",
  },
  cyan: {
    primary: "#06B6D4",
    primaryHover: "#0891B2",
    primaryLight: "#CFFAFE",
  },
};

// Get color scheme by name
export const getColorScheme = (scheme: ColorScheme): ColorSchemeConfig =>
  colorSchemes[scheme];

// Theme store with persistence
export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      theme: "light",
      colorScheme: "blue",
      setTheme: (theme) => set({ theme }),
      setColorScheme: (colorScheme) => set({ colorScheme }),
      toggleTheme: () =>
        set((state) => ({ theme: state.theme === "light" ? "dark" : "light" })),
    }),
    { name: "theme-storage" },
  ),
);

// Hook to get theme state
export const useTheme = () => {
  const theme = useThemeStore((state) => state.theme);
  const colorScheme = useThemeStore((state) => state.colorScheme);
  const setTheme = useThemeStore((state) => state.setTheme);
  const setColorScheme = useThemeStore((state) => state.setColorScheme);
  const toggleTheme = useThemeStore((state) => state.toggleTheme);

  return { theme, colorScheme, setTheme, setColorScheme, toggleTheme };
};
