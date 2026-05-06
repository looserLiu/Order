import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/warehouse_provider.dart';
import '../../providers/inventory_provider.dart';
import '../../../data/models/warehouse.dart';

class WarehousesPage extends StatefulWidget {
  const WarehousesPage({super.key});

  @override
  State<WarehousesPage> createState() => _WarehousesPageState();
}

class _WarehousesPageState extends State<WarehousesPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<WarehouseProvider>().loadWarehouses();
      context.read<InventoryProvider>().loadInventoryFlows();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('仓库管理'),
        centerTitle: true,
      ),
      body: Consumer<WarehouseProvider>(
        builder: (context, provider, child) {
          if (provider.warehouses.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.warehouse,
                    size: 64,
                    color: Colors.grey.withAlpha(100),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    '暂无仓库',
                    style: TextStyle(color: Colors.grey),
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton.icon(
                    onPressed: () => _showWarehouseDialog(),
                    icon: const Icon(Icons.add),
                    label: const Text('添加仓库'),
                  ),
                ],
              ),
            );
          }

          return ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: provider.warehouses.length,
            itemBuilder: (context, index) {
              final warehouse = provider.warehouses[index];
              return _buildWarehouseCard(warehouse);
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showWarehouseDialog(),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildWarehouseCard(Warehouse warehouse) {
    final inventoryProvider = context.read<InventoryProvider>();
    final productCount = inventoryProvider.getProductsInWarehouse(warehouse.id!);
    final totalStock = inventoryProvider.getWarehouseStock(warehouse.id!);

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                CircleAvatar(
                  backgroundColor: Colors.green.withAlpha(50),
                  child: const Icon(Icons.warehouse, color: Colors.green),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        warehouse.name,
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                      if (warehouse.location != null)
                        Text(
                          warehouse.location!,
                          style: const TextStyle(
                            color: Colors.grey,
                            fontSize: 12,
                          ),
                        ),
                    ],
                  ),
                ),
                PopupMenuButton<String>(
                  onSelected: (value) {
                    if (value == 'edit') {
                      _showWarehouseDialog(warehouse);
                    } else if (value == 'delete') {
                      _showDeleteConfirmation(warehouse);
                    } else if (value == 'toggle') {
                      _toggleWarehouseStatus(warehouse);
                    }
                  },
                  itemBuilder: (context) => [
                    const PopupMenuItem(
                      value: 'edit',
                      child: Row(
                        children: [
                          Icon(Icons.edit),
                          SizedBox(width: 8),
                          Text('编辑'),
                        ],
                      ),
                    ),
                    PopupMenuItem(
                      value: 'toggle',
                      child: Row(
                        children: [
                          Icon(warehouse.isActive == 1
                              ? Icons.pause
                              : Icons.play_arrow),
                          const SizedBox(width: 8),
                          Text(warehouse.isActive == 1 ? '停用' : '启用'),
                        ],
                      ),
                    ),
                    const PopupMenuItem(
                      value: 'delete',
                      child: Row(
                        children: [
                          Icon(Icons.delete, color: Colors.red),
                          SizedBox(width: 8),
                          Text('删除', style: TextStyle(color: Colors.red)),
                        ],
                      ),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _buildStatItem(
                    '商品种类',
                    productCount.toString(),
                    Icons.category,
                  ),
                ),
                Expanded(
                  child: _buildStatItem(
                    '总库存',
                    totalStock.toString(),
                    Icons.inventory_2,
                  ),
                ),
                Expanded(
                  child: _buildStatItem(
                    '状态',
                    warehouse.isActive == 1 ? '活跃' : '已停用',
                    warehouse.isActive == 1 ? Icons.check_circle : Icons.cancel,
                  ),
                ),
              ],
            ),
            if (warehouse.description != null && warehouse.description!.isNotEmpty) ...[
              const SizedBox(height: 12),
              const Divider(),
              const SizedBox(height: 8),
              Text(
                warehouse.description!,
                style: const TextStyle(
                  color: Colors.grey,
                  fontSize: 12,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildStatItem(String label, String value, IconData icon) {
    return Column(
      children: [
        Icon(icon, color: Colors.grey, size: 20),
        const SizedBox(height: 4),
        Text(
          value,
          style: const TextStyle(
            fontWeight: FontWeight.bold,
            fontSize: 16,
          ),
        ),
        Text(
          label,
          style: const TextStyle(
            color: Colors.grey,
            fontSize: 10,
          ),
        ),
      ],
    );
  }

  void _showWarehouseDialog([Warehouse? warehouse]) {
    final isEditing = warehouse != null;
    final nameController = TextEditingController(text: warehouse?.name ?? '');
    final locationController = TextEditingController(text: warehouse?.location ?? '');
    final descriptionController = TextEditingController(text: warehouse?.description ?? '');

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(isEditing ? '编辑仓库' : '添加仓库'),
        content: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(
                  labelText: '仓库名称 *',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: locationController,
                decoration: const InputDecoration(
                  labelText: '位置',
                  border: OutlineInputBorder(),
                  hintText: '如: 北京市朝阳区',
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: descriptionController,
                maxLines: 3,
                decoration: const InputDecoration(
                  labelText: '描述',
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              final name = nameController.text.trim();
              if (name.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('请输入仓库名称')),
                );
                return;
              }

              if (isEditing) {
                await context.read<WarehouseProvider>().updateWarehouse(
                      Warehouse(
                        id: warehouse.id,
                        name: name,
                        location: locationController.text.trim().isEmpty
                            ? null
                            : locationController.text.trim(),
                        description: descriptionController.text.trim().isEmpty
                            ? null
                            : descriptionController.text.trim(),
                        isActive: warehouse.isActive,
                        createdAt: warehouse.createdAt,
                      ),
                    );
              } else {
                await context.read<WarehouseProvider>().addWarehouse(
                      Warehouse(
                        name: name,
                        location: locationController.text.trim().isEmpty
                            ? null
                            : locationController.text.trim(),
                        description: descriptionController.text.trim().isEmpty
                            ? null
                            : descriptionController.text.trim(),
                        isActive: 1,
                        createdAt: DateTime.now().millisecondsSinceEpoch,
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

  void _showDeleteConfirmation(Warehouse warehouse) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除仓库'),
        content: Text('确定要删除仓库 "${warehouse.name}" 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              await context.read<WarehouseProvider>().deleteWarehouse(warehouse.id!);
              if (context.mounted) {
                Navigator.pop(context);
              }
            },
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('删除'),
          ),
        ],
      ),
    );
  }

  void _toggleWarehouseStatus(Warehouse warehouse) async {
    final newStatus = warehouse.isActive == 1 ? 0 : 1;
    await context.read<WarehouseProvider>().updateWarehouse(
          Warehouse(
            id: warehouse.id,
            name: warehouse.name,
            location: warehouse.location,
            description: warehouse.description,
            isActive: newStatus,
            createdAt: warehouse.createdAt,
          ),
        );
  }
}