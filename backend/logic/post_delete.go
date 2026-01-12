package logic

import (
	"errors"
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
)

// DeletePost 删除帖子
func DeletePost(postID uint, userID uint) error {
	// 1. 查询帖子是否存在
	post, err := model.GetPostDetail(int(postID))
	if err != nil {
		return errors.New("帖子不存在")
	}
	userid, err := model.UserIdToId(userID)
	if err != nil {
		return errors.New("你无权删除该帖子1")
	}
	// 2. 权限校验：必须是当前登录用户发布的帖子
	if post.PublisherId != userid {
		return errors.New("你无权删除该帖子2")
	}

	// 3. 执行删除 (GORM 会根据是否由 DeletedAt 字段决定软/硬删除)
	err = model.DB.Delete(post).Error
	if err != nil {
		return errors.New("删除失败，请稍后重试")
	}

	// 4.
	rdb := myredis.GetClient()

	// 这里的 Scan 不需要 context 参数
	iter := rdb.Scan(0, "dormgo:posts:page:*", 0).Iterator()

	for iter.Next() {
		key := iter.Val()
		// 删除也不需要 context 参数
		if err := rdb.Del(key).Err(); err != nil {
			fmt.Printf("删除缓存失败 Key: %s, Error: %v\n", key, err)
		} else {
			fmt.Printf("已清理缓存: %s\n", key)
		}
	}

	if err := iter.Err(); err != nil {
		fmt.Printf("遍历缓存Key失败: %v\n", err)
	}

	return nil
}
