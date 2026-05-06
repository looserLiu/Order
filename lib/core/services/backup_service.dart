import 'dart:io';
import 'package:csv/csv.dart';
import 'package:path_provider/path_provider.dart';
import 'package:intl/intl.dart';
import '../../data/database/database_helper.dart';
import '../../data/models/account.dart';
import '../../data/models/category.dart';
import '../../data/models/transaction.dart' as tx;
import '../../data/models/budget.dart';
import '../../data/models/product.dart';
import '../../data/models/warehouse.dart';
import '../../data/models/inventory_flow.dart';
import '../../data/repositories/account_repository.dart';
import '../../data/repositories/category_repository.dart';
import '../../data/repositories/transaction_repository.dart';
import '../../data/repositories/budget_repository.dart';
import '../../data/repositories/product_repository.dart';
import '../../data/repositories/warehouse_repository.dart';
import '../../data/repositories/inventory_flow_repository.dart';

/// Service for exporting and importing app data in CSV format.
class BackupService {
  final DatabaseHelper _dbHelper = DatabaseHelper();

  /// Export all data to a backup directory.
  /// Returns the path to the backup directory.
  Future<String> exportAllData() async {
    final backupDir = await _getBackupDirectory();
    final timestamp = DateFormat('yyyyMMdd_HHmmss').format(DateTime.now());
    final exportDir = Directory('${backupDir.path}/backup_$timestamp');
    await exportDir.create(recursive: true);

    // Export each table to CSV
    await _exportAccounts(exportDir);
    await _exportCategories(exportDir);
    await _exportTransactions(exportDir);
    await _exportBudgets(exportDir);
    await _exportProducts(exportDir);
    await _exportWarehouses(exportDir);
    await _exportInventoryFlows(exportDir);

    // Create manifest
    await _createManifest(exportDir);

    return exportDir.path;
  }

  /// Import data from a backup directory.
  Future<ImportResult> importData(String backupPath) async {
    final importDir = Directory(backupPath);
    if (!await importDir.exists()) {
      return ImportResult(success: false, message: '备份目录不存在');
    }

    int importedCount = 0;
    int errorCount = 0;

    try {
      // Import in order (respecting foreign keys)
      final accResult = await _importAccounts(importDir);
      importedCount += accResult.imported;
      errorCount += accResult.errors;

      final catResult = await _importCategories(importDir);
      importedCount += catResult.imported;
      errorCount += catResult.errors;

      final txResult = await _importTransactions(importDir);
      importedCount += txResult.imported;
      errorCount += txResult.errors;

      final budgetResult = await _importBudgets(importDir);
      importedCount += budgetResult.imported;
      errorCount += budgetResult.errors;

      final prodResult = await _importProducts(importDir);
      importedCount += prodResult.imported;
      errorCount += prodResult.errors;

      final whResult = await _importWarehouses(importDir);
      importedCount += whResult.imported;
      errorCount += whResult.errors;

      final flowResult = await _importInventoryFlows(importDir);
      importedCount += flowResult.imported;
      errorCount += flowResult.errors;

      return ImportResult(
        success: true,
        message: '导入完成: $importedCount 条成功, $errorCount 条失败',
        imported: importedCount,
        errors: errorCount,
      );
    } catch (e) {
      return ImportResult(
        success: false,
        message: '导入失败: $e',
        imported: importedCount,
        errors: errorCount,
      );
    }
  }

  /// Export transactions to CSV file.
  Future<String?> exportTransactionsCSV() async {
    try {
      final repo = TransactionRepository();
      final transactions = await repo.getAll();
      
      if (transactions.isEmpty) {
        return null;
      }

      final rows = <List<dynamic>>[
        ['id', 'account_id', 'category_id', 'amount', 'type', 'description', 'date', 'created_at'],
        ...transactions.map((t) => [
          t.id,
          t.accountId,
          t.categoryId,
          t.amount,
          t.type,
          t.description ?? '',
          t.date,
          t.createdAt,
        ]),
      ];

      final csv = const ListToCsvConverter().convert(rows);
      final dir = await _getBackupDirectory();
      final timestamp = DateFormat('yyyyMMdd_HHmmss').format(DateTime.now());
      final file = File('${dir.path}/transactions_$timestamp.csv');
      await file.writeAsString(csv);

      return file.path;
    } catch (e) {
      return null;
    }
  }

  /// Import transactions from CSV file.
  Future<ImportResult> importTransactionsCSV(String filePath) async {
    try {
      final file = File(filePath);
      if (!await file.exists()) {
        return ImportResult(success: false, message: '文件不存在');
      }

      final csv = await file.readAsString();
      final rows = const CsvToListConverter().convert(csv);
      
      if (rows.isEmpty) {
        return ImportResult(success: false, message: 'CSV 文件为空');
      }

      // Skip header row
      final dataRows = rows.skip(1).toList();
      final repo = TransactionRepository();
      int imported = 0;
      int errors = 0;

      for (final row in dataRows) {
        try {
          final transaction = tx.Transaction(
            accountId: row[1] as int,
            categoryId: row[2] as int?,
            amount: (row[3] as num).toDouble(),
            type: row[4] as String,
            description: row[5] as String?,
            date: row[6] as int,
            createdAt: row[7] as int,
          );
          await repo.insert(transaction);
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return ImportResult(
        success: true,
        message: '导入完成: $imported 条成功, $errors 条失败',
        imported: imported,
        errors: errors,
      );
    } catch (e) {
      return ImportResult(success: false, message: '导入失败: $e');
    }
  }

  /// Get list of available backups.
  Future<List<BackupInfo>> listBackups() async {
    final backupDir = await _getBackupDirectory();
    final dirs = await backupDir.list().toList();
    
    final backups = <BackupInfo>[];
    for (final entity in dirs) {
      if (entity is Directory && entity.path.contains('backup_')) {
        final name = entity.path.split('/').last;
        final stat = await entity.stat();
        backups.add(BackupInfo(
          path: entity.path,
          name: name,
          createdAt: stat.modified,
        ));
      }
    }

    backups.sort((a, b) => b.createdAt.compareTo(a.createdAt));
    return backups;
  }

  /// Delete a backup.
  Future<bool> deleteBackup(String path) async {
    try {
      final dir = Directory(path);
      if (await dir.exists()) {
        await dir.delete(recursive: true);
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  // Private helpers

  Future<Directory> _getBackupDirectory() async {
    final appDir = await getApplicationDocumentsDirectory();
    final backupDir = Directory('${appDir.path}/backups');
    if (!await backupDir.exists()) {
      await backupDir.create(recursive: true);
    }
    return backupDir;
  }

  Future<void> _exportAccounts(Directory dir) async {
    final repo = AccountRepository();
    final accounts = await repo.getAll();
    await _writeCSV(
      '${dir.path}/accounts.csv',
      ['id', 'name', 'type', 'balance', 'currency', 'icon', 'color', 'created_at', 'updated_at'],
      accounts.map((a) => [a.id, a.name, a.type, a.balance, a.currency, a.icon ?? '', a.color ?? '', a.createdAt, a.updatedAt]).toList(),
    );
  }

  Future<void> _exportCategories(Directory dir) async {
    final repo = CategoryRepository();
    final categories = await repo.getAll();
    await _writeCSV(
      '${dir.path}/categories.csv',
      ['id', 'name', 'type', 'icon', 'color', 'parent_id', 'is_system', 'usage_count', 'created_at'],
      categories.map((c) => [c.id, c.name, c.type, c.icon ?? '', c.color ?? '', c.parentId ?? '', c.isSystem, c.usageCount, c.createdAt]).toList(),
    );
  }

  Future<void> _exportTransactions(Directory dir) async {
    final repo = TransactionRepository();
    final transactions = await repo.getAll();
    await _writeCSV(
      '${dir.path}/transactions.csv',
      ['id', 'account_id', 'category_id', 'amount', 'type', 'description', 'date', 'created_at'],
      transactions.map((t) => [t.id, t.accountId, t.categoryId ?? '', t.amount, t.type, t.description ?? '', t.date, t.createdAt]).toList(),
    );
  }

  Future<void> _exportBudgets(Directory dir) async {
    final repo = BudgetRepository();
    final budgets = await repo.getAll();
    await _writeCSV(
      '${dir.path}/budgets.csv',
      ['id', 'category_id', 'amount', 'period', 'start_date', 'end_date'],
      budgets.map((b) => [b.id, b.categoryId, b.amount, b.period, b.startDate, b.endDate]).toList(),
    );
  }

  Future<void> _exportProducts(Directory dir) async {
    final repo = ProductRepository();
    final products = await repo.getAll();
    await _writeCSV(
      '${dir.path}/products.csv',
      ['id', 'name', 'sku', 'category', 'unit', 'cost_price', 'sale_price', 'image_url', 'created_at', 'updated_at'],
      products.map((p) => [p.id, p.name, p.sku ?? '', p.category ?? '', p.unit ?? '', p.costPrice, p.salePrice, p.imageUrl ?? '', p.createdAt, p.updatedAt]).toList(),
    );
  }

  Future<void> _exportWarehouses(Directory dir) async {
    final repo = WarehouseRepository();
    final warehouses = await repo.getAll();
    await _writeCSV(
      '${dir.path}/warehouses.csv',
      ['id', 'name', 'location', 'description', 'is_active', 'created_at'],
      warehouses.map((w) => [w.id, w.name, w.location ?? '', w.description ?? '', w.isActive, w.createdAt]).toList(),
    );
  }

  Future<void> _exportInventoryFlows(Directory dir) async {
    final repo = InventoryFlowRepository();
    final flows = await repo.getAll();
    await _writeCSV(
      '${dir.path}/inventory_flows.csv',
      ['id', 'product_id', 'warehouse_id', 'flow_type', 'quantity', 'batch_number', 'expiration_date', 'cost_at_flow', 'reference_id', 'date', 'created_at'],
      flows.map((f) => [f.id, f.productId, f.warehouseId, f.flowType, f.quantity, f.batchNumber ?? '', f.expirationDate ?? '', f.costAtFlow ?? '', f.referenceId ?? '', f.date, f.createdAt]).toList(),
    );
  }

  Future<void> _writeCSV(String path, List<String> headers, List<List<dynamic>> rows) async {
    final allRows = [headers, ...rows];
    final csv = const ListToCsvConverter().convert(allRows);
    await File(path).writeAsString(csv);
  }

  Future<void> _createManifest(Directory dir) async {
    final manifest = '''
Order App Backup
Created: ${DateTime.now().toIso8601String()}
Version: 1.0.0

Files:
- accounts.csv
- categories.csv
- transactions.csv
- budgets.csv
- products.csv
- warehouses.csv
- inventory_flows.csv
''';
    await File('${dir.path}/manifest.txt').writeAsString(manifest);
  }

  Future<_ImportResult> _importAccounts(Directory dir) async {
    try {
      final file = File('${dir.path}/accounts.csv');
      if (!await file.exists()) return _ImportResult();

      final rows = await _readCSV(file.path);
      final repo = AccountRepository();
      int imported = 0, errors = 0;

      for (final row in rows.skip(1)) {
        try {
          await repo.insert(Account(
            name: row[1],
            type: row[2],
            balance: (row[3] as num).toDouble(),
            currency: row[4],
            icon: row[5].toString().isEmpty ? null : row[5].toString(),
            color: row[6].toString().isEmpty ? null : int.tryParse(row[6].toString()),
            createdAt: int.parse(row[7].toString()),
            updatedAt: int.parse(row[8].toString()),
          ));
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return _ImportResult(imported: imported, errors: errors);
    } catch (e) {
      return _ImportResult(imported: 0, errors: 0);
    }
  }

  Future<_ImportResult> _importCategories(Directory dir) async {
    try {
      final file = File('${dir.path}/categories.csv');
      if (!await file.exists()) return _ImportResult();

      final rows = await _readCSV(file.path);
      final repo = CategoryRepository();
      int imported = 0, errors = 0;

      for (final row in rows.skip(1)) {
        try {
          await repo.insert(Category(
            name: row[1],
            type: row[2],
            icon: row[3].toString().isEmpty ? null : row[3].toString(),
            color: row[4].toString().isEmpty ? null : int.tryParse(row[4].toString()),
            parentId: row[5].toString().isEmpty ? null : int.tryParse(row[5].toString()),
            isSystem: int.parse(row[6].toString()) == 1,
            usageCount: int.parse(row[7].toString()),
            createdAt: int.parse(row[8].toString()),
          ));
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return _ImportResult(imported: imported, errors: errors);
    } catch (e) {
      return _ImportResult(imported: 0, errors: 0);
    }
  }

  Future<_ImportResult> _importTransactions(Directory dir) async {
    try {
      final file = File('${dir.path}/transactions.csv');
      if (!await file.exists()) return _ImportResult();

      final rows = await _readCSV(file.path);
      final repo = TransactionRepository();
      int imported = 0, errors = 0;

      for (final row in rows.skip(1)) {
        try {
          await repo.insert(tx.Transaction(
            accountId: int.parse(row[1].toString()),
            categoryId: row[2].toString().isEmpty ? null : int.parse(row[2].toString()),
            amount: (row[3] as num).toDouble(),
            type: row[4],
            description: row[5].toString().isEmpty ? null : row[5].toString(),
            date: int.parse(row[6].toString()),
            createdAt: int.parse(row[7].toString()),
          ));
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return _ImportResult(imported: imported, errors: errors);
    } catch (e) {
      return _ImportResult(imported: 0, errors: 0);
    }
  }

  Future<_ImportResult> _importBudgets(Directory dir) async {
    try {
      final file = File('${dir.path}/budgets.csv');
      if (!await file.exists()) return _ImportResult();

      final rows = await _readCSV(file.path);
      final repo = BudgetRepository();
      int imported = 0, errors = 0;

      for (final row in rows.skip(1)) {
        try {
          await repo.insert(Budget(
            categoryId: int.parse(row[1].toString()),
            amount: (row[2] as num).toDouble(),
            period: row[3],
            startDate: int.parse(row[4].toString()),
            endDate: int.parse(row[5].toString()),
          ));
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return _ImportResult(imported: imported, errors: errors);
    } catch (e) {
      return _ImportResult(imported: 0, errors: 0);
    }
  }

  Future<_ImportResult> _importProducts(Directory dir) async {
    try {
      final file = File('${dir.path}/products.csv');
      if (!await file.exists()) return _ImportResult();

      final rows = await _readCSV(file.path);
      final repo = ProductRepository();
      int imported = 0, errors = 0;

      for (final row in rows.skip(1)) {
        try {
          await repo.insert(Product(
            name: row[1],
            sku: row[2].toString().isEmpty ? null : row[2].toString(),
            category: row[3].toString().isEmpty ? null : row[3].toString(),
            unit: row[4].toString().isEmpty ? null : row[4].toString(),
            costPrice: (row[5] as num).toDouble(),
            salePrice: (row[6] as num).toDouble(),
            imageUrl: row[7].toString().isEmpty ? null : row[7].toString(),
            createdAt: int.parse(row[8].toString()),
            updatedAt: int.parse(row[9].toString()),
          ));
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return _ImportResult(imported: imported, errors: errors);
    } catch (e) {
      return _ImportResult(imported: 0, errors: 0);
    }
  }

  Future<_ImportResult> _importWarehouses(Directory dir) async {
    try {
      final file = File('${dir.path}/warehouses.csv');
      if (!await file.exists()) return _ImportResult();

      final rows = await _readCSV(file.path);
      final repo = WarehouseRepository();
      int imported = 0, errors = 0;

      for (final row in rows.skip(1)) {
        try {
          await repo.insert(Warehouse(
            name: row[1],
            location: row[2].toString().isEmpty ? null : row[2].toString(),
            description: row[3].toString().isEmpty ? null : row[3].toString(),
            isActive: int.parse(row[4].toString()),
            createdAt: int.parse(row[5].toString()),
          ));
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return _ImportResult(imported: imported, errors: errors);
    } catch (e) {
      return _ImportResult(imported: 0, errors: 0);
    }
  }

  Future<_ImportResult> _importInventoryFlows(Directory dir) async {
    try {
      final file = File('${dir.path}/inventory_flows.csv');
      if (!await file.exists()) return _ImportResult();

      final rows = await _readCSV(file.path);
      final repo = InventoryFlowRepository();
      int imported = 0, errors = 0;

      for (final row in rows.skip(1)) {
        try {
          await repo.insert(InventoryFlow(
            productId: int.parse(row[1].toString()),
            warehouseId: int.parse(row[2].toString()),
            flowType: row[3],
            quantity: (row[4] as num).toDouble(),
            batchNumber: row[5].toString().isEmpty ? null : row[5].toString(),
            expirationDate: row[6].toString().isEmpty ? null : int.parse(row[6].toString()),
            costAtFlow: row[7].toString().isEmpty ? null : (row[7] as num).toDouble(),
            referenceId: row[8].toString().isEmpty ? null : row[8].toString(),
            date: int.parse(row[9].toString()),
            createdAt: int.parse(row[10].toString()),
          ));
          imported++;
        } catch (e) {
          errors++;
        }
      }

      return _ImportResult(imported: imported, errors: errors);
    } catch (e) {
      return _ImportResult(imported: 0, errors: 0);
    }
  }

  Future<List<List<dynamic>>> _readCSV(String path) async {
    final file = File(path);
    final csv = await file.readAsString();
    return const CsvToListConverter().convert(csv);
  }
}

/// Result of an import operation.
class ImportResult {
  final bool success;
  final String message;
  final int imported;
  final int errors;

  ImportResult({
    this.success = false,
    this.message = '',
    this.imported = 0,
    this.errors = 0,
  });
}

/// Internal import result for tracking counts.
class _ImportResult {
  final int imported;
  final int errors;

  _ImportResult({this.imported = 0, this.errors = 0});
}

/// Information about a backup.
class BackupInfo {
  final String path;
  final String name;
  final DateTime createdAt;

  BackupInfo({
    required this.path,
    required this.name,
    required this.createdAt,
  });
}
