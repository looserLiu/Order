import '../../domain/entities/transaction.dart';
import '../../domain/repositories/i_transaction_repository.dart';
import '../database/database_helper.dart';
import '../models/transaction.dart';
import '../mappers/transaction_mapper.dart';

/// Implementation of ITransactionRepository using SQLite.
class TransactionRepositoryImpl implements ITransactionRepository {
  final DatabaseHelper _dbHelper;

  TransactionRepositoryImpl({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  @override
  Future<List<TransactionEntity>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('transactions', orderBy: 'date DESC');
    final transactions = maps.map((map) => Transaction.fromMap(map)).toList();
    return TransactionMapper.toEntityList(transactions);
  }

  @override
  Future<TransactionEntity?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return TransactionMapper.toEntity(Transaction.fromMap(maps.first));
  }

  @override
  Future<List<TransactionEntity>> getByAccountId(int accountId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'account_id = ?',
      whereArgs: [accountId],
      orderBy: 'date DESC',
    );
    final transactions = maps.map((map) => Transaction.fromMap(map)).toList();
    return TransactionMapper.toEntityList(transactions);
  }

  @override
  Future<List<TransactionEntity>> getByCategoryId(int categoryId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'category_id = ?',
      whereArgs: [categoryId],
      orderBy: 'date DESC',
    );
    final transactions = maps.map((map) => Transaction.fromMap(map)).toList();
    return TransactionMapper.toEntityList(transactions);
  }

  @override
  Future<List<TransactionEntity>> getByDateRange(DateTime start, DateTime end) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'transactions',
      where: 'date >= ? AND date <= ?',
      whereArgs: [start.millisecondsSinceEpoch, end.millisecondsSinceEpoch],
      orderBy: 'date DESC',
    );
    final transactions = maps.map((map) => Transaction.fromMap(map)).toList();
    return TransactionMapper.toEntityList(transactions);
  }

  @override
  Future<TransactionEntity> create(TransactionEntity transaction) async {
    final db = await _dbHelper.database;
    final model = TransactionMapper.toModel(transaction);
    final id = await db.insert('transactions', model.toMap());
    return transaction.copyWith(id: id);
  }

  @override
  Future<TransactionEntity> update(TransactionEntity transaction) async {
    final db = await _dbHelper.database;
    final model = TransactionMapper.toModel(transaction);
    await db.update(
      'transactions',
      model.toMap(),
      where: 'id = ?',
      whereArgs: [transaction.id],
    );
    return transaction;
  }

  @override
  Future<void> delete(int id) async {
    final db = await _dbHelper.database;
    await db.delete(
      'transactions',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
