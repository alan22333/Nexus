package dao

import (
	"Nuxus/configs"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var (
	// RedisClient *redis.Client
	Ctx = context.Background()
)

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{client: client}
}

func NewClient(config *configs.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB, // 0 default-DB
		PoolSize: 10,              // conn-pool size
	})
	if err := client.Ping(Ctx).Err(); err != nil {
		panic("Redis connect fail: " + err.Error())
	}

	log.Println("Redis connection successful!")
	return client
}

// func InitRedis() {
// 	RedisClient = redis.NewClient(&redis.Options{
// 		Addr:     Conf.Redis.Addr,
// 		Password: Conf.Redis.Password,
// 		DB:       Conf.Redis.DB, // 0 default-DB
// 		PoolSize: 10,            // conn-pool size
// 	})
// 	if err := RedisClient.Ping(Ctx).Err(); err != nil {
// 		panic("Redis connect fail: " + err.Error())
// 	}
// 	log.Println("Redis connection successful!")
// }

// 定义 Redis Keys 的前缀，方便管理
const (
	PrefixVerifyCode    = "nexus:verify_code:%s"   // %s 是邮箱
	PrefixSendCooldown  = "nexus:send_cooldown:%s" // %s 是邮箱
	PrefixPostViewCount = "nexus:post:view:%s"     // %d 是帖子 ID
	KeyPopularPosts     = "nexus:posts:popular"    // 热门帖子的 ZSET Key
)

// 封装需要的方法
func (r *RedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.Scan(ctx, cursor, match, count).Result()
}

func (r *RedisClient) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return r.client.GetSet(ctx, key, value).Result()
}

func (r *RedisClient) IncrementPostViewCount(postId uint) error {
	key := fmt.Sprintf(PrefixPostViewCount, postId)
	// INCR 命令：如果 key 不存在，会先创建为 0 再加 1。
	// 所以无需担心初始化问题。
	return r.client.Incr(Ctx, key).Err()
}

// 更新帖子在热门榜单上的分数
func (r *RedisClient) IncrementPostRank(postID uint, increment float64) error {
	key := KeyPopularPosts
	// ZINCRBY 命令：为一个 member 增加指定的分数
	return r.client.ZIncrBy(Ctx, key, increment, fmt.Sprint(postID)).Err()
}

// GetPopularPostIDs 从榜单获取 Top N 的帖子 ID
func (r *RedisClient) GetPopularPostIDs(limit int64) ([]string, error) {
	key := KeyPopularPosts
	// ZREVRANGE 命令：按分数从高到低返回指定区间的成员
	return r.client.ZRevRange(Ctx, key, 0, limit-1).Result()
}
func (r *RedisClient) GetPostViewCount(postId string) (int64, error) {
	key := fmt.Sprintf(PrefixPostViewCount, postId)
	// 如果 key 不存在，GET 会返回错误，但 Int() 会处理成 0。
	return r.client.Get(Ctx, key).Int64()
}

func (r *RedisClient) SetVerifyCode(req_email, code string, duration time.Duration) error {
	key := fmt.Sprintf(PrefixVerifyCode, req_email)
	return r.client.Set(Ctx, key, code, duration).Err()
}

func (r *RedisClient) GetVerificationCode(email string) (string, error) {
	key := fmt.Sprintf(PrefixVerifyCode, email)
	return r.client.Get(Ctx, key).Result()
}

// DelVerificationCode 从 Redis 中删除验证码
func (r *RedisClient) DelVerificationCode(email string) error {
	key := fmt.Sprintf(PrefixVerifyCode, email)
	return r.client.Del(Ctx, key).Err()
}

// CheckSendCooldown 检查发送冷却时间
func (r *RedisClient) CheckSendCooldown(email string) (bool, error) {
	key := fmt.Sprintf(PrefixSendCooldown, email)
	// 使用 SetNX，如果 key 不存在则设置并返回 true，如果 key 已存在则不操作并返回 false
	// 这利用了原子操作来避免竞态条件
	wasSet, err := r.client.SetNX(Ctx, key, 1, 60*time.Second).Result()
	if err != nil {
		return false, err
	}
	// 如果 wasSet 为 false，说明在60秒内已经发送过了
	return !wasSet, nil
}
