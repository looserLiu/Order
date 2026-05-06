import '../entities/category.dart';
import '../repositories/i_category_repository.dart';
import 'base_usecase.dart';

/// Use case: Get all categories.
class GetAllCategories extends NoParamsUseCase<List<CategoryEntity>> {
  final ICategoryRepository _repository;

  GetAllCategories(this._repository);

  @override
  Future<List<CategoryEntity>> call() {
    return _repository.getAll();
  }
}

/// Use case: Get category by ID.
class GetCategoryById extends UseCase<CategoryEntity?, int> {
  final ICategoryRepository _repository;

  GetCategoryById(this._repository);

  @override
  Future<CategoryEntity?> call(int id) {
    return _repository.getById(id);
  }
}

/// Use case: Get categories by type.
class GetCategoriesByType extends UseCase<List<CategoryEntity>, CategoryType> {
  final ICategoryRepository _repository;

  GetCategoriesByType(this._repository);

  @override
  Future<List<CategoryEntity>> call(CategoryType type) {
    return _repository.getByType(type);
  }
}

/// Use case: Create a new category.
class CreateCategory extends UseCase<CategoryEntity, CategoryEntity> {
  final ICategoryRepository _repository;

  CreateCategory(this._repository);

  @override
  Future<CategoryEntity> call(CategoryEntity category) {
    return _repository.create(category);
  }
}

/// Use case: Update an existing category.
class UpdateCategory extends UseCase<CategoryEntity, CategoryEntity> {
  final ICategoryRepository _repository;

  UpdateCategory(this._repository);

  @override
  Future<CategoryEntity> call(CategoryEntity category) {
    return _repository.update(category);
  }
}

/// Use case: Delete a category.
class DeleteCategory extends UseCase<void, int> {
  final ICategoryRepository _repository;

  DeleteCategory(this._repository);

  @override
  Future<void> call(int id) {
    return _repository.delete(id);
  }
}

/// Use case: Increment category usage count.
class IncrementCategoryUsage extends UseCase<void, int> {
  final ICategoryRepository _repository;

  IncrementCategoryUsage(this._repository);

  @override
  Future<void> call(int id) {
    return _repository.incrementUsageCount(id);
  }
}

/// Use case: Get top categories by usage.
class GetTopCategories extends UseCase<List<CategoryEntity>, int> {
  final ICategoryRepository _repository;

  GetTopCategories(this._repository);

  @override
  Future<List<CategoryEntity>> call(int limit) {
    return _repository.getTopByUsage(limit: limit);
  }
}

/// Use case: Search categories by name.
class SearchCategories extends UseCase<List<CategoryEntity>, String> {
  final ICategoryRepository _repository;

  SearchCategories(this._repository);

  @override
  Future<List<CategoryEntity>> call(String query) {
    return _repository.searchByName(query);
  }
}
