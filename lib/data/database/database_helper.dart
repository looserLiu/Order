import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart';

/// Database helper for SQLite operations.
/// Handles database creation, migrations, and provides access to the database instance.
class DatabaseHelper {
  static const String _databaseName = 'order.db';
  static const int _databaseVersion = 1;

  // Singleton instance
  static DatabaseHelper? _instance;
  static Database? _database;

  DatabaseHelper._internal();

  factory DatabaseHelper() {
    _instance ??= DatabaseHelper._internal();
    return _instance!;
  }

  /// Get the database instance, creating it if necessary.
  Future<Database> get database async {
    _database ??= await _initDatabase();
    return _database!;
  }

  /// Initialize the database.
  Future<Database> _initDatabase() async {
    final databasePath = await getDatabasesPath();
    final path = join(databasePath, _databaseName);

    return await openDatabase(
      path,
      version: _databaseVersion,
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
  }

  /// Create database tables.
  Future<void> _onCreate(Database db, int version) async {
    // Accounts table
    await db.execute('''
      CREATE TABLE accounts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        type TEXT NOT NULL,
        balance REAL NOT NULL DEFAULT 0.0,
        currency TEXT NOT NULL DEFAULT 'CNY',
        icon TEXT,
        color INTEGER,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');

    // Categories table
    await db.execute('''
      CREATE TABLE categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        type TEXT NOT NULL,
        icon TEXT,
        color INTEGER,
        parent_id INTEGER,
        is_system INTEGER NOT NULL DEFAULT 0,
        usage_count INTEGER NOT NULL DEFAULT 0,
        FOREIGN KEY (parent_id) REFERENCES categories (id)
      )
    ''');

    // Transactions table
    await db.execute('''
      CREATE TABLE transactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        account_id INTEGER NOT NULL,
        category_id INTEGER,
        amount REAL NOT NULL,
        type TEXT NOT NULL,
        description TEXT,
        date INTEGER NOT NULL,
        created_at INTEGER NOT NULL,
        FOREIGN KEY (account_id) REFERENCES accounts (id),
        FOREIGN KEY (category_id) REFERENCES categories (id)
      )
    ''');

    // Products table
    await db.execute('''
      CREATE TABLE products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        sku TEXT,
        category TEXT,
        unit TEXT,
        cost_price REAL NOT NULL DEFAULT 0.0,
        sale_price REAL NOT NULL DEFAULT 0.0,
        image_url TEXT,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');

    // Warehouses table
    await db.execute('''
      CREATE TABLE warehouses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        location TEXT,
        description TEXT,
        is_active INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL
      )
    ''');

    // Inventory flows table
    await db.execute('''
      CREATE TABLE inventory_flows (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        product_id INTEGER NOT NULL,
        warehouse_id INTEGER NOT NULL,
        flow_type TEXT NOT NULL,
        quantity REAL NOT NULL,
        batch_number TEXT,
        expiration_date INTEGER,
        cost_at_flow REAL,
        reference_id TEXT,
        date INTEGER NOT NULL,
        created_at INTEGER NOT NULL,
        FOREIGN KEY (product_id) REFERENCES products (id),
        FOREIGN KEY (warehouse_id) REFERENCES warehouses (id)
      )
    ''');

    // Budgets table
    await db.execute('''
      CREATE TABLE budgets (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_id INTEGER NOT NULL,
        amount REAL NOT NULL,
        period TEXT NOT NULL,
        start_date INTEGER NOT NULL,
        end_date INTEGER NOT NULL,
        FOREIGN KEY (category_id) REFERENCES categories (id)
      )
    ''');

    // Cost categories table
    await db.execute('''
      CREATE TABLE cost_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        type TEXT NOT NULL,
        description TEXT
      )
    ''');

    // Category keywords for smart categorization
    await db.execute('''
      CREATE TABLE category_keywords (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_id INTEGER NOT NULL,
        keyword TEXT NOT NULL,
        weight INTEGER NOT NULL DEFAULT 1,
        FOREIGN KEY (category_id) REFERENCES categories (id)
      )
    ''');

    await db.execute('CREATE INDEX idx_category_keywords_category_id ON category_keywords (category_id)');
    await db.execute('CREATE INDEX idx_category_keywords_keyword ON category_keywords (keyword)');

    // Create indexes for better query performance
    await db.execute('CREATE INDEX idx_transactions_account_id ON transactions (account_id)');
    await db.execute('CREATE INDEX idx_transactions_category_id ON transactions (category_id)');
    await db.execute('CREATE INDEX idx_transactions_date ON transactions (date)');
    await db.execute('CREATE INDEX idx_inventory_flows_product_id ON inventory_flows (product_id)');
    await db.execute('CREATE INDEX idx_inventory_flows_warehouse_id ON inventory_flows (warehouse_id)');
    await db.execute('CREATE INDEX idx_inventory_flows_date ON inventory_flows (date)');

    // Insert default categories
    await _insertDefaultCategories(db);
  }

  /// Handle database upgrades.
  Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    // Handle future migrations here
  }

  /// Insert default income and expense categories.
  Future<void> _insertDefaultCategories(Database db) async {
    final now = DateTime.now().millisecondsSinceEpoch;

    // Default expense categories
    final expenseCategories = [
      {'name': 'Food', 'icon': 'restaurant', 'color': 0xFFE57373},
      {'name': 'Transport', 'icon': 'directions_car', 'color': 0xFF64B5F6},
      {'name': 'Shopping', 'icon': 'shopping_bag', 'color': 0xFFBA68C8},
      {'name': 'Entertainment', 'icon': 'movie', 'color': 0xFFFFD54F},
      {'name': 'Utilities', 'icon': 'bolt', 'color': 0xFF4DB6AC},
      {'name': 'Rent', 'icon': 'home', 'color': 0xFF90A4AE},
      {'name': 'Healthcare', 'icon': 'local_hospital', 'color': 0xFFF06292},
      {'name': 'Education', 'icon': 'school', 'color': 0xFF7E57C2},
      {'name': 'Other Expense', 'icon': 'more_horiz', 'color': 0xFF9E9E9E},
    ];

    // Default income categories
    final incomeCategories = [
      {'name': 'Salary', 'icon': 'work', 'color': 0xFF81C784},
      {'name': 'Bonus', 'icon': 'card_giftcard', 'color': 0xFF64B5F6},
      {'name': 'Investment', 'icon': 'trending_up', 'color': 0xFFFFD54F},
      {'name': 'Freelance', 'icon': 'laptop', 'color': 0xFFBA68C8},
      {'name': 'Other Income', 'icon': 'more_horiz', 'color': 0xFF9E9E9E},
    ];

    for (final category in expenseCategories) {
      await db.insert('categories', {
        ...category,
        'type': 'expense',
        'is_system': 1,
        'usage_count': 0,
        'created_at': now,
      });
    }

    for (final category in incomeCategories) {
      await db.insert('categories', {
        ...category,
        'type': 'income',
        'is_system': 1,
        'usage_count': 0,
        'created_at': now,
      });
    }
  }

  /// Close the database connection.
  Future<void> close() async {
    final db = await database;
    await db.close();
    _database = null;
  }

  /// Delete the database (for testing or reset).
  Future<void> deleteDatabase() async {
    final databasePath = await getDatabasesPath();
    final path = join(databasePath, _databaseName);
    await databaseFactory.deleteDatabase(path);
    _database = null;
  }
}