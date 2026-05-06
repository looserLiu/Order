import 'package:sqflite/sqflite.dart' hide Transaction;
import '../database/database_helper.dart';
import '../models/transaction.dart';

/// Repository for Transaction data operations.
class TransactionRepository {
  final DatabaseHelper _dbHelper;

  TransactionRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new transaction.
  Future<int> insert(Transaction transaction) async {
    final db = await _dbHelper.database;
    return await db.insert('transactions', transaction.toMap());
  }

  /// Update an existing transaction.
  Future<int> update(Transaction transaction) async {
    final db = await _dbHelper.database;
    return await db.update(
      'transactions',
      transaction.toMap(),
      where: 'id = ?',
      whereArgs: [transaction.id],
    );
  }

  /// Delete a transaction by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'transactions',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all transactions.
  Future<List<Transaction>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('transactions', orderBy: 'date DESC');
    return maps.map((map) => Transaction.fromMap(map)).toList();
  }

  /// Get transaction by ID.
  Future<Transaction?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return Transaction.fromMap(maps.first);
  }

  /// Get transactions by account ID.
  Future<List<Transaction>> getByAccountId(int accountId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'account_id = ?',
      whereArgs: [accountId],
      orderBy: 'date DESC',
    );
    return maps.map((map) => Transaction.fromMap(map)).toList();
  }

  /// Get transactions by category ID.
  Future<List<Transaction>> getByCategoryId(int categoryId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'category_id = ?',
      whereArgs: [categoryId],
      orderBy: 'date DESC',
    );
    return maps.map((map) => Transaction.fromMap(map)).toList();
  }

  /// Get transactions by type (income/expense).
  Future<List<Transaction>> getByType(String type) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'type = ?',
      whereArgs: [type],
      orderBy: 'date DESC',
    );
    return maps.map((map) => Transaction.fromMap(map)).toList();
  }

  /// Get transactions within a date range.
  Future<List<Transaction>> getByDateRange(int startDate, int endDate) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'date >= ? AND date <= ?',
      whereArgs: [startDate, endDate],
      orderBy: 'date DESC',
    );
    return maps.map((map) => Transaction.fromMap(map)).toList();
  }

  /// Get transactions with filters.
  Future<List<Transaction>> getWithFilters({
    int? accountId,
    int? categoryId,
    String? type,
    int? startDate,
    int? endDate,
    int? limit,
    int? offset,
  }) async {
    final db = await _dbHelper.database;
    final whereConditions = <String>[];
    final whereArgs = <dynamic>[];

    if (accountId != null) {
      whereConditions.add('account_id = ?');
      whereArgs.add(accountId);
    }
    if (categoryId != null) {
      whereConditions.add('category_id = ?');
      whereArgs.add(categoryId);
    }
    if (type != null) {
      whereConditions.add('type = ?');
      whereArgs.add(type);
    }
    if (startDate != null) {
      whereConditions.add('date >= ?');
      whereArgs.add(startDate);
    }
    if (endDate != null) {
      whereConditions.add('date <= ?');
      whereArgs.add(endDate);
    }

    final maps = await db.query(
      'transactions',
      where: whereConditions.isEmpty ? null : whereConditions.join(' AND '),
      whereArgs: whereArgs.isEmpty ? null : whereArgs,
      orderBy: 'date DESC',
      limit: limit,
      offset: offset,
    );
    return maps.map((map) => Transaction.fromMap(map)).toList();
  }

  /// Get total amount by type within date range.
  Future<double> getTotalByTypeAndDateRange(
    String type,
    int startDate,
    int endDate,
  ) async {
    final db = await _dbHelper.database;
    final result = await db.rawQuery(
      'SELECT SUM(amount) as total FROM transactions WHERE type = ? AND date >= ? AND date <= ?',
      [type, startDate, endDate],
    );
    final total = result.first['total'];
    return total != null ? (total as num).toDouble() : 0.0;
  }

  /// Get total income.
  Future<double> getTotalIncome({int? startDate, int? endDate}) async {
    final db = await _dbHelper.database;
    String query = 'SELECT SUM(amount) as total FROM transactions WHERE type = ?';
    final args = <dynamic>['income'];

    if (startDate != null) {
      query += ' AND date >= ?';
      args.add(startDate);
    }
    if (endDate != null) {
      query += ' AND date <= ?';
      args.add(endDate);
    }

    final result = await db.rawQuery(query, args);
    final total = result.first['total'];
    return total != null ? (total as num).toDouble() : 0.0;
  }

  /// Get total expense.
  Future<double> getTotalExpense({int? startDate, int? endDate}) async {
    final db = await _dbHelper.database;
    String query = 'SELECT SUM(amount) as total FROM transactions WHERE type = ?';
    final args = <dynamic>['expense'];

    if (startDate != null) {
      query += ' AND date >= ?';
      args.add(startDate);
    }
    if (endDate != null) {
      query += ' AND date <= ?';
      args.add(endDate);
    }

    final result = await db.rawQuery(query, args);
    final total = result.first['total'];
    return total != null ? (total as num).toDouble() : 0.0;
  }

  /// Search transactions by description.
  Future<List<Transaction>> searchByDescription(String query) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'description LIKE ?',
      whereArgs: ['%$query%'],
      orderBy: 'date DESC',
    );
    return maps.map((map) => Transaction.fromMap(map)).toList();
  }
}