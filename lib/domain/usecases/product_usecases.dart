import '../entities/product.dart';
import '../repositories/i_product_repository.dart';
import 'base_usecase.dart';

/// Use case: Get all products.
class GetAllProducts extends NoParamsUseCase<List<ProductEntity>> {
  final IProductRepository _repository;

  GetAllProducts(this._repository);

  @override
  Future<List<ProductEntity>> call() {
    return _repository.getAll();
  }
}

/// Use case: Get product by ID.
class GetProductById extends UseCase<ProductEntity?, int> {
  final IProductRepository _repository;

  GetProductById(this._repository);

  @override
  Future<ProductEntity?> call(int id) {
    return _repository.getById(id);
  }
}

/// Use case: Get products by category.
class GetProductsByCategory extends UseCase<List<ProductEntity>, String> {
  final IProductRepository _repository;

  GetProductsByCategory(this._repository);

  @override
  Future<List<ProductEntity>> call(String category) {
    return _repository.getByCategory(category);
  }
}

/// Use case: Search products by name.
class SearchProducts extends UseCase<List<ProductEntity>, String> {
  final IProductRepository _repository;

  SearchProducts(this._repository);

  @override
  Future<List<ProductEntity>> call(String query) {
    return _repository.searchByName(query);
  }
}

/// Use case: Create a new product.
class CreateProduct extends UseCase<ProductEntity, ProductEntity> {
  final IProductRepository _repository;

  CreateProduct(this._repository);

  @override
  Future<ProductEntity> call(ProductEntity product) {
    return _repository.create(product);
  }
}

/// Use case: Update an existing product.
class UpdateProduct extends UseCase<ProductEntity, ProductEntity> {
  final IProductRepository _repository;

  UpdateProduct(this._repository);

  @override
  Future<ProductEntity> call(ProductEntity product) {
    return _repository.update(product);
  }
}

/// Use case: Delete a product.
class DeleteProduct extends UseCase<void, int> {
  final IProductRepository _repository;

  DeleteProduct(this._repository);

  @override
  Future<void> call(int id) {
    return _repository.delete(id);
  }
}
