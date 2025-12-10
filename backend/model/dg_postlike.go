package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type DgPostLike struct {
	ID        uint      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	PostID    uint      `gorm:"column:postid" json:"postid" binding:"required"`
	Post      DgPost    `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"post"`
	LikerID   uint      `gorm:"column:likerid" json:"likerid" binding:"required"`
	Liker     DgUser    `gorm:"foreignKey:LikerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"liker"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func PostLike(like *DgPostLike) error {
	err := DB.Create(like).Error
	if err != nil {
		return err
	}
	return nil
}

func CancelPostLike(postid uint, likerid uint) error {
	err := DB.Where("postid = ? AND likerid = ?", postid, likerid).Delete(&DgPostLike{}).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckIsLiked(postid uint, userid uint) (bool, error) {
	id, err := UserIdToId(userid)
	if err != nil {
		return false, err // 转换用户ID出错，直接返回
	}
	likerid := id
	var likeRecord DgPostLike
	err = DB.Where("postid = ? AND likerid = ?", postid, likerid).First(&likeRecord).Error
	// 检查查询结果
	if err != nil {
		// 如果是记录未找到错误，说明未点赞，返回 false
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // 未找到，但不是程序错误，返回 false, nil
		}
		// 如果是其他数据库错误，返回 error
		return false, err
	}
	// 如果 err 为 nil，说明找到了记录，表示已点赞，返回 true
	return true, nil
}
