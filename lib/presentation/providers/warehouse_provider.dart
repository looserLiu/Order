import 'package:flutter/foundation.dart';
import '../../data/database/database_helper.dart';
import '../../data/models/warehouse.dart';
import '../../data/repositories/warehouse_repository.dart';

/// Provider for warehouse operations.
class WarehouseProvider extends ChangeNotifier {
  final DatabaseHelper _db = DatabaseHelper();
  final WarehouseRepository _repository;

  List<Warehouse> _warehouses = [];
  bool _isLoading = false;

  WarehouseProvider({WarehouseRepository? repository})
      : _repository = repository ?? WarehouseRepository();

  List<Warehouse> get warehouses => _warehouses;
  bool get isLoading => _isLoading;

  /// Load all warehouses from database.
  Future<void> loadWarehouses() async {
    _isLoading = true;
    notifyListeners();

    try {
      _warehouses = await _repository.getAll();
    } catch (e) {
      debugPrint('Error loading warehouses: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  /// Add a new warehouse.
  Future<void> addWarehouse(Warehouse warehouse) async {
    try {
      await _repository.insert(warehouse);
      await loadWarehouses();
    } catch (e) {
      debugPrint('Error adding warehouse: $e');
      rethrow;
    }
  }

  /// Update an existing warehouse.
  Future<void> updateWarehouse(Warehouse warehouse) async {
    try {
      await _repository.update(warehouse);
      await loadWarehouses();
    } catch (e) {
      debugPrint('Error updating warehouse: $e');
      rethrow;
    }
  }

  /// Delete a warehouse.
  Future<void> deleteWarehouse(int id) async {
    try {
      await _repository.delete(id);
      await loadWarehouses();
    } catch (e) {
      debugPrint('Error deleting warehouse: $e');
      rethrow;
    }
  }

  /// Get warehouse by ID.
  Warehouse? getWarehouseById(int id) {
    return _warehouses.where((w) => w.id == id).firstOrNull;
  }

  /// Get active warehouses only.
  List<Warehouse> get activeWarehouses {
    return _warehouses.where((w) => w.isActive == 1).toList();
  }

  /// Toggle warehouse active status.
  Future<void> toggleWarehouseActive(int id) async {
    final warehouse = getWarehouseById(id);
    if (warehouse == null) return;

    final newStatus = warehouse.isActive == 1 ? 0 : 1;
    await updateWarehouse(Warehouse(
      id: warehouse.id,
      name: warehouse.name,
      location: warehouse.location,
      description: warehouse.description,
      isActive: newStatus,
      createdAt: warehouse.createdAt,
    ));
  }

  /// Search warehouses by name.
  List<Warehouse> searchWarehouses(String query) {
    if (query.isEmpty) return _warehouses;
    return _warehouses
        .where((w) => w.name.toLowerCase().contains(query.toLowerCase()))
        .toList();
  }

  /// Get total warehouse count.
  int get totalCount => _warehouses.length;

  /// Get active warehouse count.
  int get activeCount => _warehouses.where((w) => w.isActive == 1).length;
}