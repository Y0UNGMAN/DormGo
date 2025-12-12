package model

import (
	"time"
)

type DgMessage struct {
	ID         uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	SenderID   uint      `gorm:"column:sender_id;index" json:"sender_id"`     // 发送者
	ReceiverID uint      `gorm:"column:receiver_id;index" json:"receiver_id"` // 接收者
	Content    string    `gorm:"column:content;type:text" json:"content"`     // 聊天内容
	IsRead     bool      `gorm:"column:is_read;default:false" json:"is_read"` // 是否已读
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

type SendMessageReq struct {
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

func StoreMessage(req SendMessageReq, userId uint) (DgMessage, error) {
	userid, err := UserIdToId(userId)
	if err != nil {
		return DgMessage{}, err
	}
	msg := DgMessage{
		SenderID:   userid,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}
	if err := DB.Create(&msg).Error; err != nil {
		return msg, err
	}
	return msg, nil
}

func GetMessage(receiverID uint, userId uint) ([]DgMessage, error) {
	userid, err := UserIdToId(userId)
	if err != nil {
		return []DgMessage{}, err
	}
	var msgs []DgMessage
	err = DB.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
						userid, receiverID, receiverID, userid).
		Order("created_at asc"). // 按时间正序
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
