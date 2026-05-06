import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:provider/provider.dart';
import '../../providers/transaction_provider.dart';
import '../../providers/category_provider.dart';
import '../../providers/account_provider.dart';

class ReportsPage extends StatefulWidget {
  const ReportsPage({super.key});

  @override
  State<ReportsPage> createState() => _ReportsPageState();
}

class _ReportsPageState extends State<ReportsPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  String _selectedPeriod = 'month';

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    await context.read<TransactionProvider>().loadTransactions();
    await context.read<CategoryProvider>().loadCategories();
    await context.read<AccountProvider>().loadAccounts();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('报表统计'),
        centerTitle: true,
        elevation: 0,
        backgroundColor: Colors.transparent,
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '支出', icon: Icon(Icons.trending_down)),
            Tab(text: '收入', icon: Icon(Icons.trending_up)),
            Tab(text: '账户', icon: Icon(Icons.account_balance_wallet)),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildExpenseTab(),
          _buildIncomeTab(),
          _buildAccountTab(),
        ],
      ),
    );
  }

  Widget _buildExpenseTab() {
    return Consumer2<TransactionProvider, CategoryProvider>(
      builder: (context, txProvider, catProvider, child) {
        final expenses = txProvider.transactions
            .where((t) => t.type == 'expense')
            .toList();

        if (expenses.isEmpty) {
          return _buildEmptyState('暂无支出记录');
        }

        final expenseByCategory = <int, double>{};
        for (var tx in expenses) {
          if (tx.categoryId != null) {
            expenseByCategory[tx.categoryId!] =
                (expenseByCategory[tx.categoryId!] ?? 0) + tx.amount;
          }
        }

        final total = expenseByCategory.values.fold(0.0, (a, b) => a + b);

        return RefreshIndicator(
          onRefresh: _loadData,
          child: SingleChildScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildPeriodSelector(),
                const SizedBox(height: 24),
                _buildPieChart(expenseByCategory, catProvider, total),
                const SizedBox(height: 24),
                _buildCategoryList(expenseByCategory, catProvider, total),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildIncomeTab() {
    return Consumer2<TransactionProvider, CategoryProvider>(
      builder: (context, txProvider, catProvider, child) {
        final incomes = txProvider.transactions
            .where((t) => t.type == 'income')
            .toList();

        if (incomes.isEmpty) {
          return _buildEmptyState('暂无收入记录');
        }

        final incomeByCategory = <int, double>{};
        for (var tx in incomes) {
          if (tx.categoryId != null) {
            incomeByCategory[tx.categoryId!] =
                (incomeByCategory[tx.categoryId!] ?? 0) + tx.amount;
          }
        }

        final total = incomeByCategory.values.fold(0.0, (a, b) => a + b);

        return RefreshIndicator(
          onRefresh: _loadData,
          child: SingleChildScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildPeriodSelector(),
                const SizedBox(height: 24),
                _buildPieChart(incomeByCategory, catProvider, total, isIncome: true),
                const SizedBox(height: 24),
                _buildCategoryList(incomeByCategory, catProvider, total, isIncome: true),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildAccountTab() {
    return Consumer<AccountProvider>(
      builder: (context, accountProvider, child) {
        final accounts = accountProvider.accounts;

        if (accounts.isEmpty) {
          return _buildEmptyState('暂无账户');
        }

        final totalBalance = accounts.fold<double>(0, (sum, acc) => sum + acc.balance);

        return RefreshIndicator(
          onRefresh: _loadData,
          child: SingleChildScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildTotalBalanceCard(totalBalance),
                const SizedBox(height: 24),
                ...accounts.map((acc) => _buildAccountCard(acc)),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildPeriodSelector() {
    return SegmentedButton<String>(
      segments: const [
        ButtonSegment(value: 'week', label: Text('本周')),
        ButtonSegment(value: 'month', label: Text('本月')),
        ButtonSegment(value: 'year', label: Text('本年')),
      ],
      selected: {_selectedPeriod},
      onSelectionChanged: (selection) {
        setState(() => _selectedPeriod = selection.first);
      },
    );
  }

  Widget _buildTotalBalanceCard(double total) {
    return Card(
      elevation: 0,
      color: Theme.of(context).colorScheme.primaryContainer,
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            Text(
              '总资产',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    color: Theme.of(context).colorScheme.onPrimaryContainer,
                  ),
            ),
            const SizedBox(height: 8),
            Text(
              '¥${total.toStringAsFixed(2)}',
              style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: Theme.of(context).colorScheme.onPrimaryContainer,
                  ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAccountCard(dynamic account) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ListTile(
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: Color(account.color ?? 0xFF9E9E9E).withAlpha(30),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Icon(
            _getAccountIcon(account.type),
            color: Color(account.color ?? 0xFF9E9E9E),
          ),
        ),
        title: Text(account.name),
        subtitle: Text(_getAccountTypeName(account.type)),
        trailing: Text(
          '¥${account.balance.toStringAsFixed(2)}',
          style: const TextStyle(
            fontWeight: FontWeight.bold,
            fontSize: 16,
          ),
        ),
      ),
    );
  }

  Widget _buildPieChart(Map<int, double> data, CategoryProvider catProvider, double total, {bool isIncome = false}) {
    if (data.isEmpty) return const SizedBox();

    final colors = [
      Colors.blue,
      Colors.green,
      Colors.orange,
      Colors.purple,
      Colors.red,
      Colors.teal,
      Colors.pink,
      Colors.indigo,
    ];

    final sections = <PieChartSectionData>[];
    var colorIndex = 0;

    data.forEach((categoryId, amount) {
      final category = catProvider.categories.where((c) => c.id == categoryId).firstOrNull;
      final percentage = (amount / total * 100);
      final color = Color(category?.color ?? colors[colorIndex % colors.length].value);

      sections.add(
        PieChartSectionData(
          value: amount,
          title: '${percentage.toStringAsFixed(1)}%',
          color: color,
          radius: 80,
          titleStyle: const TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
      );
      colorIndex++;
    });

    return SizedBox(
      height: 250,
      child: PieChart(
        PieChartData(
          sections: sections,
          centerSpaceRadius: 40,
          sectionsSpace: 2,
        ),
      ),
    );
  }

  Widget _buildCategoryList(Map<int, double> data, CategoryProvider catProvider, double total, {bool isIncome = false}) {
    final sortedEntries = data.entries.toList()
      ..sort((a, b) => b.value.compareTo(a.value));

    return Column(
      children: sortedEntries.map((entry) {
        final category = catProvider.categories.where((c) => c.id == entry.key).firstOrNull;
        final percentage = (entry.value / total * 100);

        return Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: ListTile(
            leading: Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: Color(category?.color ?? 0xFF9E9E9E).withAlpha(30),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Icon(
                _getCategoryIcon(category?.icon ?? 'category'),
                color: Color(category?.color ?? 0xFF9E9E9E),
                size: 20,
              ),
            ),
            title: Text(category?.name ?? '未知'),
            subtitle: Text('${percentage.toStringAsFixed(1)}%'),
            trailing: Text(
              '¥${entry.value.toStringAsFixed(2)}',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: isIncome ? Colors.green : Colors.red,
              ),
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildEmptyState(String message) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.bar_chart,
            size: 64,
            color: Colors.grey[300],
          ),
          const SizedBox(height: 16),
          Text(
            message,
            style: TextStyle(color: Colors.grey[500]),
          ),
        ],
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

  String _getAccountTypeName(String type) {
    switch (type) {
      case 'cash':
        return '现金';
      case 'bank':
        return '银行卡';
      case 'credit_card':
        return '信用卡';
      case 'digital':
        return '数字账户';
      default:
        return type;
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
      default:
        return Icons.category;
    }
  }
}
