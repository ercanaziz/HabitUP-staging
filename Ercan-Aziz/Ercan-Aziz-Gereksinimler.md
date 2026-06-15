**1. Kullanıcı Kaydı**

   **API Metodu:** `POST /api/auth/register`

   **Açıklama:** Kullanıcıların sisteme dahil olmasını sağlar. Kullanıcı adı, e-posta ve güvenli (hashlenmiş) şifre bilgilerini alarak veritabanında yeni bir kullanıcı profili oluşturur. Kayıt başarılı olduktan sonra kullanıcı kendi alışkanlık listesini oluşturmaya hazır hale gelir.

**2. Kullanıcı Girişi**

   **API Metodu:** `POST /api/auth/login`

   **Açıklama:** Kayıtlı kullanıcıların kimlik bilgilerini doğrulayarak sisteme erişimini sağlar. Başarılı giriş sonrasında kullanıcıya sonraki isteklerinde kullanılmak üzere bir erişim anahtarı (JWT) döner. Bu sayede oturum güvenliği sağlanır.

**3. Yeni Alışkanlık Tanımlama**

   **API Metodu:** `POST /api/habits`

   **Açıklama:** Kullanıcının takip etmek istediği yeni bir hedef (Örn: "Kitap Oku") oluşturmasını sağlar. Alışkanlık adı ve açıklama bilgileri veritabanına kaydedilirken, bu veri doğrudan işlemi yapan kullanıcının ID'si ile ilişkilendirilir.

**4. Alışkanlıkları Listeleme**

   **API Metodu:** `GET /api/habits`

   **Açıklama:** Sisteme giriş yapmış olan kullanıcının oluşturduğu tüm aktif alışkanlıkları getirir. Veri izolasyonu sayesinde kullanıcı sadece kendi hedeflerini görür; başkasının verilerine erişemez.

**5. Alışkanlık Durumu Güncelleme**

   **API Metodu:** `POST /api/habits/{id}/check`

   **Açıklama:** Belirli bir alışkanlığın o gün için başarıyla tamamlandığını sisteme işler. Bu işlem, takvim üzerinde ilgili günün "yapıldı" olarak işaretlenmesini sağlar ve seri (streak) hesaplamasının temel verisini oluşturur.

**6. Alışkanlık Güncelleme**

   **API Metodu:** `PUT /api/habits/{id}`

   **Açıklama:** Kullanıcının mevcut bir alışkanlığının adını veya açıklamasını değiştirmesine olanak tanır. Hedefler zamanla evrildiğinde (Örn: "10 sayfa oku" yerine "20 sayfa oku") bu metod kullanılır.

**7. Alışkanlık Silme**

   **API Metodu:** `DELETE /api/habits/{id}`

   **Açıklama:** Kullanıcının artık takip etmek istemediği bir alışkanlığı ve ona bağlı tüm geçmiş tamamlama verilerini sistemden kalıcı olarak kaldırır.

**8. İşaretlemeyi Geri Alma**

   **API Metodu:** `DELETE /api/habits/{id}/check`

   **Açıklama:** Kullanıcının yanlışlıkla yaptığı veya iptal etmek istediği "tamamlandı" işaretini ilgili tarihten kaldırır. Bu işlem yapıldığında mevcut seri (streak) algoritma tarafından yeniden hesaplanır.

**9. İstatistik ve Seri (Streak) Takibi**

   **API Metodu:** `GET /api/habits/{id}/stats`

   **Açıklama:** Belirli bir alışkanlık için kullanıcının performans verilerini sunar. Kaç gün kesintisiz devam edildiğini (current streak), en uzun başarı serisini ve genel tamamlanma oranını hesaplayarak kullanıcıya motivasyonel veri sağlar.

**10. Oturumu Kapatma**

   **API Metodu:** `POST /api/auth/logout`

   **Açıklama:** Kullanıcının aktif oturumunu sonlandırır. Güvenlik protokolü gereği, çıkış yapıldıktan sonra eski erişim anahtarı geçersiz kılınır ve kullanıcı tekrar giriş yapana kadar kişisel verilerine erişim engellenir.
