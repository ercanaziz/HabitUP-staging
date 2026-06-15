import 'package:flutter/material.dart';
import '../models/habit.dart';
import '../services/api_service.dart';
import 'stats_screen.dart';
import 'login_screen.dart';

class HabitsScreen extends StatefulWidget {
  const HabitsScreen({super.key});

  @override
  State<HabitsScreen> createState() => _HabitsScreenState();
}

class _HabitsScreenState extends State<HabitsScreen> {
  List<Habit> _habits = [];
  final Set<String> _checkedToday = {};
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    final data = await ApiService.getHabits();
    setState(() {
      _habits = data.map((e) => Habit.fromJson(e as Map<String, dynamic>)).toList();
      _loading = false;
    });
  }

  Future<void> _toggleCheck(Habit h) async {
    if (_checkedToday.contains(h.id)) {
      final ok = await ApiService.uncheckHabit(h.id);
      if (ok) setState(() => _checkedToday.remove(h.id));
    } else {
      final ok = await ApiService.checkHabit(h.id);
      if (ok) setState(() => _checkedToday.add(h.id));
    }
  }

  Future<void> _deleteHabit(Habit h) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Alışkanlığı Sil'),
        content: Text('"${h.name}" silinsin mi?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('İptal')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Sil', style: TextStyle(color: Colors.red))),
        ],
      ),
    );
    if (confirm == true) {
      await ApiService.deleteHabit(h.id);
      _load();
    }
  }

  void _showEditDialog(Habit h) {
    final nameCtrl = TextEditingController(text: h.name);
    final descCtrl = TextEditingController(text: h.description);
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Alışkanlığı Düzenle'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: 'Ad')),
            TextField(controller: descCtrl, decoration: const InputDecoration(labelText: 'Açıklama')),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('İptal')),
          ElevatedButton(
            onPressed: () async {
              await ApiService.updateHabit(h.id, nameCtrl.text, descCtrl.text);
              if (ctx.mounted) Navigator.pop(ctx);
              _load();
            },
            child: const Text('Kaydet'),
          ),
        ],
      ),
    );
  }

  void _showCreateDialog() {
    final nameCtrl = TextEditingController();
    final descCtrl = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Yeni Alışkanlık'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: 'Ad')),
            TextField(controller: descCtrl, decoration: const InputDecoration(labelText: 'Açıklama')),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('İptal')),
          ElevatedButton(
            onPressed: () async {
              await ApiService.createHabit(nameCtrl.text, descCtrl.text);
              if (ctx.mounted) Navigator.pop(ctx);
              _load();
            },
            child: const Text('Oluştur'),
          ),
        ],
      ),
    );
  }

  Future<void> _logout() async {
    await ApiService.logout();
    if (!mounted) return;
    Navigator.pushReplacement(
      context,
      MaterialPageRoute(builder: (_) => const LoginScreen()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1A1A2E),
      appBar: AppBar(
        backgroundColor: const Color(0xFF16213E),
        title: const Text('HabitUp', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout, color: Colors.white70),
            onPressed: _logout,
            tooltip: 'Çıkış Yap',
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showCreateDialog,
        backgroundColor: const Color(0xFF6C63FF),
        child: const Icon(Icons.add, color: Colors.white),
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator(color: Color(0xFF6C63FF)))
          : _habits.isEmpty
              ? const Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.track_changes, size: 64, color: Colors.white24),
                      const SizedBox(height: 16),
                      const Text('Henüz alışkanlık yok', style: TextStyle(color: Colors.white54, fontSize: 16)),
                      const SizedBox(height: 8),
                      const Text('+ butonuna basarak başla', style: TextStyle(color: Colors.white38)),
                    ],
                  ),
                )
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: _habits.length,
                    itemBuilder: (_, i) {
                      final h = _habits[i];
                      final checked = _checkedToday.contains(h.id);
                      return Card(
                        color: const Color(0xFF16213E),
                        margin: const EdgeInsets.only(bottom: 12),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                        child: ListTile(
                          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                          leading: GestureDetector(
                            onTap: () => _toggleCheck(h),
                            child: Container(
                              width: 40,
                              height: 40,
                              decoration: BoxDecoration(
                                shape: BoxShape.circle,
                                color: checked ? const Color(0xFF6C63FF) : Colors.transparent,
                                border: Border.all(
                                  color: checked ? const Color(0xFF6C63FF) : Colors.white38,
                                  width: 2,
                                ),
                              ),
                              child: checked
                                  ? const Icon(Icons.check, color: Colors.white, size: 20)
                                  : null,
                            ),
                          ),
                          title: Text(
                            h.name,
                            style: TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.w600,
                              decoration: checked ? TextDecoration.lineThrough : null,
                              decorationColor: Colors.white54,
                            ),
                          ),
                          subtitle: h.description.isNotEmpty
                              ? Text(h.description, style: const TextStyle(color: Colors.white54))
                              : null,
                          trailing: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              IconButton(
                                icon: const Icon(Icons.bar_chart, color: Color(0xFF6C63FF)),
                                onPressed: () => Navigator.push(
                                  context,
                                  MaterialPageRoute(builder: (_) => StatsScreen(habit: h)),
                                ),
                              ),
                              IconButton(
                                icon: const Icon(Icons.edit_outlined, color: Colors.white54),
                                onPressed: () => _showEditDialog(h),
                              ),
                              IconButton(
                                icon: const Icon(Icons.delete_outline, color: Colors.redAccent),
                                onPressed: () => _deleteHabit(h),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}
