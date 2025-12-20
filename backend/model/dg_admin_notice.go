package model

import "time"

// DgNotice 通知公告表
type DgNotice struct {
	ID        uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Target    string    `gorm:"type:varchar(50)" json:"target"` // all, dorm, user
	TargetID  string    `gorm:"type:varchar(255)" json:"target_id"`
	IsPinned  bool      `gorm:"default:false" json:"is_pinned"`
	IsUrgent  bool      `gorm:"default:false" json:"is_urgent"`
	Publisher string    `gorm:"type:varchar(50)" json:"publisher_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
