import '../../domain/entities/warehouse.dart';
import '../models/warehouse.dart';

/// Mapper for converting between WarehouseEntity (domain) and Warehouse (data model).
class WarehouseMapper {
  /// Convert data model to domain entity.
  static WarehouseEntity toEntity(Warehouse model) {
    return WarehouseEntity(
      id: model.id,
      name: model.name,
      location: model.location,
      description: model.description,
      isActive: model.isActive == 1,
      createdAt: DateTime.fromMillisecondsSinceEpoch(model.createdAt),
    );
  }

  /// Convert domain entity to data model.
  static Warehouse toModel(WarehouseEntity entity) {
    return Warehouse(
      id: entity.id,
      name: entity.name,
      location: entity.location,
      description: entity.description,
      isActive: entity.isActive ? 1 : 0,
      createdAt: entity.createdAt.millisecondsSinceEpoch,
    );
  }

  /// Convert list of data models to domain entities.
  static List<WarehouseEntity> toEntityList(List<Warehouse> models) {
    return models.map(toEntity).toList();
  }
}
