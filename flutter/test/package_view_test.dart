import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/models/package.dart';
import 'package:urja_flutter/screens/packages/packages_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Package View', () {
    testWidgets('shows "no packages" when list is empty', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const PackagesScreen(),
          overrides: authenticatedOverrides(packages: []),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('No packages available'), findsOneWidget);
      expect(find.byIcon(Icons.card_membership), findsOneWidget);
    });

    testWidgets('golden — empty packages', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const PackagesScreen(),
          overrides: authenticatedOverrides(packages: []),
        ),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(PackagesScreen),
        matchesGoldenFile('goldens/packages_empty.png'),
      );
    });

    testWidgets('displays active package card', (tester) async {
      final now = DateTime.now();
      final activePkg = MemberPackage(
        id: 'mp-1',
        userId: 'user-123',
        packageId: 'pkg-1',
        organizationId: 'org-1',
        startDate: now.subtract(const Duration(days: 15)),
        endDate: now.add(const Duration(days: 15)),
        amountPaid: '3000',
        status: 'active',
        createdAt: now.subtract(const Duration(days: 15)),
        packageName: 'Monthly Plan',
        orgName: 'Urja Fitness',
      );

      await tester.pumpWidget(
        buildTestApp(
          const PackagesScreen(),
          overrides: [
            ...authenticatedOverrides(packages: [activePkg]),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Monthly Plan'), findsWidgets);
      expect(find.text('Payment History'), findsOneWidget);
    });

    testWidgets('shows payment history items', (tester) async {
      final now = DateTime.now();
      final expiredPkg = MemberPackage(
        id: 'mp-old',
        userId: 'user-123',
        packageId: 'pkg-1',
        organizationId: 'org-1',
        startDate: DateTime(2024, 1, 1),
        endDate: DateTime(2024, 2, 1),
        amountPaid: '2500',
        status: 'expired',
        createdAt: DateTime(2024, 1, 1),
        packageName: 'January Plan',
      );
      final activePkg = MemberPackage(
        id: 'mp-1',
        userId: 'user-123',
        packageId: 'pkg-2',
        organizationId: 'org-1',
        startDate: now.subtract(const Duration(days: 10)),
        endDate: now.add(const Duration(days: 20)),
        amountPaid: '3000',
        status: 'active',
        createdAt: now.subtract(const Duration(days: 10)),
        packageName: 'Monthly Plan',
      );

      await tester.pumpWidget(
        buildTestApp(
          const PackagesScreen(),
          overrides: authenticatedOverrides(
            packages: [activePkg, expiredPkg],
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Both packages should appear in payment history
      expect(find.text('Monthly Plan'), findsWidgets);
      expect(find.text('January Plan'), findsWidgets);
      expect(find.text('EXPIRED'), findsOneWidget);
      expect(find.text('ACTIVE'), findsOneWidget);
    });

    testWidgets('shows renew button when < 7 days left', (tester) async {
      final now = DateTime.now();
      final expiringPkg = MemberPackage(
        id: 'mp-1',
        userId: 'user-123',
        packageId: 'pkg-1',
        organizationId: 'org-1',
        startDate: now.subtract(const Duration(days: 25)),
        endDate: now.add(const Duration(days: 5)),
        amountPaid: '3000',
        status: 'active',
        createdAt: now.subtract(const Duration(days: 25)),
        packageName: 'Monthly Plan',
      );

      await tester.pumpWidget(
        buildTestApp(
          const PackagesScreen(),
          overrides: authenticatedOverrides(packages: [expiringPkg]),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Renew Package'), findsOneWidget);
    });
  });
}
