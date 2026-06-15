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
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	RDB = redis.NewClient(&redis.Options{Addr: addr})

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
