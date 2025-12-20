package logic

import (
	"errors"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
)

// AdminLogin 管理员登录业务逻辑
func AdminLogin(p *model.LoginRequest) (string, *model.DgAdmin, error) {
	// 1. 调用 Model 层验证账号密码
	admin, err := model.LoginAdmin(p.Username, p.Password)
	if err != nil {
		return "", nil, err
	}

	// 2. 生成 Token
	token, err := utils.GenToken(admin.ID, admin.Username)
	if err != nil {
		return "", nil, err
	}

	return token, admin, nil
}

// BanUser 封禁/解封用户逻辑
// 修改为状态反转逻辑：正常(1) -> 封禁(2)；封禁(2) -> 正常(1)
func BanUser(userId int64) error {
	// 1. 检查用户是否存在
	var user model.DgUser
	if err := model.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 2. 状态切换逻辑
	// 假设 1=正常, 2=封禁
	if user.Status == 1 {
		user.Status = 2 // 正常 -> 封禁
	} else {
		user.Status = 1 // 封禁 -> 正常 (解封)
	}

	// 3. 保存更改
	if err := model.DB.Save(&user).Error; err != nil {
		return errors.New("更新用户状态失败")
	}

	return nil
}
