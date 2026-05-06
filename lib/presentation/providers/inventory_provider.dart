import 'package:flutter/foundation.dart';
import '../../data/database/database_helper.dart';
import '../../data/models/inventory_flow.dart';
import '../../data/repositories/inventory_flow_repository.dart';

/// Provider for inventory flow operations and stock calculations.
class InventoryProvider extends ChangeNotifier {
  final DatabaseHelper _db = DatabaseHelper();
  final InventoryFlowRepository _repository;

  List<InventoryFlow> _inventoryFlows = [];
  bool _isLoading = false;

  InventoryProvider({InventoryFlowRepository? repository})
      : _repository = repository ?? InventoryFlowRepository();

  List<InventoryFlow> get inventoryFlows => _inventoryFlows;
  bool get isLoading => _isLoading;

  /// Load all inventory flows from database.
  Future<void> loadInventoryFlows() async {
    _isLoading = true;
    notifyListeners();

    try {
      _inventoryFlows = await _repository.getAll();
    } catch (e) {
      debugPrint('Error loading inventory flows: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  /// Add a new inventory flow (stock in/out).
  Future<void> addInventoryFlow(InventoryFlow flow) async {
    try {
      await _repository.insert(flow);
      await loadInventoryFlows();
    } catch (e) {
      debugPrint('Error adding inventory flow: $e');
      rethrow;
    }
  }

  /// Update an existing inventory flow.
  Future<void> updateInventoryFlow(InventoryFlow flow) async {
    try {
      await _repository.update(flow);
      await loadInventoryFlows();
    } catch (e) {
      debugPrint('Error updating inventory flow: $e');
      rethrow;
    }
  }

  /// Delete an inventory flow.
  Future<void> deleteInventoryFlow(int id) async {
    try {
      await _repository.delete(id);
      await loadInventoryFlows();
    } catch (e) {
      debugPrint('Error deleting inventory flow: $e');
      rethrow;
    }
  }

  /// Get current stock for a product in a specific warehouse.
  double getProductStockInWarehouse(int productId, int warehouseId) {
    double stock = 0;
    for (final flow in _inventoryFlows) {
      if (flow.productId == productId && flow.warehouseId == warehouseId) {
        if (flow.flowType == 'in') {
          stock += flow.quantity;
        } else if (flow.flowType == 'out') {
          stock -= flow.quantity;
        }
      }
    }
    return stock;
  }

  /// Get total stock for a product across all warehouses.
  double getProductStock(int productId) {
    double stock = 0;
    for (final flow in _inventoryFlows) {
      if (flow.productId == productId) {
        if (flow.flowType == 'in') {
          stock += flow.quantity;
        } else if (flow.flowType == 'out') {
          stock -= flow.quantity;
        }
      }
    }
    return stock;
  }

  /// Get total number of products in a warehouse.
  int getProductsInWarehouse(int warehouseId) {
    final productIds = <int>{};
    for (final flow in _inventoryFlows) {
      if (flow.warehouseId == warehouseId) {
        productIds.add(flow.productId);
      }
    }
    return productIds.length;
  }

  /// Get total stock quantity in a warehouse.
  double getWarehouseStock(int warehouseId) {
    double stock = 0;
    for (final flow in _inventoryFlows) {
      if (flow.warehouseId == warehouseId) {
        if (flow.flowType == 'in') {
          stock += flow.quantity;
        } else if (flow.flowType == 'out') {
          stock -= flow.quantity;
        }
      }
    }
    return stock;
  }

  /// Get flows for a specific product.
  List<InventoryFlow> getProductFlows(int productId) {
    return _inventoryFlows
        .where((f) => f.productId == productId)
        .toList();
  }

  /// Get flows for a specific warehouse.
  List<InventoryFlow> getWarehouseFlows(int warehouseId) {
    return _inventoryFlows
        .where((f) => f.warehouseId == warehouseId)
        .toList();
  }

  /// Get recent inventory flows with limit.
  List<InventoryFlow> getRecentFlows({int limit = 10}) {
    final sorted = List<InventoryFlow>.from(_inventoryFlows)
      ..sort((a, b) => b.date.compareTo(a.date));
    return sorted.take(limit).toList();
  }

  /// Get low stock products (stock < threshold).
  List<int> getLowStockProducts({double threshold = 10}) {
    final productStocks = <int, double>{};
    for (final flow in _inventoryFlows) {
      final current = productStocks[flow.productId] ?? 0;
      if (flow.flowType == 'in') {
        productStocks[flow.productId] = current + flow.quantity;
      } else if (flow.flowType == 'out') {
        productStocks[flow.productId] = current - flow.quantity;
      }
    }
    return productStocks.entries
        .where((e) => e.value < threshold)
        .map((e) => e.key)
        .toList();
  }

  /// Get inventory value for a product.
  double getProductInventoryValue(int productId, double costPrice) {
    return getProductStock(productId) * costPrice;
  }
}