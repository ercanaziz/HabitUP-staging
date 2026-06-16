package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Connect() {
	var opts *redis.Options
	var err error

	// Railway ve diğer cloud sağlayıcılar REDIS_URL olarak tam URL verir
	// Örnek: redis://:password@host:port veya rediss://:password@host:port
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		opts, err = redis.ParseURL(redisURL)
		if err != nil {
			log.Fatal("Redis URL parse hatası:", err)
		}
		log.Println("Redis: REDIS_URL kullanılıyor")
	} else {
		// Lokal geliştirme: REDIS_ADDR (host:port)
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}
		opts = &redis.Options{Addr: addr}
		log.Println("Redis: REDIS_ADDR kullanılıyor →", addr)
	}

	RDB = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis bağlantı hatası:", err)
	}
	log.Println("Redis bağlantısı başarılı")
}

// BlacklistToken logout sırasında JWT token'ı geçersiz kılar.
// Token Redis'e TTL ile yazılır; middleware her istekte kontrol eder.
func BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return RDB.Set(ctx, "blacklist:"+token, 1, ttl).Err()
}

func IsBlacklisted(ctx context.Context, token string) bool {
	val, err := RDB.Exists(ctx, "blacklist:"+token).Result()
	return err == nil && val > 0
}

