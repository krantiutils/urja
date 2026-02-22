// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appName => 'Urja';

  @override
  String get loading => 'Loading...';

  @override
  String get error => 'Something went wrong';

  @override
  String get retry => 'Retry';

  @override
  String get save => 'Save';

  @override
  String get cancel => 'Cancel';

  @override
  String get delete => 'Delete';

  @override
  String get search => 'Search...';

  @override
  String get noResults => 'No results found';

  @override
  String get signIn => 'Sign In';

  @override
  String get signOut => 'Sign Out';

  @override
  String get phoneLabel => 'Mobile Number';

  @override
  String get phonePlaceholder => '98XXXXXXXX';

  @override
  String get phoneNumber => 'Phone Number';

  @override
  String get otpLabel => 'Verification Code';

  @override
  String get otpPlaceholder => 'Enter 6-digit code';

  @override
  String get requestOtp => 'Send OTP';

  @override
  String get sendOtp => 'Send OTP';

  @override
  String get verifyOtp => 'Verify';

  @override
  String get enterOtp => 'Enter verification code';

  @override
  String get otpSent => 'OTP sent to your phone';

  @override
  String get invalidPhone => 'Enter a valid Nepali mobile number';

  @override
  String get invalidOtp => 'Invalid verification code';

  @override
  String get welcome => 'Welcome to Urja';

  @override
  String get subtitle => 'Gym management, simplified';

  @override
  String get resendOtp => 'Resend OTP';

  @override
  String get resend => 'Resend';

  @override
  String resendIn(int seconds) {
    return 'Resend in ${seconds}s';
  }

  @override
  String get didNotReceiveCode => 'Didn\'t receive the code?';

  @override
  String get dashboard => 'Dashboard';

  @override
  String get members => 'Members';

  @override
  String get packages => 'Packages';

  @override
  String get duePayments => 'Due Payments';

  @override
  String get attendance => 'Attendance';

  @override
  String get staff => 'Staff';

  @override
  String get accounts => 'Accounts';

  @override
  String get stories => 'Stories';

  @override
  String get sms => 'SMS';

  @override
  String get feedbacks => 'Feedbacks';

  @override
  String get nfcAccess => 'NFC Card Access';

  @override
  String get settings => 'Settings';

  @override
  String get totalMembers => 'Total Members';

  @override
  String get activeMembers => 'Active Members';

  @override
  String get todayAttendance => 'Today\'s Attendance';

  @override
  String get monthlyRevenue => 'Monthly Revenue';

  @override
  String get expiringPackages => 'Expiring Packages';

  @override
  String get expiredPackages => 'Expired Packages';

  @override
  String get noExpiringPackages => 'No expiring packages';

  @override
  String get todayActivity => 'Today\'s Activity';

  @override
  String get recentActivity => 'Recent Activity';

  @override
  String get noActivityToday => 'No activity today';

  @override
  String get unknownMember => 'Unknown Member';

  @override
  String get viewAll => 'View All';

  @override
  String get addMember => 'Add Member';

  @override
  String get editMember => 'Edit Member';

  @override
  String get removeMember => 'Remove Member';

  @override
  String get name => 'Name';

  @override
  String get nameNe => 'Name (Nepali)';

  @override
  String get phone => 'Phone';

  @override
  String get email => 'Email';

  @override
  String get status => 'Status';

  @override
  String get role => 'Role';

  @override
  String get active => 'Active';

  @override
  String get inactive => 'Inactive';

  @override
  String get suspended => 'Suspended';

  @override
  String get left => 'Left';

  @override
  String get member => 'Member';

  @override
  String get confirmDelete => 'Are you sure you want to delete this?';

  @override
  String get noMembers => 'No members found';

  @override
  String get addPackage => 'Add Package';

  @override
  String get editPackage => 'Edit Package';

  @override
  String get price => 'Price';

  @override
  String get duration => 'Duration';

  @override
  String get features => 'Features';

  @override
  String get days => 'days';

  @override
  String get maxMembers => 'Max Members';

  @override
  String get unlimited => 'Unlimited';

  @override
  String get description => 'Description';

  @override
  String get filterAll => 'All';

  @override
  String get filterActive => 'Active';

  @override
  String get filterInactive => 'Inactive';

  @override
  String get todayCheckIns => 'Today\'s Check-ins';

  @override
  String get manualCheckIn => 'Manual Check-in';

  @override
  String get checkInMember => 'Check in a member';

  @override
  String get method => 'Method';

  @override
  String get qr => 'QR';

  @override
  String get nfc => 'NFC';

  @override
  String get manual => 'Manual';

  @override
  String get noRecords => 'No records found';

  @override
  String get addStaff => 'Add Staff';

  @override
  String get editStaff => 'Edit Staff';

  @override
  String get staffRole => 'Staff Role';

  @override
  String get owner => 'Owner';

  @override
  String get manager => 'Manager';

  @override
  String get trainer => 'Trainer';

  @override
  String get receptionist => 'Receptionist';

  @override
  String get amount => 'Amount';

  @override
  String get dueDate => 'Due Date';

  @override
  String get unpaid => 'Unpaid';

  @override
  String get paid => 'Paid';

  @override
  String get waived => 'Waived';

  @override
  String get recordPayment => 'Record Payment';

  @override
  String get paymentMethod => 'Payment Method';

  @override
  String get cash => 'Cash';

  @override
  String get esewa => 'eSewa';

  @override
  String get bankTransfer => 'Bank Transfer';

  @override
  String get confirmPayment => 'Confirm Payment';

  @override
  String get noDues => 'No dues found';

  @override
  String get addTransaction => 'Add Transaction';

  @override
  String get editTransaction => 'Edit Transaction';

  @override
  String get category => 'Category';

  @override
  String get type => 'Type';

  @override
  String get paymentType => 'Payment Type';

  @override
  String get reference => 'Reference';

  @override
  String get income => 'Income';

  @override
  String get expense => 'Expense';

  @override
  String get totalIncome => 'Total Income';

  @override
  String get totalExpenses => 'Total Expenses';

  @override
  String get grossProfit => 'Gross Profit';

  @override
  String get profitPercent => 'Profit %';

  @override
  String get noTransactions => 'No transactions found';

  @override
  String get addNotice => 'Add Notice';

  @override
  String get editNotice => 'Edit Notice';

  @override
  String get noticeTitle => 'Title';

  @override
  String get titleNe => 'Title (Nepali)';

  @override
  String get content => 'Content';

  @override
  String get contentNe => 'Content (Nepali)';

  @override
  String get noNotices => 'No notices found';

  @override
  String get noFeedbacks => 'No feedbacks yet';

  @override
  String get cards => 'Cards';

  @override
  String get devices => 'Devices';

  @override
  String get registerCard => 'Register Card';

  @override
  String get cardNumber => 'Card Number';

  @override
  String get assignedTo => 'Assigned To';

  @override
  String get assign => 'Assign';

  @override
  String get unassign => 'Unassign';

  @override
  String get unassigned => 'Unassigned';

  @override
  String get available => 'Available';

  @override
  String get assigned => 'Assigned';

  @override
  String get disabled => 'Disabled';

  @override
  String get noCards => 'No NFC cards registered';

  @override
  String get deviceName => 'Device Name';

  @override
  String get online => 'Online';

  @override
  String get offline => 'Offline';

  @override
  String get locked => 'Locked';

  @override
  String get unlocked => 'Unlocked';

  @override
  String get noDevices => 'No NFC devices found';

  @override
  String get door => 'Door';

  @override
  String get uptime => 'Uptime';

  @override
  String get smsBalance => 'SMS Balance';

  @override
  String get remaining => 'Remaining';

  @override
  String get totalPurchased => 'Total Purchased';

  @override
  String get totalUsed => 'Total Used';

  @override
  String get campaigns => 'Campaigns';

  @override
  String get sendSms => 'Send SMS';

  @override
  String get message => 'Message';

  @override
  String get enterMessage => 'Enter your message';

  @override
  String get selectMembers => 'Select Members';

  @override
  String membersSelected(int count) {
    return '$count members selected';
  }

  @override
  String get send => 'Send';

  @override
  String get history => 'History';

  @override
  String get noHistory => 'No history';

  @override
  String get orgProfile => 'Organization Profile';

  @override
  String get myProfile => 'My Profile';

  @override
  String get profileInfo => 'Profile Information';

  @override
  String get noProfile => 'No profile data';

  @override
  String get privacySettings => 'Privacy Settings';

  @override
  String get showEmail => 'Show email to others';

  @override
  String get showEmailDesc => 'Allow other members to see your email';

  @override
  String get showPhone => 'Show phone to others';

  @override
  String get showPhoneDesc => 'Allow other members to see your phone';

  @override
  String get showProfile => 'Show profile publicly';

  @override
  String get showProfileDesc => 'Make your profile visible to others';

  @override
  String get showAttendance => 'Show attendance to others';

  @override
  String get showAttendanceDesc => 'Allow others to see your attendance';

  @override
  String get showOnLeaderboard => 'Show on leaderboard';

  @override
  String get showOnLeaderboardDesc => 'Display on the gym leaderboard';

  @override
  String get emergencyContact => 'Emergency Contact';

  @override
  String get contactName => 'Contact Name';

  @override
  String get contactPhone => 'Contact Phone';

  @override
  String get saveChanges => 'Save Changes';

  @override
  String get saveProfile => 'Save Profile';

  @override
  String get savePrivacy => 'Save Privacy';

  @override
  String get saveEmergencyContact => 'Save Emergency Contact';

  @override
  String get saved => 'Changes saved';

  @override
  String get myDashboard => 'My Dashboard';

  @override
  String get myAttendance => 'My Attendance';

  @override
  String get myPackages => 'My Packages';

  @override
  String get myHealth => 'My Health';

  @override
  String get myWorkouts => 'My Workouts';

  @override
  String get profile => 'Profile';

  @override
  String get welcomeBack => 'Welcome back';

  @override
  String get memberSince => 'Member since';

  @override
  String get currentPackage => 'Current Package';

  @override
  String get activePackage => 'Active Package';

  @override
  String get noActivePackage => 'No active package';

  @override
  String get noActivePackages => 'No active packages';

  @override
  String get noPastPackages => 'No past packages';

  @override
  String get daysRemaining => 'days remaining';

  @override
  String get expired => 'Expired';

  @override
  String get expires => 'Expires';

  @override
  String get currentStreak => 'Current Streak';

  @override
  String get longestStreak => 'Longest Streak';

  @override
  String get recentCheckIns => 'Recent Check-ins';

  @override
  String get noRecentCheckIns => 'No recent check-ins';

  @override
  String get noCheckIns => 'No check-ins yet';

  @override
  String get noAttendanceRecords => 'No attendance records';

  @override
  String get attendanceHistory => 'Attendance History';

  @override
  String get bmi => 'BMI';

  @override
  String get measurements => 'Body Measurements';

  @override
  String get bodyMeasurements => 'Body Measurements';

  @override
  String get noHealthData => 'No health data recorded yet';

  @override
  String get noBmiData => 'No BMI data';

  @override
  String get noMeasurementData => 'No measurement data';

  @override
  String get height => 'Height';

  @override
  String get weight => 'Weight';

  @override
  String get chest => 'Chest';

  @override
  String get waist => 'Waist';

  @override
  String get hips => 'Hips';

  @override
  String get biceps => 'Biceps';

  @override
  String get thighs => 'Thighs';

  @override
  String get bodyFat => 'Body Fat';

  @override
  String get recordedOn => 'Recorded on';

  @override
  String get workoutPlan => 'Assigned Plan';

  @override
  String get assignedPlan => 'Assigned Plan';

  @override
  String get noAssignedPlan => 'No assigned plan';

  @override
  String get recentWorkouts => 'Recent Workouts';

  @override
  String get recentWorkoutLogs => 'Recent Workout Logs';

  @override
  String get noWorkouts => 'No workouts logged yet';

  @override
  String get noWorkoutLogs => 'No workout logs';

  @override
  String get noPlan => 'No workout plan assigned';

  @override
  String get minutes => 'min';

  @override
  String get scanQr => 'Scan QR';

  @override
  String get tapNfc => 'Tap NFC';

  @override
  String get checkIn => 'Check In';

  @override
  String get checkInSuccess => 'Check-in successful!';

  @override
  String get checkInFailed => 'Check-in failed';

  @override
  String get streak => 'Streak';

  @override
  String get scanQrCode => 'Scan QR Code to Check In';

  @override
  String get tapNfcCard => 'Hold your phone near the NFC card';

  @override
  String get language => 'Language';

  @override
  String get english => 'English';

  @override
  String get nepali => 'Nepali';

  @override
  String get tapToToggleLanguage => 'Tap to switch language';

  @override
  String get logout => 'Logout';

  @override
  String get more => 'More';

  @override
  String get home => 'Home';

  @override
  String get date => 'Date';

  @override
  String get time => 'Time';

  @override
  String get actions => 'Actions';

  @override
  String get previous => 'Previous';

  @override
  String get next => 'Next';

  @override
  String get showing => 'Showing';

  @override
  String get ofLabel => 'of';

  @override
  String get joined => 'Joined';

  @override
  String get package => 'Package';

  @override
  String get address => 'Address';

  @override
  String get dateOfBirth => 'Date of Birth';

  @override
  String get gender => 'Gender';

  @override
  String get male => 'Male';

  @override
  String get female => 'Female';

  @override
  String get other => 'Other';

  @override
  String get all => 'All';

  @override
  String get past => 'Past';

  @override
  String get done => 'Done';

  @override
  String get title => 'Title';

  @override
  String get welcomeToUrja => 'Welcome to Urja';

  @override
  String get gymManagementSimplified => 'Gym management, simplified';

  @override
  String get nfcNotAvailable => 'NFC is not available on this device';

  @override
  String get nfcCardDetected => 'NFC Card Detected';

  @override
  String get cardUid => 'Card UID';

  @override
  String get scanAgain => 'Scan Again';

  @override
  String get nfcError => 'NFC Error';

  @override
  String get tryAgain => 'Try Again';

  @override
  String get holdPhoneNearCard => 'Hold your phone near the NFC card';

  @override
  String get nfcInstructions => 'Make sure NFC is enabled on your device';

  @override
  String get scanSuccess => 'Scan Successful';

  @override
  String get qrCodeScanned => 'QR Code Scanned';

  @override
  String get pointCameraAtQr => 'Point camera at QR code';

  @override
  String get qrWillAutoScan => 'QR code will be scanned automatically';

  @override
  String get memberId => 'Member ID';

  @override
  String get enterMemberId => 'Enter member ID';

  @override
  String get noAttendance => 'No attendance records';

  @override
  String get paymentRecorded => 'Payment recorded';

  @override
  String get deleteFeedbackConfirm => 'Delete this feedback?';

  @override
  String get removeMemberConfirm => 'Remove this member?';

  @override
  String get enterCardNumber => 'Enter card number';

  @override
  String get register => 'Register';

  @override
  String get cardRegistered => 'Card registered';

  @override
  String get assignCard => 'Assign Card';

  @override
  String get cardAssigned => 'Card assigned';

  @override
  String get cardUnassigned => 'Card unassigned';

  @override
  String get nfcCardAccess => 'NFC Card Access';

  @override
  String get storiesNotices => 'Stories / Notices';

  @override
  String get durationDays => 'Duration (days)';

  @override
  String get featuresHint => 'gym_access, trainer, pool';

  @override
  String get noPackages => 'No packages found';

  @override
  String get noStaff => 'No staff found';

  @override
  String get settingsSaved => 'Settings saved';

  @override
  String get logoutConfirm => 'Are you sure you want to logout?';

  @override
  String get logoutConfirmation => 'Are you sure you want to logout?';

  @override
  String get organizationProfile => 'Organization Profile';

  @override
  String get personalProfile => 'Personal Profile';

  @override
  String get smsSent => 'SMS sent successfully';

  @override
  String get selectAll => 'Select All';

  @override
  String get clearAll => 'Clear All';

  @override
  String get errorLoadingData => 'Error loading data';

  @override
  String get hello => 'Hello';

  @override
  String get letsGetMoving => 'Let\'s get you moving today!';

  @override
  String daysStreak(int count) {
    return '$count Days Streak';
  }

  @override
  String get monthlyProgress => 'Monthly Progress';

  @override
  String get top3CheckIns => 'Top 3 Check-Ins';

  @override
  String get mySubscription => 'My Subscription';

  @override
  String get subscriptionActive => 'ACTIVE';

  @override
  String daysLeft(int count) {
    return '$count days left';
  }

  @override
  String get myHomeClub => 'My Home Club';

  @override
  String get sendFeedback => 'Send Us Feedback';

  @override
  String get feedbackHint => 'Tell us about your experience...';

  @override
  String get feedbackSent => 'Thank you for your feedback!';

  @override
  String get checkIns => 'check-ins';

  @override
  String get sun => 'S';

  @override
  String get mon => 'M';

  @override
  String get tue => 'T';

  @override
  String get wed => 'W';

  @override
  String get thu => 'T';

  @override
  String get fri => 'F';

  @override
  String get sat => 'S';

  @override
  String get account => 'Account';

  @override
  String get editMyProfile => 'Edit My Profile';

  @override
  String get subscriptionHistory => 'Subscription History';

  @override
  String get participationSettings => 'Participation Settings';

  @override
  String get reviewsFeedbacks => 'Reviews & Feedbacks';

  @override
  String get contactSupport => 'Contact Support';

  @override
  String get changePhoto => 'Change Photo';

  @override
  String get camera => 'Camera';

  @override
  String get gallery => 'Gallery';

  @override
  String get removePhoto => 'Remove Photo';

  @override
  String get myNutrition => 'My Nutrition';

  @override
  String get nutrition => 'Nutrition';

  @override
  String get workoutMyPlan => 'My Plan';

  @override
  String get workoutBrowsePlans => 'Browse Plans';

  @override
  String get workoutFindMyPlan => 'Find My Plan';

  @override
  String get workoutCurrentPlan => 'Current Workout Plan';

  @override
  String get workoutNoPlan => 'No workout plan assigned yet';

  @override
  String get workoutChoosePlan => 'Choose This Plan';

  @override
  String get workoutPlanAssigned => 'Plan assigned successfully!';

  @override
  String get workoutExercises => 'Exercises';

  @override
  String get workoutSets => 'sets';

  @override
  String get workoutReps => 'reps';

  @override
  String get workoutRest => 'rest';

  @override
  String get workoutSeconds => 'sec';

  @override
  String get workoutMuscleGroups => 'Muscle Groups';

  @override
  String get workoutFilterByGoal => 'Filter by Goal';

  @override
  String get workoutFilterByDifficulty => 'Filter by Difficulty';

  @override
  String get workoutAllGoals => 'All Goals';

  @override
  String get workoutAllDifficulties => 'All Levels';

  @override
  String get workoutLoseWeight => 'Lose Weight';

  @override
  String get workoutBuildMuscle => 'Build Muscle';

  @override
  String get workoutStayFit => 'Stay Fit';

  @override
  String get workoutBeginner => 'Beginner';

  @override
  String get workoutIntermediate => 'Intermediate';

  @override
  String get workoutAdvanced => 'Advanced';

  @override
  String get workoutQuestionnaire => 'Quick Assessment';

  @override
  String get workoutWhatIsYourGoal => 'What is your fitness goal?';

  @override
  String get workoutExperienceLevel => 'What is your experience level?';

  @override
  String get workoutDaysPerWeek => 'How many days per week can you train?';

  @override
  String get workoutDays23 => '2-3 days';

  @override
  String get workoutDays45 => '4-5 days';

  @override
  String get workoutDays67 => '6-7 days';

  @override
  String get workoutGetRecommendation => 'Get My Plan';

  @override
  String get workoutRecommendedPlan => 'Recommended Plan';

  @override
  String get workoutSelfAssigned => 'Self-assigned';

  @override
  String get workoutQuestionnaireAssigned => 'Recommended';

  @override
  String get workoutStaffAssigned => 'Assigned by staff';

  @override
  String get workoutTargetMuscles => 'Target Muscles';

  @override
  String get workoutBodyMap => 'Body Map';

  @override
  String get nutritionDailyIntake => 'Daily Intake';

  @override
  String get nutritionCaloriesRemaining => 'Remaining';

  @override
  String get nutritionCaloriesConsumed => 'Consumed';

  @override
  String get nutritionCalorieGoal => 'Goal';

  @override
  String get nutritionProtein => 'Protein';

  @override
  String get nutritionCarbs => 'Carbs';

  @override
  String get nutritionFat => 'Fat';

  @override
  String get nutritionBreakfast => 'Breakfast';

  @override
  String get nutritionLunch => 'Lunch';

  @override
  String get nutritionDinner => 'Dinner';

  @override
  String get nutritionSnack => 'Snack';

  @override
  String get nutritionAddFood => 'Add Food';

  @override
  String get nutritionSearchFood => 'Search food...';

  @override
  String get nutritionQuantity => 'Quantity (g)';

  @override
  String get nutritionCalories => 'Calories';

  @override
  String get nutritionNoLogs => 'No food logged yet';

  @override
  String get nutritionLogFood => 'Log Food';

  @override
  String get nutritionWeeklyProgress => 'Weekly Progress';

  @override
  String get nutritionSetGoal => 'Set Nutrition Goal';

  @override
  String get nutritionEditGoal => 'Edit Goal';

  @override
  String get nutritionWeight => 'Weight (kg)';

  @override
  String get nutritionHeight => 'Height (cm)';

  @override
  String get nutritionAge => 'Age';

  @override
  String get nutritionGender => 'Gender';

  @override
  String get nutritionMale => 'Male';

  @override
  String get nutritionFemale => 'Female';

  @override
  String get nutritionActivityLevel => 'Activity Level';

  @override
  String get nutritionSedentary => 'Sedentary';

  @override
  String get nutritionLight => 'Lightly Active';

  @override
  String get nutritionModerate => 'Moderately Active';

  @override
  String get nutritionActive => 'Active';

  @override
  String get nutritionVeryActive => 'Very Active';

  @override
  String get nutritionGoalType => 'Goal';

  @override
  String get nutritionLoseWeight => 'Lose Weight';

  @override
  String get nutritionMaintain => 'Maintain Weight';

  @override
  String get nutritionBuildMuscle => 'Build Muscle';

  @override
  String get nutritionCalculateGoal => 'Calculate Goal';

  @override
  String get nutritionGoalSet => 'Nutrition goal set!';

  @override
  String get nutritionDailyCalories => 'Daily Calories';

  @override
  String get nutritionMealTemplates => 'Meal Templates';

  @override
  String get nutritionCreateTemplate => 'Create Template';

  @override
  String get nutritionTemplateName => 'Template Name';

  @override
  String get nutritionQuickLog => 'Quick Log';

  @override
  String get nutritionNoTemplates => 'No meal templates yet';

  @override
  String get nutritionPer100g => 'per 100g';

  @override
  String get nutritionServingSize => 'Serving';

  @override
  String get nutritionDeleteLog => 'Remove';

  @override
  String get nutritionKcal => 'kcal';

  @override
  String get nutritionGrams => 'g';

  @override
  String get nutritionNoGoalSet =>
      'Set up your nutrition goal to track progress';

  @override
  String get nutritionToday => 'Today';

  @override
  String get progressTitle => 'Progress';

  @override
  String get progressLogWeight => 'Log Weight';

  @override
  String get progressWeightKg => 'Weight (kg)';

  @override
  String get progressWeightLogged => 'Weight logged successfully!';

  @override
  String get progressCurrentWeight => 'Current';

  @override
  String get progressStartWeight => 'Start';

  @override
  String get progressChange => 'Change';

  @override
  String get progressBmi => 'BMI';

  @override
  String get progressWeightTrend => 'Weight Trend';

  @override
  String get progressNoData => 'No weight data yet. Log your first entry!';

  @override
  String get progressPeriod30 => '30 Days';

  @override
  String get progressPeriod90 => '90 Days';

  @override
  String get progressPeriod180 => '180 Days';

  @override
  String get progressNutritionStreak => 'Nutrition Streak';

  @override
  String get progressEntries => 'entries';

  @override
  String get progressKg => 'kg';
}
