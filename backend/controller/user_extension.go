package controller

import (
	"net/http"
	"time"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetUserProfile 获取个人资料
func GetUserProfile(c *gin.Context) {
	uid, _ := GetCurrentUser(c) // 假设 user.go 里已定义此辅助函数
	user, err := model.GetUserById(int(uid))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "用户不存在"})
		return
	}
	// 关联查询宿舍名
	var dormName string
	model.DB.Model(&model.DgDorm{}).Select("dormname").Where("dormid = ?", user.DormId).Scan(&dormName)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":           user.ID,
			"student_id":   user.StudentId,
			"nickname":     user.Username,
			"phone_number": "", // 数据库需加字段，这里暂空
			"intro":        "", // 数据库需加字段，这里暂空
			"avatar":       user.Avatar,
			"dorm_name":    dormName,
		},
	})
}

// UpdateUserProfile 更新个人资料
func UpdateUserProfile(c *gin.Context) {
	// 简化逻辑：只演示更新头像和昵称
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	uid, _ := GetCurrentUser(c)
	model.DB.Model(&model.DgUser{}).Where("id = ?", uid).Updates(map[string]interface{}{
		"username": req.Nickname,
		"avatar":   req.Avatar,
	})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// GetLoginInfo 获取登录信息
func GetLoginInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"last_login_time": time.Now(), // 简化处理，直接返回当前时间
			"ip":              c.ClientIP(),
		},
	})
}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	uid, _ := GetCurrentUser(c)
	user, _ := model.GetUserById(int(uid))

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "原密码错误"})
		return
	}

	// 加密新密码
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	model.DB.Model(&user).Update("password", string(hash))

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "密码修改成功"})
}
