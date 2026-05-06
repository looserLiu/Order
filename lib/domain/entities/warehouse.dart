/// Domain entity for Warehouse.
class WarehouseEntity {
  final int? id;
  final String name;
  final String? location;
  final String? description;
  final bool isActive;
  final DateTime createdAt;

  WarehouseEntity({
    this.id,
    required this.name,
    this.location,
    this.description,
    this.isActive = true,
    required this.createdAt,
  });

  WarehouseEntity copyWith({
    int? id,
    String? name,
    String? location,
    String? description,
    bool? isActive,
    DateTime? createdAt,
  }) {
    return WarehouseEntity(
      id: id ?? this.id,
      name: name ?? this.name,
      location: location ?? this.location,
      description: description ?? this.description,
      isActive: isActive ?? this.isActive,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is WarehouseEntity &&
          runtimeType == other.runtimeType &&
          id == other.id &&
          name == other.name;

  @override
  int get hashCode => id.hashCode ^ name.hashCode;
}
