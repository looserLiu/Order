import 'package:sqflite/sqflite.dart';
import '../database/database_helper.dart';
import '../models/category.dart';

/// Repository for Category data operations.
class CategoryRepository {
  final DatabaseHelper _dbHelper;

  CategoryRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new category.
  Future<int> insert(Category category) async {
    final db = await _dbHelper.database;
    return await db.insert('categories', category.toMap());
  }

  /// Update an existing category.
  Future<int> update(Category category) async {
    final db = await _dbHelper.database;
    return await db.update(
      'categories',
      category.toMap(),
      where: 'id = ?',
      whereArgs: [category.id],
    );
  }

  /// Delete a category by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'categories',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all categories.
  Future<List<Category>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('categories', orderBy: 'usage_count DESC');
    return maps.map((map) => Category.fromMap(map)).toList();
  }

  /// Get category by ID.
  Future<Category?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return Category.fromMap(maps.first);
  }

  /// Get categories by type (income/expense).
  Future<List<Category>> getByType(String type) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'type = ?',
      whereArgs: [type],
      orderBy: 'usage_count DESC',
    );
    return maps.map((map) => Category.fromMap(map)).toList();
  }

  /// Get subcategories by parent ID.
  Future<List<Category>> getSubcategories(int parentId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'parent_id = ?',
      whereArgs: [parentId],
    );
    return maps.map((map) => Category.fromMap(map)).toList();
  }

  /// Increment category usage count.
  Future<int> incrementUsageCount(int id) async {
    final db = await _dbHelper.database;
    return await db.rawUpdate(
      'UPDATE categories SET usage_count = usage_count + 1 WHERE id = ?',
      [id],
    );
  }

  /// Get top categories by usage.
  Future<List<Category>> getTopByUsage({int limit = 5}) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      orderBy: 'usage_count DESC',
      limit: limit,
    );
    return maps.map((map) => Category.fromMap(map)).toList();
  }

  /// Search categories by name.
  Future<List<Category>> searchByName(String query) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'name LIKE ?',
      whereArgs: ['%$query%'],
    );
    return maps.map((map) => Category.fromMap(map)).toList();
  }
}