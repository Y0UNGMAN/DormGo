package logic

import (
	"errors"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/utils" // 引入了 utils 包
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
		UserID:   uint(userID), // 注意：确保数据库 UserID 是 uint64 或足够大
		Username: p.Username,
		Password: p.Password,
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
