# API Tasarımı - OpenAPI Specification Örneği

**OpenAPI Spesifikasyon Dosyası:** [lamine.yaml](lamine.yaml)

Bu doküman, OpenAPI Specification (OAS) 3.0 standardına göre hazırlanmış örnek bir API tasarımını içermektedir.

## OpenAPI Specification

```yaml
openapi: 3.0.0
info:
  title: HabitUp API
  description: HabitUp alışkanlık takip uygulaması için RESTful API tasarımı.
  version: 1.0.0
servers:
  - url: http://localhost:8080
    description: Yerel Geliştirme Sunucusu

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    UserRegistration:
      type: object
      required:
        - username
        - email
        - password
      properties:
        username:
          type: string
        email:
          type: string
          format: email
        password:
          type: string
          format: password

    UserLogin:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
        password:
          type: string
          format: password

    AuthResponse:
      type: object
      properties:
        token:
          type: string
          description: JWT erişim anahtarı

    Habit:
      type: object
      required:
        - name
      properties:
        id:
          type: string
        name:
          type: string
        description:
          type: string
        createdAt:
          type: string
          format: date-time

    HabitStats:
      type: object
      properties:
        currentStreak:
          type: integer
        longestStreak:
          type: integer
        completionRate:
          type: number
          format: float

    Error:
      type: object
      properties:
        message:
          type: string

security:
  - bearerAuth: []

paths:
  /api/auth/register:
    post:
      summary: Kullanıcı Kaydı
      description: Kullanıcı adı, e-posta ve şifre bilgilerini alarak yeni bir kullanıcı profili oluşturur.
      security: [] # Bu uç nokta için yetkilendirme gerekmez
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserRegistration'
      responses:
        '201':
          description: Kullanıcı başarıyla oluşturuldu
        '400':
          description: Geçersiz istek verisi

  /api/auth/login:
    post:
      summary: Kullanıcı Girişi
      description: Kayıtlı kullanıcıların kimlik bilgilerini doğrulayarak JWT döner.
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserLogin'
      responses:
        '200':
          description: Başarılı giriş
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '401':
          description: Hatalı e-posta veya şifre

  /api/auth/logout:
    post:
      summary: Oturumu Kapatma
      description: Kullanıcının aktif oturumunu sonlandırır ve mevcut token'ı geçersiz kılar.
      responses:
        '200':
          description: Başarıyla çıkış yapıldı

  /api/habits:
    get:
      summary: Alışkanlıkları Listeleme
      description: Sisteme giriş yapmış olan kullanıcının oluşturduğu tüm aktif alışkanlıkları getirir.
      responses:
        '200':
          description: Başarılı listeleme
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Habit'
    post:
      summary: Yeni Alışkanlık Tanımlama
      description: Kullanıcının takip etmek istediği yeni bir hedef oluşturmasını sağlar.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                description:
                  type: string
      responses:
        '201':
          description: Alışkanlık başarıyla oluşturuldu
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Habit'

  /api/habits/{id}:
    put:
      summary: Alışkanlık Güncelleme
      description: Kullanıcının mevcut bir alışkanlığının adını veya açıklamasını değiştirmesine olanak tanır.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                description:
                  type: string
      responses:
        '200':
          description: Alışkanlık başarıyla güncellendi
    delete:
      summary: Alışkanlık Silme
      description: Kullanıcının alışkanlığı ve ona bağlı tüm geçmiş verilerini kalıcı olarak kaldırır.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '204':
          description: Alışkanlık başarıyla silindi

  /api/habits/{id}/check:
    post:
      summary: Alışkanlık Durumu Güncelleme
      description: Belirli bir alışkanlığın o gün için başarıyla tamamlandığını işaretler.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Başarıyla tamamlandı olarak işaretlendi
    delete:
      summary: İşaretlemeyi Geri Alma
      description: Yanlışlıkla yapılan veya iptal edilmek istenen 'tamamlandı' işaretini kaldırır.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '204':
          description: İşaretleme başarıyla geri alındı

  /api/habits/{id}/stats:
    get:
      summary: İstatistik ve Seri Takibi
      description: Alışkanlığa ait kesintisiz devam süresi (streak) ve tamamlanma oranını döner.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: İstatistik verileri başarıyla getirildi
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/HabitStats'