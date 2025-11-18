package model

import "fmt"

type DgUser struct {
	ID        uint   `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	StudentId string `gorm:"column:studentid;type:varchar(20);"json:"studentid"`
	Username  string `gorm:"column:username;type:varchar(50);" json:"username"`
	Password  string `gorm:"column:password;type:varchar(255);" json:"password"`
	DormId    uint   `gorm:"column:dormid" json:"dormid"`
	Dorm      DgDorm `gorm:"foreignkey:DormId;references:DormId" `
}

// 根据id获取用户信息
func GetUserById(id int) (*DgUser, error) {
	var user DgUser
	err := db.Where("id=?", id).First(&user).Error
	if err != nil {
		fmt.Println("get user detail error: ", err)
		return nil, err
	}
	return &user, nil
}
