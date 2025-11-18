package model

import (
	"fmt"
	"time"
)

type DgPost struct {
	ID          uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	PublisherId uint      `gorm:"column:publisherid" json:"publisherid" binding:"required"`
	AccepterId  uint      `gorm:"column:accepterid" json:"accepterid"`
	DormId      uint      `gorm:"column:dormid" json:"dormid" binding:"required"`
	Dorm        DgDorm    `gorm:"foreignkey:DormId;references:DormId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Title       string    `gorm:"column:title;type:varchar(255);" json:"title" binding:"required"`
	Content     string    `gorm:"column:content;type:text;" json:"content" binding:"required"`
	TypeId      uint      `gorm:"column:typeid" json:"typeid" binding:"required"`
	Type        DgType    `gorm:"foreignkey:TypeId;references:TypeId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Status      string    `gorm:"column:status;type:varchar(50);" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type ApiPostDetail struct {
	PublisherName string `json:"publishername"`
	*DgPost
}

// 发帖
func CreatePost(post *DgPost) (err error) {
	err = db.Create(post).Error
	if err != nil {
		fmt.Println("create post error: ", err)
		return err
	}
	return nil
}

// 根据id获取帖子
func GetPostDetail(id int) (*DgPost, error) {
	var post DgPost

	err := db.Preload("Dorm").Preload("Type").Where("id=?", id).First(&post).Error
	if err != nil {
		fmt.Println("get post detail error: ", err)
		return nil, err
	}
	return &post, nil
}

// 根据宿舍楼id返回所有该宿舍楼的帖子
func GetPostByDorm(dormid int) ([]DgPost, error) {
	var posts []DgPost
	err := db.Where("dormid=?", dormid).Find(&posts).Error
	if err != nil {
		fmt.Println("get post dorm error: ", err)
		return nil, err
	}
	return posts, nil
}
