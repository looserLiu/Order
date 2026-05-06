/// Domain entity for Account.
/// Pure business entity, no framework dependencies.
class AccountEntity {
  final int? id;
  final String name;
  final AccountType type;
  final double balance;
  final String currency;
  final String? icon;
  final int? color;
  final DateTime createdAt;
  final DateTime updatedAt;

  AccountEntity({
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

  AccountEntity copyWith({
    int? id,
    String? name,
    AccountType? type,
    double? balance,
    String? currency,
    String? icon,
    int? color,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return AccountEntity(
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

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is AccountEntity &&
          runtimeType == other.runtimeType &&
          id == other.id &&
          name == other.name;

  @override
  int get hashCode => id.hashCode ^ name.hashCode;
}

/// Account type enumeration.
enum AccountType {
  cash('cash', '现金'),
  bank('bank', '银行卡'),
  creditCard('credit_card', '信用卡'),
  digital('digital', '数字账户');

  final String value;
  final String displayName;

  const AccountType(this.value, this.displayName);

  static AccountType fromValue(String value) {
    return AccountType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => AccountType.cash,
    );
  }
}
