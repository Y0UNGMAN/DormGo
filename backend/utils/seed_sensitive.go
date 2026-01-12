package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"gorm.io/gorm/clause"
)

// SeedAllSensitiveWords 批量导入所有预定义的敏感词库
func SeedAllSensitiveWords() {
	dictFiles := map[string]string{
		"sensitive-stop-words-master/广告.txt":          "ad",
		"sensitive-stop-words-master/色情类.txt":         "porn",
		"sensitive-stop-words-master/政治类.txt":         "politics",
		"sensitive-stop-words-master/网址.txt":          "url",
		"sensitive-stop-words-master/涉枪涉爆违法信息关键词.txt": "illegal",
	}

	fmt.Println("🚀 开始批量加载敏感词库...")

	for path, cat := range dictFiles {
		SeedSensitiveWords(path, cat)
	}

	fmt.Println("✨ 所有词库加载流程结束。")
}

// SeedSensitiveWords 从文件导入敏感词
func SeedSensitiveWords(filePath string, category string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("⚠️ 跳过: 未找到文件 %s\n", filePath)
		return
	}
	defer file.Close()

	fmt.Printf("⏳ 正在导入 [%s] ... ", category)

	scanner := bufio.NewScanner(file)
	var newWords []model.DgSensitiveWord

	// 用 map 做文件内的简单去重 (防止文件里自己就有重复行)
	fileUniqueMap := make(map[string]bool)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		text = strings.Trim(text, ",") // 去掉可能的尾部逗号

		// 简单过滤：跳过空行或过短的词
		if text == "" || len(text) < 1 {
			continue
		}

		if !fileUniqueMap[text] {
			newWords = append(newWords, model.DgSensitiveWord{
				Word:     text,
				Category: category,
			})
			fileUniqueMap[text] = true
		}
	}

	if len(newWords) > 0 {
		result := model.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "word"}}, // 指定冲突检测的列
			DoNothing: true,                            // 冲突时忽略
		}).CreateInBatches(newWords, 100)

		if result.Error != nil {
			fmt.Printf("❌ 失败: %v\n", result.Error)
		} else {
			// result.RowsAffected 返回实际插入成功的行数（不包含被忽略的重复项）
			fmt.Printf("✅ 成功新增 %d 个词 (文件共 %d 个)\n", result.RowsAffected, len(newWords))
		}
	} else {
		fmt.Println("⚠️ 文件为空或无有效词汇")
	}
}
