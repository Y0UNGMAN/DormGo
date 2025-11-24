package model

import "time"

type DgComment struct {
	ID          uint      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	PostID      uint      `gorm:"column:post_id" json:"post_id"`
	Post        DgPost    `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommenterID uint      `gorm:"column:comment_id" json:"comment_id"`
	Commenter   DgUser    `gorm:"foreignKey:CommenterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Content     string    `gorm:"column:content;type:text;" json:"content"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}
