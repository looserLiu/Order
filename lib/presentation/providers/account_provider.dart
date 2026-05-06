import 'package:flutter/foundation.dart';
import '../../data/database/database_helper.dart';
import '../../data/models/account.dart';

class AccountProvider extends ChangeNotifier {
  final DatabaseHelper _db = DatabaseHelper();
  List<Account> _accounts = [];
  bool _isLoading = false;

  List<Account> get accounts => _accounts;
  bool get isLoading => _isLoading;

  Future<void> loadAccounts() async {
    _isLoading = true;
    notifyListeners();

    try {
      final db = await _db.database;
      final List<Map<String, dynamic>> maps = await db.query(
        'accounts',
        orderBy: 'created_at DESC',
      );
      _accounts = maps.map((map) => Account.fromMap(map)).toList();
    } catch (e) {
      debugPrint('Error loading accounts: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> addAccount(Account account) async {
    try {
      final db = await _db.database;
      await db.insert('accounts', account.toMap());
      await loadAccounts();
    } catch (e) {
      debugPrint('Error adding account: $e');
    }
  }

  Future<void> updateAccount(Account account) async {
    try {
      final db = await _db.database;
      await db.update(
        'accounts',
        account.toMap(),
        where: 'id = ?',
        whereArgs: [account.id],
      );
      await loadAccounts();
    } catch (e) {
      debugPrint('Error updating account: $e');
    }
  }

  Future<void> deleteAccount(int id) async {
    try {
      final db = await _db.database;
      await db.delete(
        'accounts',
        where: 'id = ?',
        whereArgs: [id],
      );
      await loadAccounts();
    } catch (e) {
      debugPrint('Error deleting account: $e');
    }
  }

  Future<void> updateBalance(int id, double newBalance) async {
    try {
      final db = await _db.database;
      await db.update(
        'accounts',
        {
          'balance': newBalance,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        },
        where: 'id = ?',
        whereArgs: [id],
      );
      await loadAccounts();
    } catch (e) {
      debugPrint('Error updating balance: $e');
    }
  }
}