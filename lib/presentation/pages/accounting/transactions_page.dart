import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/account_provider.dart';
import '../../providers/transaction_provider.dart';
import '../../providers/category_provider.dart';
import '../../../data/models/transaction.dart' as tx;
import '../../../data/models/account.dart';
import '../../../data/models/category.dart';
import 'add_transaction_page.dart';

class TransactionsPage extends StatefulWidget {
  const TransactionsPage({super.key});

  @override
  State<TransactionsPage> createState() => _TransactionsPageState();
}

class _TransactionsPageState extends State<TransactionsPage> {
  String? _selectedAccountId;
  String? _selectedType;
  DateTime? _startDate;
  DateTime? _endDate;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<TransactionProvider>().loadTransactions();
      context.read<CategoryProvider>().loadCategories();
      context.read<AccountProvider>().loadAccounts();
    });
  }

  Future<void> _applyFilters() async {
    final provider = context.read<TransactionProvider>();
    await provider.loadTransactions(
      accountId: _selectedAccountId != null ? int.tryParse(_selectedAccountId!) : null,
      startDate: _startDate,
      endDate: _endDate,
    );
  }

  void _showFilterDialog() {
    showModalBottomSheet(
      context: context,
      builder: (context) => _buildFilterSheet(),
    );
  }

  Widget _buildFilterSheet() {
    return Container(
      padding: const EdgeInsets.all(20),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '筛选条件',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 16),
          Consumer<AccountProvider>(
            builder: (context, accountProvider, child) {
              return DropdownButtonFormField<String>(
                value: _selectedAccountId,
                decoration: const InputDecoration(
                  labelText: '账户',
                  border: OutlineInputBorder(),
                ),
                items: [
                  const DropdownMenuItem(value: null, child: Text('全部账户')),
                  ...accountProvider.accounts.map((account) {
                    return DropdownMenuItem(
                      value: account.id.toString(),
                      child: Text(account.name),
                    );
                  }),
                ],
                onChanged: (value) {
                  setState(() => _selectedAccountId = value);
                },
              );
            },
          ),
          const SizedBox(height: 16),
          DropdownButtonFormField<String>(
            value: _selectedType,
            decoration: const InputDecoration(
              labelText: '类型',
              border: OutlineInputBorder(),
            ),
            items: const [
              DropdownMenuItem(value: null, child: Text('全部')),
              DropdownMenuItem(value: 'income', child: Text('收入')),
              DropdownMenuItem(value: 'expense', child: Text('支出')),
            ],
            onChanged: (value) {
              setState(() => _selectedType = value);
            },
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: () async {
                    final date = await showDatePicker(
                      context: context,
                      initialDate: _startDate ?? DateTime.now(),
                      firstDate: DateTime(2000),
                      lastDate: DateTime(2100),
                    );
                    if (date != null) {
                      setState(() => _startDate = date);
                    }
                  },
                  child: Text(_startDate != null
                      ? '${_startDate!.year}/${_startDate!.month}/${_startDate!.day}'
                      : '开始日期'),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: OutlinedButton(
                  onPressed: () async {
                    final date = await showDatePicker(
                      context: context,
                      initialDate: _endDate ?? DateTime.now(),
                      firstDate: DateTime(2000),
                      lastDate: DateTime(2100),
                    );
                    if (date != null) {
                      setState(() => _endDate = date);
                    }
                  },
                  child: Text(_endDate != null
                      ? '${_endDate!.year}/${_endDate!.month}/${_endDate!.day}'
                      : '结束日期'),
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          Row(
            children: [
              Expanded(
                child: TextButton(
                  onPressed: () {
                    setState(() {
                      _selectedAccountId = null;
                      _selectedType = null;
                      _startDate = null;
                      _endDate = null;
                    });
                    context.read<TransactionProvider>().loadTransactions();
                    Navigator.pop(context);
                  },
                  child: const Text('重置'),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: ElevatedButton(
                  onPressed: () {
                    _applyFilters();
                    Navigator.pop(context);
                  },
                  child: const Text('应用'),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('交易记录'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.filter_list),
            onPressed: _showFilterDialog,
          ),
        ],
      ),
      body: Consumer3<TransactionProvider, AccountProvider, CategoryProvider>(
        builder: (context, txProvider, accountProvider, categoryProvider, child) {
          if (txProvider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (txProvider.transactions.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.receipt_long,
                    size: 64,
                    color: Colors.grey.withAlpha(100),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    '暂无交易记录',
                    style: TextStyle(color: Colors.grey),
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton.icon(
                    onPressed: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => const AddTransactionPage()),
                    ),
                    icon: const Icon(Icons.add),
                    label: const Text('添加交易'),
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () => txProvider.loadTransactions(),
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: txProvider.transactions.length,
              itemBuilder: (context, index) {
                final transaction = txProvider.transactions[index];
                return _buildTransactionCard(
                  transaction,
                  accountProvider.accounts,
                  categoryProvider.categories,
                );
              },
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => Navigator.push(
          context,
          MaterialPageRoute(builder: (_) => const AddTransactionPage()),
        ),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildTransactionCard(
    tx.Transaction transaction,
    List<Account> accounts,
    List<Category> categories,
  ) {
    final account = accounts.where((a) => a.id == transaction.accountId).firstOrNull;
    final category = categories.where((c) => c.id == transaction.categoryId).firstOrNull;
    final isExpense = transaction.type == 'expense';
    final date = DateTime.fromMillisecondsSinceEpoch(transaction.date);

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () => _showTransactionDetail(transaction, account, category),
        onLongPress: () => _showDeleteConfirmation(transaction),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              CircleAvatar(
                backgroundColor: Color(category?.color ?? 0xFF9E9E9E).withAlpha(50),
                child: Icon(
                  _getCategoryIcon(category?.icon ?? 'category'),
                  color: Color(category?.color ?? 0xFF9E9E9E),
                  size: 20,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      category?.name ?? '未分类',
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      transaction.description ?? '',
                      style: TextStyle(color: Colors.grey[600], fontSize: 12),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '${date.year}/${date.month}/${date.day} • ${account?.name ?? '未知账户'}',
                      style: TextStyle(color: Colors.grey[400], fontSize: 11),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '${isExpense ? '-' : '+'}¥${transaction.amount.toStringAsFixed(2)}',
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                      color: isExpense ? Colors.red : Colors.green,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showTransactionDetail(
    tx.Transaction transaction,
    Account? account,
    Category? category,
  ) {
    final date = DateTime.fromMillisecondsSinceEpoch(transaction.date);
    final isExpense = transaction.type == 'expense';

    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  '交易详情',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                IconButton(
                  icon: const Icon(Icons.close),
                  onPressed: () => Navigator.pop(context),
                ),
              ],
            ),
            const Divider(),
            const SizedBox(height: 12),
            _buildDetailRow('金额', '${isExpense ? '-' : '+'}¥${transaction.amount.toStringAsFixed(2)}',
                isExpense ? Colors.red : Colors.green),
            _buildDetailRow('类型', isExpense ? '支出' : '收入'),
            _buildDetailRow('分类', category?.name ?? '未分类'),
            _buildDetailRow('账户', account?.name ?? '未知'),
            _buildDetailRow('日期', '${date.year}/${date.month}/${date.day}'),
            if (transaction.description != null && transaction.description!.isNotEmpty)
              _buildDetailRow('描述', transaction.description!),
            const SizedBox(height: 20),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: () {
                      Navigator.pop(context);
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (_) => AddTransactionPage(transaction: transaction),
                        ),
                      );
                    },
                    icon: const Icon(Icons.edit),
                    label: const Text('编辑'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () {
                      Navigator.pop(context);
                      _showDeleteConfirmation(transaction);
                    },
                    style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
                    icon: const Icon(Icons.delete),
                    label: const Text('删除'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDetailRow(String label, String value, [Color? valueColor]) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: TextStyle(color: Colors.grey[600])),
          Text(
            value,
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: valueColor,
            ),
          ),
        ],
      ),
    );
  }

  void _showDeleteConfirmation(tx.Transaction transaction) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除交易'),
        content: const Text('确定要删除这条交易记录吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              await context.read<TransactionProvider>().deleteTransaction(transaction.id!);
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