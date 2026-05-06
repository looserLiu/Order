import 'package:flutter/foundation.dart';
import '../../data/database/database_helper.dart';
import '../../data/models/budget.dart';

class BudgetProvider extends ChangeNotifier {
  final DatabaseHelper _db = DatabaseHelper();
  List<Budget> _budgets = [];
  bool _isLoading = false;

  List<Budget> get budgets => _budgets;
  bool get isLoading => _isLoading;

  Future<void> loadBudgets() async {
    _isLoading = true;
    notifyListeners();

    try {
      final db = await _db.database;
      final List<Map<String, dynamic>> maps = await db.query(
        'budgets',
        orderBy: 'start_date DESC',
      );
      _budgets = maps.map((map) => Budget.fromMap(map)).toList();
    } catch (e) {
      debugPrint('Error loading budgets: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> addBudget(Budget budget) async {
    try {
      final db = await _db.database;
      await db.insert('budgets', budget.toMap());
      await loadBudgets();
    } catch (e) {
      debugPrint('Error adding budget: $e');
    }
  }

  Future<void> updateBudget(Budget budget) async {
    try {
      final db = await _db.database;
      await db.update(
        'budgets',
        budget.toMap(),
        where: 'id = ?',
        whereArgs: [budget.id],
      );
      await loadBudgets();
    } catch (e) {
      debugPrint('Error updating budget: $e');
    }
  }

  Future<void> deleteBudget(int id) async {
    try {
      final db = await _db.database;
      await db.delete(
        'budgets',
        where: 'id = ?',
        whereArgs: [id],
      );
      await loadBudgets();
    } catch (e) {
      debugPrint('Error deleting budget: $e');
    }
  }
}