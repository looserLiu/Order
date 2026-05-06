import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/settings_provider.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<SettingsProvider>().loadSettings();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('设置'),
        centerTitle: true,
        elevation: 0,
        backgroundColor: Colors.transparent,
      ),
      body: Consumer<SettingsProvider>(
        builder: (context, settings, child) {
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              _buildSection(
                '通用',
                [
                  _buildSettingTile(
                    icon: Icons.attach_money,
                    title: '默认货币',
                    subtitle: settings.currency,
                    onTap: () => _showCurrencyPicker(),
                  ),
                  _buildSwitchTile(
                    icon: Icons.dark_mode,
                    title: '深色模式',
                    value: settings.isDarkMode,
                    onChanged: (value) => settings.setDarkMode(value),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              _buildSection(
                '数据',
                [
                  _buildSettingTile(
                    icon: Icons.file_download,
                    title: '导入订单',
                    subtitle: '从 CSV 文件导入',
                    onTap: () => _importOrders(),
                  ),
                  _buildSettingTile(
                    icon: Icons.file_upload,
                    title: '导出数据',
                    subtitle: '导出为 CSV 格式',
                    onTap: () => _exportData(),
                  ),
                  _buildSettingTile(
                    icon: Icons.backup,
                    title: '备份',
                    subtitle: '创建数据备份',
                    onTap: () => _backupData(),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              _buildSection(
                '提醒',
                [
                  _buildSwitchTile(
                    icon: Icons.notifications,
                    title: '预算提醒',
                    subtitle: '超支时通知',
                    value: settings.budgetAlert,
                    onChanged: (value) => settings.setBudgetAlert(value),
                  ),
                  _buildSwitchTile(
                    icon: Icons.inventory_2,
                    title: '库存预警',
                    subtitle: '库存不足时通知',
                    value: settings.inventoryAlert,
                    onChanged: (value) => settings.setInventoryAlert(value),
                  ),
                  _buildSwitchTile(
                    icon: Icons.schedule,
                    title: '保质期提醒',
                    subtitle: '临近过期时通知',
                    value: settings.expirationAlert,
                    onChanged: (value) => settings.setExpirationAlert(value),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              _buildSection(
                '关于',
                [
                  _buildSettingTile(
                    icon: Icons.info,
                    title: '版本',
                    subtitle: '1.0.0',
                    onTap: null,
                  ),
                  _buildSettingTile(
                    icon: Icons.description,
                    title: '隐私政策',
                    onTap: () => _showPrivacyPolicy(),
                  ),
                  _buildSettingTile(
                    icon: Icons.article,
                    title: '使用条款',
                    onTap: () => _showTerms(),
                  ),
                ],
              ),
              const SizedBox(height: 32),
              Center(
                child: Text(
                  'Order - 记账与库存管理',
                  style: TextStyle(color: Colors.grey[500]),
                ),
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _buildSection(String title, List<Widget> children) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 4, bottom: 8),
          child: Text(
            title,
            style: TextStyle(
              color: Theme.of(context).colorScheme.primary,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
        Card(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
            side: BorderSide(color: Colors.grey[200]!),
          ),
          child: Column(children: children),
        ),
      ],
    );
  }

  Widget _buildSettingTile({
    required IconData icon,
    required String title,
    String? subtitle,
    VoidCallback? onTap,
  }) {
    return ListTile(
      leading: Icon(icon, color: Theme.of(context).colorScheme.primary),
      title: Text(title),
      subtitle: subtitle != null ? Text(subtitle) : null,
      trailing: onTap != null ? const Icon(Icons.chevron_right) : null,
      onTap: onTap,
    );
  }

  Widget _buildSwitchTile({
    required IconData icon,
    required String title,
    String? subtitle,
    required bool value,
    required ValueChanged<bool> onChanged,
  }) {
    return ListTile(
      leading: Icon(icon, color: Theme.of(context).colorScheme.primary),
      title: Text(title),
      subtitle: subtitle != null ? Text(subtitle) : null,
      trailing: Switch(
        value: value,
        onChanged: onChanged,
      ),
    );
  }

  void _showCurrencyPicker() {
    final currencies = ['CNY', 'USD', 'EUR', 'JPY', 'GBP'];
    showModalBottomSheet(
      context: context,
      builder: (context) => Column(
        mainAxisSize: MainAxisSize.min,
        children: currencies.map((currency) {
          return ListTile(
            title: Text(currency),
            onTap: () {
              context.read<SettingsProvider>().setCurrency(currency);
              Navigator.pop(context);
            },
          );
        }).toList(),
      ),
    );
  }

  void _importOrders() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('导入功能开发中')),
    );
  }

  void _exportData() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('导出功能开发中')),
    );
  }

  void _backupData() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('备份功能开发中')),
    );
  }

  void _showPrivacyPolicy() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('隐私政策'),
        content: const SingleChildScrollView(
          child: Text('隐私政策内容...'),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('关闭'),
          ),
        ],
      ),
    );
  }

  void _showTerms() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('使用条款'),
        content: const SingleChildScrollView(
          child: Text('使用条款内容...'),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('关闭'),
          ),
        ],
      ),
    );
  }
}
