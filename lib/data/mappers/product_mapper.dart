import '../../domain/entities/product.dart';
import '../models/product.dart';

/// Mapper for converting between ProductEntity (domain) and Product (data model).
class ProductMapper {
  /// Convert data model to domain entity.
  static ProductEntity toEntity(Product model) {
    return ProductEntity(
      id: model.id,
      name: model.name,
      sku: model.sku,
      category: model.category,
      unit: model.unit,
      costPrice: model.costPrice,
      salePrice: model.salePrice,
      imageUrl: model.imageUrl,
      createdAt: DateTime.fromMillisecondsSinceEpoch(model.createdAt),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(model.updatedAt),
    );
  }

  /// Convert domain entity to data model.
  static Product toModel(ProductEntity entity) {
    return Product(
      id: entity.id,
      name: entity.name,
      sku: entity.sku,
      category: entity.category,
      unit: entity.unit,
      costPrice: entity.costPrice,
      salePrice: entity.salePrice,
      imageUrl: entity.imageUrl,
      createdAt: entity.createdAt.millisecondsSinceEpoch,
      updatedAt: entity.updatedAt.millisecondsSinceEpoch,
    );
  }

  /// Convert list of data models to domain entities.
  static List<ProductEntity> toEntityList(List<Product> models) {
    return models.map(toEntity).toList();
  }
}
