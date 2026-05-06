import '../entities/budget.dart';

/// Repository interface for Budget operations.
abstract class IBudgetRepository {
  /// Get all budgets.
  Future<List<BudgetEntity>> getAll();

  /// Get budget by ID.
  Future<BudgetEntity?> getById(int id);

  /// Get budgets by category ID.
  Future<List<BudgetEntity>> getByCategoryId(int categoryId);

  /// Get budgets by period.
  Future<List<BudgetEntity>> getByPeriod(BudgetPeriod period);

  /// Create a new budget.
  Future<BudgetEntity> create(BudgetEntity budget);

  /// Update an existing budget.
  Future<BudgetEntity> update(BudgetEntity budget);

  /// Delete a budget by ID.
  Future<void> delete(int id);
}
