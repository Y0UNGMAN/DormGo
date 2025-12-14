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

type ConversationItem struct {
	TargetID    uint      `json:"target_id"`
	Username    string    `json:"username"`
	AvatarURL   string    `json:"avatar_url"`
	LastContent string    `json:"last_content"`
	LastTime    time.Time `json:"last_time"`
	UnreadCount int       `json:"unread_count"`
}

type ConversationMeta struct {
	OtherID uint
	MaxID   uint
	Unread  int
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
	err = DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userid, receiverID, receiverID, userid).
		Order("created_at asc"). // 按时间正序
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// 获取会话列表
func GetConversationMetaList(myUserId uint) ([]ConversationMeta, error) {
	userid, err := UserIdToId(myUserId)
	if err != nil {
		return []ConversationMeta{}, err
	}
	query := `
        SELECT 
            CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as other_id,
            MAX(id) as max_id,
            SUM(CASE WHEN receiver_id = ? AND is_read = 0 THEN 1 ELSE 0 END) as unread
        FROM dg_messages 
        WHERE sender_id = ? OR receiver_id = ?
        GROUP BY other_id
        ORDER BY max_id DESC
    `
	var results []ConversationMeta
	err = DB.Raw(query, userid, userid, userid, userid).Scan(&results).Error
	return results, err
}

// 根据ID获取单条消息
func GetMessageByID(msgID uint) (DgMessage, error) {
	var msg DgMessage
	err := DB.First(&msg, msgID).Error
	return msg, err
}

// 标记消息为已读
func UpdateMessagesReadStatus(myUserId uint, targetUserId uint) error {
	userid, err := UserIdToId(myUserId)
	if err != nil {
		return err
	}
	return DB.Model(&DgMessage{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", targetUserId, userid, false).
		Update("is_read", true).Error
}

func CountUnreadMessages(myUserId uint) (int64, error) {
	userid, err := UserIdToId(myUserId)
	if err != nil {
		return 0, err
	}
	var count int64
	err = DB.Model(&DgMessage{}).
		Where("receiver_id = ? AND is_read = ?", userid, false).
		Count(&count).Error
	return count, err
}
