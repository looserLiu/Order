/// Domain entity for InventoryFlow.
class InventoryFlowEntity {
  final int? id;
  final int productId;
  final int warehouseId;
  final FlowType flowType;
  final double quantity;
  final String? batchNumber;
  final DateTime? expirationDate;
  final double? costAtFlow;
  final String? referenceId;
  final DateTime date;
  final DateTime createdAt;

  InventoryFlowEntity({
    this.id,
    required this.productId,
    required this.warehouseId,
    required this.flowType,
    required this.quantity,
    this.batchNumber,
    this.expirationDate,
    this.costAtFlow,
    this.referenceId,
    required this.date,
    required this.createdAt,
  });

  InventoryFlowEntity copyWith({
    int? id,
    int? productId,
    int? warehouseId,
    FlowType? flowType,
    double? quantity,
    String? batchNumber,
    DateTime? expirationDate,
    double? costAtFlow,
    String? referenceId,
    DateTime? date,
    DateTime? createdAt,
  }) {
    return InventoryFlowEntity(
      id: id ?? this.id,
      productId: productId ?? this.productId,
      warehouseId: warehouseId ?? this.warehouseId,
      flowType: flowType ?? this.flowType,
      quantity: quantity ?? this.quantity,
      batchNumber: batchNumber ?? this.batchNumber,
      expirationDate: expirationDate ?? this.expirationDate,
      costAtFlow: costAtFlow ?? this.costAtFlow,
      referenceId: referenceId ?? this.referenceId,
      date: date ?? this.date,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  bool get isStockIn => flowType == FlowType.in;
  bool get isStockOut => flowType == FlowType.out;
  bool get hasExpiration => expirationDate != null;

  /// Check if expired.
  bool get isExpired {
    if (expirationDate == null) return false;
    return DateTime.now().isAfter(expirationDate!);
  }

  /// Days until expiration (negative if expired).
  int? get daysUntilExpiration {
    if (expirationDate == null) return null;
    return expirationDate!.difference(DateTime.now()).inDays;
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is InventoryFlowEntity &&
          runtimeType == other.runtimeType &&
          id == other.id;

  @override
  int get hashCode => id.hashCode;
}

/// Inventory flow type enumeration.
enum FlowType {
  inFlow('in', '入库'),
  outFlow('out', '出库'),
  transfer('transfer', '调拨'),
  adjust('adjust', '调整');

  final String value;
  final String displayName;

  const FlowType(this.value, this.displayName);

  static FlowType fromValue(String value) {
    return FlowType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => FlowType.inFlow,
    );
  }
}
