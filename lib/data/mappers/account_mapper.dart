import '../../domain/entities/account.dart';
import '../models/account.dart';

/// Mapper for converting between AccountEntity (domain) and Account (data model).
class AccountMapper {
  /// Convert data model to domain entity.
  static AccountEntity toEntity(Account model) {
    return AccountEntity(
      id: model.id,
      name: model.name,
      type: AccountType.fromValue(model.type),
      balance: model.balance,
      currency: model.currency,
      icon: model.icon,
      color: model.color,
      createdAt: DateTime.fromMillisecondsSinceEpoch(model.createdAt),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(model.updatedAt),
    );
  }

  /// Convert domain entity to data model.
  static Account toModel(AccountEntity entity) {
    return Account(
      id: entity.id,
      name: entity.name,
      type: entity.type.value,
      balance: entity.balance,
      currency: entity.currency,
      icon: entity.icon,
      color: entity.color,
      createdAt: entity.createdAt.millisecondsSinceEpoch,
      updatedAt: entity.updatedAt.millisecondsSinceEpoch,
    );
  }

  /// Convert list of data models to domain entities.
  static List<AccountEntity> toEntityList(List<Account> models) {
    return models.map(toEntity).toList();
  }
}
