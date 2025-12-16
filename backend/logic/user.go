package logic

import (
	"errors"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
)

func SignUp(p *model.RegisterRequest) error {
	// 1 判断用户不存在
	exist, err := model.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("用户已存在")
	}

	// 2 生成UID
	userID := utils.GetID()

	var user = &model.DgUser{
		UserID:    uint(userID),
		Username:  p.Username,
		StudentId: p.StudentId, // 确保注册时写入学号
		Password:  p.Password,
		DormId:    p.DormId,
		Avatar:    "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
	}

	// 3 保存进数据库
	err = model.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

// LoginIn 修改返回值，增加 *model.DgUser
func LoginIn(p *model.LoginRequest) (token string, user *model.DgUser, err error) {
	user = new(model.DgUser)
	user, err = model.Login(p.Username, p.Password)
	if err != nil {
		return "", nil, err
	}

	token, err = utils.GenToken(user.UserID, user.Username)

	return token, user, err
}
