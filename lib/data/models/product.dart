/// Product model for inventory management.
class Product {
  final int? id;
  final String name;
  final String? sku;
  final String? category;
  final String? unit;
  final double costPrice;
  final double salePrice;
  final String? imageUrl;
  final int createdAt;
  final int updatedAt;

  Product({
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

  /// Create Product from database map.
  factory Product.fromMap(Map<String, dynamic> map) {
    return Product(
      id: map['id'] as int?,
      name: map['name'] as String,
      sku: map['sku'] as String?,
      category: map['category'] as String?,
      unit: map['unit'] as String?,
      costPrice: (map['cost_price'] as num?)?.toDouble() ?? 0.0,
      salePrice: (map['sale_price'] as num?)?.toDouble() ?? 0.0,
      imageUrl: map['image_url'] as String?,
      createdAt: map['created_at'] as int,
      updatedAt: map['updated_at'] as int,
    );
  }

  /// Convert Product to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'name': name,
      'sku': sku,
      'category': category,
      'unit': unit,
      'cost_price': costPrice,
      'sale_price': salePrice,
      'image_url': imageUrl,
      'created_at': createdAt,
      'updated_at': updatedAt,
    };
  }

  /// Create a copy of Product with optional field updates.
  Product copyWith({
    int? id,
    String? name,
    String? sku,
    String? category,
    String? unit,
    double? costPrice,
    double? salePrice,
    String? imageUrl,
    int? createdAt,
    int? updatedAt,
  }) {
    return Product(
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
}