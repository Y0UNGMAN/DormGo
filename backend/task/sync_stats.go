package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
)

// InitSyncTask 初始化定时任务
func InitSyncTask() {
	// 启动一个协程，每分钟执行一次同步
	go func() {
		ticker := time.NewTicker(time.Minute * 1) // 设置同步间隔
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("⏰ 开始执行点赞/浏览数回写任务...")
			SyncLikesToDB()
			SyncViewsToDB()
		}
	}()
}

// SyncLikesToDB 同步点赞数
func SyncLikesToDB() {
	client := myredis.GetClient()
	// 1. 扫描所有相关的 Key (dormgo:post:like_inc:*)
	// 生产环境建议用 Scan，这里为了简单演示用 Keys
	keys, err := client.Keys("dormgo:post:like_inc:*").Result()
	if err != nil {
		fmt.Println("Sync task error:", err)
		return
	}

	for _, key := range keys {
		// key 格式: dormgo:post:like_inc:101
		parts := strings.Split(key, ":")
		if len(parts) < 4 {
			continue
		}
		postIDStr := parts[3]
		postID, _ := strconv.Atoi(postIDStr)

		// 2. 获取增量值
		valStr, err := client.Get(key).Result()
		if err != nil {
			continue
		}
		incVal, _ := strconv.Atoi(valStr)

		if incVal == 0 {
			continue
		}

		// 3. 更新 MySQL
		err = model.UpdateLikeCount(uint(postID), incVal)

		if err == nil {
			client.DecrBy(key, int64(incVal))
			fmt.Printf("✅ 同步帖子[%d]点赞数: %+d\n", postID, incVal)
		} else {
			fmt.Printf("❌ 同步帖子[%d]点赞失败: %v\n", postID, err)
		}
	}
}

// SyncViewsToDB 同步浏览量 (逻辑同上)
func SyncViewsToDB() {
	client := myredis.GetClient()
	keys, err := client.Keys("dormgo:post:view_inc:*").Result()
	if err != nil {
		return
	}

	for _, key := range keys {
		parts := strings.Split(key, ":")
		if len(parts) < 4 {
			continue
		}
		postID, _ := strconv.Atoi(parts[3])

		valStr, _ := client.Get(key).Result()
		incVal, _ := strconv.Atoi(valStr)

		if incVal == 0 {
			continue
		}

		err = model.UpdateViewCount(uint(postID), incVal)
		if err == nil {
			client.DecrBy(key, int64(incVal))
		}
	}
}
