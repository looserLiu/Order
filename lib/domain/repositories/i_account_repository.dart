import '../entities/account.dart';

/// Repository interface for Account operations.
abstract class IAccountRepository {
  /// Get all accounts.
  Future<List<AccountEntity>> getAll();

  /// Get account by ID.
  Future<AccountEntity?> getById(int id);

  /// Create a new account.
  Future<AccountEntity> create(AccountEntity account);

  /// Update an existing account.
  Future<AccountEntity> update(AccountEntity account);

  /// Delete an account by ID.
  Future<void> delete(int id);

  /// Update account balance.
  Future<void> updateBalance(int id, double newBalance);
}
