import 'package:flutter/foundation.dart' hide Category;
import '../../data/database/database_helper.dart';
import '../../data/models/category.dart';

class CategoryProvider extends ChangeNotifier {
  final DatabaseHelper _db = DatabaseHelper();
  List<Category> _categories = [];
  bool _isLoading = false;

  List<Category> get categories => _categories;
  bool get isLoading => _isLoading;

  Future<void> loadCategories() async {
    _isLoading = true;
    notifyListeners();

    try {
      final db = await _db.database;
      final List<Map<String, dynamic>> maps = await db.query(
        'categories',
        orderBy: 'usage_count DESC, name ASC',
      );
      _categories = maps.map((map) => Category.fromMap(map)).toList();
    } catch (e) {
      debugPrint('Error loading categories: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> addCategory(Category category) async {
    try {
      final db = await _db.database;
      await db.insert('categories', category.toMap());
      await loadCategories();
    } catch (e) {
      debugPrint('Error adding category: $e');
    }
  }

  Future<void> updateCategory(Category category) async {
    try {
      final db = await _db.database;
      await db.update(
        'categories',
        category.toMap(),
        where: 'id = ?',
        whereArgs: [category.id],
      );
      await loadCategories();
    } catch (e) {
      debugPrint('Error updating category: $e');
    }
  }

  Future<void> deleteCategory(int id) async {
    try {
      final db = await _db.database;
      await db.delete(
        'categories',
        where: 'id = ?',
        whereArgs: [id],
      );
      await loadCategories();
    } catch (e) {
      debugPrint('Error deleting category: $e');
    }
  }

  Future<void> incrementUsageCount(int id) async {
    try {
      final db = await _db.database;
      await db.rawUpdate(
        'UPDATE categories SET usage_count = usage_count + 1 WHERE id = ?',
        [id],
      );
      await loadCategories();
    } catch (e) {
      debugPrint('Error incrementing usage count: $e');
    }
  }
}