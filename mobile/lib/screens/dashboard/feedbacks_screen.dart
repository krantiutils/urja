import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/notice_service.dart';
import '../../models/feedback.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';
import '../shared/widgets/confirm_dialog.dart';

class FeedbacksScreen extends ConsumerStatefulWidget {
  const FeedbacksScreen({super.key});

  @override
  ConsumerState<FeedbacksScreen> createState() => _FeedbacksScreenState();
}

class _FeedbacksScreenState extends ConsumerState<FeedbacksScreen> {
  bool _loading = true;
  List<MemberFeedback> _feedbacks = [];

  @override
  void initState() {
    super.initState();
    _loadFeedbacks();
  }

  Future<void> _loadFeedbacks() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final noticeService = NoticeService(ref.read(apiClientProvider));

    try {
      final result = await noticeService.getFeedbacks(orgId);
      if (mounted) {
        setState(() {
          _feedbacks = result.data;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load feedbacks: $e')),
        );
      }
    }
  }

  Future<void> _deleteFeedback(MemberFeedback feedback) async {
    final l10n = AppLocalizations.of(context)!;
    final confirmed = await showConfirmDialog(
      context,
      '${l10n.deleteFeedbackConfirm} ${feedback.memberName}?',
    );

    if (confirmed && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final noticeService = NoticeService(ref.read(apiClientProvider));
      try {
        await noticeService.deleteFeedback(orgId, feedback.id);
        _loadFeedbacks();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to delete feedback: $e')),
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
        title: Text(l10n.feedbacks),
      ),
      body: _loading
          ? const LoadingIndicator()
          : _feedbacks.isEmpty
              ? EmptyState(
                  icon: Icons.feedback_outlined,
                  message: l10n.noFeedbacks,
                )
              : RefreshIndicator(
                  onRefresh: _loadFeedbacks,
                  child: ListView.builder(
                    itemCount: _feedbacks.length,
                    padding: const EdgeInsets.all(16),
                    itemBuilder: (context, index) {
                      final feedback = _feedbacks[index];
                      final date = DateTime.tryParse(feedback.createdAt);
                      final dateStr = date != null
                          ? '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}'
                          : '';

                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: Padding(
                          padding: const EdgeInsets.all(16),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  CircleAvatar(
                                    backgroundColor: AppTheme.primary,
                                    backgroundImage:
                                        feedback.avatarUrl != null
                                            ? NetworkImage(
                                                feedback.avatarUrl!)
                                            : null,
                                    radius: 20,
                                    child: feedback.avatarUrl == null
                                        ? Text(
                                            feedback.memberName.isNotEmpty
                                                ? feedback.memberName[0]
                                                    .toUpperCase()
                                                : '?',
                                            style: const TextStyle(
                                                color: Colors.white),
                                          )
                                        : null,
                                  ),
                                  const SizedBox(width: 12),
                                  Expanded(
                                    child: Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.start,
                                      children: [
                                        Text(
                                          feedback.memberName,
                                          style: const TextStyle(
                                            fontWeight: FontWeight.w600,
                                            fontSize: 15,
                                          ),
                                        ),
                                        Text(
                                          dateStr,
                                          style: const TextStyle(
                                            color: AppTheme.textSecondary,
                                            fontSize: 12,
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                  IconButton(
                                    icon: const Icon(Icons.delete_outline,
                                        color: AppTheme.error),
                                    onPressed: () =>
                                        _deleteFeedback(feedback),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 12),
                              Text(
                                feedback.message,
                                style: const TextStyle(
                                    color: AppTheme.textSecondary),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}
