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
  login: (data: { email: string; password: string }) =>
    api.post("/auth/login", data),
  register: (data: { email: string; password: string; nickname?: string }) =>
    api.post("/auth/register", data),
  refresh: () => api.post("/auth/refresh"),
};

export const dashboardApi = {
  getStats: () => api.get("/dashboard/stats"),
};

export const accountApi = {
  list: () => api.get("/accounts"),
  create: (data: any) => api.post("/accounts", data),
  update: (id: string, data: any) => api.put(`/accounts/${id}`, data),
  delete: (id: string) => api.delete(`/accounts/${id}`),
  getTotalBalance: () => api.get("/accounts/total/balance"),
};

export const categoryApi = {
  list: () => api.get("/categories"),
  getTree: () => api.get("/categories/tree"),
  create: (data: any) => api.post("/categories", data),
  update: (id: string, data: any) => api.put(`/categories/${id}`, data),
  delete: (id: string) => api.delete(`/categories/${id}`),
};

export const transactionApi = {
  list: (params?: any) => api.get("/transactions", { params }),
  listByDate: (date: string) =>
    api.get("/transactions/by-date", { params: { date } }),
  create: (data: any) => api.post("/transactions", data),
  update: (id: string, data: any) => api.put(`/transactions/${id}`, data),
  delete: (id: string) => api.delete(`/transactions/${id}`),
  batchDelete: (ids: string[]) =>
    api.post("/transactions/batch-delete", { ids }),
};

export const tagApi = {
  list: () => api.get("/tags"),
  create: (data: any) => api.post("/tags", data),
  update: (id: string, data: any) => api.put(`/tags/${id}`, data),
  delete: (id: string) => api.delete(`/tags/${id}`),
};

export const budgetApi = {
  list: () => api.get("/budgets"),
  create: (data: any) => api.post("/budgets", data),
  update: (id: string, data: any) => api.put(`/budgets/${id}`, data),
  delete: (id: string) => api.delete(`/budgets/${id}`),
  getProgress: (id: string) => api.get(`/budgets/${id}/progress`),
};

export const assetApi = {
  list: (params?: any) => api.get("/assets", { params }),
  create: (data: any) => api.post("/assets", data),
  update: (id: string, data: any) => api.put(`/assets/${id}`, data),
  delete: (id: string) => api.delete(`/assets/${id}`),
  getSummary: () => api.get("/assets/summary"),
};

export const reminderApi = {
  list: () => api.get("/reminders"),
  create: (data: any) => api.post("/reminders", data),
  update: (id: string, data: any) => api.put(`/reminders/${id}`, data),
  delete: (id: string) => api.delete(`/reminders/${id}`),
};

export const notificationApi = {
  list: (params?: any) => api.get("/notifications", { params }),
  markAsRead: (id: string) => api.put(`/notifications/${id}/read`),
  markAllAsRead: () => api.put("/notifications/read-all"),
  delete: (id: string) => api.delete(`/notifications/${id}`),
};

export const searchApi = {
  search: (params: { q: string; type?: string }) =>
    api.get("/search", { params }),
};

export const importApi = {
  importTransactions: (data: any[]) => api.post("/import/transactions", data),
};

export const reportApi = {
  summary: (params?: any) => api.get("/reports/summary", { params }),
  trend: (params?: any) => api.get("/reports/trend", { params }),
  byCategory: (params?: any) => api.get("/reports/category", { params }),
  byAccount: (params?: any) => api.get("/reports/account", { params }),
  byMerchant: (params?: any) => api.get("/reports/merchant", { params }),
  monthlyCompare: (params?: any) => api.get("/reports/monthly", { params }),
  export: (params?: any) => api.get("/reports/export", { params }),
};

export const familyApi = {
  list: () => api.get("/families"),
  create: (data: any) => api.post("/families", data),
  get: (id: string) => api.get(`/families/${id}`),
  delete: (id: string) => api.delete(`/families/${id}`),
  addMember: (id: string, data: any) =>
    api.post(`/families/${id}/members`, data),
  removeMember: (id: string, memberId: string) =>
    api.delete(`/families/${id}/members/${memberId}`),
  getTransactions: (id: string) => api.get(`/families/${id}/transactions`),
  createTransaction: (id: string, data: any) =>
    api.post(`/families/${id}/transactions`, data),
};

export const userApi = {
  getMe: () => api.get("/users/me"),
  updateMe: (data: any) => api.put("/users/me", data),
  changePassword: (data: any) => api.post("/users/password", data),
};

export const aaGroupApi = {
  list: () => api.get("/aa-groups"),
  create: (data: any) => api.post("/aa-groups", data),
  delete: (id: string) => api.delete(`/aa-groups/${id}`),
  addExpense: (id: string, data: any) =>
    api.post(`/aa-groups/${id}/expense`, data),
  getSettlements: (id: string) => api.get(`/aa-groups/${id}/settlements`),
};

export const goalApi = {
  list: () => api.get("/goals"),
  create: (data: any) => api.post("/goals", data),
  update: (id: string, data: any) => api.put(`/goals/${id}`, data),
  addAmount: (id: string, data: any) =>
    api.post(`/goals/${id}/add-amount`, data),
  delete: (id: string) => api.delete(`/goals/${id}`),
};

export const insuranceApi = {
  list: () => api.get("/insurances"),
  create: (data: any) => api.post("/insurances", data),
  update: (id: string, data: any) => api.put(`/insurances/${id}`, data),
  delete: (id: string) => api.delete(`/insurances/${id}`),
  getSummary: () => api.get("/insurances/summary"),
};

export const netWorthApi = {
  get: () => api.get("/net-worth"),
};

export const backupApi = {
  exportAll: () => api.get("/backup/export"),
  importAll: (data: any) => api.post("/backup/import", data),
  list: () => api.get("/backup/list"),
};

export const currencyApi = {
  list: () => api.get("/currencies"),
  getRates: () => api.get("/currencies/rates"),
  convert: (data: any) => api.post("/currencies/convert", data),
};

export const csvApi = {
  importCSV: (data: any) => api.post("/import/csv", data),
};

export const statisticsApi = {
  get: () => api.get("/statistics"),
};

// File upload API
export const uploadApi = {
  upload: (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    return api.post("/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  },
};

// Cash flow projection API
export const cashFlowApi = {
  getProjection: (days = 30) =>
    api.get("/cashflow/projection", { params: { days } }),
};

// Budget alerts API
export const budgetAlertApi = {
  getAlerts: () => api.get("/budgets/alerts"),
};

export default api;
