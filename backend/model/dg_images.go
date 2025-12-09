package model

import "time"

type DgImages struct {
	ID        uint      `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	PostID    uint      `gorm:"column:post_id;" json:"post_id"`
	ImageURL  string    `gorm:"column:image_url;type:varchar(512)" json:"image_url"`
	Order     uint      `gorm:"column:image_order;" json:"order"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
