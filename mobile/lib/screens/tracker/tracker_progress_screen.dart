import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../member/member_progress_screen.dart';

/// The tracker progress tab reuses the same MemberProgressScreen content.
class TrackerProgressScreen extends ConsumerWidget {
  const TrackerProgressScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return const MemberProgressScreen();
  }
}
