import '../../domain/entities/category.dart';
import '../models/category.dart';

/// Mapper for converting between CategoryEntity (domain) and Category (data model).
class CategoryMapper {
  /// Convert data model to domain entity.
  static CategoryEntity toEntity(Category model) {
    return CategoryEntity(
      id: model.id,
      name: model.name,
      type: CategoryType.fromValue(model.type),
      icon: model.icon,
      color: model.color,
      parentId: model.parentId,
      isSystem: model.isSystem == 1, // DB stores as int, convert to bool
      usageCount: model.usageCount,
      createdAt: DateTime.fromMillisecondsSinceEpoch(model.createdAt ?? DateTime.now().millisecondsSinceEpoch),
    );
  }

  /// Convert domain entity to data model.
  static Category toModel(CategoryEntity entity) {
    return Category(
      id: entity.id,
      name: entity.name,
      type: entity.type.value,
      icon: entity.icon,
      color: entity.color,
      parentId: entity.parentId,
      isSystem: entity.isSystem ? 1 : 0, // DB stores as int
      usageCount: entity.usageCount,
      createdAt: entity.createdAt.millisecondsSinceEpoch,
    );
  }

  /// Convert list of data models to domain entities.
  static List<CategoryEntity> toEntityList(List<Category> models) {
    return models.map(toEntity).toList();
  }
}
