package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetUserProfile 获取个人资料
func GetUserProfile(c *gin.Context) {
	uid, _ := GetCurrentUser(c) // 假设 user.go 里已定义此辅助函数
	// 将 UserID 转换为数据库主键 ID
	id, err := model.UserIdToId(uint(uid))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "用户不存在"})
		return
	}
	user, err := model.GetUserById(int(id))
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
			"intro":        user.Intro,
			"avatar":       user.Avatar,
			"dorm_id":      user.DormId,
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
		DormId   uint   `json:"dorm_id"`
		Bio      string `json:"bio"`
	}
	// 为了调试 binding 问题，先读取原始 body 并尝试解析
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "读取请求体失败"})
		return
	}
	// 尝试解析
	if err := json.Unmarshal(raw, &req); err != nil {
		// 返回解析错误和原始 body 帮助调试
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败", "error": err.Error(), "raw": string(raw)})
		return
	}
	uid, _ := GetCurrentUser(c)
	// 将 UserID 转换为数据库主键 ID
	id, err := model.UserIdToId(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "用户ID转换失败"})
		return
	}
	updates := map[string]interface{}{
		"username": req.Nickname,
		"avatar":   req.Avatar,
		"dormid":   req.DormId,
		"intro":    req.Bio,
	}
	if err := model.DB.Model(&model.DgUser{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}
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
	// 将 UserID 转换为数据库主键 ID
	id, err := model.UserIdToId(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "用户ID转换失败"})
		return
	}
	user, _ := model.GetUserById(int(id))

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

// GetUserStats 获取用户统计数据
func GetUserStats(c *gin.Context) {
	uid, err := GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	// 将 UserID 转换为数据库主键 ID
	id, err := model.UserIdToId(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "用户ID转换失败"})
		return
	}

	// 获取用户帖子的总点赞数（按实际列名 publisherid）
	var totalLikes int64
	model.DB.Model(&model.DgPostLike{}).
		Joins("JOIN dg_posts ON dg_post_likes.postid = dg_posts.id").
		Where("dg_posts.publisherid = ?", id).
		Count(&totalLikes)

	// 获取用户个性签名
	var user model.DgUser
	if err := model.DB.Preload("Dorm").First(&user, id).Error; err != nil {
		// 如果没找到用户，返回空统计
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"total_likes": totalLikes, "bio": ""}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_likes": totalLikes,
			"bio":         user.Intro,
			"dorm_name":   user.Dorm.DormName,
			"dorm_id":     user.Dorm.DormId,
		},
	})
}

// GetUserCoins 获取寝友币数据
func GetUserCoins(c *gin.Context) {
	// 模拟数据，实际项目中应该查询数据库
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"balance": 128,
			"records": []gin.H{
				{
					"id":         1,
					"type":       "income",
					"title":      "新用户注册奖励",
					"amount":     50,
					"created_at": time.Now().AddDate(0, 0, -7),
				},
				{
					"id":         2,
					"type":       "income",
					"title":      "发布帖子奖励",
					"amount":     5,
					"created_at": time.Now(),
				},
			},
		},
	})
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "上传文件失败"})
		return
	}

	// 限制文件大小为 2MB
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "文件大小不能超过2MB"})
		return
	}

	url, err := utils.UploadFile(file, "avatars")
	if err != nil {
		// 尝试回退到本地保存
		localURL, saveErr := utils.SaveFileLocal(file, "avatars")
		if saveErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "上传失败"})
			return
		}
		url = localURL
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"url": url,
		},
	})
}

// GetDormList 获取宿舍列表（用于用户选择）
func GetDormList(c *gin.Context) {
	var dorms []model.DgDorm
	model.DB.Find(&dorms)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": dorms,
	})
}
