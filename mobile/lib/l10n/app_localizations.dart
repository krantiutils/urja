import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_ne.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
    : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations? of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
        delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
      ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('en'),
    Locale('ne'),
  ];

  /// No description provided for @appName.
  ///
  /// In en, this message translates to:
  /// **'Urja'**
  String get appName;

  /// No description provided for @loading.
  ///
  /// In en, this message translates to:
  /// **'Loading...'**
  String get loading;

  /// No description provided for @error.
  ///
  /// In en, this message translates to:
  /// **'Something went wrong'**
  String get error;

  /// No description provided for @retry.
  ///
  /// In en, this message translates to:
  /// **'Retry'**
  String get retry;

  /// No description provided for @save.
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get save;

  /// No description provided for @cancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get cancel;

  /// No description provided for @delete.
  ///
  /// In en, this message translates to:
  /// **'Delete'**
  String get delete;

  /// No description provided for @search.
  ///
  /// In en, this message translates to:
  /// **'Search...'**
  String get search;

  /// No description provided for @noResults.
  ///
  /// In en, this message translates to:
  /// **'No results found'**
  String get noResults;

  /// No description provided for @signIn.
  ///
  /// In en, this message translates to:
  /// **'Sign In'**
  String get signIn;

  /// No description provided for @signOut.
  ///
  /// In en, this message translates to:
  /// **'Sign Out'**
  String get signOut;

  /// No description provided for @phoneLabel.
  ///
  /// In en, this message translates to:
  /// **'Mobile Number'**
  String get phoneLabel;

  /// No description provided for @phonePlaceholder.
  ///
  /// In en, this message translates to:
  /// **'98XXXXXXXX'**
  String get phonePlaceholder;

  /// No description provided for @phoneNumber.
  ///
  /// In en, this message translates to:
  /// **'Phone Number'**
  String get phoneNumber;

  /// No description provided for @otpLabel.
  ///
  /// In en, this message translates to:
  /// **'Verification Code'**
  String get otpLabel;

  /// No description provided for @otpPlaceholder.
  ///
  /// In en, this message translates to:
  /// **'Enter 6-digit code'**
  String get otpPlaceholder;

  /// No description provided for @requestOtp.
  ///
  /// In en, this message translates to:
  /// **'Send OTP'**
  String get requestOtp;

  /// No description provided for @sendOtp.
  ///
  /// In en, this message translates to:
  /// **'Send OTP'**
  String get sendOtp;

  /// No description provided for @verifyOtp.
  ///
  /// In en, this message translates to:
  /// **'Verify'**
  String get verifyOtp;

  /// No description provided for @enterOtp.
  ///
  /// In en, this message translates to:
  /// **'Enter verification code'**
  String get enterOtp;

  /// No description provided for @otpSent.
  ///
  /// In en, this message translates to:
  /// **'OTP sent to your phone'**
  String get otpSent;

  /// No description provided for @invalidPhone.
  ///
  /// In en, this message translates to:
  /// **'Enter a valid Nepali mobile number'**
  String get invalidPhone;

  /// No description provided for @invalidOtp.
  ///
  /// In en, this message translates to:
  /// **'Invalid verification code'**
  String get invalidOtp;

  /// No description provided for @welcome.
  ///
  /// In en, this message translates to:
  /// **'Welcome to Urja'**
  String get welcome;

  /// No description provided for @subtitle.
  ///
  /// In en, this message translates to:
  /// **'Gym management, simplified'**
  String get subtitle;

  /// No description provided for @resendOtp.
  ///
  /// In en, this message translates to:
  /// **'Resend OTP'**
  String get resendOtp;

  /// No description provided for @resend.
  ///
  /// In en, this message translates to:
  /// **'Resend'**
  String get resend;

  /// No description provided for @resendIn.
  ///
  /// In en, this message translates to:
  /// **'Resend in {seconds}s'**
  String resendIn(int seconds);

  /// No description provided for @didNotReceiveCode.
  ///
  /// In en, this message translates to:
  /// **'Didn\'t receive the code?'**
  String get didNotReceiveCode;

  /// No description provided for @dashboard.
  ///
  /// In en, this message translates to:
  /// **'Dashboard'**
  String get dashboard;

  /// No description provided for @members.
  ///
  /// In en, this message translates to:
  /// **'Members'**
  String get members;

  /// No description provided for @packages.
  ///
  /// In en, this message translates to:
  /// **'Packages'**
  String get packages;

  /// No description provided for @duePayments.
  ///
  /// In en, this message translates to:
  /// **'Due Payments'**
  String get duePayments;

  /// No description provided for @attendance.
  ///
  /// In en, this message translates to:
  /// **'Attendance'**
  String get attendance;

  /// No description provided for @staff.
  ///
  /// In en, this message translates to:
  /// **'Staff'**
  String get staff;

  /// No description provided for @accounts.
  ///
  /// In en, this message translates to:
  /// **'Accounts'**
  String get accounts;

  /// No description provided for @stories.
  ///
  /// In en, this message translates to:
  /// **'Stories'**
  String get stories;

  /// No description provided for @sms.
  ///
  /// In en, this message translates to:
  /// **'SMS'**
  String get sms;

  /// No description provided for @feedbacks.
  ///
  /// In en, this message translates to:
  /// **'Feedbacks'**
  String get feedbacks;

  /// No description provided for @nfcAccess.
  ///
  /// In en, this message translates to:
  /// **'NFC Card Access'**
  String get nfcAccess;

  /// No description provided for @settings.
  ///
  /// In en, this message translates to:
  /// **'Settings'**
  String get settings;

  /// No description provided for @totalMembers.
  ///
  /// In en, this message translates to:
  /// **'Total Members'**
  String get totalMembers;

  /// No description provided for @activeMembers.
  ///
  /// In en, this message translates to:
  /// **'Active Members'**
  String get activeMembers;

  /// No description provided for @todayAttendance.
  ///
  /// In en, this message translates to:
  /// **'Today\'s Attendance'**
  String get todayAttendance;

  /// No description provided for @monthlyRevenue.
  ///
  /// In en, this message translates to:
  /// **'Monthly Revenue'**
  String get monthlyRevenue;

  /// No description provided for @expiringPackages.
  ///
  /// In en, this message translates to:
  /// **'Expiring Packages'**
  String get expiringPackages;

  /// No description provided for @expiredPackages.
  ///
  /// In en, this message translates to:
  /// **'Expired Packages'**
  String get expiredPackages;

  /// No description provided for @noExpiringPackages.
  ///
  /// In en, this message translates to:
  /// **'No expiring packages'**
  String get noExpiringPackages;

  /// No description provided for @todayActivity.
  ///
  /// In en, this message translates to:
  /// **'Today\'s Activity'**
  String get todayActivity;

  /// No description provided for @recentActivity.
  ///
  /// In en, this message translates to:
  /// **'Recent Activity'**
  String get recentActivity;

  /// No description provided for @noActivityToday.
  ///
  /// In en, this message translates to:
  /// **'No activity today'**
  String get noActivityToday;

  /// No description provided for @unknownMember.
  ///
  /// In en, this message translates to:
  /// **'Unknown Member'**
  String get unknownMember;

  /// No description provided for @viewAll.
  ///
  /// In en, this message translates to:
  /// **'View All'**
  String get viewAll;

  /// No description provided for @addMember.
  ///
  /// In en, this message translates to:
  /// **'Add Member'**
  String get addMember;

  /// No description provided for @editMember.
  ///
  /// In en, this message translates to:
  /// **'Edit Member'**
  String get editMember;

  /// No description provided for @removeMember.
  ///
  /// In en, this message translates to:
  /// **'Remove Member'**
  String get removeMember;

  /// No description provided for @name.
  ///
  /// In en, this message translates to:
  /// **'Name'**
  String get name;

  /// No description provided for @nameNe.
  ///
  /// In en, this message translates to:
  /// **'Name (Nepali)'**
  String get nameNe;

  /// No description provided for @phone.
  ///
  /// In en, this message translates to:
  /// **'Phone'**
  String get phone;

  /// No description provided for @email.
  ///
  /// In en, this message translates to:
  /// **'Email'**
  String get email;

  /// No description provided for @status.
  ///
  /// In en, this message translates to:
  /// **'Status'**
  String get status;

  /// No description provided for @role.
  ///
  /// In en, this message translates to:
  /// **'Role'**
  String get role;

  /// No description provided for @active.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get active;

  /// No description provided for @inactive.
  ///
  /// In en, this message translates to:
  /// **'Inactive'**
  String get inactive;

  /// No description provided for @suspended.
  ///
  /// In en, this message translates to:
  /// **'Suspended'**
  String get suspended;

  /// No description provided for @left.
  ///
  /// In en, this message translates to:
  /// **'Left'**
  String get left;

  /// No description provided for @member.
  ///
  /// In en, this message translates to:
  /// **'Member'**
  String get member;

  /// No description provided for @confirmDelete.
  ///
  /// In en, this message translates to:
  /// **'Are you sure you want to delete this?'**
  String get confirmDelete;

  /// No description provided for @noMembers.
  ///
  /// In en, this message translates to:
  /// **'No members found'**
  String get noMembers;

  /// No description provided for @addPackage.
  ///
  /// In en, this message translates to:
  /// **'Add Package'**
  String get addPackage;

  /// No description provided for @editPackage.
  ///
  /// In en, this message translates to:
  /// **'Edit Package'**
  String get editPackage;

  /// No description provided for @price.
  ///
  /// In en, this message translates to:
  /// **'Price'**
  String get price;

  /// No description provided for @duration.
  ///
  /// In en, this message translates to:
  /// **'Duration'**
  String get duration;

  /// No description provided for @features.
  ///
  /// In en, this message translates to:
  /// **'Features'**
  String get features;

  /// No description provided for @days.
  ///
  /// In en, this message translates to:
  /// **'days'**
  String get days;

  /// No description provided for @maxMembers.
  ///
  /// In en, this message translates to:
  /// **'Max Members'**
  String get maxMembers;

  /// No description provided for @unlimited.
  ///
  /// In en, this message translates to:
  /// **'Unlimited'**
  String get unlimited;

  /// No description provided for @description.
  ///
  /// In en, this message translates to:
  /// **'Description'**
  String get description;

  /// No description provided for @filterAll.
  ///
  /// In en, this message translates to:
  /// **'All'**
  String get filterAll;

  /// No description provided for @filterActive.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get filterActive;

  /// No description provided for @filterInactive.
  ///
  /// In en, this message translates to:
  /// **'Inactive'**
  String get filterInactive;

  /// No description provided for @todayCheckIns.
  ///
  /// In en, this message translates to:
  /// **'Today\'s Check-ins'**
  String get todayCheckIns;

  /// No description provided for @manualCheckIn.
  ///
  /// In en, this message translates to:
  /// **'Manual Check-in'**
  String get manualCheckIn;

  /// No description provided for @checkInMember.
  ///
  /// In en, this message translates to:
  /// **'Check in a member'**
  String get checkInMember;

  /// No description provided for @method.
  ///
  /// In en, this message translates to:
  /// **'Method'**
  String get method;

  /// No description provided for @qr.
  ///
  /// In en, this message translates to:
  /// **'QR'**
  String get qr;

  /// No description provided for @nfc.
  ///
  /// In en, this message translates to:
  /// **'NFC'**
  String get nfc;

  /// No description provided for @manual.
  ///
  /// In en, this message translates to:
  /// **'Manual'**
  String get manual;

  /// No description provided for @noRecords.
  ///
  /// In en, this message translates to:
  /// **'No records found'**
  String get noRecords;

  /// No description provided for @addStaff.
  ///
  /// In en, this message translates to:
  /// **'Add Staff'**
  String get addStaff;

  /// No description provided for @editStaff.
  ///
  /// In en, this message translates to:
  /// **'Edit Staff'**
  String get editStaff;

  /// No description provided for @staffRole.
  ///
  /// In en, this message translates to:
  /// **'Staff Role'**
  String get staffRole;

  /// No description provided for @owner.
  ///
  /// In en, this message translates to:
  /// **'Owner'**
  String get owner;

  /// No description provided for @manager.
  ///
  /// In en, this message translates to:
  /// **'Manager'**
  String get manager;

  /// No description provided for @trainer.
  ///
  /// In en, this message translates to:
  /// **'Trainer'**
  String get trainer;

  /// No description provided for @receptionist.
  ///
  /// In en, this message translates to:
  /// **'Receptionist'**
  String get receptionist;

  /// No description provided for @amount.
  ///
  /// In en, this message translates to:
  /// **'Amount'**
  String get amount;

  /// No description provided for @dueDate.
  ///
  /// In en, this message translates to:
  /// **'Due Date'**
  String get dueDate;

  /// No description provided for @unpaid.
  ///
  /// In en, this message translates to:
  /// **'Unpaid'**
  String get unpaid;

  /// No description provided for @paid.
  ///
  /// In en, this message translates to:
  /// **'Paid'**
  String get paid;

  /// No description provided for @waived.
  ///
  /// In en, this message translates to:
  /// **'Waived'**
  String get waived;

  /// No description provided for @recordPayment.
  ///
  /// In en, this message translates to:
  /// **'Record Payment'**
  String get recordPayment;

  /// No description provided for @paymentMethod.
  ///
  /// In en, this message translates to:
  /// **'Payment Method'**
  String get paymentMethod;

  /// No description provided for @cash.
  ///
  /// In en, this message translates to:
  /// **'Cash'**
  String get cash;

  /// No description provided for @esewa.
  ///
  /// In en, this message translates to:
  /// **'eSewa'**
  String get esewa;

  /// No description provided for @bankTransfer.
  ///
  /// In en, this message translates to:
  /// **'Bank Transfer'**
  String get bankTransfer;

  /// No description provided for @confirmPayment.
  ///
  /// In en, this message translates to:
  /// **'Confirm Payment'**
  String get confirmPayment;

  /// No description provided for @noDues.
  ///
  /// In en, this message translates to:
  /// **'No dues found'**
  String get noDues;

  /// No description provided for @addTransaction.
  ///
  /// In en, this message translates to:
  /// **'Add Transaction'**
  String get addTransaction;

  /// No description provided for @editTransaction.
  ///
  /// In en, this message translates to:
  /// **'Edit Transaction'**
  String get editTransaction;

  /// No description provided for @category.
  ///
  /// In en, this message translates to:
  /// **'Category'**
  String get category;

  /// No description provided for @type.
  ///
  /// In en, this message translates to:
  /// **'Type'**
  String get type;

  /// No description provided for @paymentType.
  ///
  /// In en, this message translates to:
  /// **'Payment Type'**
  String get paymentType;

  /// No description provided for @reference.
  ///
  /// In en, this message translates to:
  /// **'Reference'**
  String get reference;

  /// No description provided for @income.
  ///
  /// In en, this message translates to:
  /// **'Income'**
  String get income;

  /// No description provided for @expense.
  ///
  /// In en, this message translates to:
  /// **'Expense'**
  String get expense;

  /// No description provided for @totalIncome.
  ///
  /// In en, this message translates to:
  /// **'Total Income'**
  String get totalIncome;

  /// No description provided for @totalExpenses.
  ///
  /// In en, this message translates to:
  /// **'Total Expenses'**
  String get totalExpenses;

  /// No description provided for @grossProfit.
  ///
  /// In en, this message translates to:
  /// **'Gross Profit'**
  String get grossProfit;

  /// No description provided for @profitPercent.
  ///
  /// In en, this message translates to:
  /// **'Profit %'**
  String get profitPercent;

  /// No description provided for @noTransactions.
  ///
  /// In en, this message translates to:
  /// **'No transactions found'**
  String get noTransactions;

  /// No description provided for @addNotice.
  ///
  /// In en, this message translates to:
  /// **'Add Notice'**
  String get addNotice;

  /// No description provided for @editNotice.
  ///
  /// In en, this message translates to:
  /// **'Edit Notice'**
  String get editNotice;

  /// No description provided for @noticeTitle.
  ///
  /// In en, this message translates to:
  /// **'Title'**
  String get noticeTitle;

  /// No description provided for @titleNe.
  ///
  /// In en, this message translates to:
  /// **'Title (Nepali)'**
  String get titleNe;

  /// No description provided for @content.
  ///
  /// In en, this message translates to:
  /// **'Content'**
  String get content;

  /// No description provided for @contentNe.
  ///
  /// In en, this message translates to:
  /// **'Content (Nepali)'**
  String get contentNe;

  /// No description provided for @noNotices.
  ///
  /// In en, this message translates to:
  /// **'No notices found'**
  String get noNotices;

  /// No description provided for @noFeedbacks.
  ///
  /// In en, this message translates to:
  /// **'No feedbacks yet'**
  String get noFeedbacks;

  /// No description provided for @cards.
  ///
  /// In en, this message translates to:
  /// **'Cards'**
  String get cards;

  /// No description provided for @devices.
  ///
  /// In en, this message translates to:
  /// **'Devices'**
  String get devices;

  /// No description provided for @registerCard.
  ///
  /// In en, this message translates to:
  /// **'Register Card'**
  String get registerCard;

  /// No description provided for @cardNumber.
  ///
  /// In en, this message translates to:
  /// **'Card Number'**
  String get cardNumber;

  /// No description provided for @assignedTo.
  ///
  /// In en, this message translates to:
  /// **'Assigned To'**
  String get assignedTo;

  /// No description provided for @assign.
  ///
  /// In en, this message translates to:
  /// **'Assign'**
  String get assign;

  /// No description provided for @unassign.
  ///
  /// In en, this message translates to:
  /// **'Unassign'**
  String get unassign;

  /// No description provided for @unassigned.
  ///
  /// In en, this message translates to:
  /// **'Unassigned'**
  String get unassigned;

  /// No description provided for @available.
  ///
  /// In en, this message translates to:
  /// **'Available'**
  String get available;

  /// No description provided for @assigned.
  ///
  /// In en, this message translates to:
  /// **'Assigned'**
  String get assigned;

  /// No description provided for @disabled.
  ///
  /// In en, this message translates to:
  /// **'Disabled'**
  String get disabled;

  /// No description provided for @noCards.
  ///
  /// In en, this message translates to:
  /// **'No NFC cards registered'**
  String get noCards;

  /// No description provided for @deviceName.
  ///
  /// In en, this message translates to:
  /// **'Device Name'**
  String get deviceName;

  /// No description provided for @online.
  ///
  /// In en, this message translates to:
  /// **'Online'**
  String get online;

  /// No description provided for @offline.
  ///
  /// In en, this message translates to:
  /// **'Offline'**
  String get offline;

  /// No description provided for @locked.
  ///
  /// In en, this message translates to:
  /// **'Locked'**
  String get locked;

  /// No description provided for @unlocked.
  ///
  /// In en, this message translates to:
  /// **'Unlocked'**
  String get unlocked;

  /// No description provided for @noDevices.
  ///
  /// In en, this message translates to:
  /// **'No NFC devices found'**
  String get noDevices;

  /// No description provided for @door.
  ///
  /// In en, this message translates to:
  /// **'Door'**
  String get door;

  /// No description provided for @uptime.
  ///
  /// In en, this message translates to:
  /// **'Uptime'**
  String get uptime;

  /// No description provided for @smsBalance.
  ///
  /// In en, this message translates to:
  /// **'SMS Balance'**
  String get smsBalance;

  /// No description provided for @remaining.
  ///
  /// In en, this message translates to:
  /// **'Remaining'**
  String get remaining;

  /// No description provided for @totalPurchased.
  ///
  /// In en, this message translates to:
  /// **'Total Purchased'**
  String get totalPurchased;

  /// No description provided for @totalUsed.
  ///
  /// In en, this message translates to:
  /// **'Total Used'**
  String get totalUsed;

  /// No description provided for @campaigns.
  ///
  /// In en, this message translates to:
  /// **'Campaigns'**
  String get campaigns;

  /// No description provided for @sendSms.
  ///
  /// In en, this message translates to:
  /// **'Send SMS'**
  String get sendSms;

  /// No description provided for @message.
  ///
  /// In en, this message translates to:
  /// **'Message'**
  String get message;

  /// No description provided for @enterMessage.
  ///
  /// In en, this message translates to:
  /// **'Enter your message'**
  String get enterMessage;

  /// No description provided for @selectMembers.
  ///
  /// In en, this message translates to:
  /// **'Select Members'**
  String get selectMembers;

  /// No description provided for @membersSelected.
  ///
  /// In en, this message translates to:
  /// **'{count} members selected'**
  String membersSelected(int count);

  /// No description provided for @send.
  ///
  /// In en, this message translates to:
  /// **'Send'**
  String get send;

  /// No description provided for @history.
  ///
  /// In en, this message translates to:
  /// **'History'**
  String get history;

  /// No description provided for @noHistory.
  ///
  /// In en, this message translates to:
  /// **'No history'**
  String get noHistory;

  /// No description provided for @orgProfile.
  ///
  /// In en, this message translates to:
  /// **'Organization Profile'**
  String get orgProfile;

  /// No description provided for @myProfile.
  ///
  /// In en, this message translates to:
  /// **'My Profile'**
  String get myProfile;

  /// No description provided for @profileInfo.
  ///
  /// In en, this message translates to:
  /// **'Profile Information'**
  String get profileInfo;

  /// No description provided for @noProfile.
  ///
  /// In en, this message translates to:
  /// **'No profile data'**
  String get noProfile;

  /// No description provided for @privacySettings.
  ///
  /// In en, this message translates to:
  /// **'Privacy Settings'**
  String get privacySettings;

  /// No description provided for @showEmail.
  ///
  /// In en, this message translates to:
  /// **'Show email to others'**
  String get showEmail;

  /// No description provided for @showEmailDesc.
  ///
  /// In en, this message translates to:
  /// **'Allow other members to see your email'**
  String get showEmailDesc;

  /// No description provided for @showPhone.
  ///
  /// In en, this message translates to:
  /// **'Show phone to others'**
  String get showPhone;

  /// No description provided for @showPhoneDesc.
  ///
  /// In en, this message translates to:
  /// **'Allow other members to see your phone'**
  String get showPhoneDesc;

  /// No description provided for @showProfile.
  ///
  /// In en, this message translates to:
  /// **'Show profile publicly'**
  String get showProfile;

  /// No description provided for @showProfileDesc.
  ///
  /// In en, this message translates to:
  /// **'Make your profile visible to others'**
  String get showProfileDesc;

  /// No description provided for @showAttendance.
  ///
  /// In en, this message translates to:
  /// **'Show attendance to others'**
  String get showAttendance;

  /// No description provided for @showAttendanceDesc.
  ///
  /// In en, this message translates to:
  /// **'Allow others to see your attendance'**
  String get showAttendanceDesc;

  /// No description provided for @showOnLeaderboard.
  ///
  /// In en, this message translates to:
  /// **'Show on leaderboard'**
  String get showOnLeaderboard;

  /// No description provided for @showOnLeaderboardDesc.
  ///
  /// In en, this message translates to:
  /// **'Display on the gym leaderboard'**
  String get showOnLeaderboardDesc;

  /// No description provided for @emergencyContact.
  ///
  /// In en, this message translates to:
  /// **'Emergency Contact'**
  String get emergencyContact;

  /// No description provided for @contactName.
  ///
  /// In en, this message translates to:
  /// **'Contact Name'**
  String get contactName;

  /// No description provided for @contactPhone.
  ///
  /// In en, this message translates to:
  /// **'Contact Phone'**
  String get contactPhone;

  /// No description provided for @saveChanges.
  ///
  /// In en, this message translates to:
  /// **'Save Changes'**
  String get saveChanges;

  /// No description provided for @saveProfile.
  ///
  /// In en, this message translates to:
  /// **'Save Profile'**
  String get saveProfile;

  /// No description provided for @savePrivacy.
  ///
  /// In en, this message translates to:
  /// **'Save Privacy'**
  String get savePrivacy;

  /// No description provided for @saveEmergencyContact.
  ///
  /// In en, this message translates to:
  /// **'Save Emergency Contact'**
  String get saveEmergencyContact;

  /// No description provided for @saved.
  ///
  /// In en, this message translates to:
  /// **'Changes saved'**
  String get saved;

  /// No description provided for @myDashboard.
  ///
  /// In en, this message translates to:
  /// **'My Dashboard'**
  String get myDashboard;

  /// No description provided for @myAttendance.
  ///
  /// In en, this message translates to:
  /// **'My Attendance'**
  String get myAttendance;

  /// No description provided for @myPackages.
  ///
  /// In en, this message translates to:
  /// **'My Packages'**
  String get myPackages;

  /// No description provided for @myHealth.
  ///
  /// In en, this message translates to:
  /// **'My Health'**
  String get myHealth;

  /// No description provided for @myWorkouts.
  ///
  /// In en, this message translates to:
  /// **'My Workouts'**
  String get myWorkouts;

  /// No description provided for @profile.
  ///
  /// In en, this message translates to:
  /// **'Profile'**
  String get profile;

  /// No description provided for @welcomeBack.
  ///
  /// In en, this message translates to:
  /// **'Welcome back'**
  String get welcomeBack;

  /// No description provided for @memberSince.
  ///
  /// In en, this message translates to:
  /// **'Member since'**
  String get memberSince;

  /// No description provided for @currentPackage.
  ///
  /// In en, this message translates to:
  /// **'Current Package'**
  String get currentPackage;

  /// No description provided for @activePackage.
  ///
  /// In en, this message translates to:
  /// **'Active Package'**
  String get activePackage;

  /// No description provided for @noActivePackage.
  ///
  /// In en, this message translates to:
  /// **'No active package'**
  String get noActivePackage;

  /// No description provided for @noActivePackages.
  ///
  /// In en, this message translates to:
  /// **'No active packages'**
  String get noActivePackages;

  /// No description provided for @noPastPackages.
  ///
  /// In en, this message translates to:
  /// **'No past packages'**
  String get noPastPackages;

  /// No description provided for @daysRemaining.
  ///
  /// In en, this message translates to:
  /// **'days remaining'**
  String get daysRemaining;

  /// No description provided for @expired.
  ///
  /// In en, this message translates to:
  /// **'Expired'**
  String get expired;

  /// No description provided for @expires.
  ///
  /// In en, this message translates to:
  /// **'Expires'**
  String get expires;

  /// No description provided for @currentStreak.
  ///
  /// In en, this message translates to:
  /// **'Current Streak'**
  String get currentStreak;

  /// No description provided for @longestStreak.
  ///
  /// In en, this message translates to:
  /// **'Longest Streak'**
  String get longestStreak;

  /// No description provided for @recentCheckIns.
  ///
  /// In en, this message translates to:
  /// **'Recent Check-ins'**
  String get recentCheckIns;

  /// No description provided for @noRecentCheckIns.
  ///
  /// In en, this message translates to:
  /// **'No recent check-ins'**
  String get noRecentCheckIns;

  /// No description provided for @noCheckIns.
  ///
  /// In en, this message translates to:
  /// **'No check-ins yet'**
  String get noCheckIns;

  /// No description provided for @noAttendanceRecords.
  ///
  /// In en, this message translates to:
  /// **'No attendance records'**
  String get noAttendanceRecords;

  /// No description provided for @attendanceHistory.
  ///
  /// In en, this message translates to:
  /// **'Attendance History'**
  String get attendanceHistory;

  /// No description provided for @bmi.
  ///
  /// In en, this message translates to:
  /// **'BMI'**
  String get bmi;

  /// No description provided for @measurements.
  ///
  /// In en, this message translates to:
  /// **'Body Measurements'**
  String get measurements;

  /// No description provided for @bodyMeasurements.
  ///
  /// In en, this message translates to:
  /// **'Body Measurements'**
  String get bodyMeasurements;

  /// No description provided for @noHealthData.
  ///
  /// In en, this message translates to:
  /// **'No health data recorded yet'**
  String get noHealthData;

  /// No description provided for @noBmiData.
  ///
  /// In en, this message translates to:
  /// **'No BMI data'**
  String get noBmiData;

  /// No description provided for @noMeasurementData.
  ///
  /// In en, this message translates to:
  /// **'No measurement data'**
  String get noMeasurementData;

  /// No description provided for @height.
  ///
  /// In en, this message translates to:
  /// **'Height'**
  String get height;

  /// No description provided for @weight.
  ///
  /// In en, this message translates to:
  /// **'Weight'**
  String get weight;

  /// No description provided for @chest.
  ///
  /// In en, this message translates to:
  /// **'Chest'**
  String get chest;

  /// No description provided for @waist.
  ///
  /// In en, this message translates to:
  /// **'Waist'**
  String get waist;

  /// No description provided for @hips.
  ///
  /// In en, this message translates to:
  /// **'Hips'**
  String get hips;

  /// No description provided for @biceps.
  ///
  /// In en, this message translates to:
  /// **'Biceps'**
  String get biceps;

  /// No description provided for @thighs.
  ///
  /// In en, this message translates to:
  /// **'Thighs'**
  String get thighs;

  /// No description provided for @bodyFat.
  ///
  /// In en, this message translates to:
  /// **'Body Fat'**
  String get bodyFat;

  /// No description provided for @recordedOn.
  ///
  /// In en, this message translates to:
  /// **'Recorded on'**
  String get recordedOn;

  /// No description provided for @workoutPlan.
  ///
  /// In en, this message translates to:
  /// **'Assigned Plan'**
  String get workoutPlan;

  /// No description provided for @assignedPlan.
  ///
  /// In en, this message translates to:
  /// **'Assigned Plan'**
  String get assignedPlan;

  /// No description provided for @noAssignedPlan.
  ///
  /// In en, this message translates to:
  /// **'No assigned plan'**
  String get noAssignedPlan;

  /// No description provided for @recentWorkouts.
  ///
  /// In en, this message translates to:
  /// **'Recent Workouts'**
  String get recentWorkouts;

  /// No description provided for @recentWorkoutLogs.
  ///
  /// In en, this message translates to:
  /// **'Recent Workout Logs'**
  String get recentWorkoutLogs;

  /// No description provided for @noWorkouts.
  ///
  /// In en, this message translates to:
  /// **'No workouts logged yet'**
  String get noWorkouts;

  /// No description provided for @noWorkoutLogs.
  ///
  /// In en, this message translates to:
  /// **'No workout logs'**
  String get noWorkoutLogs;

  /// No description provided for @noPlan.
  ///
  /// In en, this message translates to:
  /// **'No workout plan assigned'**
  String get noPlan;

  /// No description provided for @minutes.
  ///
  /// In en, this message translates to:
  /// **'min'**
  String get minutes;

  /// No description provided for @scanQr.
  ///
  /// In en, this message translates to:
  /// **'Scan QR'**
  String get scanQr;

  /// No description provided for @tapNfc.
  ///
  /// In en, this message translates to:
  /// **'Tap NFC'**
  String get tapNfc;

  /// No description provided for @checkIn.
  ///
  /// In en, this message translates to:
  /// **'Check In'**
  String get checkIn;

  /// No description provided for @checkInSuccess.
  ///
  /// In en, this message translates to:
  /// **'Check-in successful!'**
  String get checkInSuccess;

  /// No description provided for @checkInFailed.
  ///
  /// In en, this message translates to:
  /// **'Check-in failed'**
  String get checkInFailed;

  /// No description provided for @streak.
  ///
  /// In en, this message translates to:
  /// **'Streak'**
  String get streak;

  /// No description provided for @scanQrCode.
  ///
  /// In en, this message translates to:
  /// **'Scan QR Code to Check In'**
  String get scanQrCode;

  /// No description provided for @tapNfcCard.
  ///
  /// In en, this message translates to:
  /// **'Hold your phone near the NFC card'**
  String get tapNfcCard;

  /// No description provided for @language.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get language;

  /// No description provided for @english.
  ///
  /// In en, this message translates to:
  /// **'English'**
  String get english;

  /// No description provided for @nepali.
  ///
  /// In en, this message translates to:
  /// **'Nepali'**
  String get nepali;

  /// No description provided for @tapToToggleLanguage.
  ///
  /// In en, this message translates to:
  /// **'Tap to switch language'**
  String get tapToToggleLanguage;

  /// No description provided for @logout.
  ///
  /// In en, this message translates to:
  /// **'Logout'**
  String get logout;

  /// No description provided for @more.
  ///
  /// In en, this message translates to:
  /// **'More'**
  String get more;

  /// No description provided for @home.
  ///
  /// In en, this message translates to:
  /// **'Home'**
  String get home;

  /// No description provided for @date.
  ///
  /// In en, this message translates to:
  /// **'Date'**
  String get date;

  /// No description provided for @time.
  ///
  /// In en, this message translates to:
  /// **'Time'**
  String get time;

  /// No description provided for @actions.
  ///
  /// In en, this message translates to:
  /// **'Actions'**
  String get actions;

  /// No description provided for @previous.
  ///
  /// In en, this message translates to:
  /// **'Previous'**
  String get previous;

  /// No description provided for @next.
  ///
  /// In en, this message translates to:
  /// **'Next'**
  String get next;

  /// No description provided for @showing.
  ///
  /// In en, this message translates to:
  /// **'Showing'**
  String get showing;

  /// No description provided for @ofLabel.
  ///
  /// In en, this message translates to:
  /// **'of'**
  String get ofLabel;

  /// No description provided for @joined.
  ///
  /// In en, this message translates to:
  /// **'Joined'**
  String get joined;

  /// No description provided for @package.
  ///
  /// In en, this message translates to:
  /// **'Package'**
  String get package;

  /// No description provided for @address.
  ///
  /// In en, this message translates to:
  /// **'Address'**
  String get address;

  /// No description provided for @dateOfBirth.
  ///
  /// In en, this message translates to:
  /// **'Date of Birth'**
  String get dateOfBirth;

  /// No description provided for @gender.
  ///
  /// In en, this message translates to:
  /// **'Gender'**
  String get gender;

  /// No description provided for @male.
  ///
  /// In en, this message translates to:
  /// **'Male'**
  String get male;

  /// No description provided for @female.
  ///
  /// In en, this message translates to:
  /// **'Female'**
  String get female;

  /// No description provided for @other.
  ///
  /// In en, this message translates to:
  /// **'Other'**
  String get other;

  /// No description provided for @all.
  ///
  /// In en, this message translates to:
  /// **'All'**
  String get all;

  /// No description provided for @past.
  ///
  /// In en, this message translates to:
  /// **'Past'**
  String get past;

  /// No description provided for @done.
  ///
  /// In en, this message translates to:
  /// **'Done'**
  String get done;

  /// No description provided for @title.
  ///
  /// In en, this message translates to:
  /// **'Title'**
  String get title;

  /// No description provided for @welcomeToUrja.
  ///
  /// In en, this message translates to:
  /// **'Welcome to Urja'**
  String get welcomeToUrja;

  /// No description provided for @gymManagementSimplified.
  ///
  /// In en, this message translates to:
  /// **'Gym management, simplified'**
  String get gymManagementSimplified;

  /// No description provided for @nfcNotAvailable.
  ///
  /// In en, this message translates to:
  /// **'NFC is not available on this device'**
  String get nfcNotAvailable;

  /// No description provided for @nfcCardDetected.
  ///
  /// In en, this message translates to:
  /// **'NFC Card Detected'**
  String get nfcCardDetected;

  /// No description provided for @cardUid.
  ///
  /// In en, this message translates to:
  /// **'Card UID'**
  String get cardUid;

  /// No description provided for @scanAgain.
  ///
  /// In en, this message translates to:
  /// **'Scan Again'**
  String get scanAgain;

  /// No description provided for @nfcError.
  ///
  /// In en, this message translates to:
  /// **'NFC Error'**
  String get nfcError;

  /// No description provided for @tryAgain.
  ///
  /// In en, this message translates to:
  /// **'Try Again'**
  String get tryAgain;

  /// No description provided for @holdPhoneNearCard.
  ///
  /// In en, this message translates to:
  /// **'Hold your phone near the NFC card'**
  String get holdPhoneNearCard;

  /// No description provided for @nfcInstructions.
  ///
  /// In en, this message translates to:
  /// **'Make sure NFC is enabled on your device'**
  String get nfcInstructions;

  /// No description provided for @scanSuccess.
  ///
  /// In en, this message translates to:
  /// **'Scan Successful'**
  String get scanSuccess;

  /// No description provided for @qrCodeScanned.
  ///
  /// In en, this message translates to:
  /// **'QR Code Scanned'**
  String get qrCodeScanned;

  /// No description provided for @pointCameraAtQr.
  ///
  /// In en, this message translates to:
  /// **'Point camera at QR code'**
  String get pointCameraAtQr;

  /// No description provided for @qrWillAutoScan.
  ///
  /// In en, this message translates to:
  /// **'QR code will be scanned automatically'**
  String get qrWillAutoScan;

  /// No description provided for @memberId.
  ///
  /// In en, this message translates to:
  /// **'Member ID'**
  String get memberId;

  /// No description provided for @enterMemberId.
  ///
  /// In en, this message translates to:
  /// **'Enter member ID'**
  String get enterMemberId;

  /// No description provided for @noAttendance.
  ///
  /// In en, this message translates to:
  /// **'No attendance records'**
  String get noAttendance;

  /// No description provided for @paymentRecorded.
  ///
  /// In en, this message translates to:
  /// **'Payment recorded'**
  String get paymentRecorded;

  /// No description provided for @deleteFeedbackConfirm.
  ///
  /// In en, this message translates to:
  /// **'Delete this feedback?'**
  String get deleteFeedbackConfirm;

  /// No description provided for @removeMemberConfirm.
  ///
  /// In en, this message translates to:
  /// **'Remove this member?'**
  String get removeMemberConfirm;

  /// No description provided for @enterCardNumber.
  ///
  /// In en, this message translates to:
  /// **'Enter card number'**
  String get enterCardNumber;

  /// No description provided for @register.
  ///
  /// In en, this message translates to:
  /// **'Register'**
  String get register;

  /// No description provided for @cardRegistered.
  ///
  /// In en, this message translates to:
  /// **'Card registered'**
  String get cardRegistered;

  /// No description provided for @assignCard.
  ///
  /// In en, this message translates to:
  /// **'Assign Card'**
  String get assignCard;

  /// No description provided for @cardAssigned.
  ///
  /// In en, this message translates to:
  /// **'Card assigned'**
  String get cardAssigned;

  /// No description provided for @cardUnassigned.
  ///
  /// In en, this message translates to:
  /// **'Card unassigned'**
  String get cardUnassigned;

  /// No description provided for @nfcCardAccess.
  ///
  /// In en, this message translates to:
  /// **'NFC Card Access'**
  String get nfcCardAccess;

  /// No description provided for @storiesNotices.
  ///
  /// In en, this message translates to:
  /// **'Stories / Notices'**
  String get storiesNotices;

  /// No description provided for @durationDays.
  ///
  /// In en, this message translates to:
  /// **'Duration (days)'**
  String get durationDays;

  /// No description provided for @featuresHint.
  ///
  /// In en, this message translates to:
  /// **'gym_access, trainer, pool'**
  String get featuresHint;

  /// No description provided for @noPackages.
  ///
  /// In en, this message translates to:
  /// **'No packages found'**
  String get noPackages;

  /// No description provided for @noStaff.
  ///
  /// In en, this message translates to:
  /// **'No staff found'**
  String get noStaff;

  /// No description provided for @settingsSaved.
  ///
  /// In en, this message translates to:
  /// **'Settings saved'**
  String get settingsSaved;

  /// No description provided for @logoutConfirm.
  ///
  /// In en, this message translates to:
  /// **'Are you sure you want to logout?'**
  String get logoutConfirm;

  /// No description provided for @logoutConfirmation.
  ///
  /// In en, this message translates to:
  /// **'Are you sure you want to logout?'**
  String get logoutConfirmation;

  /// No description provided for @organizationProfile.
  ///
  /// In en, this message translates to:
  /// **'Organization Profile'**
  String get organizationProfile;

  /// No description provided for @personalProfile.
  ///
  /// In en, this message translates to:
  /// **'Personal Profile'**
  String get personalProfile;

  /// No description provided for @smsSent.
  ///
  /// In en, this message translates to:
  /// **'SMS sent successfully'**
  String get smsSent;

  /// No description provided for @selectAll.
  ///
  /// In en, this message translates to:
  /// **'Select All'**
  String get selectAll;

  /// No description provided for @clearAll.
  ///
  /// In en, this message translates to:
  /// **'Clear All'**
  String get clearAll;

  /// No description provided for @errorLoadingData.
  ///
  /// In en, this message translates to:
  /// **'Error loading data'**
  String get errorLoadingData;

  /// No description provided for @hello.
  ///
  /// In en, this message translates to:
  /// **'Hello'**
  String get hello;

  /// No description provided for @letsGetMoving.
  ///
  /// In en, this message translates to:
  /// **'Let\'s get you moving today!'**
  String get letsGetMoving;

  /// No description provided for @daysStreak.
  ///
  /// In en, this message translates to:
  /// **'{count} Days Streak'**
  String daysStreak(int count);

  /// No description provided for @monthlyProgress.
  ///
  /// In en, this message translates to:
  /// **'Monthly Progress'**
  String get monthlyProgress;

  /// No description provided for @top3CheckIns.
  ///
  /// In en, this message translates to:
  /// **'Top 3 Check-Ins'**
  String get top3CheckIns;

  /// No description provided for @mySubscription.
  ///
  /// In en, this message translates to:
  /// **'My Subscription'**
  String get mySubscription;

  /// No description provided for @subscriptionActive.
  ///
  /// In en, this message translates to:
  /// **'ACTIVE'**
  String get subscriptionActive;

  /// No description provided for @daysLeft.
  ///
  /// In en, this message translates to:
  /// **'{count} days left'**
  String daysLeft(int count);

  /// No description provided for @myHomeClub.
  ///
  /// In en, this message translates to:
  /// **'My Home Club'**
  String get myHomeClub;

  /// No description provided for @sendFeedback.
  ///
  /// In en, this message translates to:
  /// **'Send Us Feedback'**
  String get sendFeedback;

  /// No description provided for @feedbackHint.
  ///
  /// In en, this message translates to:
  /// **'Tell us about your experience...'**
  String get feedbackHint;

  /// No description provided for @feedbackSent.
  ///
  /// In en, this message translates to:
  /// **'Thank you for your feedback!'**
  String get feedbackSent;

  /// No description provided for @checkIns.
  ///
  /// In en, this message translates to:
  /// **'check-ins'**
  String get checkIns;

  /// No description provided for @sun.
  ///
  /// In en, this message translates to:
  /// **'S'**
  String get sun;

  /// No description provided for @mon.
  ///
  /// In en, this message translates to:
  /// **'M'**
  String get mon;

  /// No description provided for @tue.
  ///
  /// In en, this message translates to:
  /// **'T'**
  String get tue;

  /// No description provided for @wed.
  ///
  /// In en, this message translates to:
  /// **'W'**
  String get wed;

  /// No description provided for @thu.
  ///
  /// In en, this message translates to:
  /// **'T'**
  String get thu;

  /// No description provided for @fri.
  ///
  /// In en, this message translates to:
  /// **'F'**
  String get fri;

  /// No description provided for @sat.
  ///
  /// In en, this message translates to:
  /// **'S'**
  String get sat;

  /// No description provided for @account.
  ///
  /// In en, this message translates to:
  /// **'Account'**
  String get account;

  /// No description provided for @editMyProfile.
  ///
  /// In en, this message translates to:
  /// **'Edit My Profile'**
  String get editMyProfile;

  /// No description provided for @subscriptionHistory.
  ///
  /// In en, this message translates to:
  /// **'Subscription History'**
  String get subscriptionHistory;

  /// No description provided for @participationSettings.
  ///
  /// In en, this message translates to:
  /// **'Participation Settings'**
  String get participationSettings;

  /// No description provided for @reviewsFeedbacks.
  ///
  /// In en, this message translates to:
  /// **'Reviews & Feedbacks'**
  String get reviewsFeedbacks;

  /// No description provided for @contactSupport.
  ///
  /// In en, this message translates to:
  /// **'Contact Support'**
  String get contactSupport;

  /// No description provided for @changePhoto.
  ///
  /// In en, this message translates to:
  /// **'Change Photo'**
  String get changePhoto;

  /// No description provided for @camera.
  ///
  /// In en, this message translates to:
  /// **'Camera'**
  String get camera;

  /// No description provided for @gallery.
  ///
  /// In en, this message translates to:
  /// **'Gallery'**
  String get gallery;

  /// No description provided for @removePhoto.
  ///
  /// In en, this message translates to:
  /// **'Remove Photo'**
  String get removePhoto;

  /// No description provided for @myNutrition.
  ///
  /// In en, this message translates to:
  /// **'My Nutrition'**
  String get myNutrition;

  /// No description provided for @nutrition.
  ///
  /// In en, this message translates to:
  /// **'Nutrition'**
  String get nutrition;

  /// No description provided for @workoutMyPlan.
  ///
  /// In en, this message translates to:
  /// **'My Plan'**
  String get workoutMyPlan;

  /// No description provided for @workoutBrowsePlans.
  ///
  /// In en, this message translates to:
  /// **'Browse Plans'**
  String get workoutBrowsePlans;

  /// No description provided for @workoutFindMyPlan.
  ///
  /// In en, this message translates to:
  /// **'Find My Plan'**
  String get workoutFindMyPlan;

  /// No description provided for @workoutCurrentPlan.
  ///
  /// In en, this message translates to:
  /// **'Current Workout Plan'**
  String get workoutCurrentPlan;

  /// No description provided for @workoutNoPlan.
  ///
  /// In en, this message translates to:
  /// **'No workout plan assigned yet'**
  String get workoutNoPlan;

  /// No description provided for @workoutChoosePlan.
  ///
  /// In en, this message translates to:
  /// **'Choose This Plan'**
  String get workoutChoosePlan;

  /// No description provided for @workoutPlanAssigned.
  ///
  /// In en, this message translates to:
  /// **'Plan assigned successfully!'**
  String get workoutPlanAssigned;

  /// No description provided for @workoutExercises.
  ///
  /// In en, this message translates to:
  /// **'Exercises'**
  String get workoutExercises;

  /// No description provided for @workoutSets.
  ///
  /// In en, this message translates to:
  /// **'sets'**
  String get workoutSets;

  /// No description provided for @workoutReps.
  ///
  /// In en, this message translates to:
  /// **'reps'**
  String get workoutReps;

  /// No description provided for @workoutRest.
  ///
  /// In en, this message translates to:
  /// **'rest'**
  String get workoutRest;

  /// No description provided for @workoutSeconds.
  ///
  /// In en, this message translates to:
  /// **'sec'**
  String get workoutSeconds;

  /// No description provided for @workoutMuscleGroups.
  ///
  /// In en, this message translates to:
  /// **'Muscle Groups'**
  String get workoutMuscleGroups;

  /// No description provided for @workoutFilterByGoal.
  ///
  /// In en, this message translates to:
  /// **'Filter by Goal'**
  String get workoutFilterByGoal;

  /// No description provided for @workoutFilterByDifficulty.
  ///
  /// In en, this message translates to:
  /// **'Filter by Difficulty'**
  String get workoutFilterByDifficulty;

  /// No description provided for @workoutAllGoals.
  ///
  /// In en, this message translates to:
  /// **'All Goals'**
  String get workoutAllGoals;

  /// No description provided for @workoutAllDifficulties.
  ///
  /// In en, this message translates to:
  /// **'All Levels'**
  String get workoutAllDifficulties;

  /// No description provided for @workoutLoseWeight.
  ///
  /// In en, this message translates to:
  /// **'Lose Weight'**
  String get workoutLoseWeight;

  /// No description provided for @workoutBuildMuscle.
  ///
  /// In en, this message translates to:
  /// **'Build Muscle'**
  String get workoutBuildMuscle;

  /// No description provided for @workoutStayFit.
  ///
  /// In en, this message translates to:
  /// **'Stay Fit'**
  String get workoutStayFit;

  /// No description provided for @workoutBeginner.
  ///
  /// In en, this message translates to:
  /// **'Beginner'**
  String get workoutBeginner;

  /// No description provided for @workoutIntermediate.
  ///
  /// In en, this message translates to:
  /// **'Intermediate'**
  String get workoutIntermediate;

  /// No description provided for @workoutAdvanced.
  ///
  /// In en, this message translates to:
  /// **'Advanced'**
  String get workoutAdvanced;

  /// No description provided for @workoutQuestionnaire.
  ///
  /// In en, this message translates to:
  /// **'Quick Assessment'**
  String get workoutQuestionnaire;

  /// No description provided for @workoutWhatIsYourGoal.
  ///
  /// In en, this message translates to:
  /// **'What is your fitness goal?'**
  String get workoutWhatIsYourGoal;

  /// No description provided for @workoutExperienceLevel.
  ///
  /// In en, this message translates to:
  /// **'What is your experience level?'**
  String get workoutExperienceLevel;

  /// No description provided for @workoutDaysPerWeek.
  ///
  /// In en, this message translates to:
  /// **'How many days per week can you train?'**
  String get workoutDaysPerWeek;

  /// No description provided for @workoutDays23.
  ///
  /// In en, this message translates to:
  /// **'2-3 days'**
  String get workoutDays23;

  /// No description provided for @workoutDays45.
  ///
  /// In en, this message translates to:
  /// **'4-5 days'**
  String get workoutDays45;

  /// No description provided for @workoutDays67.
  ///
  /// In en, this message translates to:
  /// **'6-7 days'**
  String get workoutDays67;

  /// No description provided for @workoutGetRecommendation.
  ///
  /// In en, this message translates to:
  /// **'Get My Plan'**
  String get workoutGetRecommendation;

  /// No description provided for @workoutRecommendedPlan.
  ///
  /// In en, this message translates to:
  /// **'Recommended Plan'**
  String get workoutRecommendedPlan;

  /// No description provided for @workoutSelfAssigned.
  ///
  /// In en, this message translates to:
  /// **'Self-assigned'**
  String get workoutSelfAssigned;

  /// No description provided for @workoutQuestionnaireAssigned.
  ///
  /// In en, this message translates to:
  /// **'Recommended'**
  String get workoutQuestionnaireAssigned;

  /// No description provided for @workoutStaffAssigned.
  ///
  /// In en, this message translates to:
  /// **'Assigned by staff'**
  String get workoutStaffAssigned;

  /// No description provided for @workoutTargetMuscles.
  ///
  /// In en, this message translates to:
  /// **'Target Muscles'**
  String get workoutTargetMuscles;

  /// No description provided for @workoutBodyMap.
  ///
  /// In en, this message translates to:
  /// **'Body Map'**
  String get workoutBodyMap;

  /// No description provided for @nutritionDailyIntake.
  ///
  /// In en, this message translates to:
  /// **'Daily Intake'**
  String get nutritionDailyIntake;

  /// No description provided for @nutritionCaloriesRemaining.
  ///
  /// In en, this message translates to:
  /// **'Remaining'**
  String get nutritionCaloriesRemaining;

  /// No description provided for @nutritionCaloriesConsumed.
  ///
  /// In en, this message translates to:
  /// **'Consumed'**
  String get nutritionCaloriesConsumed;

  /// No description provided for @nutritionCalorieGoal.
  ///
  /// In en, this message translates to:
  /// **'Goal'**
  String get nutritionCalorieGoal;

  /// No description provided for @nutritionProtein.
  ///
  /// In en, this message translates to:
  /// **'Protein'**
  String get nutritionProtein;

  /// No description provided for @nutritionCarbs.
  ///
  /// In en, this message translates to:
  /// **'Carbs'**
  String get nutritionCarbs;

  /// No description provided for @nutritionFat.
  ///
  /// In en, this message translates to:
  /// **'Fat'**
  String get nutritionFat;

  /// No description provided for @nutritionBreakfast.
  ///
  /// In en, this message translates to:
  /// **'Breakfast'**
  String get nutritionBreakfast;

  /// No description provided for @nutritionLunch.
  ///
  /// In en, this message translates to:
  /// **'Lunch'**
  String get nutritionLunch;

  /// No description provided for @nutritionDinner.
  ///
  /// In en, this message translates to:
  /// **'Dinner'**
  String get nutritionDinner;

  /// No description provided for @nutritionSnack.
  ///
  /// In en, this message translates to:
  /// **'Snack'**
  String get nutritionSnack;

  /// No description provided for @nutritionAddFood.
  ///
  /// In en, this message translates to:
  /// **'Add Food'**
  String get nutritionAddFood;

  /// No description provided for @nutritionSearchFood.
  ///
  /// In en, this message translates to:
  /// **'Search food...'**
  String get nutritionSearchFood;

  /// No description provided for @nutritionQuantity.
  ///
  /// In en, this message translates to:
  /// **'Quantity (g)'**
  String get nutritionQuantity;

  /// No description provided for @nutritionCalories.
  ///
  /// In en, this message translates to:
  /// **'Calories'**
  String get nutritionCalories;

  /// No description provided for @nutritionNoLogs.
  ///
  /// In en, this message translates to:
  /// **'No food logged yet'**
  String get nutritionNoLogs;

  /// No description provided for @nutritionLogFood.
  ///
  /// In en, this message translates to:
  /// **'Log Food'**
  String get nutritionLogFood;

  /// No description provided for @nutritionWeeklyProgress.
  ///
  /// In en, this message translates to:
  /// **'Weekly Progress'**
  String get nutritionWeeklyProgress;

  /// No description provided for @nutritionSetGoal.
  ///
  /// In en, this message translates to:
  /// **'Set Nutrition Goal'**
  String get nutritionSetGoal;

  /// No description provided for @nutritionEditGoal.
  ///
  /// In en, this message translates to:
  /// **'Edit Goal'**
  String get nutritionEditGoal;

  /// No description provided for @nutritionWeight.
  ///
  /// In en, this message translates to:
  /// **'Weight (kg)'**
  String get nutritionWeight;

  /// No description provided for @nutritionHeight.
  ///
  /// In en, this message translates to:
  /// **'Height (cm)'**
  String get nutritionHeight;

  /// No description provided for @nutritionAge.
  ///
  /// In en, this message translates to:
  /// **'Age'**
  String get nutritionAge;

  /// No description provided for @nutritionGender.
  ///
  /// In en, this message translates to:
  /// **'Gender'**
  String get nutritionGender;

  /// No description provided for @nutritionMale.
  ///
  /// In en, this message translates to:
  /// **'Male'**
  String get nutritionMale;

  /// No description provided for @nutritionFemale.
  ///
  /// In en, this message translates to:
  /// **'Female'**
  String get nutritionFemale;

  /// No description provided for @nutritionActivityLevel.
  ///
  /// In en, this message translates to:
  /// **'Activity Level'**
  String get nutritionActivityLevel;

  /// No description provided for @nutritionSedentary.
  ///
  /// In en, this message translates to:
  /// **'Sedentary'**
  String get nutritionSedentary;

  /// No description provided for @nutritionLight.
  ///
  /// In en, this message translates to:
  /// **'Lightly Active'**
  String get nutritionLight;

  /// No description provided for @nutritionModerate.
  ///
  /// In en, this message translates to:
  /// **'Moderately Active'**
  String get nutritionModerate;

  /// No description provided for @nutritionActive.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get nutritionActive;

  /// No description provided for @nutritionVeryActive.
  ///
  /// In en, this message translates to:
  /// **'Very Active'**
  String get nutritionVeryActive;

  /// No description provided for @nutritionGoalType.
  ///
  /// In en, this message translates to:
  /// **'Goal'**
  String get nutritionGoalType;

  /// No description provided for @nutritionLoseWeight.
  ///
  /// In en, this message translates to:
  /// **'Lose Weight'**
  String get nutritionLoseWeight;

  /// No description provided for @nutritionMaintain.
  ///
  /// In en, this message translates to:
  /// **'Maintain Weight'**
  String get nutritionMaintain;

  /// No description provided for @nutritionBuildMuscle.
  ///
  /// In en, this message translates to:
  /// **'Build Muscle'**
  String get nutritionBuildMuscle;

  /// No description provided for @nutritionCalculateGoal.
  ///
  /// In en, this message translates to:
  /// **'Calculate Goal'**
  String get nutritionCalculateGoal;

  /// No description provided for @nutritionGoalSet.
  ///
  /// In en, this message translates to:
  /// **'Nutrition goal set!'**
  String get nutritionGoalSet;

  /// No description provided for @nutritionDailyCalories.
  ///
  /// In en, this message translates to:
  /// **'Daily Calories'**
  String get nutritionDailyCalories;

  /// No description provided for @nutritionMealTemplates.
  ///
  /// In en, this message translates to:
  /// **'Meal Templates'**
  String get nutritionMealTemplates;

  /// No description provided for @nutritionCreateTemplate.
  ///
  /// In en, this message translates to:
  /// **'Create Template'**
  String get nutritionCreateTemplate;

  /// No description provided for @nutritionTemplateName.
  ///
  /// In en, this message translates to:
  /// **'Template Name'**
  String get nutritionTemplateName;

  /// No description provided for @nutritionQuickLog.
  ///
  /// In en, this message translates to:
  /// **'Quick Log'**
  String get nutritionQuickLog;

  /// No description provided for @nutritionNoTemplates.
  ///
  /// In en, this message translates to:
  /// **'No meal templates yet'**
  String get nutritionNoTemplates;

  /// No description provided for @nutritionPer100g.
  ///
  /// In en, this message translates to:
  /// **'per 100g'**
  String get nutritionPer100g;

  /// No description provided for @nutritionServingSize.
  ///
  /// In en, this message translates to:
  /// **'Serving'**
  String get nutritionServingSize;

  /// No description provided for @nutritionDeleteLog.
  ///
  /// In en, this message translates to:
  /// **'Remove'**
  String get nutritionDeleteLog;

  /// No description provided for @nutritionKcal.
  ///
  /// In en, this message translates to:
  /// **'kcal'**
  String get nutritionKcal;

  /// No description provided for @nutritionGrams.
  ///
  /// In en, this message translates to:
  /// **'g'**
  String get nutritionGrams;

  /// No description provided for @nutritionNoGoalSet.
  ///
  /// In en, this message translates to:
  /// **'Set up your nutrition goal to track progress'**
  String get nutritionNoGoalSet;

  /// No description provided for @nutritionToday.
  ///
  /// In en, this message translates to:
  /// **'Today'**
  String get nutritionToday;

  /// No description provided for @onboardingTitle.
  ///
  /// In en, this message translates to:
  /// **'Welcome to Urja'**
  String get onboardingTitle;

  /// No description provided for @onboardingWhatBringsYou.
  ///
  /// In en, this message translates to:
  /// **'What brings you here?'**
  String get onboardingWhatBringsYou;

  /// No description provided for @onboardingGymMember.
  ///
  /// In en, this message translates to:
  /// **'Gym Member'**
  String get onboardingGymMember;

  /// No description provided for @onboardingGymMemberDesc.
  ///
  /// In en, this message translates to:
  /// **'I go to a gym and want to track my membership'**
  String get onboardingGymMemberDesc;

  /// No description provided for @onboardingFitnessTracker.
  ///
  /// In en, this message translates to:
  /// **'Fitness Tracker'**
  String get onboardingFitnessTracker;

  /// No description provided for @onboardingFitnessTrackerDesc.
  ///
  /// In en, this message translates to:
  /// **'I work out at home and want to follow training programs'**
  String get onboardingFitnessTrackerDesc;

  /// No description provided for @onboardingCalorieTracker.
  ///
  /// In en, this message translates to:
  /// **'Calorie Tracker'**
  String get onboardingCalorieTracker;

  /// No description provided for @onboardingCalorieTrackerDesc.
  ///
  /// In en, this message translates to:
  /// **'I want to track my food intake and manage my weight'**
  String get onboardingCalorieTrackerDesc;

  /// No description provided for @onboardingWhatsYourGoal.
  ///
  /// In en, this message translates to:
  /// **'What\'s your goal?'**
  String get onboardingWhatsYourGoal;

  /// No description provided for @onboardingLoseWeight.
  ///
  /// In en, this message translates to:
  /// **'Lose Weight'**
  String get onboardingLoseWeight;

  /// No description provided for @onboardingBuildMuscle.
  ///
  /// In en, this message translates to:
  /// **'Build Muscle'**
  String get onboardingBuildMuscle;

  /// No description provided for @onboardingStayFit.
  ///
  /// In en, this message translates to:
  /// **'Stay Fit'**
  String get onboardingStayFit;

  /// No description provided for @onboardingGeneralHealth.
  ///
  /// In en, this message translates to:
  /// **'General Health'**
  String get onboardingGeneralHealth;

  /// No description provided for @onboardingAboutYou.
  ///
  /// In en, this message translates to:
  /// **'About You'**
  String get onboardingAboutYou;

  /// No description provided for @onboardingYourName.
  ///
  /// In en, this message translates to:
  /// **'Your Name'**
  String get onboardingYourName;

  /// No description provided for @onboardingNamePlaceholder.
  ///
  /// In en, this message translates to:
  /// **'Enter your full name'**
  String get onboardingNamePlaceholder;

  /// No description provided for @onboardingWeight.
  ///
  /// In en, this message translates to:
  /// **'Weight (kg)'**
  String get onboardingWeight;

  /// No description provided for @onboardingHeight.
  ///
  /// In en, this message translates to:
  /// **'Height (cm)'**
  String get onboardingHeight;

  /// No description provided for @onboardingAge.
  ///
  /// In en, this message translates to:
  /// **'Age'**
  String get onboardingAge;

  /// No description provided for @onboardingGender.
  ///
  /// In en, this message translates to:
  /// **'Gender'**
  String get onboardingGender;

  /// No description provided for @onboardingActivityLevel.
  ///
  /// In en, this message translates to:
  /// **'Activity Level'**
  String get onboardingActivityLevel;

  /// No description provided for @onboardingSedentary.
  ///
  /// In en, this message translates to:
  /// **'Sedentary'**
  String get onboardingSedentary;

  /// No description provided for @onboardingLight.
  ///
  /// In en, this message translates to:
  /// **'Lightly Active'**
  String get onboardingLight;

  /// No description provided for @onboardingModerate.
  ///
  /// In en, this message translates to:
  /// **'Moderately Active'**
  String get onboardingModerate;

  /// No description provided for @onboardingActive.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get onboardingActive;

  /// No description provided for @onboardingVeryActive.
  ///
  /// In en, this message translates to:
  /// **'Very Active'**
  String get onboardingVeryActive;

  /// No description provided for @onboardingScanGymQr.
  ///
  /// In en, this message translates to:
  /// **'Scan your gym\'s QR code to join'**
  String get onboardingScanGymQr;

  /// No description provided for @onboardingJoinGym.
  ///
  /// In en, this message translates to:
  /// **'Join Gym'**
  String get onboardingJoinGym;

  /// No description provided for @onboardingSkipForNow.
  ///
  /// In en, this message translates to:
  /// **'Skip for now'**
  String get onboardingSkipForNow;

  /// No description provided for @onboardingAllSet.
  ///
  /// In en, this message translates to:
  /// **'You\'re all set!'**
  String get onboardingAllSet;

  /// No description provided for @onboardingDailyCalories.
  ///
  /// In en, this message translates to:
  /// **'Daily Calories'**
  String get onboardingDailyCalories;

  /// No description provided for @onboardingProtein.
  ///
  /// In en, this message translates to:
  /// **'Protein'**
  String get onboardingProtein;

  /// No description provided for @onboardingCarbs.
  ///
  /// In en, this message translates to:
  /// **'Carbs'**
  String get onboardingCarbs;

  /// No description provided for @onboardingFat.
  ///
  /// In en, this message translates to:
  /// **'Fat'**
  String get onboardingFat;

  /// No description provided for @onboardingLetsGo.
  ///
  /// In en, this message translates to:
  /// **'Let\'s Go!'**
  String get onboardingLetsGo;

  /// No description provided for @onboardingNext.
  ///
  /// In en, this message translates to:
  /// **'Next'**
  String get onboardingNext;

  /// No description provided for @onboardingBack.
  ///
  /// In en, this message translates to:
  /// **'Back'**
  String get onboardingBack;

  /// No description provided for @onboardingStep.
  ///
  /// In en, this message translates to:
  /// **'Step {current} of {total}'**
  String onboardingStep(int current, int total);

  /// No description provided for @waterTitle.
  ///
  /// In en, this message translates to:
  /// **'Water Intake'**
  String get waterTitle;

  /// No description provided for @waterTodayIntake.
  ///
  /// In en, this message translates to:
  /// **'Today\'s Intake'**
  String get waterTodayIntake;

  /// No description provided for @waterGoal.
  ///
  /// In en, this message translates to:
  /// **'Goal'**
  String get waterGoal;

  /// No description provided for @waterAddWater.
  ///
  /// In en, this message translates to:
  /// **'Add Water'**
  String get waterAddWater;

  /// No description provided for @waterAmount.
  ///
  /// In en, this message translates to:
  /// **'Amount (ml)'**
  String get waterAmount;

  /// No description provided for @waterGlassSize.
  ///
  /// In en, this message translates to:
  /// **'Glass (250ml)'**
  String get waterGlassSize;

  /// No description provided for @waterBottleSize.
  ///
  /// In en, this message translates to:
  /// **'Bottle (500ml)'**
  String get waterBottleSize;

  /// No description provided for @waterLargeBottle.
  ///
  /// In en, this message translates to:
  /// **'Large (1000ml)'**
  String get waterLargeBottle;

  /// No description provided for @waterCustomAmount.
  ///
  /// In en, this message translates to:
  /// **'Custom'**
  String get waterCustomAmount;

  /// No description provided for @waterRemaining.
  ///
  /// In en, this message translates to:
  /// **'Remaining'**
  String get waterRemaining;

  /// No description provided for @waterCompleted.
  ///
  /// In en, this message translates to:
  /// **'Goal Reached!'**
  String get waterCompleted;

  /// No description provided for @waterNoLogs.
  ///
  /// In en, this message translates to:
  /// **'No water logged today'**
  String get waterNoLogs;

  /// No description provided for @waterDelete.
  ///
  /// In en, this message translates to:
  /// **'Remove'**
  String get waterDelete;

  /// No description provided for @myWater.
  ///
  /// In en, this message translates to:
  /// **'Water Intake'**
  String get myWater;

  /// No description provided for @trackerHome.
  ///
  /// In en, this message translates to:
  /// **'Home'**
  String get trackerHome;

  /// No description provided for @trackerNutrition.
  ///
  /// In en, this message translates to:
  /// **'Nutrition'**
  String get trackerNutrition;

  /// No description provided for @trackerWorkouts.
  ///
  /// In en, this message translates to:
  /// **'Workouts'**
  String get trackerWorkouts;

  /// No description provided for @trackerProgress.
  ///
  /// In en, this message translates to:
  /// **'Progress'**
  String get trackerProgress;

  /// No description provided for @trackerAccount.
  ///
  /// In en, this message translates to:
  /// **'Account'**
  String get trackerAccount;

  /// No description provided for @progressTitle.
  ///
  /// In en, this message translates to:
  /// **'Progress'**
  String get progressTitle;

  /// No description provided for @progressLogWeight.
  ///
  /// In en, this message translates to:
  /// **'Log Weight'**
  String get progressLogWeight;

  /// No description provided for @progressWeightKg.
  ///
  /// In en, this message translates to:
  /// **'Weight (kg)'**
  String get progressWeightKg;

  /// No description provided for @progressWeightLogged.
  ///
  /// In en, this message translates to:
  /// **'Weight logged successfully!'**
  String get progressWeightLogged;

  /// No description provided for @progressCurrentWeight.
  ///
  /// In en, this message translates to:
  /// **'Current'**
  String get progressCurrentWeight;

  /// No description provided for @progressStartWeight.
  ///
  /// In en, this message translates to:
  /// **'Start'**
  String get progressStartWeight;

  /// No description provided for @progressChange.
  ///
  /// In en, this message translates to:
  /// **'Change'**
  String get progressChange;

  /// No description provided for @progressBmi.
  ///
  /// In en, this message translates to:
  /// **'BMI'**
  String get progressBmi;

  /// No description provided for @progressWeightTrend.
  ///
  /// In en, this message translates to:
  /// **'Weight Trend'**
  String get progressWeightTrend;

  /// No description provided for @progressNoData.
  ///
  /// In en, this message translates to:
  /// **'No weight data yet. Log your first entry!'**
  String get progressNoData;

  /// No description provided for @progressPeriod30.
  ///
  /// In en, this message translates to:
  /// **'30 Days'**
  String get progressPeriod30;

  /// No description provided for @progressPeriod90.
  ///
  /// In en, this message translates to:
  /// **'90 Days'**
  String get progressPeriod90;

  /// No description provided for @progressPeriod180.
  ///
  /// In en, this message translates to:
  /// **'180 Days'**
  String get progressPeriod180;

  /// No description provided for @progressNutritionStreak.
  ///
  /// In en, this message translates to:
  /// **'Nutrition Streak'**
  String get progressNutritionStreak;

  /// No description provided for @progressEntries.
  ///
  /// In en, this message translates to:
  /// **'entries'**
  String get progressEntries;

  /// No description provided for @progressKg.
  ///
  /// In en, this message translates to:
  /// **'kg'**
  String get progressKg;

  /// No description provided for @gymCode.
  ///
  /// In en, this message translates to:
  /// **'Gym Code'**
  String get gymCode;

  /// No description provided for @scanToCheckIn.
  ///
  /// In en, this message translates to:
  /// **'Scan to Check In'**
  String get scanToCheckIn;

  /// No description provided for @gymQrCodeError.
  ///
  /// In en, this message translates to:
  /// **'Failed to load QR code'**
  String get gymQrCodeError;

  /// No description provided for @gymQrCodeDesc.
  ///
  /// In en, this message translates to:
  /// **'Members can scan this code to check in'**
  String get gymQrCodeDesc;

  /// No description provided for @joinAGym.
  ///
  /// In en, this message translates to:
  /// **'Join a Gym'**
  String get joinAGym;

  /// No description provided for @joinAGymDesc.
  ///
  /// In en, this message translates to:
  /// **'Connect to a gym to track attendance and access workouts'**
  String get joinAGymDesc;

  /// No description provided for @scanGymQrCode.
  ///
  /// In en, this message translates to:
  /// **'Scan Gym QR Code'**
  String get scanGymQrCode;

  /// No description provided for @orEnterCodeManually.
  ///
  /// In en, this message translates to:
  /// **'or enter code manually'**
  String get orEnterCodeManually;

  /// No description provided for @joinGym.
  ///
  /// In en, this message translates to:
  /// **'Join Gym'**
  String get joinGym;

  /// No description provided for @successfullyJoinedGym.
  ///
  /// In en, this message translates to:
  /// **'Successfully joined the gym!'**
  String get successfullyJoinedGym;

  /// No description provided for @couldNotJoinGym.
  ///
  /// In en, this message translates to:
  /// **'Could not join gym. Please check the code.'**
  String get couldNotJoinGym;

  /// No description provided for @enterGymCode.
  ///
  /// In en, this message translates to:
  /// **'Please enter your gym code'**
  String get enterGymCode;

  /// No description provided for @submitFeedback.
  ///
  /// In en, this message translates to:
  /// **'Submit Feedback'**
  String get submitFeedback;

  /// No description provided for @rateYourExperience.
  ///
  /// In en, this message translates to:
  /// **'Rate your experience'**
  String get rateYourExperience;

  /// No description provided for @submit.
  ///
  /// In en, this message translates to:
  /// **'Submit'**
  String get submit;

  /// No description provided for @feedbackSubmitted.
  ///
  /// In en, this message translates to:
  /// **'Thank you for your feedback!'**
  String get feedbackSubmitted;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['en', 'ne'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'ne':
      return AppLocalizationsNe();
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
