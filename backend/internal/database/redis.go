package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis يفتح اتصالاً بـ Redis
func ConnectRedis(redisURL string) error {
	if redisURL == "" {
		// إذا لم يكن Redis URL موجوداً، لا تعمل Redis
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("فشل تحليل Redis URL: %w", err)
	}

	RedisClient = redis.NewClient(opt)

	// اختبار الاتصال
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		RedisClient = nil
		return fmt.Errorf("فشل الاتصال بـ Redis: %w", err)
	}

	fmt.Println("✅ تم الاتصال بـ Redis بنجاح")
	return nil
}

// InitRedis يهيئ Redis من متغير البيئة
func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		fmt.Println("⚠️ REDIS_URL غير معرّف، سيتم تجاوز Rate Limiter")
		return nil
	}

	return ConnectRedis(redisURL)
}
