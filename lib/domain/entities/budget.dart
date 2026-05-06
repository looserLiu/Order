/// Domain entity for Budget.
class BudgetEntity {
  final int? id;
  final int categoryId;
  final double amount;
  final BudgetPeriod period;
  final DateTime startDate;
  final DateTime endDate;

  BudgetEntity({
    this.id,
    required this.categoryId,
    required this.amount,
    required this.period,
    required this.startDate,
    required this.endDate,
  });

  BudgetEntity copyWith({
    int? id,
    int? categoryId,
    double? amount,
    BudgetPeriod? period,
    DateTime? startDate,
    DateTime? endDate,
  }) {
    return BudgetEntity(
      id: id ?? this.id,
      categoryId: categoryId ?? this.categoryId,
      amount: amount ?? this.amount,
      period: period ?? this.period,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BudgetEntity &&
          runtimeType == other.runtimeType &&
          id == other.id;

  @override
  int get hashCode => id.hashCode;
}

/// Budget period enumeration.
enum BudgetPeriod {
  weekly('weekly', '每周'),
  monthly('monthly', '每月'),
  yearly('yearly', '每年');

  final String value;
  final String displayName;

  const BudgetPeriod(this.value, this.displayName);

  static BudgetPeriod fromValue(String value) {
    return BudgetPeriod.values.firstWhere(
      (e) => e.value == value,
      orElse: () => BudgetPeriod.monthly,
    );
  }
}
