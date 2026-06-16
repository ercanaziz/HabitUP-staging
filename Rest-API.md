# REST API Görev Dağılımı

**REST API Adresi:** [habitup-production.up.railway.app](https://habitup-production.up.railway.app)

Bu dokümanda, proje ekibindeki her üyenin geliştirmekten sorumlu olduğu REST API metotları listelenmektedir.

---

## Grup Üyelerinin REST API Metotları

1. [Ercan Aziz'in REST API Metotları](Ercan-Aziz/Ercan-Aziz-Rest-API-Gorevleri.md)

---

## Genel REST API Prensipleri

- **Framework:** Go 1.22 — Gin Web Framework
- **Veritabanı:** MongoDB
- **Kimlik Doğrulama:** JWT (HS256) — `Authorization: Bearer <token>`
- **Cache:** Redis (Upstash) — JWT blacklist
- **Mesaj Kuyruğu:** RabbitMQ (CloudAMQP)
- **Deployment:** Railway + Docker

---

## Gereksinim 1 — Kullanıcı Kaydı
**Endpoint:** `POST /api/auth/register`
Kullanıcı adı, e-posta ve şifreyi alır. Şifre bcrypt ile hashlenerek MongoDB'ye kaydedilir.

## Gereksinim 2 — Kullanıcı Girişi
**Endpoint:** `POST /api/auth/login`
E-posta ve şifre doğrulanır. Başarılı girişte 24 saatlik JWT token döndürülür.

## Gereksinim 3 — Yeni Alışkanlık Tanımlama
**Endpoint:** `POST /api/habits`
Ad ve açıklama bilgisiyle kullanıcıya bağlı yeni alışkanlık oluşturulur. `201 Created` döner.

## Gereksinim 4 — Alışkanlıkları Listeleme
**Endpoint:** `GET /api/habits`
Giriş yapmış kullanıcının yalnızca kendi alışkanlıklarını getirir (veri izolasyonu).

## Gereksinim 5 — Alışkanlık Durumu Güncelleme
**Endpoint:** `POST /api/habits/{id}/check`
Alışkanlığı o gün için tamamlandı işaretler. Aynı güne tekrar izin verilmez. RabbitMQ'ya event yayınlanır.

## Gereksinim 6 — Alışkanlık Güncelleme
**Endpoint:** `PUT /api/habits/{id}`
MongoDB `$set` ile yalnızca gönderilen alanlar güncellenir. Başkasının verisine erişilemez.

## Gereksinim 7 — Alışkanlık Silme
**Endpoint:** `DELETE /api/habits/{id}`
Alışkanlık ve bağlı tüm tamamlama kayıtları MongoDB'den kalıcı olarak silinir. `204 No Content` döner.

## Gereksinim 8 — İşaretlemeyi Geri Alma
**Endpoint:** `DELETE /api/habits/{id}/check`
Bugüne ait tamamlama kaydı `checks` koleksiyonundan silinir. `204 No Content` döner.

## Gereksinim 9 — İstatistik ve Seri Takibi
**Endpoint:** `GET /api/habits/{id}/stats`
Mevcut seri, en uzun seri, tamamlanma oranı ve toplam tamamlama sayısı hesaplanarak döndürülür.

## Gereksinim 10 — Oturumu Kapatma
**Endpoint:** `POST /api/auth/logout`
Token Redis blacklist'e eklenerek geçersiz kılınır. Sonraki isteklerde bu token reddedilir.
