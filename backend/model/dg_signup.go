package model

import "time"

type DgSignup struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	PostID    uint      `gorm:"column:post_id;index" json:"post_id"`
	UserID    uint      `gorm:"column:user_id;index" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
