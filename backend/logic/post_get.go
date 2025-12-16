package logic

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
)

// 辅助函数：安全获取用户信息
func safeGetUserInfo(publisherId uint) (string, string) {
	user, err := model.GetUserById(int(publisherId))
	if err != nil {
		// 如果找不到用户，不要报错，返回默认值
		// 这样即使有坏数据，帖子列表也能显示
		return "未知用户", "https://cube.elemecdn.com/9/c2/f0ee8a3c7c9638a54940382568c9dpng.png"
	}
	return user.Username, user.Avatar
}

// 获取 post 信息
func GetPostDetail(id int) (*model.ApiPostDetail, error) {
	post, err := model.GetPostDetail(id)
	if err != nil {
		fmt.Println("GetPostDetail error: ", err)
		return nil, err
	}

	publisherName, publisherAvatar := safeGetUserInfo(post.PublisherId)

	postDetail := &model.ApiPostDetail{
		PublisherName:   publisherName,
		PublisherAvator: publisherAvatar,
		DgPost:          post,
	}

	return postDetail, nil
}

func GetPostByDorm(dormid int) ([]*model.ApiPostDetail, error) {
	posts, err := model.GetPostByDorm(dormid)
	if err != nil {
		fmt.Println("GetPostByDorm error: ", err)
		return nil, err
	}
	postDetails := make([]*model.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 修改：使用 safeGetUserInfo
		publisherName, publisherAvatar := safeGetUserInfo(post.PublisherId)

		p := post
		postDetails = append(postDetails, &model.ApiPostDetail{
			PublisherName:   publisherName,
			PublisherAvator: publisherAvatar,
			DgPost:          p,
		})
	}
	return postDetails, nil
}

func GetPosts() ([]*model.ApiPostDetail, error) {
	posts, err := model.GetPosts()
	if err != nil {
		fmt.Println("GetPosts error: ", err)
		return nil, err
	}
	postDetails := make([]*model.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 修改：使用 safeGetUserInfo，防止因为一个坏数据导致整个列表崩塌
		publisherName, publisherAvatar := safeGetUserInfo(post.PublisherId)

		p := post
		postDetails = append(postDetails, &model.ApiPostDetail{
			PublisherName:   publisherName,
			PublisherAvator: publisherAvatar,
			DgPost:          p,
		})
	}
	return postDetails, nil
}
