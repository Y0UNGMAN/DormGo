package ai

import (
    "encoding/json"
    "strings"
)

// GenerateSummary 调用 AI 生成帖子摘要
func GenerateSummary(content string) (string, error) {
    systemPrompt := `你是文本摘要助手。请为给定的文本生成简明扼要的摘要。
要求：
1. 必须使用中文。
2. 摘要控制在 100 字以内。
3. 保留关键信息：事件、时间、人物、结论（如有）。
4. 输出格式：{"summary": "生成的摘要内容"}。

注意：只输出合法的 JSON，不要包含 Markdown 代码块标记。`

    // 调用 AI（开启 JSON 模式）
    respStr, err := CallAI(systemPrompt, content, true)
    if err != nil {
        return "", err
    }

    // 清洗可能的 Markdown 标记
    respStr = strings.ReplaceAll(respStr, "```json", "")
    respStr = strings.ReplaceAll(respStr, "```", "")

    var result struct {
        Summary string `json:"summary"`
    }
    if err := json.Unmarshal([]byte(respStr), &result); err != nil {
        return "", err
    }

    return result.Summary, nil
}