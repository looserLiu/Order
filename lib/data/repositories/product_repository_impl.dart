import '../../domain/entities/product.dart';
import '../../domain/repositories/i_product_repository.dart';
import '../database/database_helper.dart';
import '../models/product.dart';
import '../mappers/product_mapper.dart';

/// Implementation of IProductRepository using SQLite.
class ProductRepositoryImpl implements IProductRepository {
  final DatabaseHelper _dbHelper;

  ProductRepositoryImpl({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  @override
  Future<List<ProductEntity>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('products', orderBy: 'created_at DESC');
    final products = maps.map((map) => Product.fromMap(map)).toList();
    return ProductMapper.toEntityList(products);
  }

  @override
  Future<ProductEntity?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'products',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return ProductMapper.toEntity(Product.fromMap(maps.first));
  }

  @override
  Future<List<ProductEntity>> getByCategory(String category) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'products',
      where: 'category = ?',
      whereArgs: [category],
      orderBy: 'name ASC',
    );
    final products = maps.map((map) => Product.fromMap(map)).toList();
    return ProductMapper.toEntityList(products);
  }

  @override
  Future<List<ProductEntity>> searchByName(String query) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'products',
      where: 'name LIKE ?',
      whereArgs: ['%$query%'],
      orderBy: 'name ASC',
    );
    final products = maps.map((map) => Product.fromMap(map)).toList();
    return ProductMapper.toEntityList(products);
  }

  @override
  Future<ProductEntity> create(ProductEntity product) async {
    final db = await _dbHelper.database;
    final model = ProductMapper.toModel(product);
    final id = await db.insert('products', model.toMap());
    return product.copyWith(id: id);
  }

  @override
  Future<ProductEntity> update(ProductEntity product) async {
    final db = await _dbHelper.database;
    final model = ProductMapper.toModel(product);
    await db.update(
      'products',
      model.toMap(),
      where: 'id = ?',
      whereArgs: [product.id],
    );
    return product;
  }

  @override
  Future<void> delete(int id) async {
    final db = await _dbHelper.database;
    await db.delete(
      'products',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
