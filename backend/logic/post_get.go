package logic

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
)

// 获取post 信息 和 发布者user（只有name） 信息
func GetPostDetail(id int, usrId uint) (*model.ApiPostDetail, bool, error) {
	post, err := model.GetPostDetail(id)
	if err != nil {
		fmt.Println("GetPostDetail error: ", err)
		return nil, false, err
	}
	go func() {
		_ = model.AddViewCount(post.ID)
	}()

	publisherId := post.PublisherId

	user, err := model.GetUserById(int(publisherId))
	if err != nil {
		fmt.Println("GetUserById error: ", err)
		return nil, false, err
	}

	isLiked, err := model.CheckIsLiked(uint(id), usrId)
	if err != nil {
		fmt.Println("IsLiked error: ", err)
		return nil, false, err
	}

	isSignedUp := false
	if post.IsLimited {
		var err error
		isSignedUp, err = model.CheckIsSignedUp(uint(id), usrId)
		if err != nil {
			fmt.Println("CheckSignedUp error:", err)
		}
	}

	publisherName := user.Username
	publisherAvatar := user.Avatar
	postDetail := &model.ApiPostDetail{
		PublisherName:   publisherName,
		PublisherAvator: publisherAvatar,
		IsSignedUp:      isSignedUp,
		DgPost:          post,
	}

	return postDetail, isLiked, nil

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
