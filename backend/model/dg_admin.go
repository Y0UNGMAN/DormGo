package model

import (
	"errors"

	"gorm.io/gorm"
)

type DgAdmin struct {
	ID       uint   `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	Username string `gorm:"column:username;type:varchar(50);" json:"username"`
	Password string `gorm:"column:password;type:varchar(255);" json:"password"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);" json:"avatar_url"`
}

// LoginAdmin 管理员登录验证 (明文比对版)
func LoginAdmin(username, password string) (*DgAdmin, error) {
	var admin DgAdmin
	// 1. 根据用户名查询
	err := DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("管理员不存在")
		}
		return nil, err
	}

	// 2. 直接比对明文密码 (不再使用 bcrypt 加密)
	if admin.Password != password {
		return nil, errors.New("密码错误")
	}

	return &admin, nil
}
