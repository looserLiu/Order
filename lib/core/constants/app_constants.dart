// Core constants - app-wide constants
class AppConstants {
  AppConstants._();

  static const String appName = 'Order';
  static const String appVersion = '1.0.0';
  static const String defaultCurrency = 'CNY';

  // Database
  static const String databaseName = 'order.db';
  static const int databaseVersion = 1;

  // Account types
  static const List<String> accountTypes = ['cash', 'bank', 'credit_card', 'digital'];

  // Transaction types
  static const String transactionTypeIncome = 'income';
  static const String transactionTypeExpense = 'expense';

  // Flow types
  static const String flowTypeIn = 'in';
  static const String flowTypeOut = 'out';
  static const String flowTypeTransfer = 'transfer';
  static const String flowTypeAdjust = 'adjust';

  // Budget periods
  static const String budgetPeriodMonthly = 'monthly';
  static const String budgetPeriodWeekly = 'weekly';
  static const String budgetPeriodYearly = 'yearly';
}
