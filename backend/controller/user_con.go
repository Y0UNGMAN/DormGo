package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	//获取参数及参数校验
	request := new(model.RegisterRequest)
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求参数有误",
		})
		return
	}
	//手动对参数进行业务规则的校验
	if len(request.Username) == 0 || len(request.Password) == 0 || request.RePassword != request.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "用户名、密码为空，或两次密码不一样",
		})
	}

	err = logic.SignUp(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
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
	user := new(model.DgUser)
	token, user, err := logic.LoginIn(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Login Success",
		"token":   token,
		"user":    user,
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

func GetUserInfoPublic(c *gin.Context) {
	idStr := c.Query("user_id")
	id, _ := strconv.Atoi(idStr)

	user, err := logic.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Success",
		"data": gin.H{
			"id":        user.UserID,
			"username":  user.Username,
			"avatarurl": user.Avatar,
			"dorm":      user.Dorm,
		},
	})
}

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

func GetCurrentUserId(c *gin.Context) (uint, error) {
	var userId uint = 0
	if v, exists := c.Get(middleware.CtxUserIDKey); exists {
		if uid, ok := v.(uint); ok {
			userId = uid
		} else if uidInt, ok := v.(int); ok {
			userId = uint(uidInt)
		} else {
			fmt.Printf("UserID type mismatch in Context. Got: %T\n", v)
		}
	} else {
		fmt.Println("UserID not found in context")
		return userId, errors.New("UserID not found in context")
	}
	return userId, nil
}
