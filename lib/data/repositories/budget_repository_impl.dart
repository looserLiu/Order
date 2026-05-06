import '../../domain/entities/budget.dart';
import '../../domain/repositories/i_budget_repository.dart';
import '../database/database_helper.dart';
import '../models/budget.dart';
import '../mappers/budget_mapper.dart';

/// Implementation of IBudgetRepository using SQLite.
class BudgetRepositoryImpl implements IBudgetRepository {
  final DatabaseHelper _dbHelper;

  BudgetRepositoryImpl({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  @override
  Future<List<BudgetEntity>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('budgets', orderBy: 'start_date DESC');
    final budgets = maps.map((map) => Budget.fromMap(map)).toList();
    return BudgetMapper.toEntityList(budgets);
  }

  @override
  Future<BudgetEntity?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'budgets',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return BudgetMapper.toEntity(Budget.fromMap(maps.first));
  }

  @override
  Future<List<BudgetEntity>> getByCategoryId(int categoryId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'budgets',
      where: 'category_id = ?',
      whereArgs: [categoryId],
      orderBy: 'start_date DESC',
    );
    final budgets = maps.map((map) => Budget.fromMap(map)).toList();
    return BudgetMapper.toEntityList(budgets);
  }

  @override
  Future<List<BudgetEntity>> getByPeriod(BudgetPeriod period) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'budgets',
      where: 'period = ?',
      whereArgs: [period.value],
      orderBy: 'start_date DESC',
    );
    final budgets = maps.map((map) => Budget.fromMap(map)).toList();
    return BudgetMapper.toEntityList(budgets);
  }

  @override
  Future<BudgetEntity> create(BudgetEntity budget) async {
    final db = await _dbHelper.database;
    final model = BudgetMapper.toModel(budget);
    final id = await db.insert('budgets', model.toMap());
    return budget.copyWith(id: id);
  }

  @override
  Future<BudgetEntity> update(BudgetEntity budget) async {
    final db = await _dbHelper.database;
    final model = BudgetMapper.toModel(budget);
    await db.update(
      'budgets',
      model.toMap(),
      where: 'id = ?',
      whereArgs: [budget.id],
    );
    return budget;
  }

  @override
  Future<void> delete(int id) async {
    final db = await _dbHelper.database;
    await db.delete(
      'budgets',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
