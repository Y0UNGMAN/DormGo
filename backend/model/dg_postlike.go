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

// PostLike 点赞操作，使用事务确保点赞和通知的一致性
func PostLike(like *DgPostLike) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 在点赞表中增加记录
		if err := tx.Create(like).Error; err != nil {
			return err
		}

		// 2. 获取被点赞的帖子信息（为了拿到 PublisherId）
		var post DgPost
		// 只查询 publisherid 字段即可
		if err := tx.Select("publisherid").Where("id = ?", like.PostID).First(&post).Error; err != nil {
			return err
		}

		// 3. 给帖子发布者发送通知（如果点赞的不是自己的帖子）
		if post.PublisherId != like.LikerID {
			notification := DgNotification{
				ReceiverID: post.PublisherId,
				SenderID:   like.LikerID, // 点赞者即为发送通知的人
				Type:       "like",       // 通知的类型，前端可据此显示图标等
				PostID:     like.PostID,
				Content:    "赞了你的帖子", // 通知的具体内容
			}
			if err := tx.Create(&notification).Error; err != nil {
				return err
			}
		}

		return nil
	})
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
