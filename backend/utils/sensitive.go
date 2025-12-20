package utils

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model" // 替换为你的model路径
	"github.com/importcjj/sensitive"
)

// 全局过滤器实例
var WordFilter *sensitive.Filter

// InitSensitiveFilter 初始化过滤器 (在 main.go 启动时调用)
func InitSensitiveFilter() {
	WordFilter = sensitive.New()

	// 1. 从数据库加载所有敏感词
	var words []model.DgSensitiveWord
	// 注意：这里需要确保数据库连接 model.DB 已经初始化
	if err := model.DB.Find(&words).Error; err != nil {
		fmt.Println("加载敏感词库失败:", err)
		return
	}

	// 2. 添加到过滤器内存中
	for _, w := range words {
		WordFilter.AddWord(w.Word)
	}
	fmt.Printf("✅ 敏感词库加载完毕，共加载 %d 个词\n", len(words))
}

// UpdateFilter 当管理员增删词汇时，调用此方法更新内存
func UpdateFilter() {
	// 简单粗暴策略：重新初始化
	// 生产环境中可以使用更精细的 Add/Remove 接口，但重载对于毕设足够快
	InitSensitiveFilter()
}
