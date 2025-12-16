package model

import "time"

// DgSensitiveWord 敏感词表
type DgSensitiveWord struct {
	ID        uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	Word      string    `gorm:"type:varchar(100);uniqueIndex" json:"word"`
	Category  string    `gorm:"type:varchar(50)" json:"category"`
	CreatedAt time.Time `json:"created_at"`
}
