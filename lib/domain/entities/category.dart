/// Domain entity for Category.
class CategoryEntity {
  final int? id;
  final String name;
  final CategoryType type;
  final String? icon;
  final int? color;
  final int? parentId;
  final bool isSystem;
  final int usageCount;
  final DateTime createdAt;

  CategoryEntity({
    this.id,
    required this.name,
    required this.type,
    this.icon,
    this.color,
    this.parentId,
    this.isSystem = false,
    this.usageCount = 0,
    required this.createdAt,
  });

  CategoryEntity copyWith({
    int? id,
    String? name,
    CategoryType? type,
    String? icon,
    int? color,
    int? parentId,
    bool? isSystem,
    int? usageCount,
    DateTime? createdAt,
  }) {
    return CategoryEntity(
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

  bool get isExpense => type == CategoryType.expense;
  bool get isIncome => type == CategoryType.income;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is CategoryEntity &&
          runtimeType == other.runtimeType &&
          id == other.id &&
          name == other.name;

  @override
  int get hashCode => id.hashCode ^ name.hashCode;
}

/// Category type enumeration.
enum CategoryType {
  income('income', '收入'),
  expense('expense', '支出');

  final String value;
  final String displayName;

  const CategoryType(this.value, this.displayName);

  static CategoryType fromValue(String value) {
    return CategoryType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => CategoryType.expense,
    );
  }
}
