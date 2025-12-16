package model

// SystemConfig 系统配置 (简单KV存储)
type SystemConfig struct {
	ID    uint   `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	Key   string `gorm:"type:varchar(100);uniqueIndex" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}
