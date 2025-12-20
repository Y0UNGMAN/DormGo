package model

import "fmt"

type DgType struct {
	TypeId   uint   `gorm:"column:typeid;primarykey;AUTO_INCREMENT" json:"typeid"`
	TypeName string `gorm:"column:typename;type:varchar(255);" json:"typename"`
}

// 获取帖子类型列表
func GetPostType() (posttypelist []*DgType, err error) {
	err = DB.Find(&posttypelist).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return
}

func GetPostTypeById(id int) (posttype *DgType, err error) {
	err = DB.Where("TypeId=?", id).Find(&posttype).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return
}
