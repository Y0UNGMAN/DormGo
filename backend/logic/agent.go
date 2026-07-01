package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const pendingActionTTL = 5 * time.Minute
const agentSessionTTL = 24 * time.Hour

type AgentChatResponse struct {
	Reply      string                 `json:"reply"`
	Type       string                 `json:"type"`
	Candidates []*model.ApiPostDetail `json:"candidates,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

type agentPythonRequest struct {
	UserID        uint   `json:"user_id"`
	Message       string `json:"message"`
	InternalToken string `json:"internal_token"`
	GoBaseURL     string `json:"go_base_url"`
}

type PendingSignupAction struct {
	ActionID  string    `json:"action_id"`
	Action    string    `json:"action"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type AgentCandidateRef struct {
	PostID    uint   `json:"post_id"`
	Title     string `json:"title"`
	IsLimited bool   `json:"is_limited"`
}

type AgentHistoryItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AgentSession struct {
	UserID         uint                `json:"user_id"`
	LastCandidates []AgentCandidateRef `json:"last_candidates"`
	History        []AgentHistoryItem  `json:"history"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

func PendingActionKey(userID uint) string {
	return fmt.Sprintf("dormgo:agent:pending_action:user:%d", userID)
}

func AgentSessionKey(userID uint) string {
	return fmt.Sprintf("dormgo:agent:session:user:%d", userID)
}

func ChatWithAgent(userID uint, message string) (*AgentChatResponse, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, errors.New("消息不能为空")
	}

	serviceURL := strings.TrimRight(viper.GetString("agent.ServiceURL"), "/")
	if serviceURL == "" {
		serviceURL = "http://127.0.0.1:8090"
	}

	reqBody := agentPythonRequest{
		UserID:        userID,
		Message:       message,
		InternalToken: AgentInternalToken(),
		GoBaseURL:     GoPublicBaseURL(),
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, serviceURL+"/chat", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Agent 服务不可用: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Agent 服务返回异常状态: %d", resp.StatusCode)
	}

	var agentResp AgentChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		return nil, err
	}
	return &agentResp, nil
}

func AgentInternalToken() string {
	token := viper.GetString("agent.InternalToken")
	if token == "" {
		return "dormgo-agent-local-token"
	}
	return token
}

func GoPublicBaseURL() string {
	baseURL := strings.TrimRight(viper.GetString("agent.GoBaseURL"), "/")
	if baseURL == "" {
		return "http://127.0.0.1:8080"
	}
	return baseURL
}

func SearchSignupPosts(keyword string, limit int, userID uint) ([]*model.ApiPostDetail, error) {
	return buildPostDetails(model.SearchSignupPosts, keyword, limit, userID)
}

func SearchPosts(keyword string, limit int, userID uint) ([]*model.ApiPostDetail, error) {
	return buildPostDetails(model.SearchNormalPosts, keyword, limit, userID)
}

func buildPostDetails(search func(string, int) ([]*model.DgPost, error), keyword string, limit int, userID uint) ([]*model.ApiPostDetail, error) {
	if limit <= 0 || limit > 10 {
		limit = 5
	}
	posts, err := search(keyword, limit)
	if err != nil {
		return nil, err
	}

	postDetails := make([]*model.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := model.GetUserById(int(post.PublisherId))
		if err != nil {
			continue
		}
		isSignedUp := false
		if userID > 0 {
			isSignedUp, _ = model.CheckIsSignedUp(post.ID, userID)
		}
		p := post
		postDetails = append(postDetails, &model.ApiPostDetail{
			PublisherName:   user.Username,
			PublisherAvator: user.Avatar,
			PublisherIntro:  user.Intro,
			IsSignedUp:      isSignedUp,
			DgPost:          p,
		})
	}
	return postDetails, nil
}

func PrepareSignupAction(userID uint, postID uint) (*PendingSignupAction, error) {
	post, err := model.GetPostDetail(int(postID))
	if err != nil {
		return nil, err
	}
	if post.Status != "normal" {
		return nil, errors.New("帖子当前不可报名")
	}
	if !post.IsLimited {
		return nil, errors.New("该帖子未开启报名")
	}

	action := PendingSignupAction{
		ActionID:  uuid.NewString(),
		Action:    "signup_post",
		UserID:    userID,
		PostID:    postID,
		Title:     post.Title,
		CreatedAt: time.Now(),
	}
	payload, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	if err := myredis.GetClient().Set(PendingActionKey(userID), string(payload), pendingActionTTL).Err(); err != nil {
		return nil, err
	}
	return &action, nil
}

func ConfirmPendingSignup(userID uint, actionID string) (*PendingSignupAction, error) {
	key := PendingActionKey(userID)
	raw, err := myredis.GetClient().Get(key).Result()
	if err != nil {
		return nil, errors.New("没有待确认的报名动作，或确认已过期")
	}

	var action PendingSignupAction
	if err := json.Unmarshal([]byte(raw), &action); err != nil {
		return nil, err
	}
	if action.UserID != userID || action.Action != "signup_post" {
		return nil, errors.New("待确认动作不匹配")
	}
	if actionID != "" && action.ActionID != actionID {
		return nil, errors.New("待确认动作已变化，请重新确认")
	}
	if err := model.SignupPost(action.PostID, userID); err != nil {
		return nil, err
	}
	_ = myredis.GetClient().Del(key).Err()
	return &action, nil
}

func GetAgentSession(userID uint) (*AgentSession, error) {
	raw, err := myredis.GetClient().Get(AgentSessionKey(userID)).Result()
	if err != nil {
		return &AgentSession{UserID: userID}, nil
	}
	var session AgentSession
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return &AgentSession{UserID: userID}, nil
	}
	return &session, nil
}

func SaveAgentSession(session *AgentSession) error {
	if session == nil || session.UserID == 0 {
		return errors.New("invalid agent session")
	}
	if len(session.LastCandidates) > 8 {
		session.LastCandidates = session.LastCandidates[:8]
	}
	if len(session.History) > 20 {
		session.History = session.History[len(session.History)-20:]
	}
	session.UpdatedAt = time.Now()
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return myredis.GetClient().Set(AgentSessionKey(session.UserID), string(payload), agentSessionTTL).Err()
}

func PrepareSignupFromCandidate(userID uint, index int) (*PendingSignupAction, error) {
	session, err := GetAgentSession(userID)
	if err != nil {
		return nil, err
	}
	if index <= 0 {
		index = 1
	}
	if index > len(session.LastCandidates) {
		return nil, errors.New("没有可引用的候选帖子，请先搜索帖子")
	}
	candidate := session.LastCandidates[index-1]
	if !candidate.IsLimited {
		return nil, errors.New("这个帖子没有开启报名，不能报名")
	}
	return PrepareSignupAction(userID, candidate.PostID)
}
