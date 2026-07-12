import { describe, it, expect, beforeEach } from "vitest";
import {
  useThemeStore,
  getColorScheme,
  useTheme,
} from "../../stores/themeStore";

// Mock localStorage for zustand persist
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};

Object.defineProperty(window, "localStorage", {
  value: localStorageMock,
});

describe("Theme Store", () => {
  beforeEach(() => {
    // Reset store state before each test
    useThemeStore.setState({ theme: "light", colorScheme: "blue" });
    vi.clearAllMocks();
  });

  it("should have initial state with light theme and blue color scheme", () => {
    const { theme, colorScheme } = useTheme();

    expect(theme).toBe("light");
    expect(colorScheme).toBe("blue");
  });

  it("should set theme", () => {
    useThemeStore.getState().setTheme("dark");

    const { theme } = useTheme();
    expect(theme).toBe("dark");
  });

  it("should set color scheme", () => {
    useThemeStore.getState().setColorScheme("green");

    const { colorScheme } = useTheme();
    expect(colorScheme).toBe("green");
  });

  it("should toggle theme from light to dark", () => {
    useThemeStore.getState().toggleTheme();

    const { theme } = useTheme();
    expect(theme).toBe("dark");
  });

  it("should toggle theme from dark to light", () => {
    useThemeStore.setState({ theme: "dark" });
    useThemeStore.getState().toggleTheme();

    const { theme } = useTheme();
    expect(theme).toBe("light");
  });
});

describe("getColorScheme", () => {
  it("returns blue color scheme config", () => {
    const scheme = getColorScheme("blue");

    expect(scheme.primary).toBe("#3B82F6");
    expect(scheme.primaryHover).toBe("#2563EB");
    expect(scheme.primaryLight).toBe("#DBEAFE");
  });

  it("returns green color scheme config", () => {
    const scheme = getColorScheme("green");

    expect(scheme.primary).toBe("#10B981");
    expect(scheme.primaryHover).toBe("#059669");
    expect(scheme.primaryLight).toBe("#D1FAE5");
  });

  it("returns purple color scheme config", () => {
    const scheme = getColorScheme("purple");

    expect(scheme.primary).toBe("#8B5CF6");
    expect(scheme.primaryHover).toBe("#7C3AED");
    expect(scheme.primaryLight).toBe("#EDE9FE");
  });

  it("returns orange color scheme config", () => {
    const scheme = getColorScheme("orange");

    expect(scheme.primary).toBe("#F97316");
    expect(scheme.primaryHover).toBe("#EA580C");
    expect(scheme.primaryLight).toBe("#FFEDD5");
  });

  it("returns pink color scheme config", () => {
    const scheme = getColorScheme("pink");

    expect(scheme.primary).toBe("#EC4899");
    expect(scheme.primaryHover).toBe("#DB2777");
    expect(scheme.primaryLight).toBe("#FCE7F3");
  });

  it("returns cyan color scheme config", () => {
    const scheme = getColorScheme("cyan");

    expect(scheme.primary).toBe("#06B6D4");
    expect(scheme.primaryHover).toBe("#0891B2");
    expect(scheme.primaryLight).toBe("#CFFAFE");
  });
});
