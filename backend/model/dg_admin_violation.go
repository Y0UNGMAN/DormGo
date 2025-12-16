package model

import "time"

// DgViolation 违规记录表
type DgViolation struct {
	ID            uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	User          DgUser    `gorm:"foreignKey:UserID" json:"user"`
	Type          string    `gorm:"type:varchar(50)" json:"type"` // content, behavior
	Reason        string    `gorm:"type:varchar(255)" json:"reason"`
	Status        string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, processed
	Punishment    string    `gorm:"type:varchar(255)" json:"punishment"`
	Notes         string    `gorm:"type:text" json:"notes"`
	ViolationTime time.Time `json:"violation_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
