import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'

const resources = {
  en: {
    translation: {
      // Common
      save: 'Save',
      cancel: 'Cancel',
      delete: 'Delete',
      edit: 'Edit',
      add: 'Add',
      search: 'Search',
      loading: 'Loading...',
      noData: 'No data',
      
      // Navigation
      dashboard: 'Dashboard',
      transactions: 'Transactions',
      accounts: 'Accounts',
      categories: 'Categories',
      budgets: 'Budgets',
      reports: 'Reports',
      goals: 'Goals',
      debts: 'Debts',
      insurances: 'Insurances',
      families: 'Families',
      aaGroups: 'AA Groups',
      settings: 'Settings',
      
      // Auth
      login: 'Login',
      register: 'Register',
      logout: 'Logout',
      email: 'Email',
      password: 'Password',
      nickname: 'Nickname',
      
      // Dashboard
      thisMonthIncome: 'This Month Income',
      thisMonthExpense: 'This Month Expense',
      totalBalance: 'Total Balance',
      monthlyBalance: 'Monthly Balance',
      recentTransactions: 'Recent Transactions',
      budgetProgress: 'Budget Progress',
      
      // Transactions
      addTransaction: 'Add Transaction',
      expense: 'Expense',
      income: 'Income',
      transfer: 'Transfer',
      amount: 'Amount',
      date: 'Date',
      merchant: 'Merchant',
      note: 'Note',
      account: 'Account',
      category: 'Category',
      
      // Reports
      summary: 'Summary',
      trend: 'Trend',
      byCategory: 'By Category',
      byAccount: 'By Account',
      byMerchant: 'By Merchant',
      monthlyCompare: 'Monthly Compare',
      
      // Settings
      profile: 'Profile',
      appearance: 'Appearance',
      darkMode: 'Dark Mode',
      language: 'Language',
      currency: 'Currency',
      exportData: 'Export Data',
      importData: 'Import Data',
      security: 'Security',
      changePassword: 'Change Password',
    }
  },
  zh: {
    translation: {
      // Common
      save: '保存',
      cancel: '取消',
      delete: '删除',
      edit: '编辑',
      add: '添加',
      search: '搜索',
      loading: '加载中...',
      noData: '暂无数据',
      
      // Navigation
      dashboard: '仪表盘',
      transactions: '交易记录',
      accounts: '账户',
      categories: '分类',
      budgets: '预算',
      reports: '报表',
      goals: '目标',
      debts: '债务',
      insurances: '保险',
      families: '家庭',
      aaGroups: 'AA记账',
      settings: '设置',
      
      // Auth
      login: '登录',
      register: '注册',
      logout: '退出登录',
      email: '邮箱',
      password: '密码',
      nickname: '昵称',
      
      // Dashboard
      thisMonthIncome: '本月收入',
      thisMonthExpense: '本月支出',
      totalBalance: '账户余额',
      monthlyBalance: '本月结余',
      recentTransactions: '最近交易',
      budgetProgress: '预算进度',
      
      // Transactions
      addTransaction: '记一笔',
      expense: '支出',
      income: '收入',
      transfer: '转账',
      amount: '金额',
      date: '日期',
      merchant: '商家',
      note: '备注',
      account: '账户',
      category: '分类',
      
      // Reports
      summary: '摘要',
      trend: '趋势',
      byCategory: '按分类',
      byAccount: '按账户',
      byMerchant: '按商家',
      monthlyCompare: '月度对比',
      
      // Settings
      profile: '个人资料',
      appearance: '外观',
      darkMode: '深色模式',
      language: '语言',
      currency: '货币',
      exportData: '导出数据',
      importData: '导入数据',
      security: '安全设置',
      changePassword: '修改密码',
    }
  },
  ja: {
    translation: {
      // Common
      save: '保存',
      cancel: 'キャンセル',
      delete: '削除',
      edit: '編集',
      add: '追加',
      search: '検索',
      loading: '読み込み中...',
      noData: 'データなし',
      
      // Navigation
      dashboard: 'ダッシュボード',
      transactions: '取引',
      accounts: '口座',
      categories: 'カテゴリ',
      budgets: '予算',
      reports: 'レポート',
      goals: '目標',
      debts: '負債',
      insurances: '保険',
      families: 'ファミリー',
      aaGroups: '割り勘',
      settings: '設定',
      
      // Auth
      login: 'ログイン',
      register: '登録',
      logout: 'ログアウト',
      email: 'メール',
      password: 'パスワード',
      nickname: 'ニックネーム',
      
      // Dashboard
      thisMonthIncome: '今月の収入',
      thisMonthExpense: '今月の支出',
      totalBalance: '残高',
      monthlyBalance: '月間収支',
      recentTransactions: '最近の取引',
      budgetProgress: '予算進捗',
      
      // Transactions
      addTransaction: '記録',
      expense: '支出',
      income: '収入',
      transfer: '振替',
      amount: '金額',
      date: '日付',
      merchant: '店家',
      note: 'メモ',
      account: '口座',
      category: 'カテゴリ',
      
      // Reports
      summary: '概要',
      trend: 'トレンド',
      byCategory: 'カテゴリ別',
      byAccount: '口座別',
      byMerchant: '店家別',
      monthlyCompare: '月度比較',
      
      // Settings
      profile: 'プロフィール',
      appearance: '外観',
      darkMode: 'ダークモード',
      language: '言語',
      currency: '通貨',
      exportData: 'データ出力',
      importData: 'データインポート',
      security: 'セキュリティ',
      changePassword: 'パスワード変更',
    }
  }
}

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'zh',
    interpolation: {
      escapeValue: false
    },
    detection: {
      order: ['localStorage', 'navigator'],
      caches: ['localStorage']
    }
  })

export default i18n
