import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

class ApiService {
  // Production URL
  static const String _base = 'https://habitup-staging-production.up.railway.app/api';

  static Future<String?> _token() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('token');
  }

  static Future<Map<String, String>> _headers({bool auth = true}) async {
    final headers = {'Content-Type': 'application/json'};
    if (auth) {
      final t = await _token();
      if (t != null) headers['Authorization'] = 'Bearer $t';
    }
    return headers;
  }

  // Gereksinim 1: Kullanıcı Kaydı
  static Future<Map<String, dynamic>> register(
      String username, String email, String password) async {
    final res = await http.post(
      Uri.parse('$_base/auth/register'),
      headers: await _headers(auth: false),
      body: jsonEncode(
          {'username': username, 'email': email, 'password': password}),
    );
    return {'statusCode': res.statusCode, 'body': jsonDecode(res.body)};
  }

  // Gereksinim 2: Kullanıcı Girişi
  static Future<String?> login(String email, String password) async {
    final res = await http.post(
      Uri.parse('$_base/auth/login'),
      headers: await _headers(auth: false),
      body: jsonEncode({'email': email, 'password': password}),
    );
    if (res.statusCode == 200) {
      final data = jsonDecode(res.body);
      final token = data['token'] as String;
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('token', token);
      return token;
    }
    return null;
  }

  // Gereksinim 10: Oturumu Kapatma
  static Future<void> logout() async {
    await http.post(
      Uri.parse('$_base/auth/logout'),
      headers: await _headers(),
    );
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('token');
  }

  // Gereksinim 4: Alışkanlıkları Listeleme
  static Future<List<dynamic>> getHabits() async {
    final res = await http.get(
      Uri.parse('$_base/habits'),
      headers: await _headers(),
    );
    if (res.statusCode == 200) return jsonDecode(res.body) as List;
    return [];
  }

  // Gereksinim 3: Yeni Alışkanlık Tanımlama
  static Future<bool> createHabit(String name, String description) async {
    final res = await http.post(
      Uri.parse('$_base/habits'),
      headers: await _headers(),
      body: jsonEncode({'name': name, 'description': description}),
    );
    return res.statusCode == 201;
  }

  // Gereksinim 6: Alışkanlık Güncelleme
  static Future<bool> updateHabit(
      String id, String name, String description) async {
    final res = await http.put(
      Uri.parse('$_base/habits/$id'),
      headers: await _headers(),
      body: jsonEncode({'name': name, 'description': description}),
    );
    return res.statusCode == 200;
  }

  // Gereksinim 7: Alışkanlık Silme
  static Future<bool> deleteHabit(String id) async {
    final res = await http.delete(
      Uri.parse('$_base/habits/$id'),
      headers: await _headers(),
    );
    return res.statusCode == 204;
  }

  // Gereksinim 5: Alışkanlık Durumu Güncelleme
  static Future<bool> checkHabit(String id) async {
    final res = await http.post(
      Uri.parse('$_base/habits/$id/check'),
      headers: await _headers(),
    );
    return res.statusCode == 200;
  }

  // Gereksinim 8: İşaretlemeyi Geri Alma
  static Future<bool> uncheckHabit(String id) async {
    final res = await http.delete(
      Uri.parse('$_base/habits/$id/check'),
      headers: await _headers(),
    );
    return res.statusCode == 204;
  }

  // Gereksinim 9: İstatistik ve Seri Takibi
  static Future<Map<String, dynamic>?> getStats(String id) async {
    final res = await http.get(
      Uri.parse('$_base/habits/$id/stats'),
      headers: await _headers(),
    );
    if (res.statusCode == 200) return jsonDecode(res.body);
    return null;
  }
}
