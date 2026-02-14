import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:percent_indicator/circular_percent_indicator.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../providers/providers.dart';

class PackagesScreen extends ConsumerWidget {
  const PackagesScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context)!;
    final packagesAsync = ref.watch(memberPackagesProvider);

    return Scaffold(
      appBar: AppBar(title: Text(l10n.myPackage)),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(memberPackagesProvider);
        },
        child: packagesAsync.when(
          data: (packages) {
            if (packages.isEmpty) {
              return Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.card_membership,
                        size: 64, color: UrjaColors.foregroundMuted),
                    const SizedBox(height: 16),
                    Text(l10n.noPackages,
                        style: Theme.of(context).textTheme.bodyLarge),
                  ],
                ),
              );
            }

            return ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // Active package
                ...packages.where((p) => p.isActive && !p.isExpired).map(
                  (pkg) => _ActivePackageCard(pkg: pkg, l10n: l10n),
                ),
                const SizedBox(height: 16),
                // Payment history header
                Text(l10n.paymentHistory,
                    style: Theme.of(context).textTheme.titleMedium),
                const SizedBox(height: 8),
                // All packages as history
                ...packages.map(
                  (pkg) => _PaymentHistoryItem(pkg: pkg),
                ),
              ],
            );
          },
          loading: () =>
              const Center(child: CircularProgressIndicator()),
          error: (err, _) => Center(
            child: Text(l10n.error),
          ),
        ),
      ),
    );
  }
}

class _ActivePackageCard extends StatelessWidget {
  final dynamic pkg;
  final AppLocalizations l10n;

  const _ActivePackageCard({required this.pkg, required this.l10n});

  @override
  Widget build(BuildContext context) {
    final daysLeft = pkg.daysRemaining as int;
    final progress = pkg.progressFraction as double;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                CircularPercentIndicator(
                  radius: 40,
                  lineWidth: 5,
                  percent: 1.0 - progress,
                  center: Text(
                    '$daysLeft',
                    style: Theme.of(context)
                        .textTheme
                        .headlineSmall
                        ?.copyWith(color: UrjaColors.accent),
                  ),
                  progressColor: UrjaColors.accent,
                  backgroundColor: UrjaColors.borderDefault,
                  circularStrokeCap: CircularStrokeCap.round,
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        pkg.packageName ?? l10n.currentPackage,
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                      const SizedBox(height: 4),
                      Text(
                        l10n.daysLeft(daysLeft),
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                      Text(
                        '${DateFormat('MMM dd').format(pkg.startDate as DateTime)} - ${DateFormat('MMM dd, yyyy').format(pkg.endDate as DateTime)}',
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ],
                  ),
                ),
              ],
            ),
            if (daysLeft <= 7) ...[
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    // TODO: Renewal flow
                  },
                  child: Text(l10n.renewPackage),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _PaymentHistoryItem extends StatelessWidget {
  final dynamic pkg;

  const _PaymentHistoryItem({required this.pkg});

  @override
  Widget build(BuildContext context) {
    final statusColor = _statusColor(pkg.status as String);

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Icon(
          Icons.receipt_long,
          color: statusColor,
        ),
        title: Text(pkg.packageName ?? 'Package'),
        subtitle: Text(
          'Rs. ${pkg.amountPaid} - ${DateFormat('MMM dd, yyyy').format(pkg.createdAt as DateTime)}',
        ),
        trailing: Container(
          padding:
              const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
          decoration: BoxDecoration(
            color: statusColor.withAlpha(25),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Text(
            (pkg.status as String).toUpperCase(),
            style: TextStyle(
              color: statusColor,
              fontSize: 11,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ),
    );
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'active':
        return UrjaColors.success;
      case 'pending':
        return UrjaColors.warning;
      case 'expired':
      case 'cancelled':
        return UrjaColors.error;
      default:
        return UrjaColors.foregroundMuted;
    }
  }
}
