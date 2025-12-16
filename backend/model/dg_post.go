package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DgPost struct {
	ID           uint       `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	PublisherId  uint       `gorm:"column:publisherid" json:"publisherid" binding:"required"`
	AccepterId   uint       `gorm:"column:accepterid" json:"accepterid"`
	DormId       uint       `gorm:"column:dormid" json:"dormid" binding:"required"`
	Dorm         DgDorm     `gorm:"foreignkey:DormId;references:DormId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Title        string     `gorm:"column:title;type:varchar(255);" json:"title" binding:"required"`
	Content      string     `gorm:"column:content;type:text;" json:"content" binding:"required"`
	TypeId       uint       `gorm:"column:typeid" json:"typeid" binding:"required"`
	Type         DgType     `gorm:"foreignkey:TypeId;references:TypeId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Status       string     `gorm:"column:status;type:varchar(50);" json:"status"`
	IsPinned     bool       `gorm:"column:is_pinned;default:false" json:"is_pinned"` // 新增：置顶状态
	CommentCount uint       `gorm:"column:comment_count" json:"comment_count"`
	ViewCount    uint       `gorm:"column:view_count" json:"view_count"`
	LikeCount    uint       `gorm:"column:like_count" json:"like_count"`
	Images       []DgImages `gorm:"foreignKey:PostID;references:ID" json:"images"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

type ApiPostDetail struct {
	PublisherName   string `json:"publishername"`
	PublisherAvator string `json:"publisheravator"`
	*DgPost
}

// 发帖
func CreatePost(post *DgPost) (err error) {
	err = DB.Create(post).Error
	if err != nil {
		fmt.Println("create post error: ", err)
		return err
	}
	return nil
}

// 根据id获取帖子
func GetPostDetail(id int) (*DgPost, error) {
	var post DgPost
	// 增加 view_count
	DB.Model(&DgPost{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	err := DB.Preload("Images").Preload("Dorm").Preload("Type").Where("id=?", id).First(&post).Error
	if err != nil {
		fmt.Println("get post detail error: ", err)
		return nil, err
	}
	return &post, nil
}

// 根据宿舍楼id返回所有该宿舍楼的帖子
func GetPostByDorm(dormid int) ([]*DgPost, error) {
	posts := make([]*DgPost, 0)
	// 按照置顶优先，然后时间倒序排列
	err := DB.Where("dormid=?", dormid).
		Order("is_pinned desc, created_at desc").
		Find(&posts).Error
	if err != nil {
		fmt.Println("get post dorm error: ", err)
		return nil, err
	}
	return posts, nil
}

// 获取所有帖子
func GetPosts() ([]*DgPost, error) {
	posts := make([]*DgPost, 0)
	// 按照置顶优先，然后时间倒序排列
	err := DB.Preload("Images").Preload("Dorm").Preload("Type").
		Order("is_pinned desc, created_at desc").
		Find(&posts).Error
	if err != nil {
		fmt.Println("get posts error: ", err)
		return nil, err
	}
	return posts, err
}

// 点赞数加一
func AddLikeCount(postid uint) error {
	err := DB.Model(&DgPost{}).
		Where("id = ?", postid).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	if err != nil {
		fmt.Println("add like count error: ", err)
		return err
	}
	return nil
}

// 点赞数减一
func DelLikeCount(postid uint) error {
	err := DB.Model(&DgPost{}).
		Where("id = ?", postid).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	if err != nil {
		fmt.Println("del like count error: ", err)
		return err
	}
	return nil
}

// 评论数加一
func AddCommentCount(db *gorm.DB, postid uint) error {
	err := db.Model(&DgPost{}).
		Where("id = ?", postid).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	if err != nil {
		fmt.Println("add comment count error: ", err)
		return err
	}
	return nil
}
