// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Nepali (`ne`).
class AppLocalizationsNe extends AppLocalizations {
  AppLocalizationsNe([String locale = 'ne']) : super(locale);

  @override
  String get appTitle => 'ऊर्जा';

  @override
  String get login => 'लग इन';

  @override
  String get enterPhoneNumber => 'आफ्नो फोन नम्बर प्रविष्ट गर्नुहोस्';

  @override
  String get phoneNumber => 'फोन नम्बर';

  @override
  String get phoneHint => '98XXXXXXXX';

  @override
  String get sendOtp => 'OTP पठाउनुहोस्';

  @override
  String get verifyOtp => 'OTP प्रमाणित गर्नुहोस्';

  @override
  String enterOtp(String phone) {
    return '$phone मा पठाइएको ६ अङ्कको OTP प्रविष्ट गर्नुहोस्';
  }

  @override
  String get otpHint => '६ अङ्कको कोड प्रविष्ट गर्नुहोस्';

  @override
  String get verify => 'प्रमाणित गर्नुहोस्';

  @override
  String get resendOtp => 'OTP पुन: पठाउनुहोस्';

  @override
  String get invalidPhone => 'मान्य १० अङ्कको फोन नम्बर प्रविष्ट गर्नुहोस्';

  @override
  String get invalidOtp => 'मान्य ६ अङ्कको OTP प्रविष्ट गर्नुहोस्';

  @override
  String get otpSent => 'OTP सफलतापूर्वक पठाइयो';

  @override
  String get loginSuccess => 'लगइन सफल';

  @override
  String get logout => 'लग आउट';

  @override
  String get logoutAll => 'सबै उपकरणबाट लग आउट';

  @override
  String get logoutConfirm => 'के तपाईं लग आउट गर्न चाहनुहुन्छ?';

  @override
  String get cancel => 'रद्द गर्नुहोस्';

  @override
  String get confirm => 'पुष्टि गर्नुहोस्';

  @override
  String get home => 'गृहपृष्ठ';

  @override
  String get checkIn => 'चेक इन';

  @override
  String get attendance => 'उपस्थिति';

  @override
  String get myPackage => 'मेरो प्याकेज';

  @override
  String get profile => 'प्रोफाइल';

  @override
  String get settings => 'सेटिङ्स';

  @override
  String get notifications => 'सूचनाहरू';

  @override
  String get leaderboard => 'लिडरबोर्ड';

  @override
  String get achievements => 'उपलब्धिहरू';

  @override
  String get health => 'स्वास्थ्य';

  @override
  String get workouts => 'व्यायाम';

  @override
  String get guides => 'गाइडहरू';

  @override
  String greeting(String name) {
    return 'नमस्ते, $name!';
  }

  @override
  String get todayStatus => 'आजको स्थिति';

  @override
  String get checkedIn => 'चेक इन भयो';

  @override
  String get notCheckedIn => 'चेक इन भएको छैन';

  @override
  String get currentStreak => 'हालको स्ट्रीक';

  @override
  String get bestStreak => 'उत्तम स्ट्रीक';

  @override
  String daysLeft(int days) {
    return '$days दिन बाँकी';
  }

  @override
  String get days => 'दिन';

  @override
  String get packageExpired => 'प्याकेज समाप्त भयो';

  @override
  String get noActivePackage => 'कुनै सक्रिय प्याकेज छैन';

  @override
  String get scanQrCode => 'QR कोड स्क्यान गर्नुहोस्';

  @override
  String get tapNfc => 'NFC कार्ड ट्याप गर्नुहोस्';

  @override
  String get scanGymQr => 'जिमको QR कोडमा आफ्नो क्यामेरा राख्नुहोस्';

  @override
  String get checkInSuccess => 'चेक-इन सफल!';

  @override
  String get checkInFailed => 'चेक-इन असफल';

  @override
  String streakUpdated(int count) {
    return 'स्ट्रीक: $count दिन';
  }

  @override
  String get attendanceHistory => 'उपस्थिति इतिहास';

  @override
  String get noAttendance => 'अहिलेसम्म कुनै उपस्थिति रेकर्ड छैन';

  @override
  String get date => 'मिति';

  @override
  String get time => 'समय';

  @override
  String get method => 'विधि';

  @override
  String get qr => 'QR';

  @override
  String get nfc => 'NFC';

  @override
  String get manual => 'म्यानुअल';

  @override
  String get currentPackage => 'हालको प्याकेज';

  @override
  String get paymentHistory => 'भुक्तानी इतिहास';

  @override
  String get renewPackage => 'प्याकेज नवीकरण';

  @override
  String get subscribe => 'सदस्यता लिनुहोस्';

  @override
  String get noPackages => 'कुनै प्याकेज उपलब्ध छैन';

  @override
  String get editProfile => 'प्रोफाइल सम्पादन';

  @override
  String get name => 'नाम';

  @override
  String get email => 'इमेल';

  @override
  String get dateOfBirth => 'जन्म मिति';

  @override
  String get gender => 'लिङ्ग';

  @override
  String get emergencyContact => 'आपतकालीन सम्पर्क';

  @override
  String get emergencyPhone => 'आपतकालीन फोन';

  @override
  String get saveChanges => 'परिवर्तनहरू सुरक्षित गर्नुहोस्';

  @override
  String get profileUpdated => 'प्रोफाइल सफलतापूर्वक अपडेट भयो';

  @override
  String get weekly => 'साप्ताहिक';

  @override
  String get monthly => 'मासिक';

  @override
  String get allTime => 'सबै समय';

  @override
  String get rank => 'र्‍यांक';

  @override
  String get checkIns => 'चेक-इनहरू';

  @override
  String get bmi => 'BMI';

  @override
  String get whr => 'WHR';

  @override
  String get weight => 'तौल';

  @override
  String get bodyMeasurements => 'शरीरको नाप';

  @override
  String get progressPhotos => 'प्रगति फोटोहरू';

  @override
  String get logBmi => 'BMI रेकर्ड गर्नुहोस्';

  @override
  String get heightCm => 'उचाइ (से.मी.)';

  @override
  String get weightKg => 'तौल (के.जी.)';

  @override
  String get save => 'सुरक्षित गर्नुहोस्';

  @override
  String get noData => 'कुनै डाटा उपलब्ध छैन';

  @override
  String get offline => 'अफलाइन';

  @override
  String get offlineMode =>
      'तपाईं अफलाइन हुनुहुन्छ। जडान हुँदा डाटा सिंक हुनेछ।';

  @override
  String get syncing => 'सिंक हुँदैछ...';

  @override
  String get syncComplete => 'सिंक पूरा भयो';

  @override
  String get language => 'भाषा';

  @override
  String get english => 'English';

  @override
  String get nepali => 'नेपाली';

  @override
  String get privacySettings => 'गोपनीयता सेटिङ्स';

  @override
  String get showEmail => 'इमेल देखाउनुहोस्';

  @override
  String get showPhone => 'फोन देखाउनुहोस्';

  @override
  String get showProfile => 'प्रोफाइल देखाउनुहोस्';

  @override
  String get showAttendance => 'उपस्थिति देखाउनुहोस्';

  @override
  String get showOnLeaderboard => 'लिडरबोर्डमा देखाउनुहोस्';

  @override
  String get error => 'त्रुटि';

  @override
  String get retry => 'पुन: प्रयास गर्नुहोस्';

  @override
  String get loading => 'लोड हुँदैछ...';

  @override
  String get noInternet => 'इन्टरनेट जडान छैन';

  @override
  String get serverError => 'सर्भर त्रुटि। कृपया पुन: प्रयास गर्नुहोस्।';

  @override
  String get sessionExpired => 'सत्र समाप्त भयो। कृपया पुन: लग इन गर्नुहोस्।';

  @override
  String get male => 'पुरुष';

  @override
  String get female => 'महिला';

  @override
  String get other => 'अन्य';

  @override
  String get preferNotSay => 'भन्न नचाहने';
}
