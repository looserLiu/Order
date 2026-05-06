import 'package:sqflite/sqflite.dart';
import '../database/database_helper.dart';
import '../models/budget.dart';

/// Repository for Budget data operations.
class BudgetRepository {
  final DatabaseHelper _dbHelper;

  BudgetRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new budget.
  Future<int> insert(Budget budget) async {
    final db = await _dbHelper.database;
    return await db.insert('budgets', budget.toMap());
  }

  /// Update an existing budget.
  Future<int> update(Budget budget) async {
    final db = await _dbHelper.database;
    return await db.update(
      'budgets',
      budget.toMap(),
      where: 'id = ?',
      whereArgs: [budget.id],
    );
  }

  /// Delete a budget by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'budgets',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all budgets.
  Future<List<Budget>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('budgets', orderBy: 'start_date DESC');
    return maps.map((map) => Budget.fromMap(map)).toList();
  }

  /// Get budget by ID.
  Future<Budget?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'budgets',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return Budget.fromMap(maps.first);
  }

  /// Get budgets by category ID.
  Future<List<Budget>> getByCategoryId(int categoryId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'budgets',
      where: 'category_id = ?',
      whereArgs: [categoryId],
      orderBy: 'start_date DESC',
    );
    return maps.map((map) => Budget.fromMap(map)).toList();
  }

  /// Get budgets by period.
  Future<List<Budget>> getByPeriod(String period) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'budgets',
      where: 'period = ?',
      whereArgs: [period],
      orderBy: 'start_date DESC',
    );
    return maps.map((map) => Budget.fromMap(map)).toList();
  }

  /// Get active budgets (current date is within start and end date).
  Future<List<Budget>> getActive() async {
    final db = await _dbHelper.database;
    final now = DateTime.now().millisecondsSinceEpoch;
    final maps = await db.query(
      'budgets',
      where: 'start_date <= ? AND end_date >= ?',
      whereArgs: [now, now],
    );
    return maps.map((map) => Budget.fromMap(map)).toList();
  }

  /// Get budget for a specific category and month.
  Future<Budget?> getByCategoryAndMonth(int categoryId, int month, int year) async {
    final db = await _dbHelper.database;

    // Calculate start and end of month
    final startOfMonth = DateTime(year, month, 1).millisecondsSinceEpoch;
    final endOfMonth = DateTime(year, month + 1, 0, 23, 59, 59).millisecondsSinceEpoch;

    final maps = await db.query(
      'budgets',
      where: 'category_id = ? AND start_date <= ? AND end_date >= ?',
      whereArgs: [categoryId, endOfMonth, startOfMonth],
    );

    if (maps.isEmpty) return null;
    return Budget.fromMap(maps.first);
  }

  /// Get budgets within a date range.
  Future<List<Budget>> getByDateRange(int startDate, int endDate) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'budgets',
      where: 'start_date <= ? AND end_date >= ?',
      whereArgs: [endDate, startDate],
      orderBy: 'start_date DESC',
    );
    return maps.map((map) => Budget.fromMap(map)).toList();
  }

  /// Get total budget amount by period.
  Future<double> getTotalBudgetByPeriod(String period) async {
    final db = await _dbHelper.database;
    final result = await db.rawQuery(
      'SELECT SUM(amount) as total FROM budgets WHERE period = ?',
      [period],
    );
    final total = result.first['total'];
    return total != null ? (total as num).toDouble() : 0.0;
  }
}