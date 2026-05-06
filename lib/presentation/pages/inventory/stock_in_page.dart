import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import '../../providers/product_provider.dart';
import '../../providers/warehouse_provider.dart';
import '../../providers/inventory_provider.dart';
import '../../../data/models/inventory_flow.dart';

class StockInPage extends StatefulWidget {
  const StockInPage({super.key});

  @override
  State<StockInPage> createState() => _StockInPageState();
}

class _StockInPageState extends State<StockInPage> {
  final _formKey = GlobalKey<FormState>();
  final _quantityController = TextEditingController();
  final _batchController = TextEditingController();
  final _costController = TextEditingController();
  final _referenceController = TextEditingController();

  int? _selectedProductId;
  int? _selectedWarehouseId;
  DateTime _selectedDate = DateTime.now();
  DateTime? _expirationDate;

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
        }
        if (context.read<WarehouseProvider>().warehouses.isNotEmpty) {
          _selectedWarehouseId = context.read<WarehouseProvider>().warehouses.first.id;
        }
      });
    }
  }

  @override
  void dispose() {
    _quantityController.dispose();
    _batchController.dispose();
    _costController.dispose();
    _referenceController.dispose();
    super.dispose();
  }

  Future<void> _selectDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2000),
      lastDate: DateTime.now().add(const Duration(days: 365)),
    );
    if (picked != null) {
      setState(() => _selectedDate = picked);
    }
  }

  Future<void> _selectExpirationDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _expirationDate ?? DateTime.now().add(const Duration(days: 30)),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 3650)),
    );
    if (picked != null) {
      setState(() => _expirationDate = picked);
    }
  }

  Future<void> _saveStockIn() async {
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
    final cost = double.tryParse(_costController.text) ?? 0;

    final flow = InventoryFlow(
      productId: _selectedProductId!,
      warehouseId: _selectedWarehouseId!,
      flowType: 'in',
      quantity: quantity,
      batchNumber: _batchController.text.trim().isEmpty ? null : _batchController.text.trim(),
      expirationDate: _expirationDate?.millisecondsSinceEpoch,
      costAtFlow: cost,
      referenceId: _referenceController.text.trim().isEmpty ? null : _referenceController.text.trim(),
      date: _selectedDate.millisecondsSinceEpoch,
      createdAt: DateTime.now().millisecondsSinceEpoch,
    );

    await context.read<InventoryProvider>().addInventoryFlow(flow);

    if (mounted) {
      Navigator.pop(context);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('入库成功')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('入库'),
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
              _buildOptionalFields(),
              const SizedBox(height: 32),
              _buildSaveButton(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProductSelector() {
    return Consumer<ProductProvider>(
      builder: (context, provider, child) {
        return Card(
          child: ListTile(
            leading: const Icon(Icons.inventory_2),
            title: const Text('商品'),
            subtitle: Text(
              provider.products
                      .where((p) => p.id == _selectedProductId)
                      .firstOrNull
                      ?.name ??
                  '请选择商品',
            ),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showProductPicker(provider.products),
          ),
        );
      },
    );
  }

  void _showProductPicker(List products) {
    showModalBottomSheet(
      context: context,
      builder: (context) => ListView.builder(
        shrinkWrap: true,
        itemCount: products.length,
        itemBuilder: (context, index) {
          final product = products[index];
          final inventoryProvider = context.read<InventoryProvider>();
          final currentStock = inventoryProvider.getProductStock(product.id!);
          return ListTile(
            leading: const Icon(Icons.inventory_2, color: Colors.blue),
            title: Text(product.name),
            subtitle: Text('SKU: ${product.sku ?? "N/A"} | 当前库存: $currentStock ${product.unit}'),
            trailing: _selectedProductId == product.id
                ? const Icon(Icons.check, color: Colors.green)
                : null,
            onTap: () {
              setState(() {
                _selectedProductId = product.id;
                _costController.text = product.costPrice.toString();
              });
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
          return ListTile(
            leading: const Icon(Icons.warehouse, color: Colors.green),
            title: Text(warehouse.name),
            subtitle: Text(warehouse.location ?? ''),
            trailing: _selectedWarehouseId == warehouse.id
                ? const Icon(Icons.check, color: Colors.green)
                : null,
            onTap: () {
              setState(() => _selectedWarehouseId = warehouse.id);
              Navigator.pop(context);
            },
          );
        },
      ),
    );
  }

  Widget _buildQuantityInput() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '入库数量',
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey,
              ),
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
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return '请输入入库数量';
                }
                final qty = double.tryParse(value);
                if (qty == null || qty <= 0) {
                  return '请输入有效数量';
                }
                return null;
              },
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
        title: const Text('入库日期'),
        subtitle: Text(DateFormat('yyyy年MM月dd日').format(_selectedDate)),
        trailing: const Icon(Icons.chevron_right),
        onTap: _selectDate,
      ),
    );
  }

  Widget _buildOptionalFields() {
    return ExpansionTile(
      title: const Text('可选字段'),
      tilePadding: EdgeInsets.zero,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Column(
            children: [
              TextField(
                controller: _costController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: const InputDecoration(
                  labelText: '成本单价',
                  prefixText: '¥ ',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: _batchController,
                decoration: const InputDecoration(
                  labelText: '批次号',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              Card(
                child: ListTile(
                  title: const Text('有效期'),
                  subtitle: Text(
                    _expirationDate != null
                        ? DateFormat('yyyy年MM月dd日').format(_expirationDate!)
                        : '请选择有效期',
                  ),
                  trailing: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (_expirationDate != null)
                        IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: () => setState(() => _expirationDate = null),
                        ),
                      const Icon(Icons.chevron_right),
                    ],
                  ),
                  onTap: _selectExpirationDate,
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: _referenceController,
                decoration: const InputDecoration(
                  labelText: '参考单号',
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildSaveButton() {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: _saveStockIn,
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.symmetric(vertical: 16),
          backgroundColor: Colors.green,
        ),
        child: const Text(
          '确认入库',
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