// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Urja';

  @override
  String get login => 'Login';

  @override
  String get enterPhoneNumber => 'Enter your phone number';

  @override
  String get phoneNumber => 'Phone Number';

  @override
  String get phoneHint => '98XXXXXXXX';

  @override
  String get sendOtp => 'Send OTP';

  @override
  String get verifyOtp => 'Verify OTP';

  @override
  String enterOtp(String phone) {
    return 'Enter the 6-digit OTP sent to $phone';
  }

  @override
  String get otpHint => 'Enter 6-digit code';

  @override
  String get verify => 'Verify';

  @override
  String get resendOtp => 'Resend OTP';

  @override
  String get invalidPhone => 'Enter a valid 10-digit phone number';

  @override
  String get invalidOtp => 'Enter a valid 6-digit OTP';

  @override
  String get otpSent => 'OTP sent successfully';

  @override
  String get loginSuccess => 'Login successful';

  @override
  String get logout => 'Logout';

  @override
  String get logoutAll => 'Logout from all devices';

  @override
  String get logoutConfirm => 'Are you sure you want to logout?';

  @override
  String get cancel => 'Cancel';

  @override
  String get confirm => 'Confirm';

  @override
  String get home => 'Home';

  @override
  String get checkIn => 'Check In';

  @override
  String get attendance => 'Attendance';

  @override
  String get myPackage => 'My Package';

  @override
  String get profile => 'Profile';

  @override
  String get settings => 'Settings';

  @override
  String get notifications => 'Notifications';

  @override
  String get leaderboard => 'Leaderboard';

  @override
  String get achievements => 'Achievements';

  @override
  String get health => 'Health';

  @override
  String get workouts => 'Workouts';

  @override
  String get guides => 'Guides';

  @override
  String greeting(String name) {
    return 'Hello, $name!';
  }

  @override
  String get todayStatus => 'Today\'s Status';

  @override
  String get checkedIn => 'Checked In';

  @override
  String get notCheckedIn => 'Not Checked In';

  @override
  String get currentStreak => 'Current Streak';

  @override
  String get bestStreak => 'Best Streak';

  @override
  String daysLeft(int days) {
    return '$days days left';
  }

  @override
  String get days => 'days';

  @override
  String get packageExpired => 'Package Expired';

  @override
  String get noActivePackage => 'No Active Package';

  @override
  String get scanQrCode => 'Scan QR Code';

  @override
  String get tapNfc => 'Tap NFC Card';

  @override
  String get scanGymQr => 'Point your camera at the gym\'s QR code';

  @override
  String get checkInSuccess => 'Check-in successful!';

  @override
  String get checkInFailed => 'Check-in failed';

  @override
  String streakUpdated(int count) {
    return 'Streak: $count days';
  }

  @override
  String get attendanceHistory => 'Attendance History';

  @override
  String get noAttendance => 'No attendance records yet';

  @override
  String get date => 'Date';

  @override
  String get time => 'Time';

  @override
  String get method => 'Method';

  @override
  String get qr => 'QR';

  @override
  String get nfc => 'NFC';

  @override
  String get manual => 'Manual';

  @override
  String get currentPackage => 'Current Package';

  @override
  String get paymentHistory => 'Payment History';

  @override
  String get renewPackage => 'Renew Package';

  @override
  String get subscribe => 'Subscribe';

  @override
  String get noPackages => 'No packages available';

  @override
  String get editProfile => 'Edit Profile';

  @override
  String get name => 'Name';

  @override
  String get email => 'Email';

  @override
  String get dateOfBirth => 'Date of Birth';

  @override
  String get gender => 'Gender';

  @override
  String get emergencyContact => 'Emergency Contact';

  @override
  String get emergencyPhone => 'Emergency Phone';

  @override
  String get saveChanges => 'Save Changes';

  @override
  String get profileUpdated => 'Profile updated successfully';

  @override
  String get weekly => 'Weekly';

  @override
  String get monthly => 'Monthly';

  @override
  String get allTime => 'All Time';

  @override
  String get rank => 'Rank';

  @override
  String get checkIns => 'Check-ins';

  @override
  String get bmi => 'BMI';

  @override
  String get whr => 'WHR';

  @override
  String get weight => 'Weight';

  @override
  String get bodyMeasurements => 'Body Measurements';

  @override
  String get progressPhotos => 'Progress Photos';

  @override
  String get logBmi => 'Log BMI';

  @override
  String get heightCm => 'Height (cm)';

  @override
  String get weightKg => 'Weight (kg)';

  @override
  String get save => 'Save';

  @override
  String get noData => 'No data available';

  @override
  String get offline => 'Offline';

  @override
  String get offlineMode => 'You are offline. Data will sync when connected.';

  @override
  String get syncing => 'Syncing...';

  @override
  String get syncComplete => 'Sync complete';

  @override
  String get language => 'Language';

  @override
  String get english => 'English';

  @override
  String get nepali => 'Nepali';

  @override
  String get privacySettings => 'Privacy Settings';

  @override
  String get showEmail => 'Show Email';

  @override
  String get showPhone => 'Show Phone';

  @override
  String get showProfile => 'Show Profile';

  @override
  String get showAttendance => 'Show Attendance';

  @override
  String get showOnLeaderboard => 'Show on Leaderboard';

  @override
  String get error => 'Error';

  @override
  String get retry => 'Retry';

  @override
  String get loading => 'Loading...';

  @override
  String get noInternet => 'No internet connection';

  @override
  String get serverError => 'Server error. Please try again.';

  @override
  String get sessionExpired => 'Session expired. Please login again.';

  @override
  String get male => 'Male';

  @override
  String get female => 'Female';

  @override
  String get other => 'Other';

  @override
  String get preferNotSay => 'Prefer not to say';
}
