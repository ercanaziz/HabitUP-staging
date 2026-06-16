package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Connect() {
	opts := buildRedisOptions()
	RDB = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis bağlantı hatası:", err)
	}
	log.Println("Redis bağlantısı başarılı")
}

func buildRedisOptions() *redis.Options {
	// Öncelik 1: Tam scheme'li URL (redis:// veya rediss://)
	if url := os.Getenv("REDIS_URL"); url != "" {
		if strings.HasPrefix(url, "redis://") || strings.HasPrefix(url, "rediss://") {
			opts, err := redis.ParseURL(url)
			if err != nil {
				log.Fatal("Redis URL parse hatası:", err)
			}
			log.Println("Redis: REDIS_URL (scheme'li) kullanılıyor")
			return opts
		}
		// Scheme eksikse başına ekle
		log.Printf("Redis: REDIS_URL scheme'siz geldi (%s), redis:// ekleniyor", url)
		url = "redis://" + url
		opts, err := redis.ParseURL(url)
		if err != nil {
			log.Printf("Redis URL yine parse edilemedi: %v — ayrı değişkenlere bakılıyor", err)
		} else {
			return opts
		}
	}

	// Öncelik 2: Railway'in ayrı env değişkenleri (REDISHOST, REDISPORT, REDISPASSWORD)
	if host := os.Getenv("REDISHOST"); host != "" {
		port := os.Getenv("REDISPORT")
		if port == "" {
			port = "6379"
		}
		password := os.Getenv("REDISPASSWORD")
		addr := fmt.Sprintf("%s:%s", host, port)
		log.Printf("Redis: REDISHOST/REDISPORT/REDISPASSWORD kullanılıyor → %s", addr)
		return &redis.Options{
			Addr:     addr,
			Password: password,
		}
	}

	// Öncelik 3: Lokal geliştirme — REDIS_ADDR (host:port, şifresiz)
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	log.Println("Redis: REDIS_ADDR (lokal) kullanılıyor →", addr)
	return &redis.Options{Addr: addr}
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
