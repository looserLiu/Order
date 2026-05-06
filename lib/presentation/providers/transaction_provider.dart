import 'package:flutter/foundation.dart';
import '../../data/database/database_helper.dart';
import '../../data/models/transaction.dart';

class TransactionProvider extends ChangeNotifier {
  final DatabaseHelper _db = DatabaseHelper();
  List<Transaction> _transactions = [];
  bool _isLoading = false;

  List<Transaction> get transactions => _transactions;
  bool get isLoading => _isLoading;

  Future<void> loadTransactions({int? accountId, DateTime? startDate, DateTime? endDate}) async {
    _isLoading = true;
    notifyListeners();

    try {
      final db = await _db.database;
      String? where;
      List<dynamic> whereArgs = [];

      if (accountId != null) {
        where = 'account_id = ?';
        whereArgs.add(accountId);
      }

      if (startDate != null) {
        where = where != null ? '$where AND date >= ?' : 'date >= ?';
        whereArgs.add(startDate.millisecondsSinceEpoch);
      }

      if (endDate != null) {
        where = where != null ? '$where AND date <= ?' : 'date <= ?';
        whereArgs.add(endDate.millisecondsSinceEpoch);
      }

      final List<Map<String, dynamic>> maps = await db.query(
        'transactions',
        where: where,
        whereArgs: whereArgs.isEmpty ? null : whereArgs,
        orderBy: 'date DESC, created_at DESC',
      );
      _transactions = maps.map((map) => Transaction.fromMap(map)).toList();
    } catch (e) {
      debugPrint('Error loading transactions: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> addTransaction(Transaction transaction) async {
    try {
      final db = await _db.database;
      await db.insert('transactions', transaction.toMap());
      await loadTransactions();
    } catch (e) {
      debugPrint('Error adding transaction: $e');
    }
  }

  Future<void> updateTransaction(Transaction transaction) async {
    try {
      final db = await _db.database;
      await db.update(
        'transactions',
        transaction.toMap(),
        where: 'id = ?',
        whereArgs: [transaction.id],
      );
      await loadTransactions();
    } catch (e) {
      debugPrint('Error updating transaction: $e');
    }
  }

  Future<void> deleteTransaction(int id) async {
    try {
      final db = await _db.database;
      await db.delete(
        'transactions',
        where: 'id = ?',
        whereArgs: [id],
      );
      await loadTransactions();
    } catch (e) {
      debugPrint('Error deleting transaction: $e');
    }
  }
}