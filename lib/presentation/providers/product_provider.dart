import 'package:flutter/foundation.dart';
import '../../data/database/database_helper.dart';
import '../../data/models/product.dart';
import '../../data/repositories/product_repository.dart';

/// Provider for product operations.
class ProductProvider extends ChangeNotifier {
  final DatabaseHelper _db = DatabaseHelper();
  final ProductRepository _repository;

  List<Product> _products = [];
  bool _isLoading = false;

  ProductProvider({ProductRepository? repository})
      : _repository = repository ?? ProductRepository();

  List<Product> get products => _products;
  bool get isLoading => _isLoading;

  /// Load all products from database.
  Future<void> loadProducts() async {
    _isLoading = true;
    notifyListeners();

    try {
      _products = await _repository.getAll();
    } catch (e) {
      debugPrint('Error loading products: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  /// Add a new product.
  Future<void> addProduct(Product product) async {
    try {
      await _repository.insert(product);
      await loadProducts();
    } catch (e) {
      debugPrint('Error adding product: $e');
      rethrow;
    }
  }

  /// Update an existing product.
  Future<void> updateProduct(Product product) async {
    try {
      await _repository.update(product);
      await loadProducts();
    } catch (e) {
      debugPrint('Error updating product: $e');
      rethrow;
    }
  }

  /// Delete a product.
  Future<void> deleteProduct(int id) async {
    try {
      await _repository.delete(id);
      await loadProducts();
    } catch (e) {
      debugPrint('Error deleting product: $e');
      rethrow;
    }
  }

  /// Get product by ID.
  Product? getProductById(int id) {
    return _products.where((p) => p.id == id).firstOrNull;
  }

  /// Search products by name.
  List<Product> searchProducts(String query) {
    if (query.isEmpty) return _products;
    return _products
        .where((p) => p.name.toLowerCase().contains(query.toLowerCase()))
        .toList();
  }

  /// Get products by category.
  List<Product> getProductsByCategory(String category) {
    return _products.where((p) => p.category == category).toList();
  }

  /// Get all unique categories.
  List<String> get allCategories {
    return _products
        .where((p) => p.category != null)
        .map((p) => p.category!)
        .toSet()
        .toList();
  }

  /// Get total product count.
  int get totalCount => _products.length;

  /// Get product by SKU.
  Product? getProductBySku(String sku) {
    return _products.where((p) => p.sku == sku).firstOrNull;
  }
}