class Habit {
  final String id;
  final String name;
  final String description;
  final DateTime createdAt;

  Habit({
    required this.id,
    required this.name,
    required this.description,
    required this.createdAt,
  });

  factory Habit.fromJson(Map<String, dynamic> json) => Habit(
        id: json['id'] as String,
        name: json['name'] as String,
        description: (json['description'] as String?) ?? '',
        createdAt: DateTime.parse(json['createdAt'] as String),
      );
}

class HabitStats {
  final int currentStreak;
  final int longestStreak;
  final double completionRate;
  final int totalChecks;

  HabitStats({
    required this.currentStreak,
    required this.longestStreak,
    required this.completionRate,
    required this.totalChecks,
  });

  factory HabitStats.fromJson(Map<String, dynamic> json) => HabitStats(
        currentStreak: json['currentStreak'] as int,
        longestStreak: json['longestStreak'] as int,
        completionRate: (json['completionRate'] as num).toDouble(),
        totalChecks: json['totalChecks'] as int,
      );
}
