import '../../domain/entities/category.dart';
import '../../domain/repositories/i_category_repository.dart';
import '../database/database_helper.dart';
import '../models/category.dart';
import '../mappers/category_mapper.dart';

/// Implementation of ICategoryRepository using SQLite.
class CategoryRepositoryImpl implements ICategoryRepository {
  final DatabaseHelper _dbHelper;

  CategoryRepositoryImpl({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  @override
  Future<List<CategoryEntity>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('categories', orderBy: 'usage_count DESC');
    final categories = maps.map((map) => Category.fromMap(map)).toList();
    return CategoryMapper.toEntityList(categories);
  }

  @override
  Future<CategoryEntity?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return CategoryMapper.toEntity(Category.fromMap(maps.first));
  }

  @override
  Future<List<CategoryEntity>> getByType(CategoryType type) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'type = ?',
      whereArgs: [type.value],
      orderBy: 'usage_count DESC',
    );
    final categories = maps.map((map) => Category.fromMap(map)).toList();
    return CategoryMapper.toEntityList(categories);
  }

  @override
  Future<List<CategoryEntity>> getSubcategories(int parentId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'parent_id = ?',
      whereArgs: [parentId],
    );
    final categories = maps.map((map) => Category.fromMap(map)).toList();
    return CategoryMapper.toEntityList(categories);
  }

  @override
  Future<CategoryEntity> create(CategoryEntity category) async {
    final db = await _dbHelper.database;
    final model = CategoryMapper.toModel(category);
    final id = await db.insert('categories', model.toMap());
    return category.copyWith(id: id);
  }

  @override
  Future<CategoryEntity> update(CategoryEntity category) async {
    final db = await _dbHelper.database;
    final model = CategoryMapper.toModel(category);
    await db.update(
      'categories',
      model.toMap(),
      where: 'id = ?',
      whereArgs: [category.id],
    );
    return category;
  }

  @override
  Future<void> delete(int id) async {
    final db = await _dbHelper.database;
    await db.delete(
      'categories',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  @override
  Future<void> incrementUsageCount(int id) async {
    final db = await _dbHelper.database;
    await db.rawUpdate(
      'UPDATE categories SET usage_count = usage_count + 1 WHERE id = ?',
      [id],
    );
  }

  @override
  Future<List<CategoryEntity>> getTopByUsage({int limit = 5}) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      orderBy: 'usage_count DESC',
      limit: limit,
    );
    final categories = maps.map((map) => Category.fromMap(map)).toList();
    return CategoryMapper.toEntityList(categories);
  }

  @override
  Future<List<CategoryEntity>> searchByName(String query) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'categories',
      where: 'name LIKE ?',
      whereArgs: ['%$query%'],
    );
    final categories = maps.map((map) => Category.fromMap(map)).toList();
    return CategoryMapper.toEntityList(categories);
  }
}
