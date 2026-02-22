import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/nfc_service.dart';
import '../../models/nfc.dart';
import '../shared/widgets/status_badge.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';

class NfcScreen extends ConsumerStatefulWidget {
  const NfcScreen({super.key});

  @override
  ConsumerState<NfcScreen> createState() => _NfcScreenState();
}

class _NfcScreenState extends ConsumerState<NfcScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  bool _loadingCards = true;
  bool _loadingDevices = true;
  List<NfcCard> _cards = [];
  List<NfcDevice> _devices = [];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _loadCards();
    _loadDevices();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadCards() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loadingCards = true);
    final nfcService = NfcApiService(ref.read(apiClientProvider));

    try {
      final result = await nfcService.getCards(orgId);
      if (mounted) {
        setState(() {
          _cards = result.data;
          _loadingCards = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loadingCards = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load NFC cards: $e')),
        );
      }
    }
  }

  Future<void> _loadDevices() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loadingDevices = true);
    final nfcService = NfcApiService(ref.read(apiClientProvider));

    try {
      final devices = await nfcService.getDevices(orgId);
      if (mounted) {
        setState(() {
          _devices = devices;
          _loadingDevices = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loadingDevices = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load NFC devices: $e')),
        );
      }
    }
  }

  Future<void> _showRegisterCardDialog() async {
    final l10n = AppLocalizations.of(context)!;
    final cardNumberCtrl = TextEditingController();

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.registerCard),
        content: TextField(
          controller: cardNumberCtrl,
          decoration: InputDecoration(
            labelText: l10n.cardNumber,
            hintText: l10n.enterCardNumber,
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text(l10n.cancel),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            child: Text(l10n.register),
          ),
        ],
      ),
    );

    if (result == true && cardNumberCtrl.text.trim().isNotEmpty && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final nfcService = NfcApiService(ref.read(apiClientProvider));
      try {
        await nfcService.registerCard(orgId, cardNumberCtrl.text.trim());
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(l10n.cardRegistered),
              backgroundColor: AppTheme.primary,
            ),
          );
          _loadCards();
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to register card: $e')),
          );
        }
      }
    }
  }

  Future<void> _showAssignDialog(NfcCard card) async {
    final l10n = AppLocalizations.of(context)!;
    final userIdCtrl = TextEditingController();

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.assignCard),
        content: TextField(
          controller: userIdCtrl,
          decoration: InputDecoration(
            labelText: l10n.memberId,
            hintText: l10n.enterMemberId,
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text(l10n.cancel),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            child: Text(l10n.assign),
          ),
        ],
      ),
    );

    if (result == true && userIdCtrl.text.trim().isNotEmpty && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final nfcService = NfcApiService(ref.read(apiClientProvider));
      try {
        await nfcService.assignCard(orgId, card.id, userIdCtrl.text.trim());
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(l10n.cardAssigned),
              backgroundColor: AppTheme.primary,
            ),
          );
          _loadCards();
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to assign card: $e')),
          );
        }
      }
    }
  }

  Future<void> _unassignCard(NfcCard card) async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    final nfcService = NfcApiService(ref.read(apiClientProvider));
    try {
      await nfcService.unassignCard(orgId, card.id);
      if (mounted) {
        final l10n = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.cardUnassigned),
            backgroundColor: AppTheme.primary,
          ),
        );
        _loadCards();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to unassign card: $e')),
        );
      }
    }
  }

  String _formatUptime(int seconds) {
    if (seconds < 60) return '${seconds}s';
    if (seconds < 3600) return '${(seconds / 60).floor()}m';
    if (seconds < 86400) return '${(seconds / 3600).floor()}h';
    return '${(seconds / 86400).floor()}d';
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.nfcCardAccess),
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: AppTheme.primary,
          labelColor: AppTheme.primary,
          unselectedLabelColor: AppTheme.textSecondary,
          tabs: [
            Tab(text: l10n.cards),
            Tab(text: l10n.devices),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildCardsTab(l10n),
          _buildDevicesTab(l10n),
        ],
      ),
    );
  }

  Widget _buildCardsTab(AppLocalizations l10n) {
    if (_loadingCards) return const LoadingIndicator();

    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(16),
          child: SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: _showRegisterCardDialog,
              icon: const Icon(Icons.add_card),
              label: Text(l10n.registerCard),
            ),
          ),
        ),
        if (_cards.isEmpty)
          Expanded(
            child: EmptyState(
              icon: Icons.nfc_outlined,
              message: l10n.noCards,
            ),
          )
        else
          Expanded(
            child: RefreshIndicator(
              onRefresh: _loadCards,
              child: ListView.builder(
                itemCount: _cards.length,
                padding: const EdgeInsets.symmetric(horizontal: 16),
                itemBuilder: (context, index) {
                  final card = _cards[index];
                  return Card(
                    margin: const EdgeInsets.only(bottom: 8),
                    child: ListTile(
                      leading: const CircleAvatar(
                        backgroundColor: AppTheme.info,
                        child: Icon(Icons.nfc, color: Colors.white),
                      ),
                      title: Text(card.cardNumber),
                      subtitle: Text(
                        card.fullName ?? l10n.unassigned,
                        style:
                            const TextStyle(color: AppTheme.textSecondary),
                      ),
                      trailing: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          StatusBadge(label: card.status),
                          const SizedBox(width: 8),
                          if (card.assignedMember != null)
                            IconButton(
                              icon: const Icon(Icons.link_off,
                                  color: AppTheme.warning, size: 20),
                              onPressed: () => _unassignCard(card),
                              tooltip: l10n.unassign,
                            )
                          else
                            IconButton(
                              icon: const Icon(Icons.link,
                                  color: AppTheme.primary, size: 20),
                              onPressed: () => _showAssignDialog(card),
                              tooltip: l10n.assign,
                            ),
                        ],
                      ),
                    ),
                  );
                },
              ),
            ),
          ),
      ],
    );
  }

  Widget _buildDevicesTab(AppLocalizations l10n) {
    if (_loadingDevices) return const LoadingIndicator();

    if (_devices.isEmpty) {
      return EmptyState(
        icon: Icons.devices_outlined,
        message: l10n.noDevices,
      );
    }

    return RefreshIndicator(
      onRefresh: _loadDevices,
      child: ListView.builder(
        itemCount: _devices.length,
        padding: const EdgeInsets.all(16),
        itemBuilder: (context, index) {
          final device = _devices[index];
          final isOnline = device.status == 'online';

          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(
                        Icons.router,
                        color: isOnline ? AppTheme.primary : AppTheme.error,
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              device.name,
                              style: const TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 15,
                              ),
                            ),
                            Text(
                              device.deviceIdentifier,
                              style: const TextStyle(
                                  color: AppTheme.textSecondary, fontSize: 12),
                            ),
                          ],
                        ),
                      ),
                      StatusBadge(label: device.status),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      _buildDeviceInfo(
                        Icons.door_front_door,
                        '${l10n.door}: ${device.doorState}',
                      ),
                      const SizedBox(width: 16),
                      _buildDeviceInfo(
                        Icons.timer,
                        '${l10n.uptime}: ${_formatUptime(device.uptimeSeconds)}',
                      ),
                    ],
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildDeviceInfo(IconData icon, String text) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 16, color: AppTheme.textSecondary),
        const SizedBox(width: 4),
        Text(text,
            style:
                const TextStyle(color: AppTheme.textSecondary, fontSize: 13)),
      ],
    );
  }
}
