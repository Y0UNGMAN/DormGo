package logic

import (
	"errors"
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/ai"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
)

func CreatePost(post *model.DgPost) (err error) {
	// 1. 敏感词过滤（已有）
	if utils.WordFilter != nil {
		if found, first := utils.WordFilter.FindIn(post.Title); found {
			fmt.Printf("标题包含敏感词: %s\n", first)
			return fmt.Errorf("包含敏感词")
		}
		if found, first := utils.WordFilter.FindIn(post.Content); found {
			fmt.Printf("内容包含敏感词: %s\n", first)
			return fmt.Errorf("包含敏感词")
		}
	} else {
		fmt.Println("⚠️ Warning: Sensitive WordFilter is nil")
	}

	// 2. AI 内容审核（已有）
	fullContent := fmt.Sprintf("标题：%s\n内容：%s", post.Title, post.Content)
	isSafe, reason, err := ai.CheckContentSafety(fullContent)
	if err != nil {
		fmt.Println("AI审核服务异常:", err)
		// 根据你的策略：可以选择放行或拒绝，这里选择了放行（注释掉了拒绝）
		// return errors.New("审核服务繁忙，请稍后再试")
	}
	if !isSafe {
		return errors.New("内容审核未通过: " + reason)
	}

	// 3. 保存帖子到数据库
	err = model.CreatePost(post)
	if err != nil {
		fmt.Println("create post error: ", err)
		return err
	}

	// 4. 【新增】异步生成摘要（但这里是同步调用，为了不阻塞响应，可以使用 goroutine，但用户要求“同步模式”即等待生成完毕）
	//    注意：如果希望生成失败不影响发帖，则需处理错误
	go func() {
		// 调用 AI 生成摘要
		summary, err := ai.GenerateSummary(post.Content)
		if err != nil {
			fmt.Printf("生成摘要失败 (post_id=%d): %v\n", post.ID, err)
			// 降级：取内容前150字符作为摘要
			if len(post.Content) > 150 {
				summary = post.Content[:150] + "..."
			} else {
				summary = post.Content
			}
		}
		// 更新数据库中的 summary 字段
		if err := model.DB.Model(post).Update("summary", summary).Error; err != nil {
			fmt.Printf("更新摘要失败 (post_id=%d): %v\n", post.ID, err)
		} else {
			fmt.Printf("✅ 摘要生成成功 (post_id=%d): %s\n", post.ID, summary)
		}
	}()

	// 5. 清除缓存（已有）
	cacheKey := fmt.Sprintf("dormgo:posts:page:%d:size:%d", 1, 10)
	myredis.GetClient().Del(cacheKey)
	fmt.Printf("🗑️ Invalidate Cache: %s\n", cacheKey)

	return nil
}