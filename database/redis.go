package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	log.Printf("🔍 REDIS_HOST: %s", redisHost)
	log.Printf("🔍 REDIS_PORT: %s", redisPort)
	log.Printf("🔍 REDIS_PASSWORD: %s", redisPassword)

	if redisHost == "" {
		log.Printf("⚠️ REDIS_HOST boş, gps_redis atanıyor")
		redisHost = "gps_redis"
	}
	if redisPort == "" {
		log.Printf("⚠️ REDIS_PORT boş, 6379 atanıyor")
		redisPort = "6379"
	}
	if redisPassword == "" {
		log.Printf("🔓 Redis şifresi yok")
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	log.Printf("📡 Redis'e bağlanmayı deniyor: %s", addr)

	rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     redisPassword,
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	for i := 0; i < 5; i++ {
		pong, err := rdb.Ping(ctx).Result()
		if err == nil {
			log.Printf("✅ Redis bağlantısı başarılı: %s - %s", addr, pong)
			return
		}
		log.Printf("❌ Redis bağlantı denemesi #%d başarısız: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	log.Fatalf("❌ Redis bağlantısı başarısız: %s", addr)
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return rdb
}
