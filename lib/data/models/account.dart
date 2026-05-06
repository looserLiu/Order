import 'package:sqflite/sqflite.dart';

/// Account model representing a financial account.
class Account {
  final int? id;
  final String name;
  final String type;
  final double balance;
  final String currency;
  final String? icon;
  final int? color;
  final int createdAt;
  final int updatedAt;

  Account({
    this.id,
    required this.name,
    required this.type,
    this.balance = 0.0,
    this.currency = 'CNY',
    this.icon,
    this.color,
    required this.createdAt,
    required this.updatedAt,
  });

  /// Create Account from database map.
  factory Account.fromMap(Map<String, dynamic> map) {
    return Account(
      id: map['id'] as int?,
      name: map['name'] as String,
      type: map['type'] as String,
      balance: (map['balance'] as num?)?.toDouble() ?? 0.0,
      currency: map['currency'] as String? ?? 'CNY',
      icon: map['icon'] as String?,
      color: map['color'] as int?,
      createdAt: map['created_at'] as int,
      updatedAt: map['updated_at'] as int,
    );
  }

  /// Convert Account to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'name': name,
      'type': type,
      'balance': balance,
      'currency': currency,
      'icon': icon,
      'color': color,
      'created_at': createdAt,
      'updated_at': updatedAt,
    };
  }

  /// Create a copy of Account with optional field updates.
  Account copyWith({
    int? id,
    String? name,
    String? type,
    double? balance,
    String? currency,
    String? icon,
    int? color,
    int? createdAt,
    int? updatedAt,
  }) {
    return Account(
      id: id ?? this.id,
      name: name ?? this.name,
      type: type ?? this.type,
      balance: balance ?? this.balance,
      currency: currency ?? this.currency,
      icon: icon ?? this.icon,
      color: color ?? this.color,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  /// Account types enum values.
  static const String typeCash = 'cash';
  static const String typeBank = 'bank';
  static const String typeCreditCard = 'credit_card';
  static const String typeDigital = 'digital';

  /// Get all valid account types.
  static List<String> get types => [
        typeCash,
        typeBank,
        typeCreditCard,
        typeDigital,
      ];
}