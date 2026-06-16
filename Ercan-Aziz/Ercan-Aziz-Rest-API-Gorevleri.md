# Ercan Aziz'in REST API Metotları

**API Test Videosu:** [Link buraya eklenecek](https://example.com)

**Canlı API Adresi:** `https://habitup-staging-production.up.railway.app`

**Teknoloji:** Go 1.22 — Gin Web Framework, MongoDB, Redis (JWT blacklist), RabbitMQ (event queue)

---

## Gereksinim 1 — Kullanıcı Kaydı

**Endpoint:** `POST /api/auth/register`
**Yetkilendirme:** Gerekmez (public)
**Kaynak:** `internal/auth/handler.go` → `Register()`

### İstek (Request)
```json
{
  "username": "ercan",
  "email": "ercan@example.com",
  "password": "sifre123"
}
```

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `201 Created` | Kullanıcı başarıyla oluşturuldu |
| `400 Bad Request` | Eksik alan veya bu e-posta zaten kayıtlı |
| `500 Internal Server Error` | Sunucu hatası |

### Uygulama Detayı
Şifre sunucuda `bcrypt` (DefaultCost) ile hashlenerek MongoDB `users` koleksiyonuna kaydedilmektedir. E-posta benzersizliği kontrol edilmekte, mükerrer kayıt engellenmektedir.

---

## Gereksinim 2 — Kullanıcı Girişi

**Endpoint:** `POST /api/auth/login`
**Yetkilendirme:** Gerekmez (public)
**Kaynak:** `internal/auth/handler.go` → `Login()`

### İstek (Request)
```json
{
  "email": "ercan@example.com",
  "password": "sifre123"
}
```

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `200 OK` | `{ "token": "eyJhbGci..." }` |
| `401 Unauthorized` | Hatalı e-posta veya şifre |
| `500 Internal Server Error` | Sunucu hatası |

### Uygulama Detayı
MongoDB'den kullanıcı bulunarak `bcrypt.CompareHashAndPassword()` ile şifre doğrulanmaktadır. Başarılı girişte `HS256` algoritmasıyla imzalanmış, 24 saat geçerli JWT token üretilip döndürülmektedir. Token içinde `userId` claim'i bulunmaktadır.

---

## Gereksinim 3 — Yeni Alışkanlık Tanımlama

**Endpoint:** `POST /api/habits`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/habit/handler.go` → `Create()`

### İstek (Request)
```json
{
  "name": "Kitap Oku",
  "description": "Her gün en az 20 sayfa"
}
```

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `201 Created` | Oluşturulan alışkanlık objesi döner |
| `400 Bad Request` | `name` alanı zorunludur |
| `401 Unauthorized` | Geçersiz veya eksik token |
| `500 Internal Server Error` | MongoDB yazma hatası |

### Uygulama Detayı
JWT middleware'den gelen `userId` ile alışkanlık ilişkilendirilmektedir. `primitive.NewObjectID()` ile benzersiz MongoDB ID üretilmekte, `createdAt` alanı sunucu saatiyle atanmaktadır.

---

## Gereksinim 4 — Alışkanlıkları Listeleme

**Endpoint:** `GET /api/habits`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/habit/handler.go` → `List()`

### İstek (Request)
Gövde yok. Token header'dan okunur.

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `200 OK` | Alışkanlık dizisi (boşsa `[]`) |
| `401 Unauthorized` | Geçersiz token |
| `500 Internal Server Error` | MongoDB okuma hatası |

### Örnek Başarılı Yanıt
```json
[
  {
    "id": "665f1a2b3c4d5e6f7a8b9c0d",
    "userId": "665f1a2b3c4d5e6f7a8b9c0e",
    "name": "Kitap Oku",
    "description": "Her gün en az 20 sayfa",
    "createdAt": "2026-06-15T10:00:00Z"
  }
]
```

### Uygulama Detayı
MongoDB `habits` koleksiyonu `userId` alanına göre filtrelenmektedir. Kullanıcı yalnızca kendi alışkanlıklarını görmektedir (veri izolasyonu). Sonuç boş olduğunda `null` yerine `[]` döndürülmektedir.

---

## Gereksinim 5 — Alışkanlık Durumu Güncelleme (Tamamlandı İşaretle)

**Endpoint:** `POST /api/habits/{id}/check`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/habit/handler.go` → `Check()`

### İstek (Request)
Gövde yok. `{id}` URL parametresi.

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `200 OK` | Tamamlandı olarak işaretlendi |
| `200 OK` | Bugün zaten işaretlenmişti (idempotent) |
| `400 Bad Request` | Geçersiz MongoDB ID |
| `401 Unauthorized` | Geçersiz token |
| `500 Internal Server Error` | MongoDB yazma hatası |

### Uygulama Detayı
`checks` koleksiyonuna `habitId`, `userId`, `date` (YYYY-MM-DD), `checkedAt` alanlarıyla kayıt eklenmektedir. Aynı gün için tekrar işaret atılması engellenmektedir. İşlem başarılı olduğunda RabbitMQ `habit.checked` kuyruğuna `{ userId, habitId }` event'i yayınlanmaktadır.

---

## Gereksinim 6 — Alışkanlık Güncelleme

**Endpoint:** `PUT /api/habits/{id}`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/habit/handler.go` → `Update()`

### İstek (Request)
```json
{
  "name": "Kitap Oku (Güncellenmiş)",
  "description": "Her gün en az 30 sayfa"
}
```

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `200 OK` | Alışkanlık güncellendi |
| `400 Bad Request` | Geçersiz ID veya gövde |
| `401 Unauthorized` | Geçersiz token |
| `404 Not Found` | Alışkanlık bulunamadı veya başkasına ait |

### Uygulama Detayı
MongoDB `$set` operatörü kullanılmaktadır; yalnızca gönderilen alanlar güncellenir. Güncelleme `{ _id: habitID, userId: uid }` filtresiyle yapılmaktadır; böylece bir kullanıcı başkasının alışkanlığını değiştiremez.

---

## Gereksinim 7 — Alışkanlık Silme

**Endpoint:** `DELETE /api/habits/{id}`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/habit/handler.go` → `Delete()`

### İstek (Request)
Gövde yok. `{id}` URL parametresi.

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `204 No Content` | Başarıyla silindi |
| `400 Bad Request` | Geçersiz MongoDB ID |
| `401 Unauthorized` | Geçersiz token |
| `404 Not Found` | Alışkanlık bulunamadı |

### Uygulama Detayı
Önce `habits` koleksiyonundan alışkanlık silinmekte, ardından `checks` koleksiyonundan bu alışkanlığa ait tüm tamamlama kayıtları `DeleteMany()` ile toplu olarak temizlenmektedir.

---

## Gereksinim 8 — İşaretlemeyi Geri Alma

**Endpoint:** `DELETE /api/habits/{id}/check`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/habit/handler.go` → `Uncheck()`

### İstek (Request)
Gövde yok. `{id}` URL parametresi.

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `204 No Content` | İşaretleme başarıyla geri alındı |
| `400 Bad Request` | Geçersiz ID |
| `401 Unauthorized` | Geçersiz token |
| `404 Not Found` | Bugün için işaretleme kaydı bulunamadı |

### Uygulama Detayı
`checks` koleksiyonunda bugünün tarihi (`YYYY-MM-DD`), `habitId` ve `userId` ile eşleşen tek kayıt silinmektedir. Yalnızca günün işareti geri alınabilir; geçmiş günlere dokunulmaz.

---

## Gereksinim 9 — İstatistik ve Seri Takibi

**Endpoint:** `GET /api/habits/{id}/stats`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/habit/handler.go` → `Stats()`

### İstek (Request)
Gövde yok. `{id}` URL parametresi.

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `200 OK` | İstatistik objesi döner |
| `401 Unauthorized` | Geçersiz token |
| `404 Not Found` | Alışkanlık bulunamadı |

### Örnek Başarılı Yanıt
```json
{
  "currentStreak": 5,
  "longestStreak": 12,
  "completionRate": 78.5,
  "totalChecks": 47
}
```

### Uygulama Detayı
`checks` koleksiyonundan tüm tamamlama kayıtları çekilerek tarihler sıralanmaktadır. Streak algoritması (`calculateStreaks`) sıralı tarihleri geriden ileriye tarayarak ardışık günleri saymaktadır. Tamamlanma oranı, oluşturulma tarihinden bugüne geçen toplam gün sayısına bölünerek hesaplanmaktadır.

---

## Gereksinim 10 — Oturumu Kapatma

**Endpoint:** `POST /api/auth/logout`
**Yetkilendirme:** Gerekir — `Authorization: Bearer <token>`
**Kaynak:** `internal/auth/handler.go` → `Logout()`

### İstek (Request)
Gövde yok. Token header'dan okunur.

### Yanıtlar (Response)
| HTTP Kodu | Açıklama |
|-----------|----------|
| `200 OK` | Başarıyla çıkış yapıldı |

### Uygulama Detayı
Header'dan alınan JWT token parse edilerek kalan geçerlilik süresi (`ttl`) hesaplanmaktadır. Token, `blacklist:<token>` anahtarıyla TTL süresi kadar **Redis (Upstash)**'e yazılmaktadır. Middleware her istekte bu anahtarın varlığını kontrol etmekte; blacklist'teki token geçersiz sayılmaktadır. Bu sayede logout sonrası eski token'larla erişim engellenmektedir.

---

## Genel API Mimarisi

| Katman | Teknoloji | Açıklama |
|--------|-----------|----------|
| **Web Framework** | Go — Gin | HTTP routing, middleware zinciri |
| **Veritabanı** | MongoDB | Kullanıcı, alışkanlık ve check kayıtları |
| **Cache / Oturum** | Redis (Upstash) | JWT token blacklist (logout güvenliği) |
| **Mesaj Kuyruğu** | RabbitMQ (CloudAMQP) | Alışkanlık tamamlama event'leri |
| **Kimlik Doğrulama** | JWT (HS256) | 24 saatlik token, Bearer scheme |
| **Şifre Güvenliği** | bcrypt | Şifreler hash'lenerek saklanır |
| **Deployment** | Railway + Docker | Otomatik CI/CD, canlı production |

### Middleware Zinciri
```
İstek → CORS → AuthRequired (JWT doğrula + Redis blacklist kontrolü) → Handler
```

### Veri İzolasyonu
Tüm korumalı endpoint'lerde JWT'den çıkarılan `userId`, MongoDB sorgularına eklenmektedir. Bir kullanıcı asla başkasının verisine erişemez veya değiştiremez.
