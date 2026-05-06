import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/category_provider.dart';
import '../../../data/models/category.dart';

class CategoriesPage extends StatefulWidget {
  const CategoriesPage({super.key});

  @override
  State<CategoriesPage> createState() => _CategoriesPageState();
}

class _CategoriesPageState extends State<CategoriesPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<CategoryProvider>().loadCategories();
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('分类管理'),
        centerTitle: true,
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '支出'),
            Tab(text: '收入'),
          ],
        ),
      ),
      body: Consumer<CategoryProvider>(
        builder: (context, provider, child) {
          final List<Category> expenseCategories =
              provider.categories.where((c) => c.type == 'expense').toList();
          final List<Category> incomeCategories =
              provider.categories.where((c) => c.type == 'income').toList();

          return TabBarView(
            controller: _tabController,
            children: [
              _buildCategoryList(expenseCategories, 'expense'),
              _buildCategoryList(incomeCategories, 'income'),
            ],
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showCategoryDialog(),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildCategoryList(List<Category> categories, String type) {
    if (categories.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.category,
              size: 64,
              color: Colors.grey.withAlpha(100),
            ),
            const SizedBox(height: 16),
            Text(
              '暂无${type == 'expense' ? '支出' : '收入'}分类',
              style: const TextStyle(color: Colors.grey),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: categories.length,
      itemBuilder: (context, index) {
        final category = categories[index];
        return _buildCategoryCard(category);
      },
    );
  }

  Widget _buildCategoryCard(Category category) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Color(category.color ?? 0xFF9E9E9E).withAlpha(50),
          child: Icon(
            _getCategoryIcon(category.icon ?? 'category'),
            color: Color(category.color ?? 0xFF9E9E9E),
          ),
        ),
        title: Text(category.name),
        subtitle: Text(
          '使用 ${category.usageCount} 次',
          style: const TextStyle(fontSize: 12),
        ),
        trailing: category.isSystem == 1
            ? const Chip(
                label: Text('系统', style: TextStyle(fontSize: 10)),
                padding: EdgeInsets.zero,
              )
            : IconButton(
                icon: const Icon(Icons.edit),
                onPressed: () => _showCategoryDialog(category),
              ),
        onLongPress: category.isSystem == 0
            ? () => _showDeleteConfirmation(category)
            : null,
      ),
    );
  }

  void _showCategoryDialog([Category? category]) {
    final isEditing = category != null;
    final nameController = TextEditingController(text: category?.name ?? '');
    String selectedType = category?.type ?? 'expense';
    String selectedIcon = category?.icon ?? 'category';
    int selectedColor = category?.color ?? 0xFF9E9E9E;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(isEditing ? '编辑分类' : '添加分类'),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                TextField(
                  controller: nameController,
                  decoration: const InputDecoration(
                    labelText: '分类名称',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 16),
                if (!isEditing) ...[
                  const Text('分类类型'),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      Expanded(
                        child: ChoiceChip(
                          label: const Text('支出'),
                          selected: selectedType == 'expense',
                          onSelected: (selected) {
                            if (selected) {
                              setDialogState(() => selectedType = 'expense');
                            }
                          },
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: ChoiceChip(
                          label: const Text('收入'),
                          selected: selectedType == 'income',
                          onSelected: (selected) {
                            if (selected) {
                              setDialogState(() => selectedType = 'income');
                            }
                          },
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                ],
                const Text('图标'),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: _iconOptions.map((icon) {
                    final isSelected = selectedIcon == icon;
                    return GestureDetector(
                      onTap: () => setDialogState(() => selectedIcon = icon),
                      child: Container(
                        width: 40,
                        height: 40,
                        decoration: BoxDecoration(
                          color: isSelected
                              ? Theme.of(context).primaryColor
                              : Colors.grey.withAlpha(30),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Icon(
                          _getCategoryIcon(icon),
                          color: isSelected ? Colors.white : Colors.grey,
                        ),
                      ),
                    );
                  }).toList(),
                ),
                const SizedBox(height: 16),
                const Text('颜色'),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  children: [
                    0xFFE57373,
                    0xFF64B5F6,
                    0xFFBA68C8,
                    0xFFFFD54F,
                    0xFF4DB6AC,
                    0xFF90A4AE,
                    0xFF81C784,
                    0xFFF06292,
                    0xFF7E57C2,
                  ].map((color) {
                    return GestureDetector(
                      onTap: () => setDialogState(() => selectedColor = color),
                      child: Container(
                        width: 40,
                        height: 40,
                        decoration: BoxDecoration(
                          color: Color(color),
                          shape: BoxShape.circle,
                          border: selectedColor == color
                              ? Border.all(color: Colors.black, width: 3)
                              : null,
                        ),
                      ),
                    );
                  }).toList(),
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
                    const SnackBar(content: Text('请输入分类名称')),
                  );
                  return;
                }

                if (isEditing) {
                  await context.read<CategoryProvider>().updateCategory(
                        Category(
                          id: category.id,
                          name: name,
                          type: category.type,
                          icon: selectedIcon,
                          color: selectedColor,
                          parentId: category.parentId,
                          isSystem: category.isSystem,
                          usageCount: category.usageCount,
                        ),
                      );
                } else {
                  await context.read<CategoryProvider>().addCategory(
                        Category(
                          name: name,
                          type: selectedType,
                          icon: selectedIcon,
                          color: selectedColor,
                          parentId: null,
                          isSystem: false,
                          usageCount: 0,
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
      ),
    );
  }

  void _showDeleteConfirmation(Category category) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除分类'),
        content: Text('确定要删除分类 "${category.name}" 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              await context.read<CategoryProvider>().deleteCategory(category.id!);
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

  static const List<String> _iconOptions = [
    'restaurant',
    'directions_car',
    'shopping_bag',
    'movie',
    'bolt',
    'home',
    'local_hospital',
    'school',
    'work',
    'card_giftcard',
    'trending_up',
    'laptop',
    'category',
  ];
}