package model

type DgAdmin struct {
	ID       uint   `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	Username string `gorm:"column:username;type:varchar(50);" json:"username"`
	Password string `gorm:"column:password;type:varchar(255);" json:"password"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);" json:"avatar_url"`
}
