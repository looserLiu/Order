/// CostCategory model for inventory cost classification.
class CostCategory {
  final int? id;
  final String name;
  final String type;
  final String? description;

  CostCategory({
    this.id,
    required this.name,
    required this.type,
    this.description,
  });

  /// Create CostCategory from database map.
  factory CostCategory.fromMap(Map<String, dynamic> map) {
    return CostCategory(
      id: map['id'] as int?,
      name: map['name'] as String,
      type: map['type'] as String,
      description: map['description'] as String?,
    );
  }

  /// Convert CostCategory to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'name': name,
      'type': type,
      'description': description,
    };
  }

  /// Create a copy of CostCategory with optional field updates.
  CostCategory copyWith({
    int? id,
    String? name,
    String? type,
    String? description,
  }) {
    return CostCategory(
      id: id ?? this.id,
      name: name ?? this.name,
      type: type ?? this.type,
      description: description ?? this.description,
    );
  }

  /// Cost category type values.
  static const String typePurchase = 'purchase';
  static const String typeStorage = 'storage';
  static const String typeTransport = 'transport';
  static const String typeOther = 'other';

  /// Get all valid types.
  static List<String> get types => [typePurchase, typeStorage, typeTransport, typeOther];
}