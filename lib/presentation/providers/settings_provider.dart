import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SettingsProvider extends ChangeNotifier {
  static const String _keyCurrency = 'currency';
  static const String _keyDarkMode = 'dark_mode';
  static const String _keyBudgetAlert = 'budget_alert';
  static const String _keyInventoryAlert = 'inventory_alert';
  static const String _keyExpirationAlert = 'expiration_alert';
  static const String _keyLowStockThreshold = 'low_stock_threshold';
  static const String _keyExpirationDays = 'expiration_days';

  String _currency = 'CNY';
  bool _isDarkMode = false;
  bool _budgetAlert = true;
  bool _inventoryAlert = true;
  bool _expirationAlert = true;
  double _lowStockThreshold = 10.0;
  int _expirationDays = 7;
  bool _isLoaded = false;

  String get currency => _currency;
  bool get isDarkMode => _isDarkMode;
  bool get budgetAlert => _budgetAlert;
  bool get inventoryAlert => _inventoryAlert;
  bool get expirationAlert => _expirationAlert;
  double get lowStockThreshold => _lowStockThreshold;
  int get expirationDays => _expirationDays;

  Future<void> loadSettings() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      _currency = prefs.getString(_keyCurrency) ?? 'CNY';
      _isDarkMode = prefs.getBool(_keyDarkMode) ?? false;
      _budgetAlert = prefs.getBool(_keyBudgetAlert) ?? true;
      _inventoryAlert = prefs.getBool(_keyInventoryAlert) ?? true;
      _expirationAlert = prefs.getBool(_keyExpirationAlert) ?? true;
      _lowStockThreshold = prefs.getDouble(_keyLowStockThreshold) ?? 10.0;
      _expirationDays = prefs.getInt(_keyExpirationDays) ?? 7;
      _isLoaded = true;
    } catch (e) {
      debugPrint('Error loading settings: $e');
    }
    notifyListeners();
  }

  Future<void> setCurrency(String value) async {
    _currency = value;
    notifyListeners();
    await _savePreference(_keyCurrency, value);
  }

  Future<void> setDarkMode(bool value) async {
    _isDarkMode = value;
    notifyListeners();
    await _savePreference(_keyDarkMode, value);
  }

  Future<void> setBudgetAlert(bool value) async {
    _budgetAlert = value;
    notifyListeners();
    await _savePreference(_keyBudgetAlert, value);
  }

  Future<void> setInventoryAlert(bool value) async {
    _inventoryAlert = value;
    notifyListeners();
    await _savePreference(_keyInventoryAlert, value);
  }

  Future<void> setExpirationAlert(bool value) async {
    _expirationAlert = value;
    notifyListeners();
    await _savePreference(_keyExpirationAlert, value);
  }

  Future<void> setLowStockThreshold(double value) async {
    _lowStockThreshold = value;
    notifyListeners();
    await _savePreference(_keyLowStockThreshold, value);
  }

  Future<void> setExpirationDays(int value) async {
    _expirationDays = value;
    notifyListeners();
    await _savePreference(_keyExpirationDays, value);
  }

  Future<void> _savePreference(String key, dynamic value) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      if (value is String) {
        await prefs.setString(key, value);
      } else if (value is bool) {
        await prefs.setBool(key, value);
      } else if (value is double) {
        await prefs.setDouble(key, value);
      } else if (value is int) {
        await prefs.setInt(key, value);
      }
    } catch (e) {
      debugPrint('Error saving preference: $e');
    }
  }

  /// Reset all settings to defaults.
  Future<void> resetToDefaults() async {
    _currency = 'CNY';
    _isDarkMode = false;
    _budgetAlert = true;
    _inventoryAlert = true;
    _expirationAlert = true;
    _lowStockThreshold = 10.0;
    _expirationDays = 7;
    
    notifyListeners();
    
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.remove(_keyCurrency);
      await prefs.remove(_keyDarkMode);
      await prefs.remove(_keyBudgetAlert);
      await prefs.remove(_keyInventoryAlert);
      await prefs.remove(_keyExpirationAlert);
      await prefs.remove(_keyLowStockThreshold);
      await prefs.remove(_keyExpirationDays);
    } catch (e) {
      debugPrint('Error resetting settings: $e');
    }
  }
}
