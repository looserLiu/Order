import 'package:sqflite/sqflite.dart';
import '../database/database_helper.dart';
import '../models/product.dart';

/// Repository for Product data operations.
class ProductRepository {
  final DatabaseHelper _dbHelper;

  ProductRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new product.
  Future<int> insert(Product product) async {
    final db = await _dbHelper.database;
    return await db.insert('products', product.toMap());
  }

  /// Update an existing product.
  Future<int> update(Product product) async {
    final db = await _dbHelper.database;
    return await db.update(
      'products',
      product.toMap(),
      where: 'id = ?',
      whereArgs: [product.id],
    );
  }

  /// Delete a product by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'products',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all products.
  Future<List<Product>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('products', orderBy: 'created_at DESC');
    return maps.map((map) => Product.fromMap(map)).toList();
  }

  /// Get product by ID.
  Future<Product?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'products',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return Product.fromMap(maps.first);
  }

  /// Get products by category.
  Future<List<Product>> getByCategory(String category) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'products',
      where: 'category = ?',
      whereArgs: [category],
      orderBy: 'name ASC',
    );
    return maps.map((map) => Product.fromMap(map)).toList();
  }

  /// Search products by name.
  Future<List<Product>> searchByName(String query) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'products',
      where: 'name LIKE ?',
      whereArgs: ['%$query%'],
      orderBy: 'name ASC',
    );
    return maps.map((map) => Product.fromMap(map)).toList();
  }

  /// Get product by SKU.
  Future<Product?> getBySku(String sku) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'products',
      where: 'sku = ?',
      whereArgs: [sku],
    );
    if (maps.isEmpty) return null;
    return Product.fromMap(maps.first);
  }

  /// Get all unique product categories.
  Future<List<String>> getAllCategories() async {
    final db = await _dbHelper.database;
    final result = await db.rawQuery(
      'SELECT DISTINCT category FROM products WHERE category IS NOT NULL ORDER BY category ASC',
    );
    return result.map((row) => row['category'] as String).toList();
  }

  /// Get total product count.
  Future<int> getCount() async {
    final db = await _dbHelper.database;
    final result = await db.rawQuery('SELECT COUNT(*) as count FROM products');
    return (result.first['count'] as int?) ?? 0;
  }
}