/// Transaction model representing a financial transaction.
class Transaction {
  final int? id;
  final int accountId;
  final int? categoryId;
  final double amount;
  final String type;
  final String? description;
  final int date;
  final int createdAt;

  Transaction({
    this.id,
    required this.accountId,
    this.categoryId,
    required this.amount,
    required this.type,
    this.description,
    required this.date,
    required this.createdAt,
  });

  /// Create Transaction from database map.
  factory Transaction.fromMap(Map<String, dynamic> map) {
    return Transaction(
      id: map['id'] as int?,
      accountId: map['account_id'] as int,
      categoryId: map['category_id'] as int?,
      amount: (map['amount'] as num).toDouble(),
      type: map['type'] as String,
      description: map['description'] as String?,
      date: map['date'] as int,
      createdAt: map['created_at'] as int,
    );
  }

  /// Convert Transaction to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'account_id': accountId,
      'category_id': categoryId,
      'amount': amount,
      'type': type,
      'description': description,
      'date': date,
      'created_at': createdAt,
    };
  }

  /// Create a copy of Transaction with optional field updates.
  Transaction copyWith({
    int? id,
    int? accountId,
    int? categoryId,
    double? amount,
    String? type,
    String? description,
    int? date,
    int? createdAt,
  }) {
    return Transaction(
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

  /// Transaction type values.
  static const String typeIncome = 'income';
  static const String typeExpense = 'expense';

  /// Get all valid transaction types.
  static List<String> get types => [typeIncome, typeExpense];
}