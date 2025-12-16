package controller

import (
	"net/http"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/middleware"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// 注册
func Register(c *gin.Context) {
	request := new(model.RegisterRequest)
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求参数有误",
		})
		return
	}

	if len(request.Username) == 0 || len(request.Password) == 0 || request.RePassword != request.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户名、密码为空，或两次密码不一样",
		})
		return
	}

	err = logic.SignUp(request)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Sign Up Success",
	})
}

// 登录
func Login(c *gin.Context) {
	p := new(model.LoginRequest)
	err := c.ShouldBindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求参数有误",
		})
		return
	}

	// 修改：接收 user 对象
	token, user, err := logic.LoginIn(p)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 修改：返回 user_id 等详细信息
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Login Success",
		"token":   token,
		"data": gin.H{
			"id":         user.ID,     // 关键：返回数据库主键ID
			"user_id":    user.UserID, // 雪花ID
			"username":   user.Username,
			"student_id": user.StudentId,
			"avatar":     user.Avatar,
			"dorm_id":    user.DormId,
		},
	})
}

// 获取当前登录用户ID
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "用户未登录",
		})
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	return
}

// AdminLogin 管理员登录接口
func AdminLogin(c *gin.Context) {
	p := new(model.LoginRequest)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求参数有误",
		})
		return
	}

	token, admin, err := logic.AdminLogin(p)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Admin Login Success",
		"token":   token,
		"data": gin.H{
			"avatar":   admin.Avatar,
			"nickname": admin.Username,
		},
	})
}
