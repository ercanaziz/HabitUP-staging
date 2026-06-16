# Mobil Backend (REST API Bağlantısı) Görev Dağılımı

**REST API Adresi:** [habitup-staging-production.up.railway.app](https://habitup-staging-production.up.railway.app)

Bu dokümanda, mobil uygulamanın REST API ile iletişimini sağlayan backend entegrasyon görevleri listelenmektedir. Her grup üyesi, kendisine atanan API endpoint'lerinin mobil uygulamadan çağrılması ve yönetilmesinden sorumludur.

---

## Grup Üyelerinin Mobil Backend Görevleri

1. [Ercan Aziz'in Mobil Backend Görevleri](Ercan-Aziz/Ercan-Aziz-Mobil-Backend-Gorevleri.md)

---

## Genel Mobil Backend Prensipleri

- **İletişim:** HTTP paketi ile REST API çağrıları
- **Kimlik Doğrulama:** JWT — her istekte `Authorization: Bearer <token>` header'ı
- **Veri Formatı:** JSON
