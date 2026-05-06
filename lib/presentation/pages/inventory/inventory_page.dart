import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/inventory_provider.dart';
import '../../providers/product_provider.dart';
import '../../providers/warehouse_provider.dart';
import 'products_page.dart';
import 'stock_in_page.dart';
import 'stock_out_page.dart';
import 'warehouses_page.dart';
import '../../../data/models/product.dart';
import '../../../data/models/inventory_flow.dart';

class InventoryPage extends StatefulWidget {
  const InventoryPage({super.key});

  @override
  State<InventoryPage> createState() => _InventoryPageState();
}

class _InventoryPageState extends State<InventoryPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
    });
  }

  Future<void> _loadData() async {
    await context.read<ProductProvider>().loadProducts();
    await context.read<WarehouseProvider>().loadWarehouses();
    await context.read<InventoryProvider>().loadInventoryFlows();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('存库'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.warehouse),
            onPressed: () => Navigator.push(
              context,
              MaterialPageRoute(builder: (_) => const WarehousesPage()),
            ),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _loadData,
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildOverviewCards(),
              const SizedBox(height: 24),
              _buildWarnings(),
              const SizedBox(height: 24),
              _buildRecentFlows(),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showFlowOptions(),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildOverviewCards() {
    return Consumer3<ProductProvider, WarehouseProvider, InventoryProvider>(
      builder: (context, productProvider, warehouseProvider, inventoryProvider, child) {
        final totalProducts = productProvider.products.length;
        final totalWarehouses = warehouseProvider.warehouses.length;

        // Calculate total inventory value
        double totalValue = 0;
        for (final product in productProvider.products) {
          final stock = inventoryProvider.getProductStock(product.id!);
          totalValue += stock * product.costPrice;
        }

        // Calculate low stock warnings
        final lowStockProducts = productProvider.products.where((p) {
          final stock = inventoryProvider.getProductStock(p.id!);
          return stock < 10;
        }).toList();

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '库存概览',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: _buildOverviewCard(
                    '商品总数',
                    totalProducts.toString(),
                    Icons.inventory_2,
                    Colors.blue,
                    () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => const ProductsPage()),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildOverviewCard(
                    '仓库数量',
                    totalWarehouses.toString(),
                    Icons.warehouse,
                    Colors.green,
                    () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => const WarehousesPage()),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: _buildOverviewCard(
                    '库存总值',
                    '¥${totalValue.toStringAsFixed(0)}',
                    Icons.attach_money,
                    Colors.orange,
                    null,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildOverviewCard(
                    '低库存预警',
                    lowStockProducts.length.toString(),
                    Icons.warning,
                    lowStockProducts.isEmpty ? Colors.grey : Colors.red,
                    lowStockProducts.isEmpty
                        ? null
                        : () => _showLowStockDialog(lowStockProducts),
                  ),
                ),
              ],
            ),
          ],
        );
      },
    );
  }

  Widget _buildOverviewCard(
    String title,
    String value,
    IconData icon,
    Color color,
    VoidCallback? onTap,
  ) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Icon(icon, color: color, size: 24),
                  if (onTap != null)
                    Icon(Icons.chevron_right, color: Colors.grey, size: 20),
                ],
              ),
              const SizedBox(height: 12),
              Text(
                value,
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: color,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                title,
                style: const TextStyle(
                  color: Colors.grey,
                  fontSize: 12,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildWarnings() {
    return Consumer2<ProductProvider, InventoryProvider>(
      builder: (context, productProvider, inventoryProvider, child) {
        final lowStockProducts = productProvider.products.where((p) {
          final stock = inventoryProvider.getProductStock(p.id!);
          return stock < 10;
        }).toList();

        if (lowStockProducts.isEmpty) {
          return const SizedBox.shrink();
        }

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.warning, color: Colors.orange, size: 20),
                const SizedBox(width: 8),
                const Text(
                  '库存预警',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Card(
              color: Colors.orange.withAlpha(30),
              child: ListView.separated(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: lowStockProducts.take(5).length,
                separatorBuilder: (_, __) => const Divider(height: 1),
                itemBuilder: (context, index) {
                  final product = lowStockProducts[index];
                  final stock = inventoryProvider.getProductStock(product.id!);
                  return ListTile(
                    leading: const CircleAvatar(
                      backgroundColor: Colors.orange,
                      child: Icon(Icons.inventory_2, color: Colors.white, size: 20),
                    ),
                    title: Text(product.name),
                    subtitle: Text('SKU: ${product.sku ?? "N/A"}'),
                    trailing: Text(
                      '库存: $stock ${product.unit}',
                      style: const TextStyle(
                        color: Colors.red,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  );
                },
              ),
            ),
          ],
        );
      },
    );
  }

  Widget _buildRecentFlows() {
    return Consumer<InventoryProvider>(
      builder: (context, provider, child) {
        final recentFlows = provider.inventoryFlows.take(10).toList();
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  '最近出入库',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                TextButton(
                  onPressed: () {},
                  child: const Text('查看全部'),
                ),
              ],
            ),
            const SizedBox(height: 12),
            if (recentFlows.isEmpty)
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(32),
                  child: Center(
                    child: Column(
                      children: [
                        Icon(
                          Icons.history,
                          size: 48,
                          color: Colors.grey.withAlpha(100),
                        ),
                        const SizedBox(height: 12),
                        const Text(
                          '暂无出入库记录',
                          style: TextStyle(color: Colors.grey),
                        ),
                      ],
                    ),
                  ),
                ),
              )
            else
              Card(
                child: ListView.separated(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  itemCount: recentFlows.length,
                  separatorBuilder: (_, __) => const Divider(height: 1),
                  itemBuilder: (context, index) {
                    final flow = recentFlows[index];
                    final product = context
                        .read<ProductProvider>()
                        .products
                        .where((p) => p.id == flow.productId)
                        .firstOrNull;
                    final warehouse = context
                        .read<WarehouseProvider>()
                        .warehouses
                        .where((w) => w.id == flow.warehouseId)
                        .firstOrNull;
                    final isIn = flow.flowType == 'in';

                    return ListTile(
                      leading: CircleAvatar(
                        backgroundColor: isIn ? Colors.green.withAlpha(50) : Colors.red.withAlpha(50),
                        child: Icon(
                          isIn ? Icons.arrow_downward : Icons.arrow_upward,
                          color: isIn ? Colors.green : Colors.red,
                          size: 20,
                        ),
                      ),
                      title: Text(product?.name ?? 'Unknown'),
                      subtitle: Text(warehouse?.name ?? 'Unknown'),
                      trailing: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.end,
                        children: [
                          Text(
                            '${isIn ? '+' : '-'}${flow.quantity}',
                            style: TextStyle(
                              fontWeight: FontWeight.bold,
                              color: isIn ? Colors.green : Colors.red,
                            ),
                          ),
                          Text(
                            _formatDate(DateTime.fromMillisecondsSinceEpoch(flow.date)),
                            style: const TextStyle(
                              fontSize: 10,
                              color: Colors.grey,
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ),
          ],
        );
      },
    );
  }

  void _showFlowOptions() {
    showModalBottomSheet(
      context: context,
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const CircleAvatar(
                backgroundColor: Colors.green,
                child: Icon(Icons.arrow_downward, color: Colors.white),
              ),
              title: const Text('入库'),
              subtitle: const Text('商品入库操作'),
              onTap: () {
                Navigator.pop(context);
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (_) => const StockInPage()),
                );
              },
            ),
            ListTile(
              leading: const CircleAvatar(
                backgroundColor: Colors.red,
                child: Icon(Icons.arrow_upward, color: Colors.white),
              ),
              title: const Text('出库'),
              subtitle: const Text('商品出库操作'),
              onTap: () {
                Navigator.pop(context);
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (_) => const StockOutPage()),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  void _showLowStockDialog(List products) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Row(
          children: [
            Icon(Icons.warning, color: Colors.orange),
            SizedBox(width: 8),
            Text('低库存预警'),
          ],
        ),
        content: SizedBox(
          width: double.maxFinite,
          child: ListView.builder(
            shrinkWrap: true,
            itemCount: products.length,
            itemBuilder: (context, index) {
              final product = products[index];
              final stock = context
                  .read<InventoryProvider>()
                  .getProductStock(product.id!);
              return ListTile(
                title: Text(product.name),
                subtitle: Text('SKU: ${product.sku ?? "N/A"}'),
                trailing: Text(
                  '$stock ${product.unit}',
                  style: const TextStyle(
                    color: Colors.red,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              );
            },
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('关闭'),
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.month}/${date.day}';
  }
}