import '../../domain/entities/transaction.dart';
import '../models/transaction.dart';

/// Mapper for converting between TransactionEntity (domain) and Transaction (data model).
class TransactionMapper {
  /// Convert data model to domain entity.
  static TransactionEntity toEntity(Transaction model) {
    return TransactionEntity(
      id: model.id,
      accountId: model.accountId,
      categoryId: model.categoryId,
      amount: model.amount,
      type: TransactionType.fromValue(model.type),
      description: model.description,
      date: DateTime.fromMillisecondsSinceEpoch(model.date),
      createdAt: DateTime.fromMillisecondsSinceEpoch(model.createdAt),
    );
  }

  /// Convert domain entity to data model.
  static Transaction toModel(TransactionEntity entity) {
    return Transaction(
      id: entity.id,
      accountId: entity.accountId,
      categoryId: entity.categoryId,
      amount: entity.amount,
      type: entity.type.value,
      description: entity.description,
      date: entity.date.millisecondsSinceEpoch,
      createdAt: entity.createdAt.millisecondsSinceEpoch,
    );
  }

  /// Convert list of data models to domain entities.
  static List<TransactionEntity> toEntityList(List<Transaction> models) {
    return models.map(toEntity).toList();
  }
}
