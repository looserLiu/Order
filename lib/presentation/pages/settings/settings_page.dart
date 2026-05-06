import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/settings_provider.dart';
import '../../../core/services/backup_service.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  final BackupService _backupService = BackupService();
  bool _isExporting = false;
  bool _isImporting = false;

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
                    onTap: () => _showCurrencyPicker(settings),
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
                '提醒设置',
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
                  _buildSettingTile(
                    icon: Icons.warning,
                    title: '低库存阈值',
                    subtitle: '低于 ${settings.lowStockThreshold.toInt()} 件时提醒',
                    onTap: () => _showLowStockDialog(settings),
                  ),
                  _buildSwitchTile(
                    icon: Icons.schedule,
                    title: '保质期提醒',
                    subtitle: '临近过期时通知',
                    value: settings.expirationAlert,
                    onChanged: (value) => settings.setExpirationAlert(value),
                  ),
                  _buildSettingTile(
                    icon: Icons.date_range,
                    title: '过期预警天数',
                    subtitle: '提前 ${settings.expirationDays} 天提醒',
                    onTap: () => _showExpirationDaysDialog(settings),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              _buildSection(
                '数据管理',
                [
                  _buildSettingTile(
                    icon: Icons.file_upload,
                    title: '导出数据',
                    subtitle: '导出为 CSV 格式',
                    trailing: _isExporting
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : null,
                    onTap: _isExporting ? null : () => _exportData(),
                  ),
                  _buildSettingTile(
                    icon: Icons.file_download,
                    title: '导入数据',
                    subtitle: '从 CSV 文件导入',
                    trailing: _isImporting
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : null,
                    onTap: _isImporting ? null : () => _importData(),
                  ),
                  _buildSettingTile(
                    icon: Icons.backup,
                    title: '完整备份',
                    subtitle: '备份所有数据',
                    onTap: () => _fullBackup(),
                  ),
                  _buildSettingTile(
                    icon: Icons.restore,
                    title: '从备份恢复',
                    subtitle: '从备份文件恢复数据',
                    onTap: () => _showRestoreDialog(),
                  ),
                  _buildSettingTile(
                    icon: Icons.folder,
                    title: '查看备份',
                    subtitle: '管理已创建的备份',
                    onTap: () => _showBackupList(),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              _buildSection(
                '重置',
                [
                  _buildSettingTile(
                    icon: Icons.restart_alt,
                    title: '重置设置',
                    subtitle: '恢复默认设置',
                    onTap: () => _showResetSettingsDialog(settings),
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
    Widget? trailing,
    VoidCallback? onTap,
  }) {
    return ListTile(
      leading: Icon(icon, color: Theme.of(context).colorScheme.primary),
      title: Text(title),
      subtitle: subtitle != null ? Text(subtitle) : null,
      trailing: trailing ?? (onTap != null ? const Icon(Icons.chevron_right) : null),
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

  void _showCurrencyPicker(SettingsProvider settings) {
    final currencies = ['CNY', 'USD', 'EUR', 'JPY', 'GBP', 'HKD', 'KRW'];
    showModalBottomSheet(
      context: context,
      builder: (context) => Column(
        mainAxisSize: MainAxisSize.min,
        children: currencies.map((currency) {
          return ListTile(
            title: Text(currency),
            trailing: settings.currency == currency
                ? Icon(Icons.check, color: Theme.of(context).colorScheme.primary)
                : null,
            onTap: () {
              settings.setCurrency(currency);
              Navigator.pop(context);
            },
          );
        }).toList(),
      ),
    );
  }

  void _showLowStockDialog(SettingsProvider settings) {
    double tempValue = settings.lowStockThreshold;
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('低库存阈值'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('当前: ${tempValue.toInt()} 件'),
            Slider(
              value: tempValue,
              min: 5,
              max: 100,
              divisions: 19,
              label: tempValue.toInt().toString(),
              onChanged: (value) => tempValue = value,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              settings.setLowStockThreshold(tempValue);
              Navigator.pop(context);
            },
            child: const Text('保存'),
          ),
        ],
      ),
    );
  }

  void _showExpirationDaysDialog(SettingsProvider settings) {
    int tempValue = settings.expirationDays;
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('过期预警天数'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('提前 $tempValue 天提醒'),
            Slider(
              value: tempValue.toDouble(),
              min: 3,
              max: 30,
              divisions: 27,
              label: '$tempValue',
              onChanged: (value) => tempValue = value.toInt(),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              settings.setExpirationDays(tempValue);
              Navigator.pop(context);
            },
            child: const Text('保存'),
          ),
        ],
      ),
    );
  }

  Future<void> _exportData() async {
    setState(() => _isExporting = true);
    
    try {
      final path = await _backupService.exportTransactionsCSV();
      if (path != null) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('导出成功: $path')),
          );
        }
      } else {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('导出失败或暂无数据')),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('导出错误: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isExporting = false);
      }
    }
  }

  Future<void> _importData() async {
    // Note: In a real app, you would use file_picker package
    // For now, show instructions
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('导入数据'),
        content: const Text(
          '导入功能需要使用文件选择器。\n\n'
          '请将 CSV 文件放到应用文档目录下的 backups 文件夹中。',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('知道了'),
          ),
        ],
      ),
    );
  }

  Future<void> _fullBackup() async {
    setState(() => _isExporting = true);
    
    try {
      final path = await _backupService.exportAllData();
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('备份成功: $path')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('备份失败: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isExporting = false);
      }
    }
  }

  void _showRestoreDialog() async {
    final backups = await _backupService.listBackups();
    
    if (backups.isEmpty) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('暂无备份文件')),
        );
      }
      return;
    }

    if (!mounted) return;
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('选择备份'),
        content: SizedBox(
          width: double.maxFinite,
          height: 300,
          child: ListView.builder(
            itemCount: backups.length,
            itemBuilder: (context, index) {
              final backup = backups[index];
              return ListTile(
                title: Text(backup.name),
                subtitle: Text(
                  '${backup.createdAt.year}-${backup.createdAt.month.toString().padLeft(2, '0')}-${backup.createdAt.day.toString().padLeft(2, '0')} '
                  '${backup.createdAt.hour.toString().padLeft(2, '0')}:${backup.createdAt.minute.toString().padLeft(2, '0')}',
                ),
                onTap: () async {
                  Navigator.pop(context);
                  await _confirmRestore(backup.path);
                },
              );
            },
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
        ],
      ),
    );
  }

  Future<void> _confirmRestore(String backupPath) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认恢复'),
        content: const Text(
          '恢复操作将覆盖当前数据。\n\n'
          '建议先创建备份再进行恢复。\n\n'
          '是否继续？',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('继续恢复'),
          ),
        ],
      ),
    );

    if (confirmed != true) return;

    setState(() => _isImporting = true);

    try {
      final result = await _backupService.importData(backupPath);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(result.message)),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('恢复失败: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isImporting = false);
      }
    }
  }

  void _showBackupList() async {
    final backups = await _backupService.listBackups();

    if (!mounted) return;

    if (backups.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('暂无备份')),
      );
      return;
    }

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('备份列表'),
        content: SizedBox(
          width: double.maxFinite,
          height: 300,
          child: ListView.builder(
            itemCount: backups.length,
            itemBuilder: (context, index) {
              final backup = backups[index];
              return ListTile(
                title: Text(backup.name),
                subtitle: Text(
                  '${backup.createdAt.year}-${backup.createdAt.month.toString().padLeft(2, '0')}-${backup.createdAt.day.toString().padLeft(2, '0')}',
                ),
                trailing: IconButton(
                  icon: const Icon(Icons.delete_outline),
                  onPressed: () async {
                    final deleted = await _backupService.deleteBackup(backup.path);
                    if (mounted) {
                      Navigator.pop(context);
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          content: Text(deleted ? '删除成功' : '删除失败'),
                        ),
                      );
                    }
                  },
                ),
              );
            },
          ),
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

  void _showResetSettingsDialog(SettingsProvider settings) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('重置设置'),
        content: const Text('确定要恢复所有设置为默认值吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              settings.resetToDefaults();
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('设置已重置')),
              );
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _showPrivacyPolicy() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('隐私政策'),
        content: const SingleChildScrollView(
          child: Text(
            'Order 尊重并保护您的个人隐私。\n\n'
            '• 本应用所有数据均存储在本地设备\n'
            '• 不会收集或上传您的财务数据\n'
            '• 备份文件仅存储在您设备本地\n'
            '• 我们不会获取您的账户信息',
          ),
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
          child: Text(
            '使用 Order 应用即表示您同意以下条款：\n\n'
            '• 本应用按"现状"提供，不提供任何明示或暗示的保证\n'
            '• 您对数据安全负责，请定期备份重要数据\n'
            '• 请勿使用本应用进行非法交易记录\n'
            '• 我们保留随时修改条款的权利',
          ),
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
