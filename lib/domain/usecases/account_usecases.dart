import '../entities/account.dart';
import '../repositories/i_account_repository.dart';
import 'base_usecase.dart';

/// Use case: Get all accounts.
class GetAllAccounts extends NoParamsUseCase<List<AccountEntity>> {
  final IAccountRepository _repository;

  GetAllAccounts(this._repository);

  @override
  Future<List<AccountEntity>> call() {
    return _repository.getAll();
  }
}

/// Use case: Get account by ID.
class GetAccountById extends UseCase<AccountEntity?, int> {
  final IAccountRepository _repository;

  GetAccountById(this._repository);

  @override
  Future<AccountEntity?> call(int id) {
    return _repository.getById(id);
  }
}

/// Use case: Create a new account.
class CreateAccount extends UseCase<AccountEntity, AccountEntity> {
  final IAccountRepository _repository;

  CreateAccount(this._repository);

  @override
  Future<AccountEntity> call(AccountEntity account) {
    return _repository.create(account);
  }
}

/// Use case: Update an existing account.
class UpdateAccount extends UseCase<AccountEntity, AccountEntity> {
  final IAccountRepository _repository;

  UpdateAccount(this._repository);

  @override
  Future<AccountEntity> call(AccountEntity account) {
    return _repository.update(account);
  }
}

/// Use case: Delete an account.
class DeleteAccount extends UseCase<void, int> {
  final IAccountRepository _repository;

  DeleteAccount(this._repository);

  @override
  Future<void> call(int id) {
    return _repository.delete(id);
  }
}

/// Use case: Update account balance.
class UpdateAccountBalance extends UseCase<void, UpdateBalanceParams> {
  final IAccountRepository _repository;

  UpdateAccountBalance(this._repository);

  @override
  Future<void> call(UpdateBalanceParams params) {
    return _repository.updateBalance(params.id, params.newBalance);
  }
}

/// Parameters for UpdateAccountBalance.
class UpdateBalanceParams {
  final int id;
  final double newBalance;

  UpdateBalanceParams({required this.id, required this.newBalance});
}
