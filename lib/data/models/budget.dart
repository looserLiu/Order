/// Budget model for category spending limits.
class Budget {
  final int? id;
  final int categoryId;
  final double amount;
  final String period;
  final int startDate;
  final int endDate;

  Budget({
    this.id,
    required this.categoryId,
    required this.amount,
    required this.period,
    required this.startDate,
    required this.endDate,
  });

  /// Create Budget from database map.
  factory Budget.fromMap(Map<String, dynamic> map) {
    return Budget(
      id: map['id'] as int?,
      categoryId: map['category_id'] as int,
      amount: (map['amount'] as num).toDouble(),
      period: map['period'] as String,
      startDate: map['start_date'] as int,
      endDate: map['end_date'] as int,
    );
  }

  /// Convert Budget to database map.
  Map<String, dynamic> toMap() {
    return {
      if (id != null) 'id': id,
      'category_id': categoryId,
      'amount': amount,
      'period': period,
      'start_date': startDate,
      'end_date': endDate,
    };
  }

  /// Create a copy of Budget with optional field updates.
  Budget copyWith({
    int? id,
    int? categoryId,
    double? amount,
    String? period,
    int? startDate,
    int? endDate,
  }) {
    return Budget(
      id: id ?? this.id,
      categoryId: categoryId ?? this.categoryId,
      amount: amount ?? this.amount,
      period: period ?? this.period,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
    );
  }

  /// Budget period values.
  static const String periodMonthly = 'monthly';
  static const String periodWeekly = 'weekly';
  static const String periodYearly = 'yearly';

  /// Get all valid periods.
  static List<String> get periods => [periodMonthly, periodWeekly, periodYearly];
}