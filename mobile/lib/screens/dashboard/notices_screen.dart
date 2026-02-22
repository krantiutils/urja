import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/notice_service.dart';
import '../../models/notice.dart';
import '../shared/widgets/status_badge.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';
import '../shared/widgets/confirm_dialog.dart';

class NoticesScreen extends ConsumerStatefulWidget {
  const NoticesScreen({super.key});

  @override
  ConsumerState<NoticesScreen> createState() => _NoticesScreenState();
}

class _NoticesScreenState extends ConsumerState<NoticesScreen> {
  bool _loading = true;
  List<Notice> _notices = [];

  @override
  void initState() {
    super.initState();
    _loadNotices();
  }

  Future<void> _loadNotices() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final noticeService = NoticeService(ref.read(apiClientProvider));

    try {
      final result = await noticeService.getNotices(orgId);
      if (mounted) {
        setState(() {
          _notices = result.data;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load notices: $e')),
        );
      }
    }
  }

  Future<void> _showNoticeDialog({Notice? existing}) async {
    final l10n = AppLocalizations.of(context)!;
    final titleCtrl = TextEditingController(text: existing?.title ?? '');
    final titleNeCtrl = TextEditingController(text: existing?.titleNe ?? '');
    final contentCtrl = TextEditingController(text: existing?.content ?? '');
    final contentNeCtrl =
        TextEditingController(text: existing?.contentNe ?? '');
    bool isActive = existing?.isActive ?? true;

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(existing != null ? l10n.editNotice : l10n.addNotice),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(
                  controller: titleCtrl,
                  decoration: InputDecoration(labelText: l10n.title),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: titleNeCtrl,
                  decoration: InputDecoration(labelText: l10n.titleNe),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: contentCtrl,
                  decoration: InputDecoration(labelText: l10n.content),
                  maxLines: 4,
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: contentNeCtrl,
                  decoration: InputDecoration(labelText: l10n.contentNe),
                  maxLines: 4,
                ),
                const SizedBox(height: 12),
                SwitchListTile(
                  title: Text(l10n.active),
                  value: isActive,
                  onChanged: (v) => setDialogState(() => isActive = v),
                  contentPadding: EdgeInsets.zero,
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: Text(l10n.cancel),
            ),
            ElevatedButton(
              onPressed: () => Navigator.pop(context, true),
              child: Text(l10n.save),
            ),
          ],
        ),
      ),
    );

    if (result == true && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final noticeService = NoticeService(ref.read(apiClientProvider));
      final data = {
        'title': titleCtrl.text.trim(),
        'title_ne': titleNeCtrl.text.trim(),
        'content': contentCtrl.text.trim(),
        'content_ne': contentNeCtrl.text.trim(),
        'is_active': isActive,
      };

      try {
        if (existing != null) {
          await noticeService.updateNotice(orgId, existing.id, data);
        } else {
          await noticeService.addNotice(orgId, data);
        }
        _loadNotices();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to save notice: $e')),
          );
        }
      }
    }
  }

  Future<void> _deleteNotice(Notice notice) async {
    final confirmed = await showConfirmDialog(
      context,
      'Delete notice "${notice.title}"?',
    );

    if (confirmed && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final noticeService = NoticeService(ref.read(apiClientProvider));
      try {
        await noticeService.deleteNotice(orgId, notice.id);
        _loadNotices();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to delete notice: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.storiesNotices),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showNoticeDialog(),
          ),
        ],
      ),
      body: _loading
          ? const LoadingIndicator()
          : _notices.isEmpty
              ? EmptyState(
                  icon: Icons.article_outlined,
                  message: l10n.noNotices,
                )
              : RefreshIndicator(
                  onRefresh: _loadNotices,
                  child: ListView.builder(
                    itemCount: _notices.length,
                    padding: const EdgeInsets.all(16),
                    itemBuilder: (context, index) {
                      final notice = _notices[index];
                      final date = DateTime.tryParse(notice.createdAt);
                      final dateStr = date != null
                          ? '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}'
                          : '';

                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          title: Text(
                            notice.title,
                            style:
                                const TextStyle(fontWeight: FontWeight.w600),
                          ),
                          subtitle: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const SizedBox(height: 4),
                              Text(
                                notice.content,
                                maxLines: 2,
                                overflow: TextOverflow.ellipsis,
                                style: const TextStyle(
                                    color: AppTheme.textSecondary),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                dateStr,
                                style: const TextStyle(
                                    color: AppTheme.textSecondary,
                                    fontSize: 12),
                              ),
                            ],
                          ),
                          trailing: StatusBadge(
                            label: notice.isActive
                                ? l10n.active
                                : l10n.inactive,
                          ),
                          isThreeLine: true,
                          onTap: () =>
                              _showNoticeDialog(existing: notice),
                          onLongPress: () => _deleteNotice(notice),
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}
