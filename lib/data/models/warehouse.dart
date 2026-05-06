/// Warehouse model for inventory storage locations.
class Warehouse {
  final int? id;
  final String name;
  final String? location;
  final String? description;
  final int isActive; // 1 = active, 0 = inactive
  final int createdAt;

  Warehouse({
    this.id,
    required this.name,
    this.location,
    this.description,
    this.isActive = 1,
    required this.createdAt,
  });

  /// Check if warehouse is active.
  bool get isActiveStatus => isActive == 1;

  /// Create Warehouse from database map.
  factory Warehouse.fromMap(Map<String, dynamic> map) {
    return Warehouse(
      id: map['id'] as int?,
      name: map['name'] as String,
      location: map['location'] as String?,
      description: map['description'] as String?,
      isActive: (map['is_active'] as int?) ?? 1,
      createdAt: map['created_at'] as int,
    );
  }

  /// Convert Warehouse to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'name': name,
      'location': location,
      'description': description,
      'is_active': isActive,
      'created_at': createdAt,
    };
  }

  /// Create a copy of Warehouse with optional field updates.
  Warehouse copyWith({
    int? id,
    String? name,
    String? location,
    String? description,
    int? isActive,
    int? createdAt,
  }) {
    return Warehouse(
      id: id ?? this.id,
      name: name ?? this.name,
      location: location ?? this.location,
      description: description ?? this.description,
      isActive: isActive ?? this.isActive,
      createdAt: createdAt ?? this.createdAt,
    );
  }
}