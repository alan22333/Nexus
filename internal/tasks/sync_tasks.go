package tasks

import (
	"Nuxus/internal/dao"
	"context"
	"log"
	"strconv"
	"strings"
)

func SyncViewCountsToDB() {
	log.Println("开始同步浏览量到DB")

	ctx := context.Background()
	var cursor uint64
	var keys []string
	var err error

	// 1. 使用 SCAN 安全地遍历所有帖子的浏览量 Key
	// 我们每次扫描100个key，直到游标回到0
	matchPattern := "nexus:post:view:*"
	for {
		var scanResult []string
		scanResult, cursor, err = dao.RedisClient.Scan(ctx, cursor, matchPattern, 100).Result()
		if err != nil {
			log.Printf("扫描 Redis Key 失败: %v", err)
			return // 发生错误，终止本次任务
		}
		keys = append(keys, scanResult...)
		if cursor == 0 { // 遍历完成
			break
		}
	}

	if len(keys) == 0 {
		log.Println("没有需要同步的浏览量数据。")
		return
	}
	log.Printf("发现 %d 个需要同步的帖子浏览量。", len(keys))

	// 2. 遍历所有 Key，进行数据同步
	for _, key := range keys {
		// 3. 施展“乾坤挪移”大法 (GETSET)，原子性地获取并清零
		// GETSET key 0 会返回 key 的旧值，然后将 key 的值设为 "0"
		// 这确保了我们拿到的是准确的增量，并且不会丢失在操作期间的新增浏览
		incrementStr, err := dao.RedisClient.GetSet(ctx, key, "0").Result()
		if err != nil {
			log.Printf("获取并重置 Key [%s] 失败: %v", key, err)
			continue // 跳过这个 key，处理下一个
		}

		increment, _ := strconv.ParseInt(incrementStr, 10, 64)
		if increment == 0 {
			continue // 如果增量为0，无需更新数据库
		}

		// 4. 从 Key 中解析出 Post ID
		// key 的格式是 "nexus:post:view:123"
		parts := strings.Split(key, ":")
		if len(parts) != 4 {
			log.Printf("无效的 Key 格式: %s", key)
			continue
		}
		postID, _ := strconv.ParseUint(parts[3], 10, 64)
		if postID == 0 {
			continue
		}

		// 5. 将增量归入“宗门宝库” (MySQL)
		err = dao.AddPostViewCount(uint(postID), int(increment))
		if err != nil {
			log.Printf("更新帖子ID [%d] 的浏览量失败: %v。丢失的增量: %d", postID, err, increment)
			// 重要：这里需要有错误处理策略。最简单的是记录日志。
			// 复杂些可以把失败的任务放入一个“重试队列”。
			return
		}
	}
	log.Println("同步帖子浏览量任务完成。")
}
