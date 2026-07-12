import { describe, it, expect, beforeEach } from "vitest";
import { useAuthStore, useAuth } from "../../stores/authStore";
import { User } from "../../services/api";

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

describe("Auth Store", () => {
  beforeEach(() => {
    // Reset store state before each test
    useAuthStore.setState({ token: null, user: null });
    vi.clearAllMocks();
  });

  it("should have initial state with null values", () => {
    const { token, user } = useAuth();

    expect(token).toBeNull();
    expect(user).toBeNull();
  });

  it("should set auth with token and user", () => {
    const testUser: User = {
      id: "1",
      email: "test@example.com",
      nickname: "Test User",
    };
    const testToken = "test-token-123";

    useAuthStore.getState().setAuth(testToken, testUser);

    const { token, user } = useAuth();
    expect(token).toBe(testToken);
    expect(user).toEqual(testUser);
  });

  it("should logout and clear auth state", () => {
    const testUser: User = {
      id: "1",
      email: "test@example.com",
      nickname: "Test User",
    };

    // First set auth
    useAuthStore.getState().setAuth("token", testUser);
    expect(useAuthStore.getState().token).toBe("token");

    // Then logout
    useAuthStore.getState().logout();

    const { token, user } = useAuth();
    expect(token).toBeNull();
    expect(user).toBeNull();
  });

  it("should update user with partial data", () => {
    const testUser: User = {
      id: "1",
      email: "test@example.com",
      nickname: "Test User",
    };

    useAuthStore.getState().setAuth("token", testUser);

    // Update user
    useAuthStore.getState().updateUser({ nickname: "Updated User" });

    const { user } = useAuth();
    expect(user?.nickname).toBe("Updated User");
    expect(user?.email).toBe("test@example.com");
  });

  it("should not update user when user is null", () => {
    useAuthStore.getState().updateUser({ nickname: "Updated User" });

    const { user } = useAuth();
    expect(user).toBeNull();
  });
});
