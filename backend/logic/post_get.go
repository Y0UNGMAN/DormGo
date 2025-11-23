package logic

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
)

// 获取post 信息 和 发布者user（只有name） 信息
func GetPostDetail(id int) (*model.ApiPostDetail, error) {
	post, err := model.GetPostDetail(id)
	if err != nil {
		fmt.Println("GetPostDetail error: ", err)
		return nil, err
	}
	publisherId := post.PublisherId

	user, err := model.GetUserById(int(publisherId))
	if err != nil {
		fmt.Println("GetUserById error: ", err)
		return nil, err
	}

	publisherName := user.Username
	publisherAvatar := user.Avatar
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
		user, err := model.GetUserById(int(post.PublisherId))
		if err != nil {
			fmt.Println("GetUserById error: ", err)
			return nil, err
		}
		p := post
		postDetails = append(postDetails, &model.ApiPostDetail{
			PublisherName:   user.Username,
			PublisherAvator: user.Avatar,
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
		user, err := model.GetUserById(int(post.PublisherId))
		if err != nil {
			fmt.Println("GetUserById error: ", err)
			return nil, err
		}
		p := post
		postDetails = append(postDetails, &model.ApiPostDetail{
			PublisherName:   user.Username,
			PublisherAvator: user.Avatar,
			DgPost:          p,
		})
	}
	return postDetails, nil
}
