package model

type DgUser struct {
	ID        uint   `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	StudentId string `gorm:"column:studentid;type:varchar(20);"json:"studentid"`
	Username  string `gorm:"column:username;type:varchar(50);" json:"username"`
	Password  string `gorm:"column:password;type:varchar(255);" json:"password"`
	DormId    uint   `gorm:"column:dormid" json:"dormid"`
	Dorm      DgDorm `gorm:"foreignkey:DormId;references:DormId" `
}
