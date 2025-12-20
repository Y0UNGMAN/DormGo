package logic

import (
	"log"

	"github.com/Y0UNGMAN/DormGo/backend/model"
)

func SendMessage(req model.SendMessageReq, userId uint) (model.DgMessage, error) {
	return model.StoreMessage(req, userId)
}

func GetMessage(targetID uint, userId uint) ([]model.DgMessage, error) {
	return model.GetMessage(targetID, userId)
}

func GetConversationList(myUserId uint) ([]model.ConversationItem, error) {
	// 1. 从 Model 层获取聚合数据 (谁、哪条最新、几条未读)
	metaList, err := model.GetConversationMetaList(myUserId)
	if err != nil {
		return nil, err
	}

	var list []model.ConversationItem

	// 2. 遍历聚合数据，去 Model 层查详情，组装最终结果
	for _, r := range metaList {
		// 查最后一条消息内容
		msg, _ := model.GetMessageByID(r.MaxID)

		// 查对方用户信息
		user, err := model.GetUserById(int(r.OtherID))
		var targetName = "未知用户"
		var targetAvatar = "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" // 默认头像
		if err == nil {
			targetName = user.Username
			targetAvatar = user.Avatar
		} else {
			log.Printf("警告: 找不到用户ID %d 的信息 (可能已注销), 使用默认值", r.OtherID)
		}
		item := model.ConversationItem{
			TargetID:    r.OtherID,
			Username:    targetName,
			AvatarURL:   targetAvatar, // 注意大小写匹配你的 User 结构体
			LastContent: msg.Content,
			LastTime:    msg.CreatedAt,
			UnreadCount: r.Unread,
		}
		list = append(list, item)
	}

	return list, nil
}

// MarkAsRead 标记已读逻辑
func MarkAsRead(myUserId, targetUserId uint) error {
	return model.UpdateMessagesReadStatus(myUserId, targetUserId)
}

// GetUnreadCount 获取总未读数逻辑
func GetUnreadMessageCount(myUserId uint) (int64, error) {
	return model.CountUnreadMessages(myUserId)
}
