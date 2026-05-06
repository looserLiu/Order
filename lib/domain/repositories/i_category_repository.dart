import '../entities/category.dart';

/// Repository interface for Category operations.
abstract class ICategoryRepository {
  /// Get all categories.
  Future<List<CategoryEntity>> getAll();

  /// Get category by ID.
  Future<CategoryEntity?> getById(int id);

  /// Get categories by type.
  Future<List<CategoryEntity>> getByType(CategoryType type);

  /// Get subcategories by parent ID.
  Future<List<CategoryEntity>> getSubcategories(int parentId);

  /// Create a new category.
  Future<CategoryEntity> create(CategoryEntity category);

  /// Update an existing category.
  Future<CategoryEntity> update(CategoryEntity category);

  /// Delete a category by ID.
  Future<void> delete(int id);

  /// Increment category usage count.
  Future<void> incrementUsageCount(int id);

  /// Get top categories by usage.
  Future<List<CategoryEntity>> getTopByUsage({int limit = 5});

  /// Search categories by name.
  Future<List<CategoryEntity>> searchByName(String query);
}
