package model

import (
	"time"

	"gorm.io/gorm"
)

type DgComment struct {
	ID          uint    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	PostID      uint    `gorm:"column:post_id" json:"post_id"`
	Post        *DgPost `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"post,omitempty"`
	CommenterID uint    `gorm:"column:commenter_id" json:"commenter_id"`
	Commenter   *DgUser `gorm:"foreignKey:CommenterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"commenter,omitempty"`
	Content     string  `gorm:"column:content;type:text;" json:"content"`
	// ParentID: 直系父级 ID (即你点了谁的“回复”按钮)
	// 顶级评论为 0 (或 null)
	ParentID uint `gorm:"column:parent_id;index;default:0" json:"parent_id"`
	// RootID: 根评论 ID (所属楼层的 ID)
	// 如果是顶级评论，RootID = 0；如果是回复(不管几层)，RootID = 顶级评论的ID
	RootID uint `gorm:"column:root_id;index;default:0" json:"root_id"`
	// ReplyToUID: 被回复的人的 ID (用于前端展示 "回复 @张三")
	// 如果是顶级评论，此值为 0
	ReplyToUID uint `gorm:"column:reply_to_uid;default:0" json:"reply_to_uid"`
	// 关联获取被回复人的昵称/头像
	ReplyToUser *DgUser `gorm:"foreignKey:ReplyToUID;references:ID;" json:"reply_to_user,omitempty"`
	// ReplyCount: 子回复总数 (仅在 RootID=0 的顶级评论中维护此数据)
	// 作用：前端显示 "展开 15 条回复"
	ReplyCount uint `gorm:"column:reply_count;default:0" json:"reply_count"`
	// SubComments: 用于在查询时，把子回复挂载到顶级评论下面
	SubComments []DgComment `gorm:"-" json:"sub_comments"`
	CreatedAt   time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

func CreateComment(db *gorm.DB, comment *DgComment) error {
	return db.Create(comment).Error
}

func GetCommentById(db *gorm.DB, id uint) (*DgComment, error) {
	var comment DgComment
	err := db.Where("id = ?", id).First(&comment).Error
	return &comment, err
}

func IncrementReplyCount(db *gorm.DB, rootID uint) error {
	return db.Model(&DgComment{}).
		Where("id = ?", rootID).
		UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error
}

func GetRootCommentList(db *gorm.DB, postID uint, offset, limmit int) ([]*DgComment, int64, error) {
	var comments []*DgComment
	var total int64
	// 1. 构建基础查询
	query := db.Model(&DgComment{}).
		Preload("Commenter").
		Where("post_id = ? AND root_id = ?", postID, 0)
	// 2. 查总数 (用于前端分页)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. 查列表 (按时间倒序)
	err := query.Order("created_at desc").
		Offset(offset).
		Limit(limmit).
		Find(&comments).Error

	return comments, total, err

}

// GetPreviewReplies 获取某条评论下的“前N条”子回复 (用于预览)
func GetPreviewReplies(db *gorm.DB, rootID uint, limit int) ([]DgComment, error) {
	var replies []DgComment

	err := db.Model(&DgComment{}).
		Preload("Commenter").
		Preload("ReplyToUser").
		Where("root_id = ?", rootID).
		Order("created_at desc").
		Limit(limit).
		Find(&replies).Error

	return replies, err
}
