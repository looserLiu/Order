import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/account_provider.dart';
import '../../../data/models/account.dart';

class AccountsPage extends StatefulWidget {
  const AccountsPage({super.key});

  @override
  State<AccountsPage> createState() => _AccountsPageState();
}

class _AccountsPageState extends State<AccountsPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AccountProvider>().loadAccounts();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('账户管理'),
        centerTitle: true,
      ),
      body: Consumer<AccountProvider>(
        builder: (context, provider, child) {
          if (provider.accounts.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.account_balance_wallet,
                    size: 64,
                    color: Colors.grey.withAlpha(100),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    '暂无账户',
                    style: TextStyle(color: Colors.grey),
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton.icon(
                    onPressed: () => _showAccountDialog(),
                    icon: const Icon(Icons.add),
                    label: const Text('添加账户'),
                  ),
                ],
              ),
            );
          }

          final totalBalance = provider.accounts.fold<double>(
            0,
            (sum, account) => sum + account.balance,
          );

          return Column(
            children: [
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(20),
                color: Theme.of(context).primaryColor.withAlpha(30),
                child: Column(
                  children: [
                    const Text(
                      '账户总额',
                      style: TextStyle(color: Colors.grey),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '¥${totalBalance.toStringAsFixed(2)}',
                      style: const TextStyle(
                        fontSize: 32,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),
              Expanded(
                child: ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: provider.accounts.length,
                  itemBuilder: (context, index) {
                    final account = provider.accounts[index];
                    return _buildAccountCard(account);
                  },
                ),
              ),
            ],
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAccountDialog(),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildAccountCard(Account account) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Color(account.color ?? 0xFF2196F3).withAlpha(50),
          child: Icon(
            _getAccountIcon(account.type),
            color: Color(account.color ?? 0xFF2196F3),
          ),
        ),
        title: Text(account.name),
        subtitle: Text(_getAccountTypeName(account.type)),
        trailing: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Text(
              '¥${account.balance.toStringAsFixed(2)}',
              style: const TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
            Text(
              account.currency,
              style: const TextStyle(color: Colors.grey, fontSize: 12),
            ),
          ],
        ),
        onTap: () => _showAccountDialog(account),
        onLongPress: () => _showDeleteConfirmation(account),
      ),
    );
  }

  void _showAccountDialog([Account? account]) {
    final isEditing = account != null;
    final nameController = TextEditingController(text: account?.name ?? '');
    final balanceController = TextEditingController(
      text: account?.balance.toString() ?? '0',
    );
    String selectedType = account?.type ?? 'cash';
    int selectedColor = account?.color ?? 0xFF2196F3;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(isEditing ? '编辑账户' : '添加账户'),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                TextField(
                  controller: nameController,
                  decoration: const InputDecoration(
                    labelText: '账户名称',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: balanceController,
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  decoration: const InputDecoration(
                    labelText: '余额',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 16),
                const Text('账户类型'),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  children: [
                    _buildTypeChip('cash', '现金', Icons.wallet, selectedType, (type) {
                      setDialogState(() => selectedType = type);
                    }),
                    _buildTypeChip('bank', '银行', Icons.account_balance, selectedType, (type) {
                      setDialogState(() => selectedType = type);
                    }),
                    _buildTypeChip('credit_card', '信用卡', Icons.credit_card, selectedType, (type) {
                      setDialogState(() => selectedType = type);
                    }),
                    _buildTypeChip('digital', '数字钱包', Icons.phone_android, selectedType, (type) {
                      setDialogState(() => selectedType = type);
                    }),
                  ],
                ),
                const SizedBox(height: 16),
                const Text('颜色'),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  children: [
                    0xFF2196F3,
                    0xFF4CAF50,
                    0xFFFF9800,
                    0xFFE91E63,
                    0xFF9C27B0,
                    0xFF795548,
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
                    const SnackBar(content: Text('请输入账户名称')),
                  );
                  return;
                }

                final balance = double.tryParse(balanceController.text) ?? 0;

                if (isEditing) {
                  await context.read<AccountProvider>().updateAccount(
                        Account(
                          id: account.id,
                          name: name,
                          type: selectedType,
                          balance: balance,
                          currency: account.currency,
                          color: selectedColor,
                          createdAt: account.createdAt,
                          updatedAt: DateTime.now().millisecondsSinceEpoch,
                        ),
                      );
                } else {
                  await context.read<AccountProvider>().addAccount(
                        Account(
                          name: name,
                          type: selectedType,
                          balance: balance,
                          currency: 'CNY',
                          color: selectedColor,
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
      ),
    );
  }

  Widget _buildTypeChip(
    String type,
    String label,
    IconData icon,
    String selectedType,
    Function(String) onSelected,
  ) {
    final isSelected = selectedType == type;
    return GestureDetector(
      onTap: () => onSelected(type),
      child: Chip(
        avatar: Icon(icon, size: 16),
        label: Text(label),
        backgroundColor: isSelected ? Theme.of(context).primaryColor.withAlpha(50) : null,
      ),
    );
  }

  void _showDeleteConfirmation(Account account) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除账户'),
        content: Text('确定要删除账户 "${account.name}" 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              await context.read<AccountProvider>().deleteAccount(account.id!);
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

  String _getAccountTypeName(String type) {
    switch (type) {
      case 'cash':
        return '现金';
      case 'bank':
        return '银行账户';
      case 'credit_card':
        return '信用卡';
      case 'digital':
        return '数字钱包';
      default:
        return '其他';
    }
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
}