/// InventoryFlow model representing stock movements.
class InventoryFlow {
  final int? id;
  final int productId;
  final int warehouseId;
  final String flowType;
  final double quantity;
  final String? batchNumber;
  final int? expirationDate;
  final double? costAtFlow;
  final String? referenceId;
  final int date;
  final int createdAt;

  InventoryFlow({
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

  /// Create InventoryFlow from database map.
  factory InventoryFlow.fromMap(Map<String, dynamic> map) {
    return InventoryFlow(
      id: map['id'] as int?,
      productId: map['product_id'] as int,
      warehouseId: map['warehouse_id'] as int,
      flowType: map['flow_type'] as String,
      quantity: (map['quantity'] as num).toDouble(),
      batchNumber: map['batch_number'] as String?,
      expirationDate: map['expiration_date'] as int?,
      costAtFlow: (map['cost_at_flow'] as num?)?.toDouble(),
      referenceId: map['reference_id'] as String?,
      date: map['date'] as int,
      createdAt: map['created_at'] as int,
    );
  }

  /// Convert InventoryFlow to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'product_id': productId,
      'warehouse_id': warehouseId,
      'flow_type': flowType,
      'quantity': quantity,
      'batch_number': batchNumber,
      'expiration_date': expirationDate,
      'cost_at_flow': costAtFlow,
      'reference_id': referenceId,
      'date': date,
      'created_at': createdAt,
    };
  }

  /// Create a copy of InventoryFlow with optional field updates.
  InventoryFlow copyWith({
    int? id,
    int? productId,
    int? warehouseId,
    String? flowType,
    double? quantity,
    String? batchNumber,
    int? expirationDate,
    double? costAtFlow,
    String? referenceId,
    int? date,
    int? createdAt,
  }) {
    return InventoryFlow(
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

  /// Flow type values.
  static const String typeIn = 'in';
  static const String typeOut = 'out';
  static const String typeTransfer = 'transfer';
  static const String typeAdjust = 'adjust';

  /// Get all valid flow types.
  static List<String> get types => [typeIn, typeOut, typeTransfer, typeAdjust];
}