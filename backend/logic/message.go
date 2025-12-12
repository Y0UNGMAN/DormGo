package logic

import "github.com/Y0UNGMAN/DormGo/backend/model"

func SendMessage(req model.SendMessageReq, userId uint) (model.DgMessage, error) {
	return model.StoreMessage(req, userId)
}

func GetMessage(targetID uint, userId uint) ([]model.DgMessage, error) {
	return model.GetMessage(uint(targetID), userId)
}
