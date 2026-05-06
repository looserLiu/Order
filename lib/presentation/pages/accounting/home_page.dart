import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/account_provider.dart';
import '../../providers/transaction_provider.dart';
import '../../providers/category_provider.dart';
import '../../providers/budget_provider.dart';
import 'add_transaction_page.dart';
import 'accounts_page.dart';
import 'categories_page.dart';
import 'budget_page.dart';
import '../../../data/models/account.dart';
import '../../../data/models/transaction.dart' as tx;
import '../../../data/models/category.dart';
import '../../../data/models/budget.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
    });
  }

  Future<void> _loadData() async {
    await context.read<AccountProvider>().loadAccounts();
    await context.read<TransactionProvider>().loadTransactions();
    await context.read<CategoryProvider>().loadCategories();
    await context.read<BudgetProvider>().loadBudgets();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('记账'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.account_balance_wallet),
            onPressed: () => Navigator.push(
              context,
              MaterialPageRoute(builder: (_) => const AccountsPage()),
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
              _buildAccountOverview(),
              const SizedBox(height: 24),
              _buildMonthlyStats(),
              const SizedBox(height: 24),
              _buildRecentTransactions(),
              const SizedBox(height: 24),
              _buildBudgetProgress(),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => Navigator.push(
          context,
          MaterialPageRoute(builder: (_) => const AddTransactionPage()),
        ),
        icon: const Icon(Icons.add),
        label: const Text('记一笔'),
      ),
    );
  }

  Widget _buildAccountOverview() {
    return Consumer<AccountProvider>(
      builder: (context, provider, child) {
        final totalBalance = provider.accounts.fold<double>(
          0,
          (sum, account) => sum + account.balance,
        );
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      '总资产',
                      style: TextStyle(
                        fontSize: 16,
                        color: Colors.grey,
                      ),
                    ),
                    TextButton(
                      onPressed: () => Navigator.push(
                        context,
                        MaterialPageRoute(builder: (_) => const AccountsPage()),
                      ),
                      child: const Text('管理账户'),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                Text(
                  '¥${totalBalance.toStringAsFixed(2)}',
                  style: const TextStyle(
                    fontSize: 32,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                SizedBox(
                  height: 80,
                  child: ListView.builder(
                    scrollDirection: Axis.horizontal,
                    itemCount: provider.accounts.length,
                    itemBuilder: (context, index) {
                      final account = provider.accounts[index];
                      return Container(
                        width: 120,
                        margin: const EdgeInsets.only(right: 12),
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: Color(account.color ?? 0xFF2196F3).withAlpha(30),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              _getAccountIcon(account.type),
                              color: Color(account.color ?? 0xFF2196F3),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              account.name,
                              style: const TextStyle(fontSize: 12),
                              overflow: TextOverflow.ellipsis,
                            ),
                            Text(
                              '¥${account.balance.toStringAsFixed(2)}',
                              style: const TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 14,
                              ),
                            ),
                          ],
                        ),
                      );
                    },
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildMonthlyStats() {
    return Consumer<TransactionProvider>(
      builder: (context, provider, child) {
        final now = DateTime.now();
        final monthStart = DateTime(now.year, now.month, 1);
        final monthStartMs = monthStart.millisecondsSinceEpoch;
        final monthTransactions = provider.transactions.where((t) {
          return t.date >= monthStartMs;
        }).toList();

        final income = monthTransactions
            .where((t) => t.type == 'income')
            .fold<double>(0, (sum, t) => sum + t.amount);
        final expense = monthTransactions
            .where((t) => t.type == 'expense')
            .fold<double>(0, (sum, t) => sum + t.amount);

        return Card(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '${now.month}月收支',
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    TextButton(
                      onPressed: () {},
                      child: const Text('查看详情'),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: _buildStatItem(
                        '收入',
                        income,
                        Colors.green,
                        Icons.arrow_downward,
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: _buildStatItem(
                        '支出',
                        expense,
                        Colors.red,
                        Icons.arrow_upward,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.grey.withAlpha(30),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Text('本月结余: '),
                      Text(
                        '¥${(income - expense).toStringAsFixed(2)}',
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: income >= expense ? Colors.green : Colors.red,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildStatItem(String label, double amount, Color color, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: color.withAlpha(30),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(icon, color: color, size: 16),
              const SizedBox(width: 4),
              Text(label, style: TextStyle(color: color)),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            '¥${amount.toStringAsFixed(2)}',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: color,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRecentTransactions() {
    return Consumer2<TransactionProvider, CategoryProvider>(
      builder: (context, txProvider, catProvider, child) {
        final recentTx = txProvider.transactions.take(5).toList();
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      '最近交易',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    TextButton(
                      onPressed: () {},
                      child: const Text('全部'),
                    ),
                  ],
                ),
                if (recentTx.isEmpty)
                  const Padding(
                    padding: EdgeInsets.symmetric(vertical: 20),
                    child: Center(
                      child: Text('暂无交易记录'),
                    ),
                  )
                else
                  ListView.separated(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    itemCount: recentTx.length,
                    separatorBuilder: (_, __) => const Divider(),
                    itemBuilder: (context, index) {
                      final tx = recentTx[index];
                      final category = catProvider.categories
                          .where((c) => c.id == tx.categoryId)
                          .firstOrNull;
                      final isExpense = tx.type == 'expense';
                      return ListTile(
                        contentPadding: EdgeInsets.zero,
                        leading: CircleAvatar(
                          backgroundColor:
                              Color(category?.color ?? 0xFF9E9E9E).withAlpha(50),
                          child: Icon(
                            _getCategoryIcon(category?.icon ?? 'category'),
                            color: Color(category?.color ?? 0xFF9E9E9E),
                            size: 20,
                          ),
                        ),
                        title: Text(category?.name ?? '未分类'),
                        subtitle: Text(
                          tx.description ?? '',
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        trailing: Text(
                          '${isExpense ? '-' : '+'}¥${tx.amount.toStringAsFixed(2)}',
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: isExpense ? Colors.red : Colors.green,
                          ),
                        ),
                      );
                    },
                  ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildBudgetProgress() {
    return Consumer2<BudgetProvider, TransactionProvider>(
      builder: (context, budgetProvider, txProvider, child) {
        final budgets = budgetProvider.budgets;
        if (budgets.isEmpty) {
          return Card(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  const Text(
                    '预算管理',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 12),
                  const Text('暂无预算设置'),
                  TextButton(
                    onPressed: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => const BudgetPage()),
                    ),
                    child: const Text('设置预算'),
                  ),
                ],
              ),
            ),
          );
        }
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      '预算进度',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    TextButton(
                      onPressed: () => Navigator.push(
                        context,
                        MaterialPageRoute(builder: (_) => const BudgetPage()),
                      ),
                      child: const Text('管理'),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                ...budgets.take(3).map((budget) {
                  final category = context
                      .read<CategoryProvider>()
                      .categories
                      .where((c) => c.id == budget.categoryId)
                      .firstOrNull;
                  final spent = txProvider.transactions
                      .where((t) =>
                          t.categoryId == budget.categoryId &&
                          t.type == 'expense' &&
                          t.date >= budget.startDate &&
                          t.date <= budget.endDate)
                      .fold<double>(0, (sum, t) => sum + t.amount);
                  final progress = (spent / budget.amount).clamp(0.0, 1.0);
                  final isOverBudget = spent > budget.amount;
                  return Padding(
                    padding: const EdgeInsets.only(bottom: 12),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(category?.name ?? 'Unknown'),
                            Text(
                              '¥${spent.toStringAsFixed(0)}/¥${budget.amount.toStringAsFixed(0)}',
                              style: TextStyle(
                                color: isOverBudget ? Colors.red : Colors.grey,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 4),
                        LinearProgressIndicator(
                          value: progress,
                          backgroundColor: Colors.grey.withAlpha(50),
                          valueColor: AlwaysStoppedAnimation(
                            isOverBudget ? Colors.red : Colors.blue,
                          ),
                        ),
                      ],
                    ),
                  );
                }),
              ],
            ),
          ),
        );
      },
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