import '../../data/models/category.dart';
import '../../data/models/transaction.dart' as tx;

/// Smart categorization service that learns from transaction history
/// to suggest categories for new transactions.
class SmartCategorization {
  /// Suggest categories for a new transaction based on description.
  /// Returns top 3 suggested categories with confidence scores.
  List<CategorySuggestion> suggestCategories({
    required String description,
    required List<Category> categories,
    required List<tx.Transaction> transactions,
    double? amount,
    int? date,
  }) {
    if (description.isEmpty && amount == null) {
      return _getTopUsedCategories(categories, transactions, limit: 3);
    }

    final suggestions = <CategorySuggestion>[];
    final descLower = description.toLowerCase();

    // 1. Keyword matching
    final keywordMatches = _matchByKeywords(descLower, categories, transactions);
    suggestions.addAll(keywordMatches);

    // 2. Amount-based matching (±10% similarity)
    if (amount != null) {
      final amountMatches = _matchByAmount(amount, categories, transactions);
      suggestions.addAll(amountMatches);
    }

    // 3. Time-based pattern (same day of month ±3 days)
    if (date != null) {
      final timeMatches = _matchByTime(date, categories, transactions);
      suggestions.addAll(timeMatches);
    }

    // 4. Fallback to top used categories
    if (suggestions.isEmpty) {
      return _getTopUsedCategories(categories, transactions, limit: 3);
    }

    // Deduplicate and sort by score, return top 3
    final deduped = _deduplicateSuggestions(suggestions);
    deduped.sort((a, b) => b.score.compareTo(a.score));
    return deduped.take(3).toList();
  }

  /// Match categories by keyword extraction from description.
  List<CategorySuggestion> _matchByKeywords(
    String description,
    List<Category> categories,
    List<tx.Transaction> transactions,
  ) {
    final suggestions = <CategorySuggestion>[];
    final words = description.split(RegExp(r'\s+')).where((w) => w.length > 2).toList();

    if (words.isEmpty) return suggestions;

    for (final transaction in transactions) {
      if (transaction.description == null) continue;

      final txDescLower = transaction.description!.toLowerCase();
      final matchCount = words.where((w) => txDescLower.contains(w)).length;

      if (matchCount > 0) {
        final category = categories.where((c) => c.id == transaction.categoryId).firstOrNull;
        if (category != null) {
          final score = (matchCount / words.length) * 0.6; // 60% weight for keyword
          suggestions.add(CategorySuggestion(
            category: category,
            score: score,
            reason: '关键词匹配: ${words.take(2).join(", ")}',
          ));
        }
      }
    }

    return suggestions;
  }

  /// Match categories by similar amount (±10%).
  List<CategorySuggestion> _matchByAmount(
    double amount,
    List<Category> categories,
    List<tx.Transaction> transactions,
  ) {
    final suggestions = <CategorySuggestion>[];
    final tolerance = amount * 0.1;

    for (final transaction in transactions) {
      if (transaction.categoryId == null) continue;

      final diff = (transaction.amount - amount).abs();
      if (diff <= tolerance) {
        final category = categories.where((c) => c.id == transaction.categoryId).firstOrNull;
        if (category != null) {
          final score = (1 - diff / amount) * 0.25; // 25% weight for amount
          suggestions.add(CategorySuggestion(
            category: category,
            score: score,
            reason: '金额相似: ¥${transaction.amount.toStringAsFixed(2)}',
          ));
        }
      }
    }

    return suggestions;
  }

  /// Match categories by recurring time pattern (same day of month).
  List<CategorySuggestion> _matchByTime(
    int date,
    List<Category> categories,
    List<tx.Transaction> transactions,
  ) {
    final suggestions = <CategorySuggestion>[];
    final txDate = DateTime.fromMillisecondsSinceEpoch(date);
    final targetDay = txDate.day;

    for (final transaction in transactions) {
      if (transaction.categoryId == null) continue;

      final otherDate = DateTime.fromMillisecondsSinceEpoch(transaction.date);
      final dayDiff = (otherDate.day - targetDay).abs();

      // Same day ±3 days
      if (dayDiff <= 3 && dayDiff > 0) {
        final category = categories.where((c) => c.id == transaction.categoryId).firstOrNull;
        if (category != null) {
          final score = (1 - dayDiff / 3) * 0.15; // 15% weight for time
          suggestions.add(CategorySuggestion(
            category: category,
            score: score,
            reason: '周期模式: 每月 ${otherDate.day} 日',
          ));
        }
      }
    }

    return suggestions;
  }

  /// Get top categories by usage count.
  List<CategorySuggestion> _getTopUsedCategories(
    List<Category> categories,
    List<tx.Transaction> transactions,
    {int limit = 3}
  ) {
    final usageCount = <int, int>{};
    for (final t in transactions) {
      if (t.categoryId != null) {
        usageCount[t.categoryId!] = (usageCount[t.categoryId!] ?? 0) + 1;
      }
    }

    final sorted = categories.where((c) => usageCount.containsKey(c.id)).toList()
      ..sort((a, b) => (usageCount[b.id] ?? 0).compareTo(usageCount[a.id] ?? 0));

    return sorted.take(limit).map((c) => CategorySuggestion(
      category: c,
      score: 0.3,
      reason: '常用分类 (${usageCount[c.id]}次)',
    )).toList();
  }

  /// Deduplicate suggestions keeping highest score per category.
  List<CategorySuggestion> _deduplicateSuggestions(List<CategorySuggestion> suggestions) {
    final Map<int, CategorySuggestion> deduped = {};
    for (final s in suggestions) {
      final existing = deduped[s.category.id];
      if (existing == null || s.score > existing.score) {
        deduped[s.category.id!] = s;
      }
    }
    return deduped.values.toList();
  }

  /// Update category usage count after a transaction is created.
  /// Call this after user confirms a transaction with a category.
  void recordCategoryUsage({
    required int categoryId,
    required List<Category> categories,
  }) {
    final index = categories.indexWhere((c) => c.id == categoryId);
    if (index != -1) {
      // Usage count is incremented via CategoryRepository.incrementUsageCount()
      // This method is for any future ML/advanced logic
    }
  }
}

/// A category suggestion with confidence score and reason.
class CategorySuggestion {
  final Category category;
  final double score;
  final String reason;

  CategorySuggestion({
    required this.category,
    required this.score,
    required this.reason,
  });

  /// Get confidence level as a string.
  String get confidenceLevel {
    if (score >= 0.7) return '高';
    if (score >= 0.4) return '中';
    return '低';
  }
}
