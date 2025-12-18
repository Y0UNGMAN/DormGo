package ai

import (
	"encoding/json"
	"strings"
)

type AuditResult struct {
	Status string `json:"status"` // PASS 或 BLOCK
	Reason string `json:"reason"`
}

// CheckContentSafety 审核内容
func CheckContentSafety(text string) (bool, string, error) {
	systemPrompt := `你是一个校园社区管理员。请审核用户的帖子内容。
	如果内容包含：严重脏话辱骂、诈骗广告、色情描述、暴力恐吓、政治敏感等违规信息，请拒绝。
	如果是正常的校园生活分享、闲置交易、交友、吐槽（非人身攻击），请通过。
	请务必返回合法的 JSON 格式，不要包含 Markdown 格式（如 '''json ），格式如下：
	{"status": "BLOCK", "reason": "包含涉黄信息"} 或 {"status": "PASS", "reason": "内容正常"}`

	// 调用 AI (开启 JSON 模式)
	respStr, err := CallAI(systemPrompt, text, true)
	if err != nil {
		// 如果 AI 挂了，为了不影响用户发帖，可以选择默认通过，或者报错
		// 这里选择报错，安全第一
		return true, "", err
	}

	// 清洗数据：有些模型即使开了 JSON 模式也可能带 ```json
	respStr = strings.ReplaceAll(respStr, "```json", "")
	respStr = strings.ReplaceAll(respStr, "```", "")

	var res AuditResult
	if err := json.Unmarshal([]byte(respStr), &res); err != nil {
		return false, "", err
	}

	if res.Status == "BLOCK" {
		return false, res.Reason, nil
	}

	return true, "", nil
}
