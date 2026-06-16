# Ercan Aziz'in Mobil Frontend Görevleri

**Mobile Front-end Demo Videosu:** [Link buraya eklenecek](https://example.com)

---

## Gereksinim 1 — Kullanıcı Kaydı Ekranı

**Ekran:** `register_screen.dart`

Kullanıcının sisteme yeni hesap oluşturmasını sağlayan kayıt ekranı. Kullanıcı adı, e-posta ve şifre alanlarından oluşan form yapısı bulunmaktadır. Form doğrulama (validation) mantığı uygulanmış olup eksik veya hatalı girişlerde kullanıcıya uyarı mesajı gösterilmektedir. Kayıt başarılı olduğunda kullanıcı otomatik olarak giriş ekranına yönlendirilmektedir.

**Kullanılan Widget'lar:** `TextField`, `ElevatedButton`, `SnackBar`

---

## Gereksinim 2 — Kullanıcı Girişi Ekranı

**Ekran:** `login_screen.dart`

Kayıtlı kullanıcıların sisteme erişimini sağlayan giriş ekranı. E-posta ve şifre alanları içermektedir. Başarılı girişte JWT token `SharedPreferences` üzerine kaydedilmekte ve kullanıcı alışkanlıklar ekranına yönlendirilmektedir. Uygulama açıldığında daha önce oturum açılmışsa otomatik olarak alışkanlıklar ekranına geçiş yapılmaktadır (Splash ekranı ile kontrol).

**Kullanılan Widget'lar:** `TextField`, `ElevatedButton`, `SharedPreferences`

---

## Gereksinim 3 — Yeni Alışkanlık Tanımlama

**Ekran:** `habits_screen.dart` → `_showCreateDialog()`

Ana ekrandaki `+` (FloatingActionButton) butonuna basıldığında açılan dialog penceresi aracılığıyla kullanıcı yeni alışkanlık oluşturabilmektedir. Ad ve açıklama alanlarından oluşmaktadır. Oluşturma başarılı olduğunda liste otomatik olarak yenilenmektedir.

**Kullanılan Widget'lar:** `FloatingActionButton`, `AlertDialog`, `TextField`

---

## Gereksinim 4 — Alışkanlıkları Listeleme Ekranı

**Ekran:** `habits_screen.dart`

Kullanıcının oluşturduğu tüm alışkanlıkları kart (Card) yapısında listeleyen ana ekrandır. Her kartta alışkanlık adı, açıklaması ve işlem butonları (tamamla, düzenle, sil, istatistik) yer almaktadır. Liste boş olduğunda yönlendirici bir boş durum (empty state) mesajı gösterilmektedir. Pull-to-refresh (aşağı çekme ile yenileme) desteği mevcuttur.

**Kullanılan Widget'lar:** `ListView.builder`, `Card`, `RefreshIndicator`

---

## Gereksinim 5 — Alışkanlık Durumu Güncelleme (Tamamlandı İşaretle)

**Ekran:** `habits_screen.dart` → `_toggleCheck()`

Her alışkanlık kartının solunda yer alan dairesel checkbox ile kullanıcı o günkü alışkanlığını tamamlandı olarak işaretleyebilmektedir. İşaretlenen alışkanlıklar mor renk ile dolup üzeri çizili (strikethrough) görünüm almaktadır. Durum `_checkedToday` kümesinde tutulmaktadır.

**Kullanılan Widget'lar:** `GestureDetector`, `Container` (dairesel), `Icon(Icons.check)`

---

## Gereksinim 6 — Alışkanlık Güncelleme

**Ekran:** `habits_screen.dart` → `_showEditDialog()`

Her kartın sağındaki kalem (edit) ikonuna basıldığında açılan dialog ile kullanıcı mevcut alışkanlığın adını ve açıklamasını güncelleyebilmektedir. Mevcut değerler forma otomatik olarak yüklenmektedir. Güncelleme sonrasında liste yenilenmektedir.

**Kullanılan Widget'lar:** `AlertDialog`, `TextEditingController`, `ElevatedButton`

---

## Gereksinim 7 — Alışkanlık Silme

**Ekran:** `habits_screen.dart` → `_deleteHabit()`

Her kartın sağındaki çöp kutusu ikonuna basıldığında onay dialogu açılmaktadır. Kullanıcı silme işlemini onaylarsa alışkanlık ve ona bağlı tüm tamamlama kayıtları sistemden kaldırılmakta, liste güncellenmektedir. Yanlışlıkla silmenin önüne geçmek için iki adımlı onay mekanizması kullanılmıştır.

**Kullanılan Widget'lar:** `AlertDialog`, `TextButton` (İptal/Sil)

---

## Gereksinim 8 — İşaretlemeyi Geri Alma

**Ekran:** `habits_screen.dart` → `_toggleCheck()`

Daha önce tamamlandı olarak işaretlenen bir alışkanlığa tekrar tıklanıldığında işaretleme geri alınmaktadır. `_checkedToday` kümesinden ilgili ID kaldırılmakta ve daire boş (çizgisiz) görünümüne geri dönmektedir.

**Kullanılan Widget'lar:** `GestureDetector`, durum yönetimi (`setState`)

---

## Gereksinim 9 — İstatistik ve Seri Takibi Ekranı

**Ekran:** `stats_screen.dart`

Her alışkanlık kartındaki grafik (bar_chart) ikonuna basıldığında açılan ekrandır. Seçilen alışkanlığa ait şu istatistikleri göstermektedir:

- **Mevcut Seri (Current Streak):** Kesintisiz devam eden gün sayısı
- **En Uzun Seri (Longest Streak):** Tüm zamanlardaki en iyi seri
- **Tamamlanma Oranı:** Oluşturulma tarihinden itibaren yüzde kaç tamamlandığı
- **Toplam Tamamlama:** Toplam tamamlandı işaretleme sayısı

**Kullanılan Widget'lar:** `Card`, `Column`, `CircularProgressIndicator` (yükleme)

---

## Gereksinim 10 — Oturumu Kapatma

**Ekran:** `habits_screen.dart` → `_logout()`

AppBar'ın sağ üst köşesindeki çıkış (logout) ikonuna basıldığında oturum sonlandırılmaktadır. `SharedPreferences`'tan JWT token silinmekte ve kullanıcı giriş ekranına yönlendirilmektedir. Sunucu tarafında da token geçersiz kılınmaktadır (Redis blacklist).

**Kullanılan Widget'lar:** `IconButton(Icons.logout)`, `Navigator.pushReplacement`

---

## Genel UI/UX Tercihleri

- **Tema:** Dark mode — `ThemeData.dark()` bazlı özel renk paleti
- **Ana Renk:** `#6C63FF` (mor) — buton, checkbox, ikon vurguları
- **Arka Plan:** `#1A1A2E` (koyu lacivert) 
- **Kart Arka Planı:** `#16213E`
- **Yazı Tipi:** Flutter varsayılan (Material Design)
- **Animasyon:** `RefreshIndicator`, `CircularProgressIndicator` yükleme göstergeleri
