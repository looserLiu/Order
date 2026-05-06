import '../../domain/entities/account.dart';
import '../../domain/entities/category.dart';
import '../../domain/entities/transaction.dart';
import '../../domain/entities/budget.dart';
import '../../domain/entities/product.dart';
import '../../domain/entities/warehouse.dart';
import '../../domain/entities/inventory_flow.dart';

/// Mapper for converting between Domain entities and Data models.
class EntityMappers {
  // ==================== Account ====================
  static AccountEntity accountFromMap(Map<String, dynamic> map) {
    return AccountEntity(
      id: map['id'] as int?,
      name: map['name'] as String,
      type: AccountType.fromValue(map['type'] as String),
      balance: (map['balance'] as num?)?.toDouble() ?? 0.0,
      currency: map['currency'] as String? ?? 'CNY',
      icon: map['icon'] as String?,
      color: map['color'] as int?,
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at'] as int),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(map['updated_at'] as int),
    );
  }

  static Map<String, dynamic> accountToMap(AccountEntity entity) {
    return {
      if (entity.id != null) 'id': entity.id,
      'name': entity.name,
      'type': entity.type.value,
      'balance': entity.balance,
      'currency': entity.currency,
      'icon': entity.icon,
      'color': entity.color,
      'created_at': entity.createdAt.millisecondsSinceEpoch,
      'updated_at': entity.updatedAt.millisecondsSinceEpoch,
    };
  }

  // ==================== Category ====================
  static CategoryEntity categoryFromMap(Map<String, dynamic> map) {
    return CategoryEntity(
      id: map['id'] as int?,
      name: map['name'] as String,
      type: CategoryType.fromValue(map['type'] as String),
      icon: map['icon'] as String?,
      color: map['color'] as int?,
      parentId: map['parent_id'] as int?,
      isSystem: (map['is_system'] as int?) == 1,
      usageCount: map['usage_count'] as int? ?? 0,
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at'] as int),
    );
  }

  static Map<String, dynamic> categoryToMap(CategoryEntity entity) {
    return {
      if (entity.id != null) 'id': entity.id,
      'name': entity.name,
      'type': entity.type.value,
      'icon': entity.icon,
      'color': entity.color,
      'parent_id': entity.parentId,
      'is_system': entity.isSystem ? 1 : 0,
      'usage_count': entity.usageCount,
      'created_at': entity.createdAt.millisecondsSinceEpoch,
    };
  }

  // ==================== Transaction ====================
  static TransactionEntity transactionFromMap(Map<String, dynamic> map) {
    return TransactionEntity(
      id: map['id'] as int?,
      accountId: map['account_id'] as int,
      categoryId: map['category_id'] as int?,
      amount: (map['amount'] as num).toDouble(),
      type: TransactionType.fromValue(map['type'] as String),
      description: map['description'] as String?,
      date: DateTime.fromMillisecondsSinceEpoch(map['date'] as int),
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at'] as int),
    );
  }

  static Map<String, dynamic> transactionToMap(TransactionEntity entity) {
    return {
      if (entity.id != null) 'id': entity.id,
      'account_id': entity.accountId,
      'category_id': entity.categoryId,
      'amount': entity.amount,
      'type': entity.type.value,
      'description': entity.description,
      'date': entity.date.millisecondsSinceEpoch,
      'created_at': entity.createdAt.millisecondsSinceEpoch,
    };
  }

  // ==================== Budget ====================
  static BudgetEntity budgetFromMap(Map<String, dynamic> map) {
    return BudgetEntity(
      id: map['id'] as int?,
      categoryId: map['category_id'] as int,
      amount: (map['amount'] as num).toDouble(),
      period: BudgetPeriod.fromValue(map['period'] as String),
      startDate: DateTime.fromMillisecondsSinceEpoch(map['start_date'] as int),
      endDate: DateTime.fromMillisecondsSinceEpoch(map['end_date'] as int),
    );
  }

  static Map<String, dynamic> budgetToMap(BudgetEntity entity) {
    return {
      if (entity.id != null) 'id': entity.id,
      'category_id': entity.categoryId,
      'amount': entity.amount,
      'period': entity.period.value,
      'start_date': entity.startDate.millisecondsSinceEpoch,
      'end_date': entity.endDate.millisecondsSinceEpoch,
    };
  }

  // ==================== Product ====================
  static ProductEntity productFromMap(Map<String, dynamic> map) {
    return ProductEntity(
      id: map['id'] as int?,
      name: map['name'] as String,
      sku: map['sku'] as String?,
      category: map['category'] as String?,
      unit: map['unit'] as String?,
      costPrice: (map['cost_price'] as num).toDouble(),
      salePrice: (map['sale_price'] as num).toDouble(),
      imageUrl: map['image_url'] as String?,
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at'] as int),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(map['updated_at'] as int),
    );
  }

  static Map<String, dynamic> productToMap(ProductEntity entity) {
    return {
      if (entity.id != null) 'id': entity.id,
      'name': entity.name,
      'sku': entity.sku,
      'category': entity.category,
      'unit': entity.unit,
      'cost_price': entity.costPrice,
      'sale_price': entity.salePrice,
      'image_url': entity.imageUrl,
      'created_at': entity.createdAt.millisecondsSinceEpoch,
      'updated_at': entity.updatedAt.millisecondsSinceEpoch,
    };
  }

  // ==================== Warehouse ====================
  static WarehouseEntity warehouseFromMap(Map<String, dynamic> map) {
    return WarehouseEntity(
      id: map['id'] as int?,
      name: map['name'] as String,
      location: map['location'] as String?,
      description: map['description'] as String?,
      isActive: (map['is_active'] as int?) == 1,
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at'] as int),
    );
  }

  static Map<String, dynamic> warehouseToMap(WarehouseEntity entity) {
    return {
      if (entity.id != null) 'id': entity.id,
      'name': entity.name,
      'location': entity.location,
      'description': entity.description,
      'is_active': entity.isActive ? 1 : 0,
      'created_at': entity.createdAt.millisecondsSinceEpoch,
    };
  }

  // ==================== InventoryFlow ====================
  static InventoryFlowEntity inventoryFlowFromMap(Map<String, dynamic> map) {
    return InventoryFlowEntity(
      id: map['id'] as int?,
      productId: map['product_id'] as int,
      warehouseId: map['warehouse_id'] as int,
      flowType: FlowType.fromValue(map['flow_type'] as String),
      quantity: (map['quantity'] as num).toDouble(),
      batchNumber: map['batch_number'] as String?,
      expirationDate: map['expiration_date'] != null
          ? DateTime.fromMillisecondsSinceEpoch(map['expiration_date'] as int)
          : null,
      costAtFlow: (map['cost_at_flow'] as num?)?.toDouble(),
      referenceId: map['reference_id'] as String?,
      date: DateTime.fromMillisecondsSinceEpoch(map['date'] as int),
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at'] as int),
    );
  }

  static Map<String, dynamic> inventoryFlowToMap(InventoryFlowEntity entity) {
    return {
      if (entity.id != null) 'id': entity.id,
      'product_id': entity.productId,
      'warehouse_id': entity.warehouseId,
      'flow_type': entity.flowType.value,
      'quantity': entity.quantity,
      'batch_number': entity.batchNumber,
      'expiration_date': entity.expirationDate?.millisecondsSinceEpoch,
      'cost_at_flow': entity.costAtFlow,
      'reference_id': entity.referenceId,
      'date': entity.date.millisecondsSinceEpoch,
      'created_at': entity.createdAt.millisecondsSinceEpoch,
    };
  }
}
