package logic

import (
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
)

// AdminLogin 管理员登录业务逻辑
func AdminLogin(p *model.LoginRequest) (string, *model.DgAdmin, error) {
	// 1. 调用 Model 层验证账号密码
	// 注意：这里调用的是我们在第一步 model/dg_admin.go 中新增的 LoginAdmin 方法
	admin, err := model.LoginAdmin(p.Username, p.Password)
	if err != nil {
		return "", nil, err
	}

	// 2. 生成 Token
	// 复用 utils 包中的 GenToken 方法，传入管理员ID和用户名
	token, err := utils.GenToken(admin.ID, admin.Username)
	if err != nil {
		return "", nil, err
	}

	return token, admin, nil
}
