import '../../domain/entities/inventory_flow.dart';
import '../models/inventory_flow.dart';

/// Mapper for converting between InventoryFlowEntity (domain) and InventoryFlow (data model).
class InventoryFlowMapper {
  /// Convert data model to domain entity.
  static InventoryFlowEntity toEntity(InventoryFlow model) {
    return InventoryFlowEntity(
      id: model.id,
      productId: model.productId,
      warehouseId: model.warehouseId,
      flowType: FlowType.fromValue(model.flowType),
      quantity: model.quantity,
      batchNumber: model.batchNumber,
      expirationDate: model.expirationDate != null
          ? DateTime.fromMillisecondsSinceEpoch(model.expirationDate!)
          : null,
      costAtFlow: model.costAtFlow,
      referenceId: model.referenceId,
      date: DateTime.fromMillisecondsSinceEpoch(model.date),
      createdAt: DateTime.fromMillisecondsSinceEpoch(model.createdAt),
    );
  }

  /// Convert domain entity to data model.
  static InventoryFlow toModel(InventoryFlowEntity entity) {
    return InventoryFlow(
      id: entity.id,
      productId: entity.productId,
      warehouseId: entity.warehouseId,
      flowType: entity.flowType.value,
      quantity: entity.quantity,
      batchNumber: entity.batchNumber,
      expirationDate: entity.expirationDate?.millisecondsSinceEpoch,
      costAtFlow: entity.costAtFlow,
      referenceId: entity.referenceId,
      date: entity.date.millisecondsSinceEpoch,
      createdAt: entity.createdAt.millisecondsSinceEpoch,
    );
  }

  /// Convert list of data models to domain entities.
  static List<InventoryFlowEntity> toEntityList(List<InventoryFlow> models) {
    return models.map(toEntity).toList();
  }
}
