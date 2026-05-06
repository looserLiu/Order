import '../../data/models/inventory_flow.dart';
import '../../data/models/product.dart';
import '../../data/repositories/inventory_flow_repository.dart';
import '../../data/repositories/product_repository.dart';

/// Service for checking and reporting expiring inventory items.
class ExpirationAlert {
  final InventoryFlowRepository _flowRepo;
  final ProductRepository _productRepo;

  ExpirationAlert({
    InventoryFlowRepository? flowRepo,
    ProductRepository? productRepo,
  })  : _flowRepo = flowRepo ?? InventoryFlowRepository(),
        _productRepo = productRepo ?? ProductRepository();

  /// Get all products expiring within [days] days.
  Future<List<ExpiringItem>> getExpiringItems({int days = 7}) async {
    final flows = await _flowRepo.getExpiringSoon(days);
    final items = <ExpiringItem>[];

    for (final flow in flows) {
      if (flow.expirationDate == null) continue;

      final product = await _productRepo.getById(flow.productId);
      if (product == null) continue;

      final daysUntilExpiration = _calculateDaysUntilExpiration(flow.expirationDate!);
      final urgency = _getUrgency(daysUntilExpiration);

      items.add(ExpiringItem(
        flow: flow,
        product: product,
        daysUntilExpiration: daysUntilExpiration,
        urgency: urgency,
      ));
    }

    // Sort by urgency (most urgent first)
    items.sort((a, b) => a.daysUntilExpiration.compareTo(b.daysUntilExpiration));
    return items;
  }

  /// Get all expired items.
  Future<List<ExpiringItem>> getExpiredItems() async {
    final now = DateTime.now().millisecondsSinceEpoch;
    final items = <ExpiringItem>[];

    // Get all inventory flows with expiration dates
    final flows = await _flowRepo.getAll();
    
    for (final flow in flows) {
      if (flow.expirationDate == null) continue;
      if (flow.expirationDate! > now) continue;

      final product = await _productRepo.getById(flow.productId);
      if (product == null) continue;

      final daysExpired = _calculateDaysExpired(flow.expirationDate!);

      items.add(ExpiringItem(
        flow: flow,
        product: product,
        daysUntilExpiration: -daysExpired,
        urgency: ExpirationUrgency.expired,
      ));
    }

    items.sort((a, b) => a.daysUntilExpiration.compareTo(b.daysUntilExpiration));
    return items;
  }

  /// Get a summary of expiration status.
  Future<ExpirationSummary> getExpirationSummary({int days = 30}) async {
    final expiringItems = await getExpiringItems(days: days);
    final expiredItems = await getExpiredItems();

    int expiredCount = expiredItems.length;
    int urgentCount = 0;  // ≤ 3 days
    int warningCount = 0;  // ≤ 7 days
    int noticeCount = 0;   // ≤ 30 days

    for (final item in expiringItems) {
      switch (item.urgency) {
        case ExpirationUrgency.expired:
          expiredCount++;
          break;
        case ExpirationUrgency.urgent:
          urgentCount++;
          break;
        case ExpirationUrgency.warning:
          warningCount++;
          break;
        case ExpirationUrgency.notice:
          noticeCount++;
          break;
        case ExpirationUrgency.normal:
          break;
      }
    }

    return ExpirationSummary(
      expiredCount: expiredCount,
      urgentCount: urgentCount,
      warningCount: warningCount,
      noticeCount: noticeCount,
      totalAtRisk: expiredCount + urgentCount + warningCount,
    );
  }

  /// Calculate days until expiration (negative if already expired).
  int _calculateDaysUntilExpiration(int expirationTimestamp) {
    final now = DateTime.now();
    final expiration = DateTime.fromMillisecondsSinceEpoch(expirationTimestamp);
    return expiration.difference(now).inDays;
  }

  /// Calculate days since expiration (how long ago it expired).
  int _calculateDaysExpired(int expirationTimestamp) {
    final now = DateTime.now();
    final expiration = DateTime.fromMillisecondsSinceEpoch(expirationTimestamp);
    return now.difference(expiration).inDays;
  }

  /// Get urgency level based on days until expiration.
  ExpirationUrgency _getUrgency(int days) {
    if (days < 0) return ExpirationUrgency.expired;
    if (days <= 3) return ExpirationUrgency.urgent;
    if (days <= 7) return ExpirationUrgency.warning;
    if (days <= 30) return ExpirationUrgency.notice;
    return ExpirationUrgency.normal;
  }
}

/// An item that is expiring or expired.
class ExpiringItem {
  final InventoryFlow flow;
  final Product product;
  final int daysUntilExpiration;
  final ExpirationUrgency urgency;

  ExpiringItem({
    required this.flow,
    required this.product,
    required this.daysUntilExpiration,
    required this.urgency,
  });

  /// Human-readable expiration text.
  String get expirationText {
    if (daysUntilExpiration < 0) {
      return '已过期 ${-daysUntilExpiration} 天';
    } else if (daysUntilExpiration == 0) {
      return '今天过期';
    } else if (daysUntilExpiration == 1) {
      return '明天过期';
    } else {
      return '还有 $daysUntilExpiration 天过期';
    }
  }
}

/// Summary of expiration status.
class ExpirationSummary {
  final int expiredCount;
  final int urgentCount;
  final int warningCount;
  final int noticeCount;
  final int totalAtRisk;

  ExpirationSummary({
    required this.expiredCount,
    required this.urgentCount,
    required this.warningCount,
    required this.noticeCount,
    required this.totalAtRisk,
  });

  /// Whether there are any items that need attention.
  bool get hasAlerts => expiredCount > 0 || urgentCount > 0 || warningCount > 0;

  /// Get the highest urgency level present.
  ExpirationUrgency get highestUrgency {
    if (expiredCount > 0) return ExpirationUrgency.expired;
    if (urgentCount > 0) return ExpirationUrgency.urgent;
    if (warningCount > 0) return ExpirationUrgency.warning;
    if (noticeCount > 0) return ExpirationUrgency.notice;
    return ExpirationUrgency.normal;
  }
}

/// Urgency level for expiration.
enum ExpirationUrgency {
  normal,    // > 30 days
  notice,    // ≤ 30 days
  warning,   // ≤ 7 days
  urgent,    // ≤ 3 days
  expired,   // Already expired
}
