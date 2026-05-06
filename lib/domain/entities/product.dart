/// Domain entity for Product.
class ProductEntity {
  final int? id;
  final String name;
  final String? sku;
  final String? category;
  final String? unit;
  final double costPrice;
  final double salePrice;
  final String? imageUrl;
  final DateTime createdAt;
  final DateTime updatedAt;

  ProductEntity({
    this.id,
    required this.name,
    this.sku,
    this.category,
    this.unit,
    this.costPrice = 0.0,
    this.salePrice = 0.0,
    this.imageUrl,
    required this.createdAt,
    required this.updatedAt,
  });

  ProductEntity copyWith({
    int? id,
    String? name,
    String? sku,
    String? category,
    String? unit,
    double? costPrice,
    double? salePrice,
    String? imageUrl,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return ProductEntity(
      id: id ?? this.id,
      name: name ?? this.name,
      sku: sku ?? this.sku,
      category: category ?? this.category,
      unit: unit ?? this.unit,
      costPrice: costPrice ?? this.costPrice,
      salePrice: salePrice ?? this.salePrice,
      imageUrl: imageUrl ?? this.imageUrl,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  /// Calculate profit margin.
  double get profitMargin {
    if (salePrice == 0) return 0;
    return (salePrice - costPrice) / salePrice;
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ProductEntity &&
          runtimeType == other.runtimeType &&
          id == other.id &&
          name == other.name;

  @override
  int get hashCode => id.hashCode ^ name.hashCode;
}
