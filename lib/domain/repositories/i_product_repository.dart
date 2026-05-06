import '../entities/product.dart';

/// Repository interface for Product operations.
abstract class IProductRepository {
  /// Get all products.
  Future<List<ProductEntity>> getAll();

  /// Get product by ID.
  Future<ProductEntity?> getById(int id);

  /// Get products by category.
  Future<List<ProductEntity>> getByCategory(String category);

  /// Search products by name.
  Future<List<ProductEntity>> searchByName(String query);

  /// Create a new product.
  Future<ProductEntity> create(ProductEntity product);

  /// Update an existing product.
  Future<ProductEntity> update(ProductEntity product);

  /// Delete a product by ID.
  Future<void> delete(int id);
}
