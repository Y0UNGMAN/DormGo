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
		UserID:      uint(userID),
		Username:    p.Username,
		StudentId:   p.StudentId,
		Password:    p.Password,
		DormId:      p.DormId,
		Avatar:      "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
		Status:      1,   // 【新增】注册时默认为正常状态
		CreditScore: 100, // 【新增】初始信用分
	}

	// 3 保存进数据库
	err = model.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

func LoginIn(p *model.LoginRequest) (token string, user *model.DgUser, err error) {
	user = new(model.DgUser)
	// model.Login 内部已经包含了状态检测逻辑
	user, err = model.Login(p.Username, p.Password)
	if err != nil {
		return "", nil, err
	}

	token, err = utils.GenToken(user.UserID, user.Username)

	return token, user, err
}

func GetUserById(id int) (*model.DgUser, error) {
	return model.GetUserById(id)
}
