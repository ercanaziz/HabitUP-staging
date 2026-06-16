# Mobil Backend (REST API Bağlantısı) Görev Dağılımı

**REST API Adresi:** [habitup-production.up.railway.app](https://habitup-production.up.railway.app)

Bu dokümanda, mobil uygulamanın REST API ile iletişimini sağlayan backend entegrasyon görevleri listelenmektedir. Her grup üyesi, kendisine atanan API endpoint'lerinin mobil uygulamadan çağrılması ve yönetilmesinden sorumludur.

---

## Grup Üyelerinin Mobil Backend Görevleri

1. [Ercan Aziz'in Mobil Backend Görevleri](Ercan-Aziz/Ercan-Aziz-Mobil-Backend-Gorevleri.md)

---

## Genel Mobil Backend Prensipleri

- **İletişim Katmanı:** `lib/services/api_service.dart` — tüm HTTP istekleri bu sınıf üzerinden yapılır
- **Kimlik Doğrulama:** JWT — `Authorization: Bearer <token>` header'ı otomatik eklenir
- **Token Saklama:** `shared_preferences` paketi
- **Veri Formatı:** JSON (`dart:convert`)
- **Güvenli İletişim:** HTTPS

---

## Gereksinim 1 — Kullanıcı Kaydı
**Endpoint:** `POST /api/auth/register`
Kullanıcı adı, e-posta ve şifre API'ye gönderilir. Şifre sunucuda bcrypt ile hashlenerek MongoDB'ye kaydedilir.

## Gereksinim 2 — Kullanıcı Girişi
**Endpoint:** `POST /api/auth/login`
E-posta ve şifre doğrulanır. Başarılı yanıtta dönen JWT token SharedPreferences'a kaydedilir.

## Gereksinim 3 — Yeni Alışkanlık Tanımlama
**Endpoint:** `POST /api/habits`
Ad ve açıklama JWT ile birlikte gönderilir. `201 Created` dönmesi durumunda liste yenilenir.

## Gereksinim 4 — Alışkanlıkları Listeleme
**Endpoint:** `GET /api/habits`
JWT ile korunan endpoint, kullanıcıya ait tüm alışkanlıkları JSON dizisi olarak döndürür.

## Gereksinim 5 — Alışkanlık Durumu Güncelleme
**Endpoint:** `POST /api/habits/{id}/check`
Seçilen alışkanlık o gün için tamamlandı olarak işaretlenir. Sunucu aynı güne tekrar izin vermez.

## Gereksinim 6 — Alışkanlık Güncelleme
**Endpoint:** `PUT /api/habits/{id}`
Güncellenecek ad/açıklama gönderilir. Yalnızca gönderilen alanlar değiştirilir (partial update).

## Gereksinim 7 — Alışkanlık Silme
**Endpoint:** `DELETE /api/habits/{id}`
Alışkanlık ve bağlı tüm tamamlama kayıtları sunucudan kalıcı olarak silinir. `204 No Content` döner.

## Gereksinim 8 — İşaretlemeyi Geri Alma
**Endpoint:** `DELETE /api/habits/{id}/check`
Bugüne ait tamamlama kaydı silinir. `204 No Content` başarı yanıtıdır.

## Gereksinim 9 — İstatistik ve Seri Takibi
**Endpoint:** `GET /api/habits/{id}/stats`
`currentStreak`, `longestStreak`, `completionRate`, `totalChecks` alanlarını içeren JSON döner.

## Gereksinim 10 — Oturumu Kapatma
**Endpoint:** `POST /api/auth/logout`
Token sunucuda Redis blacklist'e eklenerek geçersiz kılınır. Mobilde SharedPreferences'tan silinir.
