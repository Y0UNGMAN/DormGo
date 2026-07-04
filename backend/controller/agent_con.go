package controller

import (
	"net/http"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/gin-gonic/gin"
)

type AgentChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type internalSearchPostsRequest struct {
	UserID  uint   `json:"user_id"`
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
}

type internalPrepareSignupRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	PostID uint `json:"post_id" binding:"required"`
}

type internalConfirmSignupRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	ActionID string `json:"action_id"`
}

type internalSessionRequest struct {
	UserID         uint                      `json:"user_id" binding:"required"`
	LastCandidates []logic.AgentCandidateRef `json:"last_candidates"`
	History        []logic.AgentHistoryItem  `json:"history"`
	SelectedPost   map[string]interface{}    `json:"selected_post"`
	PendingAction  map[string]interface{}    `json:"pending_action"`
}

type internalPrepareFromCandidateRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Index  int  `json:"index"`
}

type internalPostDetailRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	PostID uint `json:"post_id" binding:"required"`
}

func AgentChat(c *gin.Context) {
	var req AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "消息不能为空"})
		return
	}

	userID, err := GetCurrentUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
		return
	}

	resp, err := logic.ChatWithAgent(userID, req.Message)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": resp})
}

func AgentChatStream(c *gin.Context) {
	var req AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "消息不能为空"})
		return
	}

	userID, err := GetCurrentUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
		return
	}

	if err := logic.ChatWithAgentStream(userID, req.Message, c.Writer); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
}

func InternalSearchPosts(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalSearchPostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	posts, err := logic.SearchPosts(req.Keyword, req.Limit, req.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": posts})
}

func InternalSearchSignupPosts(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalSearchPostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	posts, err := logic.SearchSignupPosts(req.Keyword, req.Limit, req.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": posts})
}

func InternalPrepareSignup(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalPrepareSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	action, err := logic.PrepareSignupAction(req.UserID, req.PostID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": action})
}

func InternalPrepareSignupFromCandidate(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalPrepareFromCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	action, err := logic.PrepareSignupFromCandidate(req.UserID, req.Index)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": action})
}

func InternalPostDetail(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalPostDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	post, err := logic.GetAgentPostDetail(req.UserID, req.PostID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": post})
}

func InternalPostDetailFromCandidate(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalPrepareFromCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	post, err := logic.GetPostDetailFromCandidate(req.UserID, req.Index)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": post})
}

func InternalConfirmSignup(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalConfirmSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	action, err := logic.ConfirmPendingSignup(req.UserID, req.ActionID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "报名成功", "data": action})
}

func InternalGetAgentSession(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	session, err := logic.GetAgentSession(req.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": session})
}

func InternalSaveAgentSession(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req internalSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	session := &logic.AgentSession{
		UserID:         req.UserID,
		LastCandidates: req.LastCandidates,
		History:        req.History,
		SelectedPost:   req.SelectedPost,
		PendingAction:  req.PendingAction,
	}
	if err := logic.SaveAgentSession(session); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": session})
}

func InternalAllPosts(c *gin.Context) {
	if !checkAgentInternalToken(c) {
		return
	}
	var req struct {
		Limit int `json:"limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Limit = 100
	}
	posts, err := logic.GetAllPostsForIndex(req.Limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": posts})
}

func checkAgentInternalToken(c *gin.Context) bool {
	if c.GetHeader("X-Agent-Token") != logic.AgentInternalToken() {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid agent token"})
		return false
	}
	return true
}
