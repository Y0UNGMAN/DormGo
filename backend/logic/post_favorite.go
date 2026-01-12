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

	// 调用 Model 查询，获取帖子列表（DgPost）
	posts, total, err := model.SelectFavoritesByUserID(userID, offset, size)
	if err != nil {
		return nil, err
	}

	// 将 DgPost 转换为 ApiPostDetail，补充发布者信息（与 GetPosts/ GetUserPosts 保持一致）
	postDetails := make([]*model.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := model.GetUserById(int(post.PublisherId))
		if err != nil {
			// 若获取用户失败，跳过该条记录但继续返回其他记录
			continue
		}
		p := post
		postDetails = append(postDetails, &model.ApiPostDetail{
			PublisherName:   user.Username,
			PublisherAvator: user.Avatar,
			PublisherIntro:  user.Intro,
			DgPost:          p,
		})
	}

	// 组装返回数据结构，保持与其他接口兼容的字段名
	result := map[string]interface{}{
		"list":  postDetails,
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
