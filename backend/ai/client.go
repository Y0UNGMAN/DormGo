package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// OpenAI 标准请求格式
type ChatRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	Temperature    float64         `json:"temperature"`               // 0-2，越低越严谨，越高越发散
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"` // 用于强制 JSON 输出
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI 标准响应格式
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// 通用 AI 调用函数
func CallAI(systemPrompt string, userContent string, jsonMode bool) (string, error) {
	reqBody := ChatRequest{
		Model: viper.GetString("api.ApiModel"),
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userContent},
		},
		Temperature: 1.3,
	}

	if jsonMode {
		reqBody.ResponseFormat = &ResponseFormat{Type: "json_object"}
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", viper.GetString("api.ApiUrl"), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+viper.GetString("api.ApiKey"))

	client := &http.Client{Timeout: 30 * time.Second} // 设置超时
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// 解析响应
	var aiResp ChatResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", err
	}

	if aiResp.Error.Message != "" {
		return "", fmt.Errorf("AI API Error: %s", aiResp.Error.Message)
	}

	if len(aiResp.Choices) == 0 {
		return "", errors.New("AI 没有返回内容")
	}

	return aiResp.Choices[0].Message.Content, nil
}
