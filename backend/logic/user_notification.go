package logic

import "github.com/Y0UNGMAN/DormGo/backend/model"

func GetNotifacation(userId uint) ([]model.DgNotification, error) {
	return model.GetNotifacation(userId)
}

func GetUnreadCount(userId uint) (int64, error) {
	return model.GetUnreadCount(userId)
}

func ReadAllNotification(userId uint) error {
	return model.ReadAllNotification(userId)
}
