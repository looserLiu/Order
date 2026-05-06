import 'package:sqflite/sqflite.dart';
import '../database/database_helper.dart';
import '../models/warehouse.dart';

/// Repository for Warehouse data operations.
class WarehouseRepository {
  final DatabaseHelper _dbHelper;

  WarehouseRepository({DatabaseHelper? dbHelper})
      : _dbHelper = dbHelper ?? DatabaseHelper();

  /// Insert a new warehouse.
  Future<int> insert(Warehouse warehouse) async {
    final db = await _dbHelper.database;
    return await db.insert('warehouses', warehouse.toMap());
  }

  /// Update an existing warehouse.
  Future<int> update(Warehouse warehouse) async {
    final db = await _dbHelper.database;
    return await db.update(
      'warehouses',
      warehouse.toMap(),
      where: 'id = ?',
      whereArgs: [warehouse.id],
    );
  }

  /// Delete a warehouse by ID.
  Future<int> delete(int id) async {
    final db = await _dbHelper.database;
    return await db.delete(
      'warehouses',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get all warehouses.
  Future<List<Warehouse>> getAll() async {
    final db = await _dbHelper.database;
    final maps = await db.query('warehouses', orderBy: 'created_at DESC');
    return maps.map((map) => Warehouse.fromMap(map)).toList();
  }

  /// Get active warehouses only.
  Future<List<Warehouse>> getActive() async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'warehouses',
      where: 'is_active = ?',
      whereArgs: [1],
      orderBy: 'name ASC',
    );
    return maps.map((map) => Warehouse.fromMap(map)).toList();
  }

  /// Get warehouse by ID.
  Future<Warehouse?> getById(int id) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'warehouses',
      where: 'id = ?',
      whereArgs: [id],
    );
    if (maps.isEmpty) return null;
    return Warehouse.fromMap(maps.first);
  }

  /// Search warehouses by name.
  Future<List<Warehouse>> searchByName(String query) async {
    final db = await _dbHelper.database;
    final maps = await db.query(
      'warehouses',
      where: 'name LIKE ?',
      whereArgs: ['%$query%'],
      orderBy: 'name ASC',
    );
    return maps.map((map) => Warehouse.fromMap(map)).toList();
  }

  /// Toggle warehouse active status.
  Future<int> toggleActive(int id, bool isActive) async {
    final db = await _dbHelper.database;
    return await db.update(
      'warehouses',
      {'is_active': isActive ? 1 : 0},
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// Get total warehouse count.
  Future<int> getCount() async {
    final db = await _dbHelper.database;
    final result = await db.rawQuery('SELECT COUNT(*) as count FROM warehouses');
    return (result.first['count'] as int?) ?? 0;
  }
}