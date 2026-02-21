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
