# Mobil Frontend Görev Dağılımı

Bu dokümanda, mobil uygulamanın kullanıcı arayüzü (UI) ve kullanıcı deneyimi (UX) görevleri listelenmektedir. Her grup üyesi, kendisine atanan ekranların tasarımı, implementasyonu ve kullanıcı etkileşimlerinden sorumludur.

---

## Grup Üyelerinin Mobil Frontend Görevleri

1. [Ercan Aziz'in Mobil Frontend Görevleri](Ercan-Aziz/Ercan-Aziz-Mobil-Frontend-Gorevleri.md)

---

## Genel Mobil Frontend Prensipleri

- **Framework:** Flutter (Dart) — cross-platform mobil geliştirme
- **Mimari:** Tek servis katmanı (`ApiService`) üzerinden merkezi API iletişimi
- **Tema:** Dark mode — `#1A1A2E` arka plan, `#6C63FF` mor vurgu rengi
- **Durum Yönetimi:** `StatefulWidget` + `setState()` ile lokal durum yönetimi
- **Oturum Yönetimi:** `SharedPreferences` ile JWT token kalıcı saklama
- **Navigasyon:** `Navigator.push` / `Navigator.pushReplacement` ile ekran geçişleri

### Ekran Listesi

| Ekran | Dosya | Sorumluluk |
|-------|-------|------------|
| Splash (Açılış) | `main.dart` | Token kontrolü, otomatik yönlendirme |
| Kayıt | `register_screen.dart` | Yeni kullanıcı kaydı formu |
| Giriş | `login_screen.dart` | E-posta/şifre giriş formu |
| Alışkanlıklar | `habits_screen.dart` | Ana liste, oluştur, düzenle, sil, tamamla |
| İstatistikler | `stats_screen.dart` | Streak ve tamamlanma oranı görüntüleme |

### Uygulanan Gereksinimler

| Gereksinim | Ekran / Fonksiyon |
|------------|-------------------|
| 1 — Kullanıcı Kaydı | `register_screen.dart` |
| 2 — Kullanıcı Girişi | `login_screen.dart` |
| 3 — Yeni Alışkanlık | `habits_screen.dart` → `_showCreateDialog()` |
| 4 — Listeleme | `habits_screen.dart` → `ListView.builder` |
| 5 — Tamamlandı İşaretle | `habits_screen.dart` → `_toggleCheck()` |
| 6 — Alışkanlık Güncelle | `habits_screen.dart` → `_showEditDialog()` |
| 7 — Alışkanlık Sil | `habits_screen.dart` → `_deleteHabit()` |
| 8 — İşareti Geri Al | `habits_screen.dart` → `_toggleCheck()` |
| 9 — İstatistikler | `stats_screen.dart` |
| 10 — Oturum Kapat | `habits_screen.dart` → `_logout()` |
