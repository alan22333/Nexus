package dao

import (
	. "Nuxus/configs"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Conf.Redis.Addr,
		Password: Conf.Redis.Password,
		DB:       Conf.Redis.DB, // 0 default-DB
		PoolSize: 10,            // conn-pool size
	})
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		panic("Redis connect fail: " + err.Error())
	}
	log.Println("Redis connection successful!")
}

// 定义 Redis Keys 的前缀，方便管理
const (
	PrefixVerifyCode   = "nexus:verify_code:%s"   // %s 是邮箱
	PrefixSendCooldown = "nexus:send_cooldown:%s" // %s 是邮箱
)

func SetVerifyCode(req_email, code string, duration time.Duration) error {
	key := fmt.Sprintf(PrefixVerifyCode, req_email)
	return RedisClient.Set(Ctx, key, code, duration).Err()
}

func GetVerificationCode(email string) (string, error) {
	key := fmt.Sprintf(PrefixVerifyCode, email)
	return RedisClient.Get(Ctx, key).Result()
}

// DelVerificationCode 从 Redis 中删除验证码
func DelVerificationCode(email string) error {
	key := fmt.Sprintf(PrefixVerifyCode, email)
	return RedisClient.Del(Ctx, key).Err()
}

// CheckSendCooldown 检查发送冷却时间
func CheckSendCooldown(email string) (bool, error) {
	key := fmt.Sprintf(PrefixSendCooldown, email)
	// 使用 SetNX，如果 key 不存在则设置并返回 true，如果 key 已存在则不操作并返回 false
	// 这利用了原子操作来避免竞态条件
	wasSet, err := RedisClient.SetNX(Ctx, key, 1, 60*time.Second).Result()
	if err != nil {
		return false, err
	}
	// 如果 wasSet 为 false，说明在60秒内已经发送过了
	return !wasSet, nil
}
