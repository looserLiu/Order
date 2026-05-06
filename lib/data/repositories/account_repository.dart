import 'package:sqflite/sqflite.dart';
import '../database/database_helper.dart';
import '../models/account.dart';

/// Repository for Account data operations.
class AccountRepository {
  final DatabaseHelper _dbHelper;

  AccountRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new account.
  Future<int> insert(Account account) async {
    final db = await _dbHelper.database;
    return await db.insert('accounts', account.toMap());
  }

  /// Update an existing account.
  Future<int> update(Account account) async {
    final db = await _dbHelper.database;
    return await db.update(
      'accounts',
      account.toMap(),
      where: 'id = ?',
      whereArgs: [account.id],
    );
  }

  /// Delete an account by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'accounts',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all accounts.
  Future<List<Account>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('accounts', orderBy: 'created_at DESC');
    return maps.map((map) => Account.fromMap(map)).toList();
  }

  /// Get account by ID.
  Future<Account?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'accounts',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return Account.fromMap(maps.first);
  }

  /// Get accounts by type.
  Future<List<Account>> getByType(String type) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'accounts',
      where: 'type = ?',
      whereArgs: [type],
      orderBy: 'created_at DESC',
    );
    return maps.map((map) => Account.fromMap(map)).toList();
  }

  /// Get total balance across all accounts.
  Future<double> getTotalBalance() async {
    final db = await _dbHelper.database;
    final result = await db.rawQuery('SELECT SUM(balance) as total FROM accounts');
    final total = result.first['total'];
    return total != null ? (total as num).toDouble() : 0.0;
  }

  /// Update account balance.
  Future<int> updateBalance(int id, double newBalance) async {
    final db = await _dbHelper.database;
    return await db.update(
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