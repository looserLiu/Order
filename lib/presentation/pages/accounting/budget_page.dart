import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/budget_provider.dart';
import '../../providers/category_provider.dart';
import '../../providers/transaction_provider.dart';
import '../../../data/models/budget.dart';

class BudgetPage extends StatefulWidget {
  const BudgetPage({super.key});

  @override
  State<BudgetPage> createState() => _BudgetPageState();
}

class _BudgetPageState extends State<BudgetPage> with SingleTickerProviderStateMixin {
  late AnimationController _animationController;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 800),
    );
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
      _animationController.forward();
    });
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    await context.read<BudgetProvider>().loadBudgets();
    await context.read<CategoryProvider>().loadCategories();
    await context.read<TransactionProvider>().loadTransactions();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface,
      appBar: AppBar(
        title: const Text('预算管理'),
        centerTitle: true,
        elevation: 0,
        backgroundColor: Colors.transparent,
      ),
      body: RefreshIndicator(
        onRefresh: _loadData,
        color: Theme.of(context).colorScheme.primary,
        child: Consumer3<BudgetProvider, CategoryProvider, TransactionProvider>(
          builder: (context, budgetProvider, catProvider, txProvider, child) {
            final budgets = budgetProvider.budgets;

            if (budgets.isEmpty) {
              return _buildEmptyState(context);
            }

            return _buildBudgetList(budgets, catProvider, txProvider);
          },
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _showBudgetDialog(),
        icon: const Icon(Icons.add),
        label: const Text('添加预算'),
      ),
    );
  }

  Widget _buildEmptyState(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 120,
            height: 120,
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primaryContainer.withAlpha(50),
              shape: BoxShape.circle,
            ),
            child: Icon(
              Icons.pie_chart_outline,
              size: 64,
              color: Theme.of(context).colorScheme.primary.withAlpha(150),
            ),
          ),
          const SizedBox(height: 24),
          Text(
            '暂无预算',
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  color: Colors.grey[600],
                ),
          ),
          const SizedBox(height: 8),
          Text(
            '设置预算帮助您合理消费',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.grey[500],
                ),
          ),
          const SizedBox(height: 24),
          FilledButton.icon(
            onPressed: () => _showBudgetDialog(),
            icon: const Icon(Icons.add),
            label: const Text('添加第一个预算'),
          ),
        ],
      ),
    );
  }

  Widget _buildBudgetList(
    List<Budget> budgets,
    CategoryProvider catProvider,
    TransactionProvider txProvider,
  ) {
    return ListView.builder(
      padding: const EdgeInsets.fromLTRB(16, 8, 16, 100),
      itemCount: budgets.length,
      itemBuilder: (context, index) {
        final budget = budgets[index];
        final category = catProvider.categories
            .where((c) => c.id == budget.categoryId)
            .firstOrNull;
        final spent = txProvider.transactions
            .where((t) =>
                t.categoryId == budget.categoryId &&
                t.type == 'expense' &&
                t.date >= budget.startDate - const Duration(days: 1).inMilliseconds &&
                t.date <= budget.endDate + const Duration(days: 1).inMilliseconds)
            .fold<double>(0, (sum, t) => sum + t.amount);
        final remaining = budget.amount - spent;
        final progress = (spent / budget.amount).clamp(0.0, 1.0);
        final isOverBudget = spent > budget.amount;

        return TweenAnimationBuilder<double>(
          tween: Tween(begin: 0.0, end: 1.0),
          duration: Duration(milliseconds: 300 + (index * 100)),
          curve: Curves.easeOutCubic,
          builder: (context, value, child) {
            return Transform.translate(
              offset: Offset(0, 20 * (1 - value)),
              child: Opacity(opacity: value, child: child),
            );
          },
          child: _buildBudgetCard(
            context,
            budget,
            category,
            spent,
            remaining,
            progress,
            isOverBudget,
          ),
        );
      },
    );
  }

  Widget _buildBudgetCard(
    BuildContext context,
    Budget budget,
    dynamic category,
    double spent,
    double remaining,
    double progress,
    bool isOverBudget,
  ) {
    final categoryColor = Color(category?.color ?? 0xFF9E9E9E);
    final periodColor = _getPeriodColor(budget.period);

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 2,
      shadowColor: categoryColor.withAlpha(50),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      child: InkWell(
        borderRadius: BorderRadius.circular(16),
        onTap: () => _showBudgetDialog(budget),
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: categoryColor.withAlpha(30),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      _getCategoryIcon(category?.icon ?? 'category'),
                      color: categoryColor,
                      size: 24,
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          category?.name ?? '未知分类',
                          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                fontWeight: FontWeight.bold,
                              ),
                        ),
                        const SizedBox(height: 4),
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 8,
                            vertical: 2,
                          ),
                          decoration: BoxDecoration(
                            color: periodColor.withAlpha(20),
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: Text(
                            _getPeriodText(budget.period),
                            style: TextStyle(
                              color: periodColor,
                              fontSize: 11,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  PopupMenuButton<String>(
                    onSelected: (value) {
                      if (value == 'edit') {
                        _showBudgetDialog(budget);
                      } else if (value == 'delete') {
                        _showDeleteConfirmation(budget);
                      }
                    },
                    itemBuilder: (context) => [
                      const PopupMenuItem(
                        value: 'edit',
                        child: Row(
                          children: [
                            Icon(Icons.edit, size: 20),
                            SizedBox(width: 8),
                            Text('编辑'),
                          ],
                        ),
                      ),
                      const PopupMenuItem(
                        value: 'delete',
                        child: Row(
                          children: [
                            Icon(Icons.delete, color: Colors.red, size: 20),
                            SizedBox(width: 8),
                            Text('删除', style: TextStyle(color: Colors.red)),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              const SizedBox(height: 20),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '已花费',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                              color: Colors.grey[600],
                            ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        '¥${spent.toStringAsFixed(2)}',
                        style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                              fontWeight: FontWeight.bold,
                              color: isOverBudget ? Colors.red : null,
                            ),
                      ),
                    ],
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text(
                        '预算',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                              color: Colors.grey[600],
                            ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        '¥${budget.amount.toStringAsFixed(2)}',
                        style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                      ),
                    ],
                  ),
                ],
              ),
              const SizedBox(height: 16),
              TweenAnimationBuilder<double>(
                tween: Tween(begin: 0.0, end: progress),
                duration: const Duration(milliseconds: 1000),
                curve: Curves.easeOutCubic,
                builder: (context, value, child) {
                  return Stack(
                    children: [
                      Container(
                        height: 10,
                        decoration: BoxDecoration(
                          color: Colors.grey[200],
                          borderRadius: BorderRadius.circular(5),
                        ),
                      ),
                      Container(
                        height: 10,
                        width: MediaQuery.of(context).size.width * 0.85 * value,
                        decoration: BoxDecoration(
                          gradient: LinearGradient(
                            colors: isOverBudget
                                ? [Colors.red[400]!, Colors.red[600]!]
                                : [categoryColor.withAlpha(180), categoryColor],
                          ),
                          borderRadius: BorderRadius.circular(5),
                          boxShadow: [
                            BoxShadow(
                              color: (isOverBudget ? Colors.red : categoryColor).withAlpha(50),
                              blurRadius: 4,
                              offset: const Offset(0, 2),
                            ),
                          ],
                        ),
                      ),
                    ],
                  );
                },
              ),
              const SizedBox(height: 12),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '${(progress * 100).toStringAsFixed(0)}% 已使用',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: isOverBudget ? Colors.red : Colors.grey[600],
                          fontWeight: isOverBudget ? FontWeight.bold : null,
                        ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 10,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: (isOverBudget ? Colors.red : Colors.green).withAlpha(20),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      isOverBudget
                          ? '超出 ¥${(-remaining).toStringAsFixed(2)}'
                          : '剩余 ¥${remaining.toStringAsFixed(2)}',
                      style: TextStyle(
                        color: isOverBudget ? Colors.red : Colors.green[700],
                        fontWeight: FontWeight.bold,
                        fontSize: 12,
                      ),
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

  void _showBudgetDialog([Budget? budget]) {
    final isEditing = budget != null;
    final amountController = TextEditingController(
      text: budget?.amount.toString() ?? '',
    );
    String selectedPeriod = budget?.period ?? 'monthly';
    int? selectedCategoryId = budget?.categoryId;

    final expenseCategories = context
        .read<CategoryProvider>()
        .categories
        .where((c) => c.type == 'expense')
        .toList();

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        padding: EdgeInsets.only(
          bottom: MediaQuery.of(context).viewInsets.bottom,
        ),
        decoration: BoxDecoration(
          color: Theme.of(context).colorScheme.surface,
          borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
        ),
        child: StatefulBuilder(
          builder: (context, setDialogState) {
            if (selectedCategoryId == null && expenseCategories.isNotEmpty) {
              WidgetsBinding.instance.addPostFrameCallback((_) {
                setDialogState(() {
                  selectedCategoryId = expenseCategories.first.id;
                });
              });
            }

            return Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Center(
                    child: Container(
                      width: 40,
                      height: 4,
                      decoration: BoxDecoration(
                        color: Colors.grey[300],
                        borderRadius: BorderRadius.circular(2),
                      ),
                    ),
                  ),
                  const SizedBox(height: 20),
                  Text(
                    isEditing ? '编辑预算' : '添加预算',
                    style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                  const SizedBox(height: 24),
                  if (!isEditing) ...[
                    Text(
                      '选择分类',
                      style: Theme.of(context).textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                    ),
                    const SizedBox(height: 12),
                    Wrap(
                      spacing: 8,
                      runSpacing: 8,
                      children: expenseCategories.map((cat) {
                        final isSelected = selectedCategoryId == cat.id;
                        return ChoiceChip(
                          label: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(
                                _getCategoryIcon(cat.icon ?? 'category'),
                                size: 16,
                                color: isSelected
                                    ? Colors.white
                                    : Color(cat.color ?? 0xFF9E9E9E),
                              ),
                              const SizedBox(width: 4),
                              Text(cat.name),
                            ],
                          ),
                          selected: isSelected,
                          selectedColor: Color(cat.color ?? 0xFF9E9E9E),
                          onSelected: (selected) {
                            if (selected) {
                              setDialogState(() => selectedCategoryId = cat.id);
                            }
                          },
                        );
                      }).toList(),
                    ),
                    const SizedBox(height: 20),
                  ],
                  Text(
                    '预算金额',
                    style: Theme.of(context).textTheme.titleSmall?.copyWith(
                          fontWeight: FontWeight.w600,
                        ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: amountController,
                    keyboardType:
                        const TextInputType.numberWithOptions(decimal: true),
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                    decoration: InputDecoration(
                      prefixText: '¥ ',
                      prefixStyle: TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.bold,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                    ),
                  ),
                  const SizedBox(height: 20),
                  Text(
                    '预算周期',
                    style: Theme.of(context).textTheme.titleSmall?.copyWith(
                          fontWeight: FontWeight.w600,
                        ),
                  ),
                  const SizedBox(height: 12),
                  SegmentedButton<String>(
                    segments: const [
                      ButtonSegment(
                        value: 'weekly',
                        label: Text('每周'),
                        icon: Icon(Icons.view_week),
                      ),
                      ButtonSegment(
                        value: 'monthly',
                        label: Text('每月'),
                        icon: Icon(Icons.calendar_month),
                      ),
                      ButtonSegment(
                        value: 'yearly',
                        label: Text('每年'),
                        icon: Icon(Icons.calendar_today),
                      ),
                    ],
                    selected: {selectedPeriod},
                    onSelectionChanged: (selection) {
                      setDialogState(() => selectedPeriod = selection.first);
                    },
                  ),
                  const SizedBox(height: 32),
                  Row(
                    children: [
                      Expanded(
                        child: OutlinedButton(
                          onPressed: () => Navigator.pop(context),
                          style: OutlinedButton.styleFrom(
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                          child: const Text('取消'),
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        flex: 2,
                        child: FilledButton(
                          onPressed: () async {
                            final amountText = amountController.text.trim();
                            if (amountText.isEmpty) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                  content: Text('请输入预算金额'),
                                  behavior: SnackBarBehavior.floating,
                                ),
                              );
                              return;
                            }

                            final amount = double.tryParse(amountText);
                            if (amount == null || amount <= 0) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                  content: Text('请输入有效金额'),
                                  behavior: SnackBarBehavior.floating,
                                ),
                              );
                              return;
                            }

                            if (selectedCategoryId == null) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                  content: Text('请选择分类'),
                                  behavior: SnackBarBehavior.floating,
                                ),
                              );
                              return;
                            }

                            final now = DateTime.now();
                            DateTime startDate;
                            DateTime endDate;

                            switch (selectedPeriod) {
                              case 'weekly':
                                startDate = now.subtract(Duration(days: now.weekday - 1));
                                endDate = startDate.add(const Duration(days: 6));
                                break;
                              case 'monthly':
                                startDate = DateTime(now.year, now.month, 1);
                                endDate = DateTime(now.year, now.month + 1, 0);
                                break;
                              case 'yearly':
                                startDate = DateTime(now.year, 1, 1);
                                endDate = DateTime(now.year, 12, 31);
                                break;
                              default:
                                startDate = DateTime(now.year, now.month, 1);
                                endDate = DateTime(now.year, now.month + 1, 0);
                            }

                            if (isEditing) {
                              await context.read<BudgetProvider>().updateBudget(
                                    Budget(
                                      id: budget.id,
                                      categoryId: budget.categoryId,
                                      amount: amount,
                                      period: selectedPeriod,
                                      startDate: startDate.millisecondsSinceEpoch,
                                      endDate: endDate.millisecondsSinceEpoch,
                                    ),
                                  );
                            } else {
                              await context.read<BudgetProvider>().addBudget(
                                    Budget(
                                      categoryId: selectedCategoryId!,
                                      amount: amount,
                                      period: selectedPeriod,
                                      startDate: startDate.millisecondsSinceEpoch,
                                      endDate: endDate.millisecondsSinceEpoch,
                                    ),
                                  );
                            }

                            if (context.mounted) {
                              Navigator.pop(context);
                            }
                          },
                          style: FilledButton.styleFrom(
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                          child: Text(isEditing ? '保存' : '添加'),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                ],
              ),
            );
          },
        ),
      ),
    );
  }

  void _showDeleteConfirmation(Budget budget) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
        title: const Text('删除预算'),
        content: const Text('确定要删除此预算吗？此操作无法撤销。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          FilledButton(
            onPressed: () async {
              await context.read<BudgetProvider>().deleteBudget(budget.id!);
              if (context.mounted) {
                Navigator.pop(context);
              }
            },
            style: FilledButton.styleFrom(
              backgroundColor: Colors.red,
            ),
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
      default:
        return Icons.category;
    }
  }

  String _getPeriodText(String period) {
    switch (period) {
      case 'weekly':
        return '每周';
      case 'monthly':
        return '每月';
      case 'yearly':
        return '每年';
      default:
        return period;
    }
  }

  Color _getPeriodColor(String period) {
    switch (period) {
      case 'weekly':
        return Colors.blue;
      case 'monthly':
        return Colors.orange;
      case 'yearly':
        return Colors.purple;
      default:
        return Colors.grey;
    }
  }
}
