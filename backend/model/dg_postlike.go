package model

import "time"

type DgPostLike struct {
	ID        uint      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	PostID    uint      `gorm:"column:postid" json:"commentid"`
	Post      DgPost    `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"post"`
	LikerID   uint      `gorm:"column:likerid" json:"likerid"`
	Liker     DgUser    `gorm:"foreignKey:LikerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"liker"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
