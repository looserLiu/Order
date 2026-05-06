import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/product_provider.dart';
import '../../providers/inventory_provider.dart';
import '../../../data/models/product.dart';

class ProductsPage extends StatefulWidget {
  const ProductsPage({super.key});

  @override
  State<ProductsPage> createState() => _ProductsPageState();
}

class _ProductsPageState extends State<ProductsPage> {
  String _searchQuery = '';
  String _sortBy = 'name';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ProductProvider>().loadProducts();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('商品管理'),
        centerTitle: true,
        actions: [
          PopupMenuButton<String>(
            icon: const Icon(Icons.sort),
            onSelected: (value) => setState(() => _sortBy = value),
            itemBuilder: (context) => [
              const PopupMenuItem(value: 'name', child: Text('按名称')),
              const PopupMenuItem(value: 'price', child: Text('按价格')),
              const PopupMenuItem(value: 'stock', child: Text('按库存')),
            ],
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              decoration: InputDecoration(
                hintText: '搜索商品...',
                prefixIcon: const Icon(Icons.search),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                filled: true,
              ),
              onChanged: (value) => setState(() => _searchQuery = value),
            ),
          ),
          Expanded(
            child: Consumer<ProductProvider>(
              builder: (context, provider, child) {
                var products = provider.products;

                // Filter
                if (_searchQuery.isNotEmpty) {
                  products = products.where((p) {
                    return p.name.toLowerCase().contains(_searchQuery.toLowerCase()) ||
                        (p.sku?.toLowerCase().contains(_searchQuery.toLowerCase()) ?? false);
                  }).toList();
                }

                // Sort
                switch (_sortBy) {
                  case 'price':
                    products.sort((a, b) => a.salePrice.compareTo(b.salePrice));
                    break;
                  case 'stock':
                    final inventoryProvider = context.read<InventoryProvider>();
                    products.sort((a, b) {
                      final stockA = inventoryProvider.getProductStock(a.id!);
                      final stockB = inventoryProvider.getProductStock(b.id!);
                      return stockA.compareTo(stockB);
                    });
                    break;
                  default:
                    products.sort((a, b) => a.name.compareTo(b.name));
                }

                if (products.isEmpty) {
                  return Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.inventory_2,
                          size: 64,
                          color: Colors.grey.withAlpha(100),
                        ),
                        const SizedBox(height: 16),
                        Text(
                          _searchQuery.isEmpty ? '暂无商品' : '未找到商品',
                          style: const TextStyle(color: Colors.grey),
                        ),
                      ],
                    ),
                  );
                }

                return ListView.builder(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  itemCount: products.length,
                  itemBuilder: (context, index) {
                    final product = products[index];
                    return _buildProductCard(product);
                  },
                );
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showProductDialog(),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildProductCard(Product product) {
    final inventoryProvider = context.read<InventoryProvider>();
    final stock = inventoryProvider.getProductStock(product.id!);
    final isLowStock = stock < 10;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () => _showProductDialog(product),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                width: 60,
                height: 60,
                decoration: BoxDecoration(
                  color: Colors.grey.withAlpha(30),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: product.imageUrl != null
                    ? ClipRRect(
                        borderRadius: BorderRadius.circular(8),
                        child: Image.network(
                          product.imageUrl!,
                          fit: BoxFit.cover,
                          errorBuilder: (_, __, ___) => const Icon(
                            Icons.inventory_2,
                            color: Colors.grey,
                          ),
                        ),
                      )
                    : const Icon(
                        Icons.inventory_2,
                        color: Colors.grey,
                      ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      product.name,
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'SKU: ${product.sku ?? "N/A"}',
                      style: const TextStyle(
                        color: Colors.grey,
                        fontSize: 12,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        Text(
                          '¥${product.salePrice.toStringAsFixed(2)}',
                          style: const TextStyle(
                            color: Colors.green,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        if (product.category != null) ...[
                          const SizedBox(width: 8),
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.blue.withAlpha(30),
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              product.category!,
                              style: const TextStyle(
                                fontSize: 10,
                                color: Colors.blue,
                              ),
                            ),
                          ),
                        ],
                      ],
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '$stock ${product.unit ?? "件"}',
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      color: isLowStock ? Colors.red : null,
                    ),
                  ),
                  if (isLowStock)
                    const Icon(
                      Icons.warning,
                      color: Colors.orange,
                      size: 16,
                    ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showProductDialog([Product? product]) {
    final isEditing = product != null;
    final nameController = TextEditingController(text: product?.name ?? '');
    final skuController = TextEditingController(text: product?.sku ?? '');
    final categoryController = TextEditingController(text: product?.category ?? '');
    final unitController = TextEditingController(text: product?.unit ?? '件');
    final costPriceController = TextEditingController(
      text: product?.costPrice.toString() ?? '0',
    );
    final salePriceController = TextEditingController(
      text: product?.salePrice.toString() ?? '0',
    );

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(isEditing ? '编辑商品' : '添加商品'),
        content: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(
                  labelText: '商品名称 *',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: skuController,
                decoration: const InputDecoration(
                  labelText: 'SKU编码',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: categoryController,
                decoration: const InputDecoration(
                  labelText: '商品分类',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: unitController,
                decoration: const InputDecoration(
                  labelText: '单位',
                  border: OutlineInputBorder(),
                  hintText: '如: 件、箱、个',
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: costPriceController,
                      keyboardType:
                          const TextInputType.numberWithOptions(decimal: true),
                      decoration: const InputDecoration(
                        labelText: '成本价',
                        prefixText: '¥ ',
                        border: OutlineInputBorder(),
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: TextField(
                      controller: salePriceController,
                      keyboardType:
                          const TextInputType.numberWithOptions(decimal: true),
                      decoration: const InputDecoration(
                        labelText: '售价',
                        prefixText: '¥ ',
                        border: OutlineInputBorder(),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
        actions: [
          if (isEditing)
            TextButton(
              onPressed: () async {
                final confirm = await showDialog<bool>(
                  context: context,
                  builder: (context) => AlertDialog(
                    title: const Text('删除商品'),
                    content: Text('确定要删除 "${product.name}" 吗？'),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.pop(context, false),
                        child: const Text('取消'),
                      ),
                      ElevatedButton(
                        onPressed: () => Navigator.pop(context, true),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.red,
                        ),
                        child: const Text('删除'),
                      ),
                    ],
                  ),
                );
                if (confirm == true && context.mounted) {
                  await context.read<ProductProvider>().deleteProduct(product.id!);
                  if (context.mounted) Navigator.pop(context);
                }
              },
              child: const Text('删除', style: TextStyle(color: Colors.red)),
            ),
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              final name = nameController.text.trim();
              if (name.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('请输入商品名称')),
                );
                return;
              }

              final costPrice = double.tryParse(costPriceController.text) ?? 0;
              final salePrice = double.tryParse(salePriceController.text) ?? 0;

              if (isEditing) {
                await context.read<ProductProvider>().updateProduct(
                      Product(
                        id: product.id,
                        name: name,
                        sku: skuController.text.trim().isEmpty
                            ? null
                            : skuController.text.trim(),
                        category: categoryController.text.trim().isEmpty
                            ? null
                            : categoryController.text.trim(),
                        unit: unitController.text.trim().isEmpty
                            ? '件'
                            : unitController.text.trim(),
                        costPrice: costPrice,
                        salePrice: salePrice,
                        imageUrl: product.imageUrl,
                        createdAt: product.createdAt,
                        updatedAt: DateTime.now().millisecondsSinceEpoch,
                      ),
                    );
              } else {
                await context.read<ProductProvider>().addProduct(
                      Product(
                        name: name,
                        sku: skuController.text.trim().isEmpty
                            ? null
                            : skuController.text.trim(),
                        category: categoryController.text.trim().isEmpty
                            ? null
                            : categoryController.text.trim(),
                        unit: unitController.text.trim().isEmpty
                            ? '件'
                            : unitController.text.trim(),
                        costPrice: costPrice,
                        salePrice: salePrice,
                        createdAt: DateTime.now().millisecondsSinceEpoch,
                        updatedAt: DateTime.now().millisecondsSinceEpoch,
                      ),
                    );
              }

              if (context.mounted) {
                Navigator.pop(context);
              }
            },
            child: Text(isEditing ? '保存' : '添加'),
          ),
        ],
      ),
    );
  }
}