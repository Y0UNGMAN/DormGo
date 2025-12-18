package logic

import (
	"errors"
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/ai"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
)

func CreatePost(post *model.DgPost) (err error) {
	fullContent := fmt.Sprintf("标题：%s\n内容：%s", post.Title, post.Content)
	isSafe, reason, err := ai.CheckContentSafety(fullContent)
	if err != nil {
		// AI 调用失败时的处理策略：
		// 策略A: 报错，暂时无法发帖 (更安全)
		// 策略B: 放行，记录日志 (更可用)
		fmt.Println("AI审核服务异常:", err)
		// return errors.New("审核服务繁忙，请稍后再试") // 选用策略A
	}
	if !isSafe {
		return errors.New("内容审核未通过: " + reason)
	}
	//1 保存到数据库
	err = model.CreatePost(post)
	//2 返回
	if err != nil {
		fmt.Println("create post error: ", err)
		return err
	}

	cacheKey := fmt.Sprintf("dormgo:posts:page:%d:size:%d", 1, 10)
	myredis.GetClient().Del(cacheKey)
	fmt.Printf("🗑️ Invalidate Cache: %s\n", cacheKey)

	return nil
}
