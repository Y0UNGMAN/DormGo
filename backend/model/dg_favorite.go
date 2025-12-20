package model

import (
	"time"
)

// DgFavorite 数据库实体
type DgFavorite struct {
	ID        uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	UserID    uint      `gorm:"column:user_id;index;not null" json:"user_id"` // 添加索引优化查询
	PostID    uint      `gorm:"column:post_id;index;not null" json:"post_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// CheckFavoriteExist 检查是否存在收藏记录
func CheckFavoriteExist(userID uint, postID uint) (bool, error) {
	var count int64
	err := DB.Model(&DgFavorite{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

// InsertFavorite 插入一条收藏记录
func InsertFavorite(fav *DgFavorite) error {
	return DB.Create(fav).Error
}

// DeleteFavoriteByKeys 根据用户ID和帖子ID物理删除
func DeleteFavoriteByKeys(userID uint, postID uint) error {
	return DB.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&DgFavorite{}).Error
}

// SelectFavoritesByUserID 查询用户的收藏列表 (包含分页)
// 这里使用了 Join 查询，把 DgPost 的信息一起查出来
func SelectFavoritesByUserID(userID uint, offset int, limit int) ([]*DgPost, int64, error) {
	var posts []*DgPost
	var total int64

	// 1. 统计总数
	countQuery := DB.Model(&DgPost{}).
		Joins("JOIN dg_favorites ON dg_favorites.post_id = dg_posts.id").
		Where("dg_favorites.user_id = ?", userID)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 2. 查询数据 (预加载关联表)
	err := DB.Model(&DgPost{}).
		Joins("JOIN dg_favorites ON dg_favorites.post_id = dg_posts.id").
		Where("dg_favorites.user_id = ?", userID).
		Preload("Images").
		Preload("Dorm").
		Preload("Type").
		Order("dg_favorites.created_at desc"). // 按收藏时间倒序
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	return posts, total, err
}
