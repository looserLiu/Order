import '../../domain/entities/account.dart';
import '../../domain/repositories/i_account_repository.dart';
import '../database/database_helper.dart';
import '../models/account.dart';
import '../mappers/account_mapper.dart';

/// Implementation of IAccountRepository using SQLite.
class AccountRepositoryImpl implements IAccountRepository {
  final DatabaseHelper _dbHelper;

  AccountRepositoryImpl({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  @override
  Future<List<AccountEntity>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('accounts', orderBy: 'created_at DESC');
    final accounts = maps.map((map) => Account.fromMap(map)).toList();
    return AccountMapper.toEntityList(accounts);
  }

  @override
  Future<AccountEntity?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'accounts',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return AccountMapper.toEntity(Account.fromMap(maps.first));
  }

  @override
  Future<AccountEntity> create(AccountEntity account) async {
    final db = await _dbHelper.database;
    final model = AccountMapper.toModel(account);
    final id = await db.insert('accounts', model.toMap());
    return account.copyWith(id: id);
  }

  @override
  Future<AccountEntity> update(AccountEntity account) async {
    final db = await _dbHelper.database;
    final model = AccountMapper.toModel(account);
    await db.update(
      'accounts',
      model.toMap(),
      where: 'id = ?',
      whereArgs: [account.id],
    );
    return account;
  }

  @override
  Future<void> delete(int id) async {
    final db = await _dbHelper.database;
    await db.delete(
      'accounts',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  @override
  Future<void> updateBalance(int id, double newBalance) async {
    final db = await _dbHelper.database;
    await db.update(
      'accounts',
      {
        'balance': newBalance,
        'updated_at': DateTime.now().millisecondsSinceEpoch,
      },
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
