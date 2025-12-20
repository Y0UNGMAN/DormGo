package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DgPost struct {
	ID                uint       `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	PublisherId       uint       `gorm:"column:publisherid" json:"publisherid" binding:"required"`
	AccepterId        uint       `gorm:"column:accepterid" json:"accepterid"`
	DormId            uint       `gorm:"column:dormid" json:"dormid" binding:"required"`
	Dorm              DgDorm     `gorm:"foreignkey:DormId;references:DormId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Title             string     `gorm:"column:title;type:varchar(255);" json:"title" binding:"required"`
	Content           string     `gorm:"column:content;type:text;" json:"content" binding:"required"`
	TypeId            uint       `gorm:"column:typeid" json:"typeid" binding:"required"`
	Type              DgType     `gorm:"foreignkey:TypeId;references:TypeId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" `
	Status            string     `gorm:"column:status;type:varchar(50);" json:"status"`
	CommentCount      uint       `gorm:"column:comment_count" json:"comment_count"`
	ViewCount         uint       `gorm:"column:view_count" json:"view_count"`
	LikeCount         uint       `gorm:"column:like_count" json:"like_count"`
	Images            []DgImages `gorm:"foreignKey:PostID;references:ID" json:"images"`
	IsLimited         bool       `gorm:"column:is_limited" json:"is_limited"`
	IsPinned          bool       `gorm:"column:is_pinned;default:false" json:"is_pinned"`
	Deadline          *time.Time `gorm:"column:deadline" json:"deadline"`
	MaxEnrollment     int        `gorm:"column:max_enrollment" json:"max_enrollment"`
	CurrentEnrollment int        `gorm:"column:current_enrollment" json:"current_enrollment"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

type ApiPostDetail struct {
	PublisherName   string `json:"publishername"`
	PublisherAvator string `json:"publisheravator"`
	PublisherIntro  string `json:"publisherintro"`
	IsSignedUp      bool   `json:"is_signed_up"`
	IsFavorited     bool   `json:"is_favorited"`
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
	err := DB.Where("dormid=?", dormid).Find(&posts).Error
	if err != nil {
		fmt.Println("get post dorm error: ", err)
		return nil, err
	}
	return posts, nil
}

// 获取所有帖子
func GetPosts(page int, pageSize int) ([]*DgPost, int64, error) {
	posts := make([]*DgPost, 0)
	var total int64
	err := DB.Model(&DgPost{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err = DB.Preload("Images").
		Preload("Dorm").
		Preload("Type").
		Order("created_at desc"). // <--- 关键：时间倒序
		Offset(offset).
		Limit(pageSize). // <--- 关键：限制数量
		Find(&posts).Error
	if err != nil {
		fmt.Println("get posts error: ", err)
		return nil, 0, err
	}
	return posts, total, err
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

// 浏览量加一
func AddViewCount(postid uint) error {
	err := DB.Model(&DgPost{}).
		Where("id = ?", postid).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	if err != nil {
		fmt.Println("add view count error: ", err)
		return err
	}
	return nil
}

// 用户报名
func SignupPost(postID uint, userID uint) error {
	userid, err := UserIdToId(userID)
	if err != nil {
		return err
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 检查帖子是否存在、未过期、名额未满
		var post DgPost
		if err := tx.Where("id = ?", postID).First(&post).Error; err != nil {
			return err
		}
		// 检查是否开启限时
		if !post.IsLimited {
			return fmt.Errorf("该帖子未开启报名")
		}
		// 检查时间
		if post.Deadline != nil && time.Now().After(*post.Deadline) {
			return fmt.Errorf("报名已截止")
		}
		// 检查人数
		if post.MaxEnrollment > 0 && post.CurrentEnrollment >= post.MaxEnrollment {
			return fmt.Errorf("名额已满")
		}

		// 2. 检查是否重复报名
		var count int64
		tx.Model(&DgSignup{}).Where("post_id = ? AND user_id = ?", postID, userid).Count(&count)
		if count > 0 {
			return fmt.Errorf("你已经报过名了")
		}

		// 3. 创建报名记录
		signup := DgSignup{PostID: postID, UserID: userid}
		if err := tx.Create(&signup).Error; err != nil {
			return err
		}

		// 4. 帖子当前报名人数 +1
		if err := tx.Model(&DgPost{}).Where("id = ?", postID).UpdateColumn("current_enrollment", gorm.Expr("current_enrollment + ?", 1)).Error; err != nil {
			return err
		}

		notification := DgNotification{
			ReceiverID: post.PublisherId,
			SenderID:   userid,
			Type:       "signup",
			PostID:     postID,
			Content:    "报名了你的活动", // 前端可以拼接成 "用户A 报名了你的活动"
		}
		if err := tx.Create(&notification).Error; err != nil {
			return err
		}

		return nil
	})
}

// 检查用户是否已报名
func CheckIsSignedUp(postID uint, userID uint) (bool, error) {
	userid, err := UserIdToId(userID)
	if err != nil {
		return false, err
	}
	var count int64
	err = DB.Model(&DgSignup{}).Where("post_id = ? AND user_id = ?", postID, userid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetPostsByUserID 根据用户ID获取该用户发布的所有帖子
func GetPostsByUserID(userID int) ([]*DgPost, error) {
	posts := make([]*DgPost, 0)
	err := DB.Where("publisherid=?", userID).
		Preload("Images").
		Preload("Dorm").
		Preload("Type").
		Order("created_at desc").
		Find(&posts).Error
	if err != nil {
		fmt.Println("get user posts error: ", err)
		return nil, err
	}
	return posts, nil
}
