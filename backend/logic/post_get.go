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

	postDetail := &model.ApiPostDetail{
		PublisherName: publisherName,
		DgPost:        post,
	}

	return postDetail, nil

}

func GetPostByDorm(dormid int) ([]model.ApiPostDetail, error) {
	posts, err := model.GetPostByDorm(dormid)
	if err != nil {
		fmt.Println("GetPostByDorm error: ", err)
		return nil, err
	}
	var postDetails []model.ApiPostDetail
	for _, post := range posts {
		user, err := model.GetUserById(int(post.PublisherId))
		if err != nil {
			fmt.Println("GetUserById error: ", err)
			return nil, err
		}
		p := post
		postDetails = append(postDetails, model.ApiPostDetail{
			PublisherName: user.Username,
			DgPost:        &p,
		})
	}
	return postDetails, nil
}
