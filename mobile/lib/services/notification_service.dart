import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:timezone/timezone.dart' as tz;
import 'package:timezone/data/latest_all.dart' as tz;
import 'package:shared_preferences/shared_preferences.dart';

class NotificationService {
  static final NotificationService _instance = NotificationService._();
  factory NotificationService() => _instance;
  NotificationService._();

  final FlutterLocalNotificationsPlugin _plugin =
      FlutterLocalNotificationsPlugin();

  static const _waterChannelId = 'water_reminders';
  static const _foodChannelId = 'food_reminders';
  static const _prefsKey = 'notifications_enabled';

  Future<void> init() async {
    if (kIsWeb) return;
    tz.initializeTimeZones();

    const androidSettings =
        AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings(
      requestAlertPermission: true,
      requestBadgePermission: true,
      requestSoundPermission: true,
    );

    await _plugin.initialize(
      const InitializationSettings(
        android: androidSettings,
        iOS: iosSettings,
      ),
    );

    // Create notification channels
    await _plugin
        .resolvePlatformSpecificImplementation<
            AndroidFlutterLocalNotificationsPlugin>()
        ?.createNotificationChannel(const AndroidNotificationChannel(
          _waterChannelId,
          'Water Reminders',
          description: 'Periodic reminders to drink water',
          importance: Importance.defaultImportance,
        ));

    await _plugin
        .resolvePlatformSpecificImplementation<
            AndroidFlutterLocalNotificationsPlugin>()
        ?.createNotificationChannel(const AndroidNotificationChannel(
          _foodChannelId,
          'Meal Reminders',
          description: 'Reminders to log your meals',
          importance: Importance.defaultImportance,
        ));
  }

  Future<bool> requestPermission() async {
    final android = _plugin.resolvePlatformSpecificImplementation<
        AndroidFlutterLocalNotificationsPlugin>();
    if (android != null) {
      return await android.requestNotificationsPermission() ?? false;
    }
    return true;
  }

  Future<void> scheduleAllReminders() async {
    if (kIsWeb) return;
    final granted = await requestPermission();
    if (!granted) return;

    await cancelAll();
    await _scheduleWaterReminders();
    await _scheduleMealReminders();

    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(_prefsKey, true);
  }

  Future<void> cancelAll() async {
    await _plugin.cancelAll();
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(_prefsKey, false);
  }

  Future<bool> get isEnabled async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool(_prefsKey) ?? false;
  }

  // Water reminders every 2 hours from 7AM to 9PM
  Future<void> _scheduleWaterReminders() async {
    const waterTimes = [7, 9, 11, 13, 15, 17, 19, 21];
    const messages = [
      'Start your day hydrated! Drink a glass of water.',
      'Mid-morning reminder: time to drink water!',
      'Stay hydrated before lunch!',
      'Afternoon water break!',
      'Keep sipping! You\'re doing great.',
      'Evening hydration check!',
      'Don\'t forget to drink water before dinner.',
      'Last water reminder for today. Stay hydrated!',
    ];

    for (var i = 0; i < waterTimes.length; i++) {
      await _scheduleDailyNotification(
        id: 100 + i,
        channelId: _waterChannelId,
        channelName: 'Water Reminders',
        title: 'Time to Drink Water',
        body: messages[i],
        hour: waterTimes[i],
        minute: 0,
      );
    }
  }

  // Meal reminders: breakfast (7:30), lunch (12:00), snack (15:30), dinner (19:00)
  Future<void> _scheduleMealReminders() async {
    final meals = [
      (hour: 7, minute: 30, title: 'Log Breakfast', body: 'Have you had breakfast? Don\'t forget to log it!'),
      (hour: 12, minute: 0, title: 'Log Lunch', body: 'Lunchtime! Remember to log what you eat.'),
      (hour: 15, minute: 30, title: 'Afternoon Snack', body: 'Had a snack? Log it to stay on track.'),
      (hour: 19, minute: 0, title: 'Log Dinner', body: 'Dinner time! Log your meal to complete your day.'),
    ];

    for (var i = 0; i < meals.length; i++) {
      await _scheduleDailyNotification(
        id: 200 + i,
        channelId: _foodChannelId,
        channelName: 'Meal Reminders',
        title: meals[i].title,
        body: meals[i].body,
        hour: meals[i].hour,
        minute: meals[i].minute,
      );
    }
  }

  Future<void> _scheduleDailyNotification({
    required int id,
    required String channelId,
    required String channelName,
    required String title,
    required String body,
    required int hour,
    required int minute,
  }) async {
    final now = tz.TZDateTime.now(tz.local);
    var scheduled = tz.TZDateTime(tz.local, now.year, now.month, now.day, hour, minute);
    if (scheduled.isBefore(now)) {
      scheduled = scheduled.add(const Duration(days: 1));
    }

    await _plugin.zonedSchedule(
      id,
      title,
      body,
      scheduled,
      NotificationDetails(
        android: AndroidNotificationDetails(
          channelId,
          channelName,
          icon: '@mipmap/ic_launcher',
          priority: Priority.defaultPriority,
          importance: Importance.defaultImportance,
        ),
        iOS: const DarwinNotificationDetails(),
      ),
      androidScheduleMode: AndroidScheduleMode.inexactAllowWhileIdle,
      uiLocalNotificationDateInterpretation:
          UILocalNotificationDateInterpretation.absoluteTime,
      matchDateTimeComponents: DateTimeComponents.time,
    );
  }
}
