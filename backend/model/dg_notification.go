package model

import "time"

type DgNotification struct {
	ID         uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	ReceiverID uint      `gorm:"column:receiver_id;index" json:"receiver_id"`
	SenderID   uint      `gorm:"column:sender_id" json:"sender_id"`
	Sender     DgUser    `gorm:"foreignKey:SenderID" json:"sender"`
	Type       string    `gorm:"column:type;type:varchar(20)" json:"type"`
	Content    string    `gorm:"column:content;type:varchar(255)" json:"content"`
	PostID     uint      `gorm:"column:post_id" json:"post_id"`
	Post       DgPost    `gorm:"foreignKey:PostID" json:"post"`
	IsRead     bool      `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func GetNotifacation(userId uint) ([]DgNotification, error) {
	userid, err := UserIdToId(userId)
	if err != nil {
		return nil, err
	}
	var notifs []DgNotification
	err = DB.Preload("Sender").Preload("Post").
		Where("receiver_id = ?", userid).
		Order("created_at desc").
		Find(&notifs).Error

	if err != nil {
		return nil, err
	}
	return notifs, nil
}

func GetUnreadCount(userId uint) (int64, error) {
	userid, err := UserIdToId(userId)
	if err != nil {
		return 0, err
	}
	var count int64
	err = DB.Model(&DgNotification{}).
		Where("receiver_id = ? AND is_read = ?", userid, false).
		Count(&count).Error

	return count, err
}

func ReadAllNotification(userId uint) error {
	userid, err := UserIdToId(userId)
	if err != nil {
		return err
	}
	err = DB.Model(&DgNotification{}).
		Where("receiver_id = ?", userid).
		Update("is_read", true).Error
	if err != nil {
		return err
	}
	return nil
}
