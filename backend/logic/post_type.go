package logic

import (
	"github.com/Y0UNGMAN/DormGo/backend/model"
)

func GetPostType() (posttypelist []*model.DgType, err error) {
	//查数据库，找到所有的类型并返回
	return model.GetPostType()
}

func GetPostTypeById(id int) (posttype *model.DgType, err error) {
	return model.GetPostTypeById(id)
}
