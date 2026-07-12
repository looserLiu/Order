import axios from "axios";
import { useAuthStore } from "../stores/authStore";

// Type definitions
export interface User {
  id: string;
  email: string;
  nickname?: string;
  avatar_url?: string;
  currency?: string;
  timezone?: string;
  phone?: string;
}

export interface Account {
  id: string;
  user_id: string;
  name: string;
  type: string;
  balance: number;
  currency?: string;
  icon?: string;
  color?: string;
  is_default: boolean;
  sort_order: number;
}

export interface Category {
  id: string;
  user_id?: string;
  parent_id?: string;
  name: string;
  icon?: string;
  color?: string;
  type: string;
  sort_order: number;
  is_system: boolean;
}

export interface Transaction {
  id: string;
  user_id: string;
  account_id: string;
  target_account_id?: string;
  category_id?: string;
  type: string;
  amount: number;
  currency?: string;
  merchant?: string;
  note?: string;
  bill_date: string;
  is_recurring: boolean;
  recurring_rule?: string;
  next_date?: string;
  tags?: Tag[];
  category?: Category;
  account?: Account;
}

export interface Tag {
  id: string;
  user_id: string;
  name: string;
  color?: string;
}

export interface Budget {
  id: string;
  user_id: string;
  category_id?: string;
  amount: number;
  period: string;
  start_date: string;
  end_date?: string;
  alert_threshold: number;
}

export interface AssetChange {
  id: string;
  user_id: string;
  asset_type: string;
  related_user?: string;
  name: string;
  amount: number;
  interest_rate?: number;
  start_date?: string;
  end_date?: string;
  status: string;
  note?: string;
}

export interface Reminder {
  id: string;
  user_id: string;
  title: string;
  content?: string;
  remind_time: string;
  repeat_type?: string;
  category_id?: string;
  is_active: boolean;
}

export interface Notification {
  id: string;
  user_id: string;
  title: string;
  content?: string;
  is_read: boolean;
  created_at: string;
  type?: string;
}

export interface AAGroup {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  total_amount: number;
  members?: AAMember[];
}

export interface AAMember {
  id: string;
  group_id: string;
  name: string;
  paid: number;
  owe: number;
}

export interface FinancialGoal {
  id: string;
  user_id: string;
  name: string;
  target_amount: number;
  current_amount: number;
  deadline?: string;
  category?: string;
  note?: string;
  status?: string;
}

export interface Insurance {
  id: string;
  user_id: string;
  name: string;
  insurance_type?: string;
  company?: string;
  premium?: number;
  payment_type?: string;
  start_date?: string;
  end_date?: string;
  coverage?: number;
  beneficiary?: string;
  note?: string;
  status?: string;
  next_payment_date?: string;
}

export interface Currency {
  code: string;
  name: string;
  symbol: string;
  rate: number;
}

export interface CashFlowProjection {
  date: string;
  projected_balance: number;
  income: number;
  expense: number;
  recurring_transactions: number;
}

export interface BudgetAlert {
  budget_id: string;
  category_name: string;
  budget_amount: number;
  spent_amount: number;
  remaining: number;
  alert_type: string;
  percentage: number;
}

// Request types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  nickname?: string;
}

export interface CreateAccountRequest {
  name: string;
  type: string;
  balance: number;
  color?: string;
}

export interface CreateCategoryRequest {
  name: string;
  type: string;
  color?: string;
}

export interface CreateTransactionRequest {
  type: string;
  account_id: string;
  target_account_id?: string;
  category_id?: string;
  amount: number;
  merchant?: string;
  note?: string;
  bill_date: string;
}

export interface CreateBudgetRequest {
  category_id?: string;
  amount: number;
  period: string;
  start_date: string;
  end_date?: string;
  alert_threshold?: number;
}

export interface CreateTagRequest {
  name: string;
  color?: string;
}

export interface UpdateUserRequest {
  nickname?: string;
  currency?: string;
  timezone?: string;
}

export interface CreateGoalRequest {
  name: string;
  target_amount: number;
  deadline?: string;
  category?: string;
  note?: string;
}

export interface CreateInsuranceRequest {
  name: string;
  insurance_type?: string;
  company?: string;
  premium?: number;
  payment_type?: string;
  start_date?: string;
  end_date?: string;
  coverage?: number;
  beneficiary?: string;
  note?: string;
}

export interface CreateAAGroupRequest {
  name: string;
  description?: string;
}

export interface AddExpenseRequest {
  amount: number;
  description?: string;
  payer_id?: string;
  splits?: { member_id: string; amount: number }[];
}

// Response types
export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface SummaryResponse {
  income: number;
  expense: number;
  balance: number;
}

export interface CashFlowResponse {
  total: number;
  projections: CashFlowProjection[];
}

export interface ImportResult {
  imported: number;
  failed: number;
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

export const authApi = {
  login: (data: LoginRequest) =>
    api.post<ApiResponse<AuthResponse>>("/auth/login", data),
  register: (data: RegisterRequest) =>
    api.post<ApiResponse<AuthResponse>>("/auth/register", data),
  refresh: () => api.post<ApiResponse<{ token: string }>>("/auth/refresh"),
};

export const dashboardApi = {
  getStats: () => api.get<ApiResponse<SummaryResponse>>("/dashboard/stats"),
};

export const accountApi = {
  list: () => api.get<ApiResponse<Account[]>>("/accounts"),
  create: (data: CreateAccountRequest) =>
    api.post<ApiResponse<Account>>("/accounts", data),
  update: (id: string, data: Partial<CreateAccountRequest>) =>
    api.put<ApiResponse<Account>>(`/accounts/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/accounts/${id}`),
  getTotalBalance: () =>
    api.get<ApiResponse<{ total: number }>>("/accounts/total/balance"),
};

export const categoryApi = {
  list: () => api.get<ApiResponse<Category[]>>("/categories"),
  getTree: () => api.get<ApiResponse<Category[]>>("/categories/tree"),
  create: (data: CreateCategoryRequest) =>
    api.post<ApiResponse<Category>>("/categories", data),
  update: (id: string, data: Partial<CreateCategoryRequest>) =>
    api.put<ApiResponse<Category>>(`/categories/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/categories/${id}`),
};

export const transactionApi = {
  list: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<PaginatedResponse<Transaction>>>("/transactions", {
      params,
    }),
  listByDate: (date: string) =>
    api.get<ApiResponse<Transaction[]>>("/transactions/by-date", {
      params: { date },
    }),
  create: (data: CreateTransactionRequest) =>
    api.post<ApiResponse<Transaction>>("/transactions", data),
  update: (id: string, data: Partial<CreateTransactionRequest>) =>
    api.put<ApiResponse<Transaction>>(`/transactions/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/transactions/${id}`),
  batchDelete: (ids: string[]) =>
    api.post<ApiResponse<void>>("/transactions/batch-delete", { ids }),
};

export const tagApi = {
  list: () => api.get<ApiResponse<Tag[]>>("/tags"),
  create: (data: CreateTagRequest) => api.post<ApiResponse<Tag>>("/tags", data),
  update: (id: string, data: Partial<CreateTagRequest>) =>
    api.put<ApiResponse<Tag>>(`/tags/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/tags/${id}`),
};

export const budgetApi = {
  list: () => api.get<ApiResponse<Budget[]>>("/budgets"),
  create: (data: CreateBudgetRequest) =>
    api.post<ApiResponse<Budget>>("/budgets", data),
  update: (id: string, data: Partial<CreateBudgetRequest>) =>
    api.put<ApiResponse<Budget>>(`/budgets/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/budgets/${id}`),
  getProgress: (id: string) =>
    api.get<
      ApiResponse<{ spent: number; remaining: number; percentage: number }>
    >(`/budgets/${id}/progress`),
};

export const assetApi = {
  list: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<AssetChange[]>>("/assets", { params }),
  create: (data: Partial<AssetChange>) =>
    api.post<ApiResponse<AssetChange>>("/assets", data),
  update: (id: string, data: Partial<AssetChange>) =>
    api.put<ApiResponse<AssetChange>>(`/assets/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/assets/${id}`),
  getSummary: () =>
    api.get<ApiResponse<{ total: number; byType: Record<string, number> }>>(
      "/assets/summary",
    ),
};

export const reminderApi = {
  list: () => api.get<ApiResponse<Reminder[]>>("/reminders"),
  create: (data: Partial<Reminder>) =>
    api.post<ApiResponse<Reminder>>("/reminders", data),
  update: (id: string, data: Partial<Reminder>) =>
    api.put<ApiResponse<Reminder>>(`/reminders/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/reminders/${id}`),
};

export const notificationApi = {
  list: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<PaginatedResponse<Notification>>>("/notifications", {
      params,
    }),
  markAsRead: (id: string) =>
    api.put<ApiResponse<void>>(`/notifications/${id}/read`),
  markAllAsRead: () => api.put<ApiResponse<void>>("/notifications/read-all"),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/notifications/${id}`),
};

export const searchApi = {
  search: (params: { q: string; type?: string }) =>
    api.get<ApiResponse<unknown[]>>("/search", { params }),
};

export const importApi = {
  importTransactions: (data: unknown[]) =>
    api.post<ApiResponse<ImportResult>>("/import/transactions", data),
};

export const reportApi = {
  summary: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<SummaryResponse>>("/reports/summary", { params }),
  trend: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<unknown[]>>("/reports/trend", { params }),
  byCategory: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<unknown[]>>("/reports/category", { params }),
  byAccount: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<unknown[]>>("/reports/account", { params }),
  byMerchant: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<unknown[]>>("/reports/merchant", { params }),
  monthlyCompare: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<unknown[]>>("/reports/monthly", { params }),
  export: (params?: Record<string, unknown>) =>
    api.get<ApiResponse<unknown>>("/reports/export", { params }),
};

export const familyApi = {
  list: () => api.get<ApiResponse<unknown[]>>("/families"),
  create: (data: unknown) => api.post<ApiResponse<unknown>>("/families", data),
  get: (id: string) => api.get<ApiResponse<unknown>>(`/families/${id}`),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/families/${id}`),
  addMember: (id: string, data: unknown) =>
    api.post<ApiResponse<unknown>>(`/families/${id}/members`, data),
  removeMember: (id: string, memberId: string) =>
    api.delete<ApiResponse<void>>(`/families/${id}/members/${memberId}`),
  getTransactions: (id: string) =>
    api.get<ApiResponse<Transaction[]>>(`/families/${id}/transactions`),
  createTransaction: (id: string, data: unknown) =>
    api.post<ApiResponse<Transaction>>(`/families/${id}/transactions`, data),
};

export const userApi = {
  getMe: () => api.get<ApiResponse<User>>("/users/me"),
  updateMe: (data: UpdateUserRequest) =>
    api.put<ApiResponse<User>>("/users/me", data),
  changePassword: (data: { old_password: string; new_password: string }) =>
    api.post<ApiResponse<void>>("/users/password", data),
};

export const aaGroupApi = {
  list: () => api.get<ApiResponse<AAGroup[]>>("/aa-groups"),
  create: (data: CreateAAGroupRequest) =>
    api.post<ApiResponse<AAGroup>>("/aa-groups", data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/aa-groups/${id}`),
  addExpense: (id: string, data: AddExpenseRequest) =>
    api.post<ApiResponse<unknown>>(`/aa-groups/${id}/expense`, data),
  getSettlements: (id: string) =>
    api.get<ApiResponse<unknown[]>>(`/aa-groups/${id}/settlements`),
};

export const goalApi = {
  list: () => api.get<ApiResponse<FinancialGoal[]>>("/goals"),
  create: (data: CreateGoalRequest) =>
    api.post<ApiResponse<FinancialGoal>>("/goals", data),
  update: (id: string, data: Partial<CreateGoalRequest>) =>
    api.put<ApiResponse<FinancialGoal>>(`/goals/${id}`, data),
  addAmount: (id: string, data: { amount: number }) =>
    api.post<ApiResponse<FinancialGoal>>(`/goals/${id}/add-amount`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/goals/${id}`),
};

export const insuranceApi = {
  list: () => api.get<ApiResponse<Insurance[]>>("/insurances"),
  create: (data: CreateInsuranceRequest) =>
    api.post<ApiResponse<Insurance>>("/insurances", data),
  update: (id: string, data: Partial<CreateInsuranceRequest>) =>
    api.put<ApiResponse<Insurance>>(`/insurances/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/insurances/${id}`),
  getSummary: () =>
    api.get<ApiResponse<{ total_premium: number; total_coverage: number }>>(
      "/insurances/summary",
    ),
};

export const netWorthApi = {
  get: () =>
    api.get<
      ApiResponse<{ total: number; assets: number; liabilities: number }>
    >("/net-worth"),
};

export const backupApi = {
  exportAll: () => api.get<ApiResponse<unknown>>("/backup/export"),
  importAll: (data: unknown) =>
    api.post<ApiResponse<void>>("/backup/import", data),
  list: () => api.get<ApiResponse<unknown[]>>("/backup/list"),
};

export const currencyApi = {
  list: () => api.get<ApiResponse<Currency[]>>("/currencies"),
  getRates: () =>
    api.get<ApiResponse<Record<string, number>>>("/currencies/rates"),
  convert: (data: { from: string; to: string; amount: number }) =>
    api.post<ApiResponse<{ result: number }>>("/currencies/convert", data),
};

export const csvApi = {
  importCSV: (data: { transactions: unknown[] }) =>
    api.post<ApiResponse<ImportResult>>("/import/csv", data),
};

export const statisticsApi = {
  get: () => api.get<ApiResponse<unknown>>("/statistics"),
};

// File upload API
export const uploadApi = {
  upload: (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    return api.post<ApiResponse<{ url: string }>>("/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  },
};

// Cash flow projection API
export const cashFlowApi = {
  getProjection: (days = 30) =>
    api.get<ApiResponse<CashFlowResponse>>("/cashflow/projection", {
      params: { days },
    }),
};

// Budget alerts API
export const budgetAlertApi = {
  getAlerts: () => api.get<ApiResponse<BudgetAlert[]>>("/budgets/alerts"),
};

export default api;
