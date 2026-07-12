import { describe, it, expect, vi, beforeEach } from "vitest";
import axios from "axios";
import {
  authApi,
  accountApi,
  transactionApi,
  categoryApi,
  budgetApi,
} from "../../services/api";

// Mock axios
vi.mock("axios", () => ({
  create: vi.fn(() => ({
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  })),
}));

// Mock auth store
vi.mock("../../stores/authStore", () => ({
  useAuthStore: {
    getState: () => ({
      token: "test-token",
      logout: vi.fn(),
    }),
  },
}));

describe("API Service", () => {
  describe("authApi", () => {
    it("should have login method", () => {
      expect(authApi.login).toBeDefined();
      expect(typeof authApi.login).toBe("function");
    });

    it("should have register method", () => {
      expect(authApi.register).toBeDefined();
      expect(typeof authApi.register).toBe("function");
    });

    it("should have refresh method", () => {
      expect(authApi.refresh).toBeDefined();
      expect(typeof authApi.refresh).toBe("function");
    });
  });

  describe("accountApi", () => {
    it("should have list method", () => {
      expect(accountApi.list).toBeDefined();
      expect(typeof accountApi.list).toBe("function");
    });

    it("should have create method", () => {
      expect(accountApi.create).toBeDefined();
      expect(typeof accountApi.create).toBe("function");
    });

    it("should have update method", () => {
      expect(accountApi.update).toBeDefined();
      expect(typeof accountApi.update).toBe("function");
    });

    it("should have delete method", () => {
      expect(accountApi.delete).toBeDefined();
      expect(typeof accountApi.delete).toBe("function");
    });

    it("should have getTotalBalance method", () => {
      expect(accountApi.getTotalBalance).toBeDefined();
      expect(typeof accountApi.getTotalBalance).toBe("function");
    });
  });

  describe("transactionApi", () => {
    it("should have list method", () => {
      expect(transactionApi.list).toBeDefined();
      expect(typeof transactionApi.list).toBe("function");
    });

    it("should have create method", () => {
      expect(transactionApi.create).toBeDefined();
      expect(typeof transactionApi.create).toBe("function");
    });

    it("should have update method", () => {
      expect(transactionApi.update).toBeDefined();
      expect(typeof transactionApi.update).toBe("function");
    });

    it("should have delete method", () => {
      expect(transactionApi.delete).toBeDefined();
      expect(typeof transactionApi.delete).toBe("function");
    });

    it("should have batchDelete method", () => {
      expect(transactionApi.batchDelete).toBeDefined();
      expect(typeof transactionApi.batchDelete).toBe("function");
    });
  });

  describe("categoryApi", () => {
    it("should have list method", () => {
      expect(categoryApi.list).toBeDefined();
      expect(typeof categoryApi.list).toBe("function");
    });

    it("should have create method", () => {
      expect(categoryApi.create).toBeDefined();
      expect(typeof categoryApi.create).toBe("function");
    });

    it("should have update method", () => {
      expect(categoryApi.update).toBeDefined();
      expect(typeof categoryApi.update).toBe("function");
    });

    it("should have delete method", () => {
      expect(categoryApi.delete).toBeDefined();
      expect(typeof categoryApi.delete).toBe("function");
    });
  });

  describe("budgetApi", () => {
    it("should have list method", () => {
      expect(budgetApi.list).toBeDefined();
      expect(typeof budgetApi.list).toBe("function");
    });

    it("should have create method", () => {
      expect(budgetApi.create).toBeDefined();
      expect(typeof budgetApi.create).toBe("function");
    });

    it("should have update method", () => {
      expect(budgetApi.update).toBeDefined();
      expect(typeof budgetApi.update).toBe("function");
    });

    it("should have delete method", () => {
      expect(budgetApi.delete).toBeDefined();
      expect(typeof budgetApi.delete).toBe("function");
    });

    it("should have getProgress method", () => {
      expect(budgetApi.getProgress).toBeDefined();
      expect(typeof budgetApi.getProgress).toBe("function");
    });
  });
});
