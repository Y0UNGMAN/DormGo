package utils

import (
	"testing"

	"github.com/importcjj/sensitive"
)

// TestSensitiveFilter 测试敏感词过滤器的核心逻辑
func TestSensitiveFilter(t *testing.T) {
	// ==========================================
	// 1. 准备阶段 (Setup)
	// ==========================================
	// 注意：我们这里不调用 InitSensitiveFilter()，因为它会尝试连接数据库。
	// 在单元测试中，我们手动初始化过滤器并添加几个“假”的敏感词用于测试。
	WordFilter = sensitive.New()

	// 模拟加载敏感词
	mockWords := []string{"笨蛋", "垃圾", "诈骗", "加我V"}
	for _, w := range mockWords {
		WordFilter.AddWord(w)
	}
	t.Logf("单元测试环境初始化：已手动添加 %d 个测试敏感词", len(mockWords))

	// ==========================================
	// 2. 定义测试用例 (Table-Driven Tests)
	// ==========================================
	tests := []struct {
		name      string // 用例名称
		input     string // 输入的文本
		wantFound bool   // 期望结果：是否包含敏感词
		wantWord  string // 期望找到的那个敏感词（如果是 false 则为空）
	}{
		{
			name:      "正常文本",
			input:     "今天天气真好，我们去图书馆学习吧。",
			wantFound: false,
			wantWord:  "",
		},
		{
			name:      "包含单个敏感词",
			input:     "你就是一个大笨蛋，真的。",
			wantFound: true,
			wantWord:  "笨蛋",
		},
		{
			name:      "包含其他敏感词",
			input:     "这简直是垃圾项目。",
			wantFound: true,
			wantWord:  "垃圾",
		},
		{
			name:      "隐晦的广告诈骗",
			input:     "想要兼职吗？加我V私聊。",
			wantFound: true,
			wantWord:  "加我V",
		},
		{
			name:      "混合文本",
			input:     "DormGo社区很棒", // 正常
			wantFound: false,
			wantWord:  "",
		},
	}

	// ==========================================
	// 3. 执行测试 (Run)
	// ==========================================
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用 sensitive 库的 FindIn 方法
			// FindIn 返回两个值：是否发现(bool)，发现的第一个词(string)
			found, firstWord := WordFilter.FindIn(tt.input)

			// 验证 1: 是否发现结果符合预期
			if found != tt.wantFound {
				t.Errorf("测试失败 [%s]: 输入内容 '%s', 期望 found=%v, 实际 found=%v",
					tt.name, tt.input, tt.wantFound, found)
			}

			// 验证 2: 发现的词是否正确（仅在期望发现时验证）
			if tt.wantFound && firstWord != tt.wantWord {
				t.Errorf("测试失败 [%s]: 期望找到敏感词 '%s', 但实际找到了 '%s'",
					tt.name, tt.wantWord, firstWord)
			}
		})
	}
}
