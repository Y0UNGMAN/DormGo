package logic

import "github.com/Y0UNGMAN/DormGo/backend/model"

func GetDorms() ([]*model.DgDorm, error) {
	return model.GetDorms()
}
