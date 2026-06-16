# Ercan Aziz'in Mobil Backend Görevleri

**Mobil Front-end ile Back-end Bağlanmış Test Videosu:** [Link buraya eklenecek](https://example.com)

**REST API Adresi:** `https://habitup-production.up.railway.app`

---

## Genel Mimari

Mobil uygulama, Railway platformunda çalışan Go (Gin) tabanlı REST API ile HTTP üzerinden haberleşmektedir. Tüm API istekleri `lib/services/api_service.dart` dosyasındaki `ApiService` sınıfı üzerinden merkezi olarak yönetilmektedir. Kimlik doğrulama gerektiren endpoint'lerde JWT token otomatik olarak `Authorization: Bearer <token>` başlığına eklenmektedir. Token, oturum açıldığında `SharedPreferences`'ta saklanmakta, oturum kapatıldığında silinmektedir.

---

## Gereksinim 1 — Kullanıcı Kaydı API Entegrasyonu

**Endpoint:** `POST /api/auth/register`

**Servis Metodu:** `ApiService.register(username, email, password)`

Kullanıcının girdiği ad, e-posta ve şifre bilgileri JSON formatında API'ye gönderilmektedir. Şifre, sunucu tarafında bcrypt ile hashlenerek MongoDB'ye kaydedilmektedir. Dönüş değeri olarak HTTP durum kodu ve mesaj işlenmektedir. `201 Created` dönerse başarılı kabul edilip giriş ekranına yönlendirme yapılmaktadır.

```dart
final res = await ApiService.register(username, email, password);
// res['statusCode'] == 201 → başarılı
```

---

## Gereksinim 2 — Kullanıcı Girişi API Entegrasyonu

**Endpoint:** `POST /api/auth/login`

**Servis Metodu:** `ApiService.login(email, password)`

E-posta ve şifre API'ye gönderilmektedir. Başarılı yanıtta (`200 OK`) dönen JWT token `SharedPreferences`'a kaydedilmekte ve sonraki tüm isteklerde otomatik olarak kullanılmaktadır. Token geçerlilik süresi 24 saattir.

```dart
final token = await ApiService.login(email, password);
// token != null → başarılı giriş, habits ekranına geç
```

---

## Gereksinim 3 — Yeni Alışkanlık Oluşturma API Entegrasyonu

**Endpoint:** `POST /api/habits`

**Servis Metodu:** `ApiService.createHabit(name, description)`

Alışkanlık adı ve açıklama JSON gövdesiyle API'ye gönderilmektedir. `Authorization` başlığı ile JWT doğrulaması yapılmaktadır. `201 Created` dönmesi durumunda işlem başarılı sayılmakta ve alışkanlık listesi yeniden yüklenmektedir.

```dart
final ok = await ApiService.createHabit(name, description);
// ok == true → liste yenile
```

---

## Gereksinim 4 — Alışkanlık Listeleme API Entegrasyonu

**Endpoint:** `GET /api/habits`

**Servis Metodu:** `ApiService.getHabits()`

JWT ile korunan bu endpoint, giriş yapmış kullanıcıya ait tüm alışkanlıkları JSON dizisi olarak döndürmektedir. Dönen veri `Habit.fromJson()` ile model nesnelerine dönüştürülerek `ListView`'da gösterilmektedir. Kullanıcı yalnızca kendi verilerini görmektedir (sunucu tarafı filtreleme).

```dart
final data = await ApiService.getHabits();
_habits = data.map((e) => Habit.fromJson(e)).toList();
```

---

## Gereksinim 5 — Alışkanlık Tamamlama API Entegrasyonu

**Endpoint:** `POST /api/habits/{id}/check`

**Servis Metodu:** `ApiService.checkHabit(id)`

Seçilen alışkanlığın ID'si URL'ye eklenerek POST isteği gönderilmektedir. Sunucu, aynı gün için tekrar işaretlemeyi engellemektedir (idempotent davranış). Başarılı yanıtta RabbitMQ üzerinden `habit.checked` kuyruğuna event yayınlanmaktadır.

```dart
final ok = await ApiService.checkHabit(habit.id);
// ok == true → UI'da checkbox işaretli göster
```

---

## Gereksinim 6 — Alışkanlık Güncelleme API Entegrasyonu

**Endpoint:** `PUT /api/habits/{id}`

**Servis Metodu:** `ApiService.updateHabit(id, name, description)`

Alışkanlık ID'si URL'de, güncellenecek alanlar JSON gövdesinde gönderilmektedir. Sunucu yalnızca gönderilen alanları güncellemekte (partial update), gönderilmeyen alanları korumaktadır. `200 OK` dönmesi durumunda liste yeniden yüklenmektedir.

```dart
final ok = await ApiService.updateHabit(id, newName, newDesc);
```

---

## Gereksinim 7 — Alışkanlık Silme API Entegrasyonu

**Endpoint:** `DELETE /api/habits/{id}`

**Servis Metodu:** `ApiService.deleteHabit(id)`

Alışkanlık ID'si URL'ye eklenerek DELETE isteği gönderilmektedir. Sunucu, alışkanlığı ve ona bağlı tüm tamamlama kayıtlarını (`checks` koleksiyonu) MongoDB'den silmektedir. `204 No Content` dönmesi durumunda işlem başarılı kabul edilmektedir.

```dart
final ok = await ApiService.deleteHabit(habit.id);
```

---

## Gereksinim 8 — İşaretleme Geri Alma API Entegrasyonu

**Endpoint:** `DELETE /api/habits/{id}/check`

**Servis Metodu:** `ApiService.uncheckHabit(id)`

O günkü tamamlama kaydının silinmesi için DELETE isteği gönderilmektedir. Sunucu bugünün tarihine ait kaydı `checks` koleksiyonundan kaldırmakta, streak algoritması yeniden hesaplanmaktadır. `204 No Content` başarı yanıtıdır.

```dart
final ok = await ApiService.uncheckHabit(habit.id);
```

---

## Gereksinim 9 — İstatistik API Entegrasyonu

**Endpoint:** `GET /api/habits/{id}/stats`

**Servis Metodu:** `ApiService.getStats(id)`

Seçilen alışkanlığa ait istatistikler API'den alınmaktadır. Dönen JSON nesnesi şu alanları içermektedir:

| Alan | Tip | Açıklama |
|------|-----|----------|
| `currentStreak` | int | Mevcut kesintisiz gün serisi |
| `longestStreak` | int | Tüm zamanlardaki en uzun seri |
| `completionRate` | double | Yüzde tamamlanma oranı |
| `totalChecks` | int | Toplam tamamlama sayısı |

```dart
final stats = await ApiService.getStats(habit.id);
// stats['currentStreak'], stats['completionRate'] vb.
```

---

## Gereksinim 10 — Oturum Kapatma API Entegrasyonu

**Endpoint:** `POST /api/auth/logout`

**Servis Metodu:** `ApiService.logout()`

Logout isteği gönderildiğinde sunucu JWT token'ı Redis blacklist'e eklemektedir. Mobil tarafta ise `SharedPreferences`'tan token silinerek yerel oturum da sonlandırılmaktadır. Bu çift taraflı yaklaşım sayesinde token çalınsa bile geçersiz sayılmaktadır.

```dart
await ApiService.logout();
// SharedPreferences token silindi + Redis blacklist eklendi
```

---

## Kullanılan Teknolojiler

| Teknoloji | Kullanım Amacı |
|-----------|----------------|
| `http` paketi | REST API istekleri (GET, POST, PUT, DELETE) |
| `shared_preferences` | JWT token'ın cihazda kalıcı saklanması |
| `dart:convert` | JSON encode/decode işlemleri |
| JWT (sunucu tarafı) | Kimlik doğrulama ve yetkilendirme |
| HTTPS | Güvenli API iletişimi (Railway SSL) |
