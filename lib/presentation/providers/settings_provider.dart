import 'package:flutter/foundation.dart';

class SettingsProvider extends ChangeNotifier {
  String _currency = 'CNY';
  bool _isDarkMode = false;
  bool _budgetAlert = true;
  bool _inventoryAlert = true;
  bool _expirationAlert = true;

  String get currency => _currency;
  bool get isDarkMode => _isDarkMode;
  bool get budgetAlert => _budgetAlert;
  bool get inventoryAlert => _inventoryAlert;
  bool get expirationAlert => _expirationAlert;

  Future<void> loadSettings() async {
    // TODO: Load from SharedPreferences or database
    notifyListeners();
  }

  void setCurrency(String value) {
    _currency = value;
    notifyListeners();
  }

  void setDarkMode(bool value) {
    _isDarkMode = value;
    notifyListeners();
  }

  void setBudgetAlert(bool value) {
    _budgetAlert = value;
    notifyListeners();
  }

  void setInventoryAlert(bool value) {
    _inventoryAlert = value;
    notifyListeners();
  }

  void setExpirationAlert(bool value) {
    _expirationAlert = value;
    notifyListeners();
  }
}
