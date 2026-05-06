import '../entities/transaction.dart';
import '../repositories/i_transaction_repository.dart';
import 'base_usecase.dart';

/// Use case: Get all transactions.
class GetAllTransactions extends NoParamsUseCase<List<TransactionEntity>> {
  final ITransactionRepository _repository;

  GetAllTransactions(this._repository);

  @override
  Future<List<TransactionEntity>> call() {
    return _repository.getAll();
  }
}

/// Use case: Get transaction by ID.
class GetTransactionById extends UseCase<TransactionEntity?, int> {
  final ITransactionRepository _repository;

  GetTransactionById(this._repository);

  @override
  Future<TransactionEntity?> call(int id) {
    return _repository.getById(id);
  }
}

/// Use case: Get transactions by account ID.
class GetTransactionsByAccount extends UseCase<List<TransactionEntity>, int> {
  final ITransactionRepository _repository;

  GetTransactionsByAccount(this._repository);

  @override
  Future<List<TransactionEntity>> call(int accountId) {
    return _repository.getByAccountId(accountId);
  }
}

/// Use case: Get transactions by category ID.
class GetTransactionsByCategory extends UseCase<List<TransactionEntity>, int> {
  final ITransactionRepository _repository;

  GetTransactionsByCategory(this._repository);

  @override
  Future<List<TransactionEntity>> call(int categoryId) {
    return _repository.getByCategoryId(categoryId);
  }
}

/// Use case: Get transactions by date range.
class GetTransactionsByDateRange extends UseCase<List<TransactionEntity>, DateRangeParams> {
  final ITransactionRepository _repository;

  GetTransactionsByDateRange(this._repository);

  @override
  Future<List<TransactionEntity>> call(DateRangeParams params) {
    return _repository.getByDateRange(params.start, params.end);
  }
}

/// Parameters for date range queries.
class DateRangeParams {
  final DateTime start;
  final DateTime end;

  DateRangeParams({required this.start, required this.end});
}

/// Use case: Create a new transaction.
class CreateTransaction extends UseCase<TransactionEntity, TransactionEntity> {
  final ITransactionRepository _repository;

  CreateTransaction(this._repository);

  @override
  Future<TransactionEntity> call(TransactionEntity transaction) {
    return _repository.create(transaction);
  }
}

/// Use case: Update an existing transaction.
class UpdateTransaction extends UseCase<TransactionEntity, TransactionEntity> {
  final ITransactionRepository _repository;

  UpdateTransaction(this._repository);

  @override
  Future<TransactionEntity> call(TransactionEntity transaction) {
    return _repository.update(transaction);
  }
}

/// Use case: Delete a transaction.
class DeleteTransaction extends UseCase<void, int> {
  final ITransactionRepository _repository;

  DeleteTransaction(this._repository);

  @override
  Future<void> call(int id) {
    return _repository.delete(id);
  }
}
