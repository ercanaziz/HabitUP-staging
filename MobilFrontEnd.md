# Mobil Frontend Görev Dağılımı

Bu dokümanda, mobil uygulamanın kullanıcı arayüzü (UI) ve kullanıcı deneyimi (UX) görevleri listelenmektedir. Her grup üyesi, kendisine atanan ekranların tasarımı, implementasyonu ve kullanıcı etkileşimlerinden sorumludur.

---

## Grup Üyelerinin Mobil Frontend Görevleri

1. [Ercan Aziz'in Mobil Frontend Görevleri](Ercan-Aziz/Ercan-Aziz-Mobil-Frontend-Gorevleri.md)

---

## Genel Mobil Frontend Prensipleri

- **Framework:** Flutter (Dart)
- **Tema:** Dark mode — `#1A1A2E` arka plan, `#6C63FF` mor vurgu rengi
- **Oturum Yönetimi:** `SharedPreferences` ile JWT token kalıcı saklama
- **Navigasyon:** `Navigator.push` / `Navigator.pushReplacement`

---

## Gereksinim 1 — Kullanıcı Kaydı Ekranı
**Ekran:** `register_screen.dart`
Kullanıcı adı, e-posta ve şifre formu. Başarılı kayıt sonrası giriş ekranına yönlendirme.

## Gereksinim 2 — Kullanıcı Girişi Ekranı
**Ekran:** `login_screen.dart`
E-posta/şifre formu. JWT token SharedPreferences'a kaydedilerek alışkanlıklar ekranına geçiş yapılır.

## Gereksinim 3 — Yeni Alışkanlık Tanımlama
**Ekran:** `habits_screen.dart` → `_showCreateDialog()`
`+` butonuyla açılan dialog üzerinden ad ve açıklama girilerek yeni alışkanlık oluşturulur.

## Gereksinim 4 — Alışkanlıkları Listeleme
**Ekran:** `habits_screen.dart`
Kullanıcının tüm alışkanlıkları Card bileşenleriyle listelenir. Pull-to-refresh desteği mevcuttur.

## Gereksinim 5 — Alışkanlık Durumu Güncelleme
**Ekran:** `habits_screen.dart` → `_toggleCheck()`
Kart solundaki dairesel checkbox ile o günkü alışkanlık tamamlandı işaretlenir, kart mor renge döner.

## Gereksinim 6 — Alışkanlık Güncelleme
**Ekran:** `habits_screen.dart` → `_showEditDialog()`
Kalem ikonuyla açılan dialog üzerinden ad ve açıklama güncellenir.

## Gereksinim 7 — Alışkanlık Silme
**Ekran:** `habits_screen.dart` → `_deleteHabit()`
Çöp kutusu ikonu → onay dialogu → silme işlemi. İki adımlı onay mekanizması.

## Gereksinim 8 — İşaretlemeyi Geri Alma
**Ekran:** `habits_screen.dart` → `_toggleCheck()`
İşaretli alışkanlığa tekrar tıklanarak günün tamamlama kaydı geri alınır.

## Gereksinim 9 — İstatistik ve Seri Takibi
**Ekran:** `stats_screen.dart`
Grafik ikonuyla açılır. Mevcut seri, en uzun seri, tamamlanma oranı ve toplam tamamlama gösterilir.

## Gereksinim 10 — Oturumu Kapatma
**Ekran:** `habits_screen.dart` → `_logout()`
AppBar'daki çıkış ikonu ile token silinir, kullanıcı giriş ekranına yönlendirilir.
