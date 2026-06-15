import 'package:flutter/material.dart';
import '../models/habit.dart';
import '../services/api_service.dart';

class StatsScreen extends StatefulWidget {
  final Habit habit;

  const StatsScreen({super.key, required this.habit});

  @override
  State<StatsScreen> createState() => _StatsScreenState();
}

class _StatsScreenState extends State<StatsScreen> {
  HabitStats? _stats;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final data = await ApiService.getStats(widget.habit.id);
    if (data != null && mounted) {
      setState(() {
        _stats = HabitStats.fromJson(data);
        _loading = false;
      });
    } else if (mounted) {
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1A1A2E),
      appBar: AppBar(
        backgroundColor: const Color(0xFF16213E),
        iconTheme: const IconThemeData(color: Colors.white),
        title: Text(widget.habit.name,
            style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator(color: Color(0xFF6C63FF)))
          : _stats == null
              ? const Center(child: Text('İstatistik yüklenemedi', style: TextStyle(color: Colors.white54)))
              : Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('İstatistikler',
                          style: TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.bold)),
                      const SizedBox(height: 24),
                      Row(
                        children: [
                          _statCard('🔥 Mevcut Seri', '${_stats!.currentStreak} gün'),
                          const SizedBox(width: 12),
                          _statCard('🏆 En Uzun Seri', '${_stats!.longestStreak} gün'),
                        ],
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          _statCard('✅ Toplam', '${_stats!.totalChecks} kez'),
                          const SizedBox(width: 12),
                          _statCard('📊 Tamamlanma', '%${_stats!.completionRate.toStringAsFixed(1)}'),
                        ],
                      ),
                      const SizedBox(height: 32),
                      LinearProgressIndicator(
                        value: _stats!.completionRate / 100,
                        backgroundColor: Colors.white12,
                        color: const Color(0xFF6C63FF),
                        minHeight: 8,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Tamamlanma oranı: %${_stats!.completionRate.toStringAsFixed(1)}',
                        style: const TextStyle(color: Colors.white54),
                      ),
                    ],
                  ),
                ),
    );
  }

  Widget _statCard(String label, String value) => Expanded(
        child: Container(
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: const Color(0xFF16213E),
            borderRadius: BorderRadius.circular(16),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(label, style: const TextStyle(color: Colors.white54, fontSize: 13)),
              const SizedBox(height: 8),
              Text(value,
                  style: const TextStyle(
                      color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold)),
            ],
          ),
        ),
      );
}
