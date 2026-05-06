/// Category model for transaction categorization.
class Category {
  final int? id;
  final String name;
  final String type;
  final String? icon;
  final int? color;
  final int? parentId;
  final bool isSystem;
  final int usageCount;
  final int? createdAt;

  Category({
    this.id,
    required this.name,
    required this.type,
    this.icon,
    this.color,
    this.parentId,
    this.isSystem = false,
    this.usageCount = 0,
    this.createdAt,
  });

  /// Create Category from database map.
  factory Category.fromMap(Map<String, dynamic> map) {
    return Category(
      id: map['id'] as int?,
      name: map['name'] as String,
      type: map['type'] as String,
      icon: map['icon'] as String?,
      color: map['color'] as int?,
      parentId: map['parent_id'] as int?,
      isSystem: (map['is_system'] as int?) == 1,
      usageCount: map['usage_count'] as int? ?? 0,
      createdAt: map['created_at'] as int?,
    );
  }

  /// Convert Category to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'name': name,
      'type': type,
      'icon': icon,
      'color': color,
      'parent_id': parentId,
      'is_system': isSystem ? 1 : 0,
      'usage_count': usageCount,
      'created_at': createdAt,
    };
  }

  /// Create a copy of Category with optional field updates.
  Category copyWith({
    int? id,
    String? name,
    String? type,
    String? icon,
    int? color,
    int? parentId,
    bool? isSystem,
    int? usageCount,
    int? createdAt,
  }) {
    return Category(
      id: id ?? this.id,
      name: name ?? this.name,
      type: type ?? this.type,
      icon: icon ?? this.icon,
      color: color ?? this.color,
      parentId: parentId ?? this.parentId,
      isSystem: isSystem ?? this.isSystem,
      usageCount: usageCount ?? this.usageCount,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  /// Category type values.
  static const String typeIncome = 'income';
  static const String typeExpense = 'expense';

  /// Get all valid category types.
  static List<String> get types => [typeIncome, typeExpense];
}