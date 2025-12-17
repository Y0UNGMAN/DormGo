package logic

import (
	"errors"

	"github.com/Y0UNGMAN/DormGo/backend/model"
)

// DoFavorite 收藏帖子业务逻辑
func DoFavorite(userID uint, postID uint) error {
	// 1. (可选) 检查帖子是否存在，防止收藏不存在的帖子
	// post, err := model.GetPostDetail(int(postID))
	// if err != nil || post.ID == 0 { return errors.New("帖子不存在") }

	// 2. 检查是否已经收藏
	exists, err := model.CheckFavoriteExist(userID, postID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("你已经收藏过该帖子")
	}

	// 3. 构建模型并保存
	fav := &model.DgFavorite{
		UserID: userID,
		PostID: postID,
	}

	return model.InsertFavorite(fav)
}

// CancelFavorite 取消收藏业务逻辑
func CancelFavorite(userID uint, postID uint) error {
	return model.DeleteFavoriteByKeys(userID, postID)
}

// GetUserFavoriteList 获取用户收藏列表业务逻辑
func GetUserFavoriteList(userID uint, page int, size int) (map[string]interface{}, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset := (page - 1) * size

	// 调用 Model 查询
	posts, total, err := model.SelectFavoritesByUserID(userID, offset, size)
	if err != nil {
		return nil, err
	}

	// 组装返回数据结构
	result := map[string]interface{}{
		"list":  posts,
		"total": total,
		"page":  page,
		"size":  size,
	}
	return result, nil
}

// CheckIsFavorited 辅助逻辑：查询帖子详情时，判断当前用户是否收藏
func CheckIsFavorited(userID uint, postID uint) bool {
	exists, err := model.CheckFavoriteExist(userID, postID)
	if err != nil {
		return false
	}
	return exists
}
