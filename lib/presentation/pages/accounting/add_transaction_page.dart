import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import '../../providers/transaction_provider.dart';
import '../../providers/account_provider.dart';
import '../../providers/category_provider.dart';
import '../../../data/models/transaction.dart';
import '../../../core/services/smart_categorization.dart';

class AddTransactionPage extends StatefulWidget {
  final Transaction? transaction;

  const AddTransactionPage({super.key, this.transaction});

  @override
  State<AddTransactionPage> createState() => _AddTransactionPageState();
}

class _AddTransactionPageState extends State<AddTransactionPage> {
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();
  final _descriptionController = TextEditingController();

  String _type = 'expense';
  int? _selectedAccountId;
  int? _selectedCategoryId;
  DateTime _selectedDate = DateTime.now();
  List<CategorySuggestion> _suggestedCategories = [];
  final SmartCategorization _smartCategorization = SmartCategorization();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
    });
  }

  Future<void> _loadData() async {
    await context.read<AccountProvider>().loadAccounts();
    await context.read<CategoryProvider>().loadCategories();
    await context.read<TransactionProvider>().loadTransactions();
    if (mounted) {
      setState(() {
        if (context.read<AccountProvider>().accounts.isNotEmpty) {
          _selectedAccountId = context.read<AccountProvider>().accounts.first.id;
        }
      });
    }
  }

  void _updateSuggestions() {
    final desc = _descriptionController.text;
    final amount = double.tryParse(_amountController.text);
    final date = _selectedDate.millisecondsSinceEpoch;

    final categories = context.read<CategoryProvider>().categories
        .where((c) => c.type == _type)
        .toList();
    final transactions = context.read<TransactionProvider>().transactions;

    final suggestions = _smartCategorization.suggestCategories(
      description: desc,
      categories: categories,
      transactions: transactions,
      amount: amount,
      date: date,
    );

    setState(() => _suggestedCategories = suggestions);
  }

  @override
  void dispose() {
    _amountController.dispose();
    _descriptionController.dispose();
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

  Future<void> _saveTransaction() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedAccountId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请选择账户')),
      );
      return;
    }
    if (_selectedCategoryId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请选择分类')),
      );
      return;
    }

    final amount = double.tryParse(_amountController.text) ?? 0;
    final provider = context.read<TransactionProvider>();

    final transaction = Transaction(
      accountId: _selectedAccountId!,
      categoryId: _selectedCategoryId,
      amount: amount,
      type: _type,
      description: _descriptionController.text.trim(),
      date: _selectedDate.millisecondsSinceEpoch,
      createdAt: DateTime.now().millisecondsSinceEpoch,
    );

    await provider.addTransaction(transaction);

    // Update category usage count
    await context.read<CategoryProvider>().incrementUsageCount(_selectedCategoryId!);

    // Update account balance
    final accountProvider = context.read<AccountProvider>();
    final account = accountProvider.accounts.firstWhere((a) => a.id == _selectedAccountId);
    final newBalance = _type == 'income'
        ? account.balance + amount
        : account.balance - amount;
    await accountProvider.updateBalance(_selectedAccountId!, newBalance);

    if (mounted) {
      Navigator.pop(context);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('记录已保存')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('记一笔'),
        centerTitle: true,
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildTypeSelector(),
              const SizedBox(height: 24),
              _buildAmountInput(),
              const SizedBox(height: 24),
              _buildAccountSelector(),
              const SizedBox(height: 24),
              _buildCategorySelector(),
              const SizedBox(height: 24),
              _buildDateSelector(),
              const SizedBox(height: 24),
              _buildDescriptionInput(),
              const SizedBox(height: 32),
              _buildSaveButton(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTypeSelector() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(8),
        child: Row(
          children: [
            Expanded(
              child: _buildTypeButton('expense', '支出', Colors.red),
            ),
            Expanded(
              child: _buildTypeButton('income', '收入', Colors.green),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTypeButton(String type, String label, Color color) {
    final isSelected = _type == type;
    return GestureDetector(
      onTap: () {
        setState(() {
          _type = type;
          _selectedCategoryId = null;
        });
        _updateSuggestions();
      },
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: isSelected ? color : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Center(
          child: Text(
            label,
            style: TextStyle(
              color: isSelected ? Colors.white : color,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildAmountInput() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '金额',
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey,
              ),
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: _amountController,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              style: const TextStyle(
                fontSize: 32,
                fontWeight: FontWeight.bold,
              ),
              decoration: const InputDecoration(
                prefixText: '¥ ',
                prefixStyle: TextStyle(
                  fontSize: 32,
                  fontWeight: FontWeight.bold,
                ),
                border: InputBorder.none,
                hintText: '0.00',
              ),
              onChanged: (_) => _updateSuggestions(),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return '请输入金额';
                }
                final amount = double.tryParse(value);
                if (amount == null || amount <= 0) {
                  return '请输入有效金额';
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAccountSelector() {
    return Consumer<AccountProvider>(
      builder: (context, provider, child) {
        return Card(
          child: ListTile(
            leading: const Icon(Icons.account_balance_wallet),
            title: const Text('账户'),
            subtitle: Text(
              provider.accounts
                      .where((a) => a.id == _selectedAccountId)
                      .firstOrNull
                      ?.name ??
                  '请选择账户',
            ),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showAccountPicker(provider.accounts),
          ),
        );
      },
    );
  }

  void _showAccountPicker(List accounts) {
    showModalBottomSheet(
      context: context,
      builder: (context) => ListView.builder(
        shrinkWrap: true,
        itemCount: accounts.length,
        itemBuilder: (context, index) {
          final account = accounts[index];
          return ListTile(
            leading: Icon(
              _getAccountIcon(account.type),
              color: Color(account.color ?? 0xFF2196F3),
            ),
            title: Text(account.name),
            subtitle: Text('¥${account.balance.toStringAsFixed(2)}'),
            trailing: _selectedAccountId == account.id
                ? const Icon(Icons.check, color: Colors.green)
                : null,
            onTap: () {
              setState(() => _selectedAccountId = account.id);
              Navigator.pop(context);
            },
          );
        },
      ),
    );
  }

  Widget _buildCategorySelector() {
    return Consumer<CategoryProvider>(
      builder: (context, provider, child) {
        final categories = provider.categories.where((c) => c.type == _type).toList();
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '分类',
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.grey,
                  ),
                ),
                const SizedBox(height: 12),
                // Smart categorization suggestions
                if (_suggestedCategories.isNotEmpty && _selectedCategoryId == null) ...[
                  const Text(
                    '✨ 智能推荐',
                    style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 8),
                  Wrap(
                    spacing: 8,
                    children: _suggestedCategories.map((suggestion) {
                      return ActionChip(
                        avatar: Icon(
                          _getCategoryIcon(suggestion.category.icon ?? 'category'),
                          size: 16,
                          color: Color(suggestion.category.color ?? 0xFF9E9E9E),
                        ),
                        label: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(suggestion.category.name),
                            Text(
                              suggestion.reason,
                              style: const TextStyle(fontSize: 9, color: Colors.grey),
                            ),
                          ],
                        ),
                        onPressed: () => setState(() => _selectedCategoryId = suggestion.category.id),
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 12),
                ],
                // All categories
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: categories.map((cat) {
                    final isSelected = _selectedCategoryId == cat.id;
                    return GestureDetector(
                      onTap: () => setState(() => _selectedCategoryId = cat.id),
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                        decoration: BoxDecoration(
                          color: isSelected
                              ? Color(cat.color ?? 0xFF9E9E9E)
                              : Colors.grey.withAlpha(30),
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(
                              _getCategoryIcon(cat.icon ?? 'category'),
                              size: 16,
                              color: isSelected ? Colors.white : Color(cat.color ?? 0xFF9E9E9E),
                            ),
                            const SizedBox(width: 4),
                            Text(
                              cat.name,
                              style: TextStyle(
                                color: isSelected ? Colors.white : null,
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  }).toList(),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildDateSelector() {
    return Card(
      child: ListTile(
        leading: const Icon(Icons.calendar_today),
        title: const Text('日期'),
        subtitle: Text(DateFormat('yyyy年MM月dd日').format(_selectedDate)),
        trailing: const Icon(Icons.chevron_right),
        onTap: _selectDate,
      ),
    );
  }

  Widget _buildDescriptionInput() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '备注',
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey,
              ),
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: _descriptionController,
              maxLines: 3,
              decoration: const InputDecoration(
                hintText: '添加备注...',
                border: InputBorder.none,
              ),
              onChanged: (_) => _updateSuggestions(),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSaveButton() {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: _saveTransaction,
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.symmetric(vertical: 16),
          backgroundColor: _type == 'income' ? Colors.green : Colors.red,
        ),
        child: const Text(
          '保存',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
      ),
    );
  }

  IconData _getAccountIcon(String type) {
    switch (type) {
      case 'cash':
        return Icons.wallet;
      case 'bank':
        return Icons.account_balance;
      case 'credit_card':
        return Icons.credit_card;
      case 'digital':
        return Icons.phone_android;
      default:
        return Icons.account_balance_wallet;
    }
  }

  IconData _getCategoryIcon(String icon) {
    switch (icon) {
      case 'restaurant':
        return Icons.restaurant;
      case 'directions_car':
        return Icons.directions_car;
      case 'shopping_bag':
        return Icons.shopping_bag;
      case 'movie':
        return Icons.movie;
      case 'bolt':
        return Icons.bolt;
      case 'home':
        return Icons.home;
      case 'local_hospital':
        return Icons.local_hospital;
      case 'school':
        return Icons.school;
      case 'work':
        return Icons.work;
      case 'card_giftcard':
        return Icons.card_giftcard;
      case 'trending_up':
        return Icons.trending_up;
      case 'laptop':
        return Icons.laptop;
      default:
        return Icons.category;
    }
  }
}
