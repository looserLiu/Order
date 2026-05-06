import 'package:sqflite/sqflite.dart';
import '../database/database_helper.dart';
import '../models/inventory_flow.dart';

/// Repository for InventoryFlow data operations.
class InventoryFlowRepository {
  final DatabaseHelper _dbHelper;

  InventoryFlowRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new inventory flow.
  Future<int> insert(InventoryFlow flow) async {
    final db = await _dbHelper.database;
    return await db.insert('inventory_flows', flow.toMap());
  }

  /// Update an existing inventory flow.
  Future<int> update(InventoryFlow flow) async {
    final db = await _dbHelper.database;
    return await db.update(
      'inventory_flows',
      flow.toMap(),
      where: 'id = ?',
      whereArgs: [flow.id],
    );
  }

  /// Delete an inventory flow by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'inventory_flows',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all inventory flows.
  Future<List<InventoryFlow>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('inventory_flows', orderBy: 'date DESC');
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }

  /// Get inventory flow by ID.
  Future<InventoryFlow?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'inventory_flows',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return InventoryFlow.fromMap(maps.first);
  }

  /// Get flows by product ID.
  Future<List<InventoryFlow>> getByProductId(int productId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'inventory_flows',
      where: 'product_id = ?',
      whereArgs: [productId],
      orderBy: 'date DESC',
    );
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }

  /// Get flows by warehouse ID.
  Future<List<InventoryFlow>> getByWarehouseId(int warehouseId) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'inventory_flows',
      where: 'warehouse_id = ?',
      whereArgs: [warehouseId],
      orderBy: 'date DESC',
    );
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }

  /// Get flows by flow type.
  Future<List<InventoryFlow>> getByType(String flowType) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'inventory_flows',
      where: 'flow_type = ?',
      whereArgs: [flowType],
      orderBy: 'date DESC',
    );
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }

  /// Get flows within a date range.
  Future<List<InventoryFlow>> getByDateRange(int startDate, int endDate) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'inventory_flows',
      where: 'date >= ? AND date <= ?',
      whereArgs: [startDate, endDate],
      orderBy: 'date DESC',
    );
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }

  /// Get flows expiring soon (within days).
  Future<List<InventoryFlow>> getExpiringSoon(int days) async {
    final db = await _dbHelper.database;
    final now = DateTime.now().millisecondsSinceEpoch;
    final futureDate = DateTime.now()
        .add(Duration(days: days))
        .millisecondsSinceEpoch;

    final maps = await db.query(
      'inventory_flows',
      where: 'expiration_date IS NOT NULL AND expiration_date > ? AND expiration_date <= ?',
      whereArgs: [now, futureDate],
      orderBy: 'expiration_date ASC',
    );
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }

  /// Get flows by product and warehouse.
  Future<List<InventoryFlow>> getByProductAndWarehouse(
    int productId,
    int warehouseId,
  ) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'inventory_flows',
      where: 'product_id = ? AND warehouse_id = ?',
      whereArgs: [productId, warehouseId],
      orderBy: 'date DESC',
    );
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }

  /// Calculate current stock for a product in a warehouse.
  Future<double> getCurrentStock(int productId, int warehouseId) async {
    final db = await _dbHelper.database;

    // Sum of 'in' flows - Sum of 'out' flows
    final result = await db.rawQuery('''
      SELECT
        COALESCE(SUM(CASE WHEN flow_type = 'in' THEN quantity ELSE 0 END), 0) -
        COALESCE(SUM(CASE WHEN flow_type = 'out' THEN quantity ELSE 0 END), 0) as stock
      FROM inventory_flows
      WHERE product_id = ? AND warehouse_id = ?
    ''', [productId, warehouseId]);

    final stock = result.first['stock'];
    return stock != null ? (stock as num).toDouble() : 0.0;
  }

  /// Get total stock for a product across all warehouses.
  Future<double> getTotalStock(int productId) async {
    final db = await _dbHelper.database;

    final result = await db.rawQuery('''
      SELECT
        COALESCE(SUM(CASE WHEN flow_type = 'in' THEN quantity ELSE 0 END), 0) -
        COALESCE(SUM(CASE WHEN flow_type = 'out' THEN quantity ELSE 0 END), 0) as stock
      FROM inventory_flows
      WHERE product_id = ?
    ''', [productId]);

    final stock = result.first['stock'];
    return stock != null ? (stock as num).toDouble() : 0.0;
  }

  /// Get flows with filters.
  Future<List<InventoryFlow>> getWithFilters({
    int? productId,
    int? warehouseId,
    String? flowType,
    int? startDate,
    int? endDate,
    int? limit,
    int? offset,
  }) async {
    final db = await _dbHelper.database;
    final whereConditions = <String>[];
    final whereArgs = <dynamic>[];

    if (productId != null) {
      whereConditions.add('product_id = ?');
      whereArgs.add(productId);
    }
    if (warehouseId != null) {
      whereConditions.add('warehouse_id = ?');
      whereArgs.add(warehouseId);
    }
    if (flowType != null) {
      whereConditions.add('flow_type = ?');
      whereArgs.add(flowType);
    }
    if (startDate != null) {
      whereConditions.add('date >= ?');
      whereArgs.add(startDate);
    }
    if (endDate != null) {
      whereConditions.add('date <= ?');
      whereArgs.add(endDate);
    }

    final maps = await db.query(
      'inventory_flows',
      where: whereConditions.isEmpty ? null : whereConditions.join(' AND '),
      whereArgs: whereArgs.isEmpty ? null : whereArgs,
      orderBy: 'date DESC',
      limit: limit,
      offset: offset,
    );
    return maps.map((map) => InventoryFlow.fromMap(map)).toList();
  }
}