/// Domain entity for Transaction.
class TransactionEntity {
  final int? id;
  final int accountId;
  final int? categoryId;
  final double amount;
  final TransactionType type;
  final String? description;
  final DateTime date;
  final DateTime createdAt;

  TransactionEntity({
    this.id,
    required this.accountId,
    this.categoryId,
    required this.amount,
    required this.type,
    this.description,
    required this.date,
    required this.createdAt,
  });

  TransactionEntity copyWith({
    int? id,
    int? accountId,
    int? categoryId,
    double? amount,
    TransactionType? type,
    String? description,
    DateTime? date,
    DateTime? createdAt,
  }) {
    return TransactionEntity(
      id: id ?? this.id,
      accountId: accountId ?? this.accountId,
      categoryId: categoryId ?? this.categoryId,
      amount: amount ?? this.amount,
      type: type ?? this.type,
      description: description ?? this.description,
      date: date ?? this.date,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  bool get isExpense => type == TransactionType.expense;
  bool get isIncome => type == TransactionType.income;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is TransactionEntity &&
          runtimeType == other.runtimeType &&
          id == other.id;

  @override
  int get hashCode => id.hashCode;
}

/// Transaction type enumeration.
enum TransactionType {
  income('income', '收入'),
  expense('expense', '支出');

  final String value;
  final String displayName;

  const TransactionType(this.value, this.displayName);

  static TransactionType fromValue(String value) {
    return TransactionType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => TransactionType.expense,
    );
  }
}
