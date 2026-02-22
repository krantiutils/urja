// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Nepali (`ne`).
class AppLocalizationsNe extends AppLocalizations {
  AppLocalizationsNe([String locale = 'ne']) : super(locale);

  @override
  String get appName => 'ऊर्जा';

  @override
  String get loading => 'लोड हुँदैछ...';

  @override
  String get error => 'केही गलत भयो';

  @override
  String get retry => 'पुनः प्रयास';

  @override
  String get save => 'सेभ गर्नुहोस्';

  @override
  String get cancel => 'रद्द गर्नुहोस्';

  @override
  String get delete => 'मेट्नुहोस्';

  @override
  String get search => 'खोज्नुहोस्...';

  @override
  String get noResults => 'कुनै नतिजा भेटिएन';

  @override
  String get signIn => 'साइन इन';

  @override
  String get signOut => 'साइन आउट';

  @override
  String get phoneLabel => 'मोबाइल नम्बर';

  @override
  String get phonePlaceholder => '98XXXXXXXX';

  @override
  String get phoneNumber => 'फोन नम्बर';

  @override
  String get otpLabel => 'प्रमाणीकरण कोड';

  @override
  String get otpPlaceholder => '६ अंकको कोड हाल्नुहोस्';

  @override
  String get requestOtp => 'OTP पठाउनुहोस्';

  @override
  String get sendOtp => 'OTP पठाउनुहोस्';

  @override
  String get verifyOtp => 'प्रमाणित गर्नुहोस्';

  @override
  String get enterOtp => 'प्रमाणीकरण कोड हाल्नुहोस्';

  @override
  String get otpSent => 'तपाईंको फोनमा OTP पठाइएको छ';

  @override
  String get invalidPhone => 'वैध नेपाली मोबाइल नम्बर हाल्नुहोस्';

  @override
  String get invalidOtp => 'अवैध प्रमाणीकरण कोड';

  @override
  String get welcome => 'ऊर्जामा स्वागत छ';

  @override
  String get subtitle => 'जिम व्यवस्थापन, सरलीकृत';

  @override
  String get resendOtp => 'OTP पुन: पठाउनुहोस्';

  @override
  String get resend => 'पुन: पठाउनुहोस्';

  @override
  String resendIn(int seconds) {
    return 'पुन: पठाउनुहोस् ${seconds}s';
  }

  @override
  String get didNotReceiveCode => 'कोड प्राप्त भएन?';

  @override
  String get dashboard => 'ड्यासबोर्ड';

  @override
  String get members => 'सदस्यहरू';

  @override
  String get packages => 'प्याकेजहरू';

  @override
  String get duePayments => 'बाँकी भुक्तानी';

  @override
  String get attendance => 'उपस्थिति';

  @override
  String get staff => 'कर्मचारी';

  @override
  String get accounts => 'लेखा';

  @override
  String get stories => 'कथाहरू';

  @override
  String get sms => 'SMS';

  @override
  String get feedbacks => 'प्रतिक्रिया';

  @override
  String get nfcAccess => 'NFC कार्ड पहुँच';

  @override
  String get settings => 'सेटिङ्स';

  @override
  String get totalMembers => 'कुल सदस्यहरू';

  @override
  String get activeMembers => 'सक्रिय सदस्यहरू';

  @override
  String get todayAttendance => 'आजको उपस्थिति';

  @override
  String get monthlyRevenue => 'मासिक आम्दानी';

  @override
  String get expiringPackages => 'समाप्त हुन लागेका प्याकेजहरू';

  @override
  String get expiredPackages => 'म्याद सकिएका प्याकेजहरू';

  @override
  String get noExpiringPackages => 'कुनै समाप्त हुन लागेका प्याकेजहरू छैनन्';

  @override
  String get todayActivity => 'आजको गतिविधि';

  @override
  String get recentActivity => 'हालको गतिविधि';

  @override
  String get noActivityToday => 'आज कुनै गतिविधि छैन';

  @override
  String get unknownMember => 'अज्ञात सदस्य';

  @override
  String get viewAll => 'सबै हेर्नुहोस्';

  @override
  String get addMember => 'सदस्य थप्नुहोस्';

  @override
  String get editMember => 'सदस्य सम्पादन';

  @override
  String get removeMember => 'सदस्य हटाउनुहोस्';

  @override
  String get name => 'नाम';

  @override
  String get nameNe => 'नाम (नेपाली)';

  @override
  String get phone => 'फोन';

  @override
  String get email => 'इमेल';

  @override
  String get status => 'स्थिति';

  @override
  String get role => 'भूमिका';

  @override
  String get active => 'सक्रिय';

  @override
  String get inactive => 'निष्क्रिय';

  @override
  String get suspended => 'निलम्बित';

  @override
  String get left => 'छोडेको';

  @override
  String get member => 'सदस्य';

  @override
  String get confirmDelete => 'के तपाईं यसलाई मेट्न चाहनुहुन्छ?';

  @override
  String get noMembers => 'कुनै सदस्य भेटिएन';

  @override
  String get addPackage => 'प्याकेज थप्नुहोस्';

  @override
  String get editPackage => 'प्याकेज सम्पादन';

  @override
  String get price => 'मूल्य';

  @override
  String get duration => 'अवधि';

  @override
  String get features => 'सुविधाहरू';

  @override
  String get days => 'दिन';

  @override
  String get maxMembers => 'अधिकतम सदस्य';

  @override
  String get unlimited => 'असीमित';

  @override
  String get description => 'विवरण';

  @override
  String get filterAll => 'सबै';

  @override
  String get filterActive => 'सक्रिय';

  @override
  String get filterInactive => 'निष्क्रिय';

  @override
  String get todayCheckIns => 'आजको चेक-इनहरू';

  @override
  String get manualCheckIn => 'म्यानुअल चेक-इन';

  @override
  String get checkInMember => 'सदस्यलाई चेक-इन गर्नुहोस्';

  @override
  String get method => 'विधि';

  @override
  String get qr => 'QR';

  @override
  String get nfc => 'NFC';

  @override
  String get manual => 'म्यानुअल';

  @override
  String get noRecords => 'कुनै रेकर्ड भेटिएन';

  @override
  String get addStaff => 'कर्मचारी थप्नुहोस्';

  @override
  String get editStaff => 'कर्मचारी सम्पादन';

  @override
  String get staffRole => 'कर्मचारी भूमिका';

  @override
  String get owner => 'मालिक';

  @override
  String get manager => 'प्रबन्धक';

  @override
  String get trainer => 'प्रशिक्षक';

  @override
  String get receptionist => 'रिसेप्सनिस्ट';

  @override
  String get amount => 'रकम';

  @override
  String get dueDate => 'भुक्तानी मिति';

  @override
  String get unpaid => 'भुक्तानी नभएको';

  @override
  String get paid => 'भुक्तानी भएको';

  @override
  String get waived => 'माफ गरिएको';

  @override
  String get recordPayment => 'भुक्तानी रेकर्ड';

  @override
  String get paymentMethod => 'भुक्तानी विधि';

  @override
  String get cash => 'नगद';

  @override
  String get esewa => 'इसेवा';

  @override
  String get bankTransfer => 'बैंक ट्रान्सफर';

  @override
  String get confirmPayment => 'भुक्तानी पुष्टि';

  @override
  String get noDues => 'कुनै बाँकी भुक्तानी छैन';

  @override
  String get addTransaction => 'कारोबार थप्नुहोस्';

  @override
  String get editTransaction => 'कारोबार सम्पादन';

  @override
  String get category => 'वर्ग';

  @override
  String get type => 'प्रकार';

  @override
  String get paymentType => 'भुक्तानी प्रकार';

  @override
  String get reference => 'सन्दर्भ';

  @override
  String get income => 'आम्दानी';

  @override
  String get expense => 'खर्च';

  @override
  String get totalIncome => 'कुल आम्दानी';

  @override
  String get totalExpenses => 'कुल खर्च';

  @override
  String get grossProfit => 'कुल नाफा';

  @override
  String get profitPercent => 'नाफा %';

  @override
  String get noTransactions => 'कुनै कारोबार भेटिएन';

  @override
  String get addNotice => 'सूचना थप्नुहोस्';

  @override
  String get editNotice => 'सूचना सम्पादन';

  @override
  String get noticeTitle => 'शीर्षक';

  @override
  String get titleNe => 'शीर्षक (नेपाली)';

  @override
  String get content => 'सामग्री';

  @override
  String get contentNe => 'सामग्री (नेपाली)';

  @override
  String get noNotices => 'कुनै सूचना भेटिएन';

  @override
  String get noFeedbacks => 'अहिलेसम्म कुनै प्रतिक्रिया छैन';

  @override
  String get cards => 'कार्डहरू';

  @override
  String get devices => 'उपकरणहरू';

  @override
  String get registerCard => 'कार्ड दर्ता';

  @override
  String get cardNumber => 'कार्ड नम्बर';

  @override
  String get assignedTo => 'तोकिएको';

  @override
  String get assign => 'तोक्नुहोस्';

  @override
  String get unassign => 'हटाउनुहोस्';

  @override
  String get unassigned => 'तोकिएको छैन';

  @override
  String get available => 'उपलब्ध';

  @override
  String get assigned => 'तोकिएको';

  @override
  String get disabled => 'निष्क्रिय';

  @override
  String get noCards => 'कुनै NFC कार्ड दर्ता छैन';

  @override
  String get deviceName => 'उपकरण नाम';

  @override
  String get online => 'अनलाइन';

  @override
  String get offline => 'अफलाइन';

  @override
  String get locked => 'बन्द';

  @override
  String get unlocked => 'खुला';

  @override
  String get noDevices => 'कुनै NFC उपकरण भेटिएन';

  @override
  String get door => 'ढोका';

  @override
  String get uptime => 'अपटाइम';

  @override
  String get smsBalance => 'SMS ब्यालेन्स';

  @override
  String get remaining => 'बाँकी';

  @override
  String get totalPurchased => 'कुल खरिद';

  @override
  String get totalUsed => 'कुल प्रयोग';

  @override
  String get campaigns => 'अभियानहरू';

  @override
  String get sendSms => 'SMS पठाउनुहोस्';

  @override
  String get message => 'सन्देश';

  @override
  String get enterMessage => 'सन्देश लेख्नुहोस्';

  @override
  String get selectMembers => 'सदस्य छान्नुहोस्';

  @override
  String membersSelected(int count) {
    return '$count सदस्य छानिएका';
  }

  @override
  String get send => 'पठाउनुहोस्';

  @override
  String get history => 'इतिहास';

  @override
  String get noHistory => 'कुनै इतिहास छैन';

  @override
  String get orgProfile => 'संस्थाको प्रोफाइल';

  @override
  String get myProfile => 'मेरो प्रोफाइल';

  @override
  String get profileInfo => 'प्रोफाइल जानकारी';

  @override
  String get noProfile => 'कुनै प्रोफाइल डाटा छैन';

  @override
  String get privacySettings => 'गोपनीयता सेटिङ्स';

  @override
  String get showEmail => 'अरूलाई इमेल देखाउनुहोस्';

  @override
  String get showEmailDesc => 'अरू सदस्यहरूलाई इमेल देख्न दिनुहोस्';

  @override
  String get showPhone => 'अरूलाई फोन देखाउनुहोस्';

  @override
  String get showPhoneDesc => 'अरू सदस्यहरूलाई फोन देख्न दिनुहोस्';

  @override
  String get showProfile => 'प्रोफाइल सार्वजनिक देखाउनुहोस्';

  @override
  String get showProfileDesc => 'प्रोफाइल अरूलाई देख्न दिनुहोस्';

  @override
  String get showAttendance => 'अरूलाई उपस्थिति देखाउनुहोस्';

  @override
  String get showAttendanceDesc => 'अरूलाई उपस्थिति देख्न दिनुहोस्';

  @override
  String get showOnLeaderboard => 'लिडरबोर्डमा देखाउनुहोस्';

  @override
  String get showOnLeaderboardDesc => 'जिम लिडरबोर्डमा देखाउनुहोस्';

  @override
  String get emergencyContact => 'आपतकालीन सम्पर्क';

  @override
  String get contactName => 'सम्पर्क नाम';

  @override
  String get contactPhone => 'सम्पर्क फोन';

  @override
  String get saveChanges => 'परिवर्तनहरू सेभ गर्नुहोस्';

  @override
  String get saveProfile => 'प्रोफाइल सेभ गर्नुहोस्';

  @override
  String get savePrivacy => 'गोपनीयता सेभ गर्नुहोस्';

  @override
  String get saveEmergencyContact => 'आपतकालीन सम्पर्क सेभ गर्नुहोस्';

  @override
  String get saved => 'परिवर्तनहरू सेभ भयो';

  @override
  String get myDashboard => 'मेरो ड्यासबोर्ड';

  @override
  String get myAttendance => 'मेरो उपस्थिति';

  @override
  String get myPackages => 'मेरा प्याकेजहरू';

  @override
  String get myHealth => 'मेरो स्वास्थ्य';

  @override
  String get myWorkouts => 'मेरा व्यायामहरू';

  @override
  String get profile => 'प्रोफाइल';

  @override
  String get welcomeBack => 'स्वागत छ';

  @override
  String get memberSince => 'सदस्य मिति';

  @override
  String get currentPackage => 'हालको प्याकेज';

  @override
  String get activePackage => 'सक्रिय प्याकेज';

  @override
  String get noActivePackage => 'कुनै सक्रिय प्याकेज छैन';

  @override
  String get noActivePackages => 'कुनै सक्रिय प्याकेज छैन';

  @override
  String get noPastPackages => 'कुनै विगतका प्याकेज छैन';

  @override
  String get daysRemaining => 'दिन बाँकी';

  @override
  String get expired => 'म्याद सकिएको';

  @override
  String get expires => 'म्याद';

  @override
  String get currentStreak => 'हालको स्ट्रिक';

  @override
  String get longestStreak => 'सबैभन्दा लामो स्ट्रिक';

  @override
  String get recentCheckIns => 'हालको चेक-इनहरू';

  @override
  String get noRecentCheckIns => 'कुनै हालको चेक-इन छैन';

  @override
  String get noCheckIns => 'अहिलेसम्म चेक-इन छैन';

  @override
  String get noAttendanceRecords => 'कुनै उपस्थिति रेकर्ड छैन';

  @override
  String get attendanceHistory => 'उपस्थिति इतिहास';

  @override
  String get bmi => 'BMI';

  @override
  String get measurements => 'शरीर नाप';

  @override
  String get bodyMeasurements => 'शरीर नाप';

  @override
  String get noHealthData => 'अहिलेसम्म कुनै स्वास्थ्य डाटा छैन';

  @override
  String get noBmiData => 'कुनै BMI डाटा छैन';

  @override
  String get noMeasurementData => 'कुनै नाप डाटा छैन';

  @override
  String get height => 'उचाइ';

  @override
  String get weight => 'तौल';

  @override
  String get chest => 'छाती';

  @override
  String get waist => 'कम्मर';

  @override
  String get hips => 'कुल्हा';

  @override
  String get biceps => 'बाइसेप्स';

  @override
  String get thighs => 'जाँघ';

  @override
  String get bodyFat => 'शरीर बोसो';

  @override
  String get recordedOn => 'रेकर्ड गरिएको';

  @override
  String get workoutPlan => 'तोकिएको योजना';

  @override
  String get assignedPlan => 'तोकिएको योजना';

  @override
  String get noAssignedPlan => 'कुनै तोकिएको योजना छैन';

  @override
  String get recentWorkouts => 'हालको व्यायामहरू';

  @override
  String get recentWorkoutLogs => 'हालको व्यायाम लगहरू';

  @override
  String get noWorkouts => 'अहिलेसम्म कुनै व्यायाम छैन';

  @override
  String get noWorkoutLogs => 'कुनै व्यायाम लग छैन';

  @override
  String get noPlan => 'कुनै व्यायाम योजना तोकिएको छैन';

  @override
  String get minutes => 'मिनेट';

  @override
  String get scanQr => 'QR स्क्यान';

  @override
  String get tapNfc => 'NFC ट्याप';

  @override
  String get checkIn => 'चेक इन';

  @override
  String get checkInSuccess => 'चेक-इन सफल भयो!';

  @override
  String get checkInFailed => 'चेक-इन असफल';

  @override
  String get streak => 'स्ट्रिक';

  @override
  String get scanQrCode => 'चेक-इन गर्न QR कोड स्क्यान गर्नुहोस्';

  @override
  String get tapNfcCard => 'NFC कार्ड नजिक फोन राख्नुहोस्';

  @override
  String get language => 'भाषा';

  @override
  String get english => 'English';

  @override
  String get nepali => 'नेपाली';

  @override
  String get tapToToggleLanguage => 'भाषा परिवर्तन गर्न ट्याप गर्नुहोस्';

  @override
  String get logout => 'लग आउट';

  @override
  String get more => 'थप';

  @override
  String get home => 'गृहपृष्ठ';

  @override
  String get date => 'मिति';

  @override
  String get time => 'समय';

  @override
  String get actions => 'कार्यहरू';

  @override
  String get previous => 'अघिल्लो';

  @override
  String get next => 'अर्को';

  @override
  String get showing => 'देखाउँदै';

  @override
  String get ofLabel => 'मध्ये';

  @override
  String get joined => 'भर्ना मिति';

  @override
  String get package => 'प्याकेज';

  @override
  String get address => 'ठेगाना';

  @override
  String get dateOfBirth => 'जन्म मिति';

  @override
  String get gender => 'लिङ्ग';

  @override
  String get male => 'पुरुष';

  @override
  String get female => 'महिला';

  @override
  String get other => 'अन्य';

  @override
  String get all => 'सबै';

  @override
  String get past => 'विगत';

  @override
  String get done => 'सकियो';

  @override
  String get title => 'शीर्षक';

  @override
  String get welcomeToUrja => 'ऊर्जामा स्वागत छ';

  @override
  String get gymManagementSimplified => 'जिम व्यवस्थापन, सरलीकृत';

  @override
  String get nfcNotAvailable => 'यो उपकरणमा NFC उपलब्ध छैन';

  @override
  String get nfcCardDetected => 'NFC कार्ड भेटियो';

  @override
  String get cardUid => 'कार्ड UID';

  @override
  String get scanAgain => 'फेरि स्क्यान गर्नुहोस्';

  @override
  String get nfcError => 'NFC त्रुटि';

  @override
  String get tryAgain => 'पुनः प्रयास गर्नुहोस्';

  @override
  String get holdPhoneNearCard => 'NFC कार्ड नजिक फोन राख्नुहोस्';

  @override
  String get nfcInstructions => 'NFC सक्षम छ भनी सुनिश्चित गर्नुहोस्';

  @override
  String get scanSuccess => 'स्क्यान सफल';

  @override
  String get qrCodeScanned => 'QR कोड स्क्यान भयो';

  @override
  String get pointCameraAtQr => 'क्यामेरा QR कोडमा देखाउनुहोस्';

  @override
  String get qrWillAutoScan => 'QR कोड स्वचालित स्क्यान हुनेछ';

  @override
  String get memberId => 'सदस्य ID';

  @override
  String get enterMemberId => 'सदस्य ID हाल्नुहोस्';

  @override
  String get noAttendance => 'कुनै उपस्थिति रेकर्ड छैन';

  @override
  String get paymentRecorded => 'भुक्तानी रेकर्ड भयो';

  @override
  String get deleteFeedbackConfirm => 'यो प्रतिक्रिया मेट्ने?';

  @override
  String get removeMemberConfirm => 'यो सदस्यलाई हटाउने?';

  @override
  String get enterCardNumber => 'कार्ड नम्बर हाल्नुहोस्';

  @override
  String get register => 'दर्ता गर्नुहोस्';

  @override
  String get cardRegistered => 'कार्ड दर्ता भयो';

  @override
  String get assignCard => 'कार्ड तोक्नुहोस्';

  @override
  String get cardAssigned => 'कार्ड तोकियो';

  @override
  String get cardUnassigned => 'कार्ड हटाइयो';

  @override
  String get nfcCardAccess => 'NFC कार्ड पहुँच';

  @override
  String get storiesNotices => 'कथा / सूचनाहरू';

  @override
  String get durationDays => 'अवधि (दिन)';

  @override
  String get featuresHint => 'gym_access, trainer, pool';

  @override
  String get noPackages => 'कुनै प्याकेज भेटिएन';

  @override
  String get noStaff => 'कुनै कर्मचारी भेटिएन';

  @override
  String get settingsSaved => 'सेटिङ सेभ भयो';

  @override
  String get logoutConfirm => 'के तपाईं लग आउट गर्न चाहनुहुन्छ?';

  @override
  String get logoutConfirmation => 'के तपाईं लग आउट गर्न चाहनुहुन्छ?';

  @override
  String get organizationProfile => 'संस्थाको प्रोफाइल';

  @override
  String get personalProfile => 'व्यक्तिगत प्रोफाइल';

  @override
  String get smsSent => 'SMS सफलतापूर्वक पठाइयो';

  @override
  String get selectAll => 'सबै छान्नुहोस्';

  @override
  String get clearAll => 'सबै हटाउनुहोस्';

  @override
  String get errorLoadingData => 'डाटा लोड गर्न त्रुटि';

  @override
  String get hello => 'नमस्ते';

  @override
  String get letsGetMoving => 'आज व्यायाम गरौं!';

  @override
  String daysStreak(int count) {
    return '$count दिन स्ट्रिक';
  }

  @override
  String get monthlyProgress => 'मासिक प्रगति';

  @override
  String get top3CheckIns => 'शीर्ष ३ चेक-इन';

  @override
  String get mySubscription => 'मेरो सदस्यता';

  @override
  String get subscriptionActive => 'सक्रिय';

  @override
  String daysLeft(int count) {
    return '$count दिन बाँकी';
  }

  @override
  String get myHomeClub => 'मेरो जिम';

  @override
  String get sendFeedback => 'हामीलाई प्रतिक्रिया दिनुहोस्';

  @override
  String get feedbackHint => 'तपाईंको अनुभव बारेमा बताउनुहोस्...';

  @override
  String get feedbackSent => 'तपाईंको प्रतिक्रियाको लागि धन्यवाद!';

  @override
  String get checkIns => 'चेक-इन';

  @override
  String get sun => 'आ';

  @override
  String get mon => 'सो';

  @override
  String get tue => 'मं';

  @override
  String get wed => 'बु';

  @override
  String get thu => 'वि';

  @override
  String get fri => 'शु';

  @override
  String get sat => 'श';

  @override
  String get account => 'खाता';

  @override
  String get editMyProfile => 'मेरो प्रोफाइल सम्पादन';

  @override
  String get subscriptionHistory => 'सदस्यता इतिहास';

  @override
  String get participationSettings => 'सहभागिता सेटिङ्स';

  @override
  String get reviewsFeedbacks => 'समीक्षा र प्रतिक्रिया';

  @override
  String get contactSupport => 'सहयोग सम्पर्क';

  @override
  String get changePhoto => 'फोटो परिवर्तन';

  @override
  String get camera => 'क्यामेरा';

  @override
  String get gallery => 'ग्यालरी';

  @override
  String get removePhoto => 'फोटो हटाउनुहोस्';

  @override
  String get myNutrition => 'मेरो पोषण';

  @override
  String get nutrition => 'पोषण';

  @override
  String get workoutMyPlan => 'मेरो योजना';

  @override
  String get workoutBrowsePlans => 'योजनाहरू हेर्नुहोस्';

  @override
  String get workoutFindMyPlan => 'मेरो योजना खोज्नुहोस्';

  @override
  String get workoutCurrentPlan => 'हालको व्यायाम योजना';

  @override
  String get workoutNoPlan => 'अहिलेसम्म कुनै व्यायाम योजना तोकिएको छैन';

  @override
  String get workoutChoosePlan => 'यो योजना छान्नुहोस्';

  @override
  String get workoutPlanAssigned => 'योजना सफलतापूर्वक तोकियो!';

  @override
  String get workoutExercises => 'व्यायामहरू';

  @override
  String get workoutSets => 'सेट';

  @override
  String get workoutReps => 'रेप्स';

  @override
  String get workoutRest => 'आराम';

  @override
  String get workoutSeconds => 'सेकेन्ड';

  @override
  String get workoutMuscleGroups => 'मांसपेशी समूहहरू';

  @override
  String get workoutFilterByGoal => 'लक्ष्यले छान्नुहोस्';

  @override
  String get workoutFilterByDifficulty => 'कठिनाइले छान्नुहोस्';

  @override
  String get workoutAllGoals => 'सबै लक्ष्य';

  @override
  String get workoutAllDifficulties => 'सबै स्तर';

  @override
  String get workoutLoseWeight => 'तौल घटाउने';

  @override
  String get workoutBuildMuscle => 'मांसपेशी बनाउने';

  @override
  String get workoutStayFit => 'फिट रहने';

  @override
  String get workoutBeginner => 'सुरुवाती';

  @override
  String get workoutIntermediate => 'मध्यम';

  @override
  String get workoutAdvanced => 'उन्नत';

  @override
  String get workoutQuestionnaire => 'छिटो मूल्यांकन';

  @override
  String get workoutWhatIsYourGoal => 'तपाईंको फिटनेस लक्ष्य के हो?';

  @override
  String get workoutExperienceLevel => 'तपाईंको अनुभव स्तर के हो?';

  @override
  String get workoutDaysPerWeek => 'हप्तामा कति दिन प्रशिक्षण गर्न सक्नुहुन्छ?';

  @override
  String get workoutDays23 => '२-३ दिन';

  @override
  String get workoutDays45 => '४-५ दिन';

  @override
  String get workoutDays67 => '६-७ दिन';

  @override
  String get workoutGetRecommendation => 'मेरो योजना पाउनुहोस्';

  @override
  String get workoutRecommendedPlan => 'सिफारिस गरिएको योजना';

  @override
  String get workoutSelfAssigned => 'आफैं छानेको';

  @override
  String get workoutQuestionnaireAssigned => 'सिफारिस गरिएको';

  @override
  String get workoutStaffAssigned => 'कर्मचारीले तोकेको';

  @override
  String get workoutTargetMuscles => 'लक्षित मांसपेशीहरू';

  @override
  String get workoutBodyMap => 'शरीर नक्सा';

  @override
  String get nutritionDailyIntake => 'दैनिक सेवन';

  @override
  String get nutritionCaloriesRemaining => 'बाँकी';

  @override
  String get nutritionCaloriesConsumed => 'खपत';

  @override
  String get nutritionCalorieGoal => 'लक्ष्य';

  @override
  String get nutritionProtein => 'प्रोटिन';

  @override
  String get nutritionCarbs => 'कार्ब्स';

  @override
  String get nutritionFat => 'बोसो';

  @override
  String get nutritionBreakfast => 'बिहानको खाना';

  @override
  String get nutritionLunch => 'दिउँसोको खाना';

  @override
  String get nutritionDinner => 'बेलुकाको खाना';

  @override
  String get nutritionSnack => 'खाजा';

  @override
  String get nutritionAddFood => 'खाना थप्नुहोस्';

  @override
  String get nutritionSearchFood => 'खाना खोज्नुहोस्...';

  @override
  String get nutritionQuantity => 'परिमाण (ग्राम)';

  @override
  String get nutritionCalories => 'क्यालोरी';

  @override
  String get nutritionNoLogs => 'अहिलेसम्म खाना लग गरिएको छैन';

  @override
  String get nutritionLogFood => 'खाना लग गर्नुहोस्';

  @override
  String get nutritionWeeklyProgress => 'साप्ताहिक प्रगति';

  @override
  String get nutritionSetGoal => 'पोषण लक्ष्य सेट गर्नुहोस्';

  @override
  String get nutritionEditGoal => 'लक्ष्य सम्पादन';

  @override
  String get nutritionWeight => 'तौल (किलो)';

  @override
  String get nutritionHeight => 'उचाइ (से.मी.)';

  @override
  String get nutritionAge => 'उमेर';

  @override
  String get nutritionGender => 'लिङ्ग';

  @override
  String get nutritionMale => 'पुरुष';

  @override
  String get nutritionFemale => 'महिला';

  @override
  String get nutritionActivityLevel => 'गतिविधि स्तर';

  @override
  String get nutritionSedentary => 'निष्क्रिय';

  @override
  String get nutritionLight => 'हल्का सक्रिय';

  @override
  String get nutritionModerate => 'मध्यम सक्रिय';

  @override
  String get nutritionActive => 'सक्रिय';

  @override
  String get nutritionVeryActive => 'धेरै सक्रिय';

  @override
  String get nutritionGoalType => 'लक्ष्य';

  @override
  String get nutritionLoseWeight => 'तौल घटाउने';

  @override
  String get nutritionMaintain => 'तौल कायम';

  @override
  String get nutritionBuildMuscle => 'मांसपेशी बनाउने';

  @override
  String get nutritionCalculateGoal => 'लक्ष्य गणना गर्नुहोस्';

  @override
  String get nutritionGoalSet => 'पोषण लक्ष्य सेट भयो!';

  @override
  String get nutritionDailyCalories => 'दैनिक क्यालोरी';

  @override
  String get nutritionMealTemplates => 'भोजन टेम्प्लेटहरू';

  @override
  String get nutritionCreateTemplate => 'टेम्प्लेट बनाउनुहोस्';

  @override
  String get nutritionTemplateName => 'टेम्प्लेट नाम';

  @override
  String get nutritionQuickLog => 'छिटो लग';

  @override
  String get nutritionNoTemplates => 'अहिलेसम्म कुनै भोजन टेम्प्लेट छैन';

  @override
  String get nutritionPer100g => 'प्रति १००ग्राम';

  @override
  String get nutritionServingSize => 'सर्भिङ';

  @override
  String get nutritionDeleteLog => 'हटाउनुहोस्';

  @override
  String get nutritionKcal => 'क्यालोरी';

  @override
  String get nutritionGrams => 'ग्राम';

  @override
  String get nutritionNoGoalSet =>
      'प्रगति ट्र्याक गर्न पोषण लक्ष्य सेट गर्नुहोस्';

  @override
  String get nutritionToday => 'आज';

  @override
  String get progressTitle => 'प्रगति';

  @override
  String get progressLogWeight => 'तौल लग गर्नुहोस्';

  @override
  String get progressWeightKg => 'तौल (किलो)';

  @override
  String get progressWeightLogged => 'तौल सफलतापूर्वक लग भयो!';

  @override
  String get progressCurrentWeight => 'हालको';

  @override
  String get progressStartWeight => 'सुरु';

  @override
  String get progressChange => 'परिवर्तन';

  @override
  String get progressBmi => 'BMI';

  @override
  String get progressWeightTrend => 'तौल प्रवृत्ति';

  @override
  String get progressNoData =>
      'अहिलेसम्म कुनै तौल डाटा छैन। पहिलो प्रविष्टि लग गर्नुहोस्!';

  @override
  String get progressPeriod30 => '३० दिन';

  @override
  String get progressPeriod90 => '९० दिन';

  @override
  String get progressPeriod180 => '१८० दिन';

  @override
  String get progressNutritionStreak => 'पोषण स्ट्रिक';

  @override
  String get progressEntries => 'प्रविष्टिहरू';

  @override
  String get progressKg => 'किलो';
}
