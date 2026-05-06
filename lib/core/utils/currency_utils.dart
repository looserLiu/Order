import 'package:intl/intl.dart';

class CurrencyUtils {
  static final NumberFormat _currencyFormat = NumberFormat.currency(
    symbol: '\$',
    decimalDigits: 2,
  );

  static final NumberFormat _compactFormat = NumberFormat.compact();

  static String format(double amount) {
    return _currencyFormat.format(amount);
  }

  static String formatCompact(double amount) {
    return _compactFormat.format(amount);
  }

  static String formatWithSymbol(double amount, String symbol) {
    return NumberFormat.currency(symbol: symbol, decimalDigits: 2).format(amount);
  }

  static double? parse(String? value) {
    if (value == null || value.isEmpty) return null;
    try {
      final cleanValue = value.replaceAll(RegExp(r'[^\d.]'), '');
      return double.tryParse(cleanValue);
    } catch (_) {
      return null;
    }
  }
}
