import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import '../../providers/product_provider.dart';
import '../../providers/warehouse_provider.dart';
import '../../providers/inventory_provider.dart';
import '../../../data/models/inventory_flow.dart';

class StockOutPage extends StatefulWidget {
  const StockOutPage({super.key});

  @override
  State<StockOutPage> createState() => _StockOutPageState();
}

class _StockOutPageState extends State<StockOutPage> {
  final _formKey = GlobalKey<FormState>();
  final _quantityController = TextEditingController();
  final _referenceController = TextEditingController();

  int? _selectedProductId;
  int? _selectedWarehouseId;
  DateTime _selectedDate = DateTime.now();
  double _availableStock = 0;

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
    if (mounted) {
      setState(() {
        if (context.read<ProductProvider>().products.isNotEmpty) {
          _selectedProductId = context.read<ProductProvider>().products.first.id;
          _updateAvailableStock();
        }
        if (context.read<WarehouseProvider>().warehouses.isNotEmpty) {
          _selectedWarehouseId = context.read<WarehouseProvider>().warehouses.first.id;
        }
      });
    }
  }

  void _updateAvailableStock() {
    if (_selectedProductId != null && _selectedWarehouseId != null) {
      final inventoryProvider = context.read<InventoryProvider>();
      _availableStock = inventoryProvider.getProductStockInWarehouse(
        _selectedProductId!,
        _selectedWarehouseId!,
      );
    }
  }

  @override
  void dispose() {
    _quantityController.dispose();
    _referenceController.dispose();
    super.dispose();
  }

  Future<void> _selectDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2000),
      lastDate: DateTime.now(),
    );
    if (picked != null) {
      setState(() => _selectedDate = picked);
    }
  }

  Future<void> _saveStockOut() async {
    if (!_formKey.currentState!.validate()) return;

    if (_selectedProductId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请选择商品')),
      );
      return;
    }

    if (_selectedWarehouseId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请选择仓库')),
      );
      return;
    }

    final quantity = double.tryParse(_quantityController.text) ?? 0;

    if (quantity > _availableStock) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('库存不足，当前库存: $_availableStock')),
      );
      return;
    }

    final flow = InventoryFlow(
      productId: _selectedProductId!,
      warehouseId: _selectedWarehouseId!,
      flowType: 'out',
      quantity: quantity,
      referenceId: _referenceController.text.trim().isEmpty
          ? null
          : _referenceController.text.trim(),
      date: _selectedDate.millisecondsSinceEpoch,
      createdAt: DateTime.now().millisecondsSinceEpoch,
    );

    await context.read<InventoryProvider>().addInventoryFlow(flow);

    if (mounted) {
      Navigator.pop(context);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('出库成功')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('出库'),
        centerTitle: true,
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildProductSelector(),
              const SizedBox(height: 16),
              _buildWarehouseSelector(),
              const SizedBox(height: 16),
              _buildQuantityInput(),
              const SizedBox(height: 16),
              _buildDateSelector(),
              const SizedBox(height: 16),
              _buildReferenceInput(),
              const SizedBox(height: 32),
              _buildSaveButton(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProductSelector() {
    return Consumer2<ProductProvider, InventoryProvider>(
      builder: (context, productProvider, inventoryProvider, child) {
        return Card(
          child: ListTile(
            leading: const Icon(Icons.inventory_2),
            title: const Text('商品'),
            subtitle: Text(
              productProvider.products
                      .where((p) => p.id == _selectedProductId)
                      .firstOrNull
                      ?.name ??
                  '请选择商品',
            ),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showProductPicker(productProvider.products, inventoryProvider),
          ),
        );
      },
    );
  }

  void _showProductPicker(List products, InventoryProvider inventoryProvider) {
    showModalBottomSheet(
      context: context,
      builder: (ctx) => ListView.builder(
        shrinkWrap: true,
        itemCount: products.length,
        itemBuilder: (context, index) {
          final product = products[index];
          final totalStock = inventoryProvider.getProductStock(product.id!);
          return ListTile(
            leading: const Icon(Icons.inventory_2, color: Colors.red),
            title: Text(product.name),
            subtitle: Text('SKU: ${product.sku ?? "N/A"} | 总库存: $totalStock ${product.unit}'),
            trailing: _selectedProductId == product.id
                ? const Icon(Icons.check, color: Colors.green)
                : null,
            onTap: () {
              setState(() {
                _selectedProductId = product.id;
              });
              _updateAvailableStock();
              Navigator.pop(context);
            },
          );
        },
      ),
    );
  }

  Widget _buildWarehouseSelector() {
    return Consumer<WarehouseProvider>(
      builder: (context, provider, child) {
        return Card(
          child: ListTile(
            leading: const Icon(Icons.warehouse),
            title: const Text('仓库'),
            subtitle: Text(
              provider.warehouses
                      .where((w) => w.id == _selectedWarehouseId)
                      .firstOrNull
                      ?.name ??
                  '请选择仓库',
            ),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showWarehousePicker(provider.warehouses),
          ),
        );
      },
    );
  }

  void _showWarehousePicker(List warehouses) {
    showModalBottomSheet(
      context: context,
      builder: (context) => ListView.builder(
        shrinkWrap: true,
        itemCount: warehouses.length,
        itemBuilder: (context, index) {
          final warehouse = warehouses[index];
          final inventoryProvider = context.read<InventoryProvider>();
          final stock = _selectedProductId != null
              ? inventoryProvider.getProductStockInWarehouse(
                  _selectedProductId!,
                  warehouse.id!,
                )
              : 0;
          return ListTile(
            leading: const Icon(Icons.warehouse, color: Colors.orange),
            title: Text(warehouse.name),
            subtitle: Text('${warehouse.location ?? ""} | 库存: $stock'),
            trailing: _selectedWarehouseId == warehouse.id
                ? const Icon(Icons.check, color: Colors.green)
                : null,
            onTap: () {
              setState(() => _selectedWarehouseId = warehouse.id);
              _updateAvailableStock();
              Navigator.pop(context);
            },
          );
        },
      ),
    );
  }

  Widget _buildQuantityInput() {
    final product = context
        .read<ProductProvider>()
        .products
        .where((p) => p.id == _selectedProductId)
        .firstOrNull;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  '出库数量',
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.grey,
                  ),
                ),
                Text(
                  '可用库存: $_availableStock ${product?.unit ?? "件"}',
                  style: const TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: _quantityController,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              style: const TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
              decoration: const InputDecoration(
                border: InputBorder.none,
                hintText: '0',
              ),
              onChanged: (_) => setState(() {}),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return '请输入出库数量';
                }
                final qty = double.tryParse(value);
                if (qty == null || qty <= 0) {
                  return '请输入有效数量';
                }
                if (qty > _availableStock) {
                  return '库存不足';
                }
                return null;
              },
            ),
            if (_availableStock < 10)
              Padding(
                padding: const EdgeInsets.only(top: 8),
                child: Row(
                  children: [
                    const Icon(Icons.warning, color: Colors.orange, size: 16),
                    const SizedBox(width: 4),
                    Text(
                      '库存不足，请及时补货',
                      style: TextStyle(color: Colors.orange.shade700, fontSize: 12),
                    ),
                  ],
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildDateSelector() {
    return Card(
      child: ListTile(
        leading: const Icon(Icons.calendar_today),
        title: const Text('出库日期'),
        subtitle: Text(DateFormat('yyyy年MM月dd日').format(_selectedDate)),
        trailing: const Icon(Icons.chevron_right),
        onTap: _selectDate,
      ),
    );
  }

  Widget _buildReferenceInput() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '参考单号',
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey,
              ),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: _referenceController,
              decoration: const InputDecoration(
                border: InputBorder.none,
                hintText: '可选',
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSaveButton() {
    final isValid = _availableStock > 0;
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: isValid ? _saveStockOut : null,
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.symmetric(vertical: 16),
          backgroundColor: Colors.red,
          disabledBackgroundColor: Colors.grey.shade300,
        ),
        child: const Text(
          '确认出库',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
      ),
    );
  }
}