# Mobil Backend (REST API Bağlantısı) Görev Dağılımı

**REST API Adresi:** `https://habitup-staging-production.up.railway.app`

Bu dokümanda, mobil uygulamanın REST API ile iletişimini sağlayan backend entegrasyon görevleri listelenmektedir. Her grup üyesi, kendisine atanan API endpoint'lerinin mobil uygulamadan çağrılması ve yönetilmesinden sorumludur.

---

## Grup Üyelerinin Mobil Backend Görevleri

1. [Ercan Aziz'in Mobil Backend Görevleri](Ercan-Aziz/Ercan-Aziz-Mobil-Backend-Gorevleri.md)

---

## Genel Mobil Backend Prensipleri

- **API İletişim Katmanı:** `lib/services/api_service.dart` — tüm HTTP istekleri bu sınıf üzerinden yapılır
- **HTTP Kütüphanesi:** `http: ^1.2.1` (Dart resmi paketi)
- **Kimlik Doğrulama:** JWT token — her korumalı istekte `Authorization: Bearer <token>` header'ı otomatik eklenir
- **Token Saklama:** `shared_preferences: ^2.2.3` — token cihazda kalıcı olarak saklanır
- **Veri Formatı:** JSON — tüm istek ve yanıtlar `dart:convert` ile encode/decode edilir
- **Güvenli İletişim:** HTTPS (Railway SSL sertifikası)

### Endpoint Tablosu

| Gereksinim | Metot | Endpoint | Yetkilendirme |
|------------|-------|----------|---------------|
| 1 — Kullanıcı Kaydı | `POST` | `/api/auth/register` | Yok |
| 2 — Kullanıcı Girişi | `POST` | `/api/auth/login` | Yok |
| 3 — Alışkanlık Oluştur | `POST` | `/api/habits` | JWT |
| 4 — Alışkanlıkları Listele | `GET` | `/api/habits` | JWT |
| 5 — Tamamlandı İşaretle | `POST` | `/api/habits/{id}/check` | JWT |
| 6 — Alışkanlık Güncelle | `PUT` | `/api/habits/{id}` | JWT |
| 7 — Alışkanlık Sil | `DELETE` | `/api/habits/{id}` | JWT |
| 8 — İşareti Geri Al | `DELETE` | `/api/habits/{id}/check` | JWT |
| 9 — İstatistikler | `GET` | `/api/habits/{id}/stats` | JWT |
| 10 — Oturum Kapat | `POST` | `/api/auth/logout` | JWT |

### Backend Altyapısı

| Bileşen | Teknoloji | Açıklama |
|---------|-----------|----------|
| **API Sunucu** | Go 1.22 — Gin | REST endpoint'leri karşılar |
| **Veritabanı** | MongoDB (Railway) | Kullanıcı ve alışkanlık verileri |
| **Cache** | Redis — Upstash | JWT blacklist (logout güvenliği) |
| **Mesaj Kuyruğu** | RabbitMQ — CloudAMQP | Alışkanlık tamamlama event'leri |
| **Deployment** | Railway + Docker | Canlı production ortamı |
| **CI/CD** | GitHub Actions | Otomatik build ve test |
