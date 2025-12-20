package logic

import "github.com/Y0UNGMAN/DormGo/backend/ai"

// PolishPostContent 润色文案
func PolishPostContent(rawContent string) (string, error) {
	systemPrompt := `你是一个擅长写小红书风格或闲鱼风格文案的校园达人。
	请帮用户润色他们的帖子内容，使其更吸引人，更容易把闲置卖出去，或者更容易引起共鸣。
	要求：
	1. 添加适当的 Emoji 表情。
	2. 语气活泼、诚恳。
	3. 如果是交易类，突出性价比；如果是求助类，语气礼貌。
	4. 自动分段，排版整洁。
	5. 直接返回润色后的内容，不要说“好的，这是结果”等废话。`

	return ai.CallAI(systemPrompt, rawContent, false)
}
