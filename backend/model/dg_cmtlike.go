package model

import "time"

type DgCmtLike struct {
	ID        uint      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CommentID uint      `gorm:"column:commentid" json:"commentid"`
	Comment   DgComment `gorm:"foreignKey:CommentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comment"`
	LikerID   uint      `gorm:"column:likerid" json:"likerid"`
	Liker     DgUser    `gorm:"foreignKey:LikerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"liker"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
