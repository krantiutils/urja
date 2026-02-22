import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../config/theme.dart';

enum HiitPhase { idle, work, rest, finished }

class HiitTimerScreen extends StatefulWidget {
  const HiitTimerScreen({super.key});

  @override
  State<HiitTimerScreen> createState() => _HiitTimerScreenState();
}

class _HiitTimerScreenState extends State<HiitTimerScreen> {
  // Configuration
  int _workSeconds = 30;
  int _restSeconds = 15;
  int _totalRounds = 8;

  // Timer state
  HiitPhase _phase = HiitPhase.idle;
  int _currentRound = 1;
  int _remainingSeconds = 0;
  Timer? _timer;
  bool _isPaused = false;

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  void _start() {
    setState(() {
      _phase = HiitPhase.work;
      _currentRound = 1;
      _remainingSeconds = _workSeconds;
      _isPaused = false;
    });
    _startCountdown();
  }

  void _startCountdown() {
    _timer?.cancel();
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (_isPaused) return;

      setState(() {
        _remainingSeconds--;
      });

      if (_remainingSeconds <= 0) {
        _onPhaseComplete();
      } else if (_remainingSeconds <= 3) {
        // Countdown beep for last 3 seconds
        SystemSound.play(SystemSoundType.click);
      }
    });
  }

  void _onPhaseComplete() {
    SystemSound.play(SystemSoundType.click);

    if (_phase == HiitPhase.work) {
      if (_currentRound >= _totalRounds) {
        // All rounds complete
        _timer?.cancel();
        setState(() {
          _phase = HiitPhase.finished;
        });
        return;
      }
      // Switch to rest
      setState(() {
        _phase = HiitPhase.rest;
        _remainingSeconds = _restSeconds;
      });
    } else if (_phase == HiitPhase.rest) {
      // Switch to next work round
      setState(() {
        _currentRound++;
        _phase = HiitPhase.work;
        _remainingSeconds = _workSeconds;
      });
    }
  }

  void _togglePause() {
    setState(() {
      _isPaused = !_isPaused;
    });
  }

  void _reset() {
    _timer?.cancel();
    setState(() {
      _phase = HiitPhase.idle;
      _currentRound = 1;
      _remainingSeconds = 0;
      _isPaused = false;
    });
  }

  String _formatTime(int seconds) {
    final mins = seconds ~/ 60;
    final secs = seconds % 60;
    return '${mins.toString().padLeft(2, '0')}:${secs.toString().padLeft(2, '0')}';
  }

  Color get _phaseColor {
    switch (_phase) {
      case HiitPhase.work:
        return AppTheme.primary;
      case HiitPhase.rest:
        return AppTheme.warning;
      case HiitPhase.finished:
        return AppTheme.info;
      case HiitPhase.idle:
        return AppTheme.textSecondary;
    }
  }

  String get _phaseLabel {
    switch (_phase) {
      case HiitPhase.work:
        return 'WORK';
      case HiitPhase.rest:
        return 'REST';
      case HiitPhase.finished:
        return 'FINISHED';
      case HiitPhase.idle:
        return 'READY';
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(
        title: const Text('HIIT Timer'),
      ),
      body: SafeArea(
        child: _phase == HiitPhase.idle
            ? _buildConfigView()
            : _buildTimerView(),
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Configuration view
  // ---------------------------------------------------------------------------
  Widget _buildConfigView() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          const SizedBox(height: 20),
          // Timer icon
          Container(
            width: 100,
            height: 100,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: AppTheme.primary.withAlpha(20),
            ),
            child: const Icon(
              Icons.timer,
              size: 48,
              color: AppTheme.primary,
            ),
          ),
          const SizedBox(height: 24),
          const Text(
            'Configure Timer',
            style: TextStyle(
              fontSize: 22,
              fontWeight: FontWeight.bold,
              color: AppTheme.textPrimary,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            'Set your work, rest, and round parameters',
            style: TextStyle(
              color: AppTheme.textSecondary,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 32),

          // Work interval
          _buildSliderCard(
            label: 'Work Interval',
            value: _workSeconds,
            min: 20,
            max: 60,
            color: AppTheme.primary,
            icon: Icons.directions_run,
            unit: 'sec',
            onChanged: (val) => setState(() => _workSeconds = val.round()),
          ),
          const SizedBox(height: 12),

          // Rest interval
          _buildSliderCard(
            label: 'Rest Interval',
            value: _restSeconds,
            min: 10,
            max: 30,
            color: AppTheme.warning,
            icon: Icons.pause_circle_outline,
            unit: 'sec',
            onChanged: (val) => setState(() => _restSeconds = val.round()),
          ),
          const SizedBox(height: 12),

          // Rounds
          _buildSliderCard(
            label: 'Rounds',
            value: _totalRounds,
            min: 4,
            max: 12,
            color: AppTheme.info,
            icon: Icons.repeat,
            unit: '',
            onChanged: (val) => setState(() => _totalRounds = val.round()),
          ),
          const SizedBox(height: 16),

          // Total time estimate
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  const Icon(Icons.schedule,
                      color: AppTheme.textSecondary, size: 20),
                  const SizedBox(width: 10),
                  const Text(
                    'Total Time: ',
                    style: TextStyle(
                      color: AppTheme.textSecondary,
                      fontSize: 14,
                    ),
                  ),
                  Text(
                    _formatTime(
                        (_workSeconds * _totalRounds) +
                            (_restSeconds * (_totalRounds - 1))),
                    style: const TextStyle(
                      color: AppTheme.textPrimary,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 32),

          // Start button
          SizedBox(
            width: double.infinity,
            height: 56,
            child: ElevatedButton.icon(
              onPressed: _start,
              icon: const Icon(Icons.play_arrow, size: 28),
              label: const Text(
                'START',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  letterSpacing: 2,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSliderCard({
    required String label,
    required int value,
    required int min,
    required int max,
    required Color color,
    required IconData icon,
    required String unit,
    required ValueChanged<double> onChanged,
  }) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Column(
          children: [
            Row(
              children: [
                Icon(icon, color: color, size: 20),
                const SizedBox(width: 10),
                Text(
                  label,
                  style: const TextStyle(
                    color: AppTheme.textPrimary,
                    fontSize: 15,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const Spacer(),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                  decoration: BoxDecoration(
                    color: color.withAlpha(25),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    unit.isNotEmpty ? '$value $unit' : '$value',
                    style: TextStyle(
                      color: color,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            SliderTheme(
              data: SliderThemeData(
                activeTrackColor: color,
                thumbColor: color,
                inactiveTrackColor: color.withAlpha(40),
                overlayColor: color.withAlpha(30),
              ),
              child: Slider(
                value: value.toDouble(),
                min: min.toDouble(),
                max: max.toDouble(),
                divisions: max - min,
                onChanged: onChanged,
              ),
            ),
          ],
        ),
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Timer view
  // ---------------------------------------------------------------------------
  Widget _buildTimerView() {
    final isFinished = _phase == HiitPhase.finished;

    return Column(
      children: [
        const Spacer(flex: 1),

        // Phase label
        Text(
          _phaseLabel,
          style: TextStyle(
            fontSize: 28,
            fontWeight: FontWeight.w900,
            color: _phaseColor,
            letterSpacing: 6,
          ),
        ),
        const SizedBox(height: 24),

        // Timer circle
        Container(
          width: 250,
          height: 250,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            border: Border.all(
              color: _phaseColor.withAlpha(60),
              width: 6,
            ),
          ),
          child: Container(
            margin: const EdgeInsets.all(10),
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              border: Border.all(
                color: _phaseColor,
                width: 4,
              ),
            ),
            child: Center(
              child: isFinished
                  ? Icon(
                      Icons.check_circle,
                      size: 80,
                      color: _phaseColor,
                    )
                  : Text(
                      _formatTime(_remainingSeconds),
                      style: TextStyle(
                        fontSize: 64,
                        fontWeight: FontWeight.w300,
                        color: _phaseColor,
                        fontFamily: 'monospace',
                      ),
                    ),
            ),
          ),
        ),
        const SizedBox(height: 32),

        // Round counter
        Text(
          isFinished
              ? 'All $_totalRounds rounds complete!'
              : 'Round $_currentRound of $_totalRounds',
          style: const TextStyle(
            fontSize: 18,
            color: AppTheme.textPrimary,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 8),

        // Round indicators
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32),
          child: Wrap(
            spacing: 6,
            runSpacing: 6,
            alignment: WrapAlignment.center,
            children: List.generate(_totalRounds, (index) {
              final roundNum = index + 1;
              final isCompleted = roundNum < _currentRound ||
                  (isFinished && roundNum <= _totalRounds);
              final isCurrent = roundNum == _currentRound && !isFinished;

              return Container(
                width: 24,
                height: 24,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: isCompleted
                      ? AppTheme.primary
                      : isCurrent
                          ? _phaseColor.withAlpha(40)
                          : AppTheme.surfaceLight,
                  border: isCurrent
                      ? Border.all(color: _phaseColor, width: 2)
                      : null,
                ),
                child: Center(
                  child: isCompleted
                      ? const Icon(Icons.check, color: Colors.white, size: 14)
                      : Text(
                          '$roundNum',
                          style: TextStyle(
                            color: isCurrent
                                ? _phaseColor
                                : AppTheme.textSecondary,
                            fontSize: 11,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                ),
              );
            }),
          ),
        ),

        const Spacer(flex: 2),

        // Control buttons
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24),
          child: isFinished
              ? SizedBox(
                  width: double.infinity,
                  height: 56,
                  child: ElevatedButton.icon(
                    onPressed: _reset,
                    icon: const Icon(Icons.replay, size: 24),
                    label: const Text(
                      'START OVER',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        letterSpacing: 2,
                      ),
                    ),
                  ),
                )
              : Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    // Reset button
                    SizedBox(
                      width: 56,
                      height: 56,
                      child: OutlinedButton(
                        onPressed: _reset,
                        style: OutlinedButton.styleFrom(
                          side: const BorderSide(
                              color: AppTheme.textSecondary, width: 1),
                          shape: const CircleBorder(),
                          padding: EdgeInsets.zero,
                        ),
                        child: const Icon(
                          Icons.stop,
                          color: AppTheme.textSecondary,
                          size: 24,
                        ),
                      ),
                    ),
                    const SizedBox(width: 24),
                    // Pause/Resume button
                    SizedBox(
                      width: 80,
                      height: 80,
                      child: ElevatedButton(
                        onPressed: _togglePause,
                        style: ElevatedButton.styleFrom(
                          shape: const CircleBorder(),
                          padding: EdgeInsets.zero,
                          backgroundColor: _phaseColor,
                        ),
                        child: Icon(
                          _isPaused ? Icons.play_arrow : Icons.pause,
                          color: Colors.white,
                          size: 36,
                        ),
                      ),
                    ),
                    const SizedBox(width: 24),
                    // Skip button
                    SizedBox(
                      width: 56,
                      height: 56,
                      child: OutlinedButton(
                        onPressed: () {
                          setState(() => _remainingSeconds = 0);
                          _onPhaseComplete();
                        },
                        style: OutlinedButton.styleFrom(
                          side: const BorderSide(
                              color: AppTheme.textSecondary, width: 1),
                          shape: const CircleBorder(),
                          padding: EdgeInsets.zero,
                        ),
                        child: const Icon(
                          Icons.skip_next,
                          color: AppTheme.textSecondary,
                          size: 24,
                        ),
                      ),
                    ),
                  ],
                ),
        ),
        const SizedBox(height: 32),
      ],
    );
  }
}
