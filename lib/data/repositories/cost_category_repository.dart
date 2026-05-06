import 'package:sqflite/sqflite.dart';
import '../database/database_helper.dart';
import '../models/cost_category.dart';

/// Repository for CostCategory data operations.
class CostCategoryRepository {
  final DatabaseHelper _dbHelper;

  CostCategoryRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new cost category.
  Future<int> insert(CostCategory costCategory) async {
    final db = await _dbHelper.database;
    return await db.insert('cost_categories', costCategory.toMap());
  }

  /// Update an existing cost category.
  Future<int> update(CostCategory costCategory) async {
    final db = await _dbHelper.database;
    return await db.update(
      'cost_categories',
      costCategory.toMap(),
      where: 'id = ?',
      whereArgs: [costCategory.id],
    );
  }

  /// Delete a cost category by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'cost_categories',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all cost categories.
  Future<List<CostCategory>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('cost_categories', orderBy: 'name ASC');
    return maps.map((map) => CostCategory.fromMap(map)).toList();
  }

  /// Get cost category by ID.
  Future<CostCategory?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'cost_categories',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return CostCategory.fromMap(maps.first);
  }

  /// Get cost categories by type.
  Future<List<CostCategory>> getByType(String type) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'cost_categories',
      where: 'type = ?',
      whereArgs: [type],
      orderBy: 'name ASC',
    );
    return maps.map((map) => CostCategory.fromMap(map)).toList();
  }

  /// Search cost categories by name.
  Future<List<CostCategory>> searchByName(String query) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'cost_categories',
      where: 'name LIKE ?',
      whereArgs: ['%$query%'],
      orderBy: 'name ASC',
    );
    return maps.map((map) => CostCategory.fromMap(map)).toList();
  }
}