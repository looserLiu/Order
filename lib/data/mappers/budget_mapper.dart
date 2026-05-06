import '../../domain/entities/budget.dart';
import '../models/budget.dart';

/// Mapper for converting between BudgetEntity (domain) and Budget (data model).
class BudgetMapper {
  /// Convert data model to domain entity.
  static BudgetEntity toEntity(Budget model) {
    return BudgetEntity(
      id: model.id,
      categoryId: model.categoryId,
      amount: model.amount,
      period: BudgetPeriod.fromValue(model.period),
      startDate: DateTime.fromMillisecondsSinceEpoch(model.startDate),
      endDate: DateTime.fromMillisecondsSinceEpoch(model.endDate),
    );
  }

  /// Convert domain entity to data model.
  static Budget toModel(BudgetEntity entity) {
    return Budget(
      id: entity.id,
      categoryId: entity.categoryId,
      amount: entity.amount,
      period: entity.period.value,
      startDate: entity.startDate.millisecondsSinceEpoch,
      endDate: entity.endDate.millisecondsSinceEpoch,
    );
  }

  /// Convert list of data models to domain entities.
  static List<BudgetEntity> toEntityList(List<Budget> models) {
    return models.map(toEntity).toList();
  }
}
