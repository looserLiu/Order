import '../entities/transaction.dart';

/// Repository interface for Transaction operations.
abstract class ITransactionRepository {
  /// Get all transactions.
  Future<List<TransactionEntity>> getAll();

  /// Get transaction by ID.
  Future<TransactionEntity?> getById(int id);

  /// Get transactions by account ID.
  Future<List<TransactionEntity>> getByAccountId(int accountId);

  /// Get transactions by category ID.
  Future<List<TransactionEntity>> getByCategoryId(int categoryId);

  /// Get transactions within a date range.
  Future<List<TransactionEntity>> getByDateRange(DateTime start, DateTime end);

  /// Create a new transaction.
  Future<TransactionEntity> create(TransactionEntity transaction);

  /// Update an existing transaction.
  Future<TransactionEntity> update(TransactionEntity transaction);

  /// Delete a transaction by ID.
  Future<void> delete(int id);
}
