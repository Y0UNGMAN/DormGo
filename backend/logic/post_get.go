package logic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
)

type CachedPostsData struct {
	List  []*model.ApiPostDetail `json:"list"`
	Total int64                  `json:"total"`
}

// 获取post 信息 和 发布者user（只有name） 信息
func GetPostDetail(id int, usrId uint) (*model.ApiPostDetail, bool, bool, error) {
	post, err := model.GetPostDetail(id)
	if err != nil {
		fmt.Println("GetPostDetail error: ", err)
		return nil, false, false, err
	}
	go func() {
		viewKey := fmt.Sprintf("dormgo:post:view_inc:%d", post.ID)
		myredis.GetClient().Incr(viewKey)
	}()
	client := myredis.GetClient()
	// 合并点赞数
	likeIncKey := fmt.Sprintf("dormgo:post:like_inc:%d", post.ID)
	if likeIncStr, err := client.Get(likeIncKey).Result(); err == nil {
		if likeInc, _ := strconv.Atoi(likeIncStr); likeInc != 0 {
			// 注意：LikeCount 是 uint，要做防溢出处理
			newCount := int(post.LikeCount) + likeInc
			if newCount < 0 {
				newCount = 0
			}
			post.LikeCount = uint(newCount)
		}
	}
	// 合并浏览量 (同理)
	viewIncKey := fmt.Sprintf("dormgo:post:view_inc:%d", post.ID)
	if viewIncStr, err := client.Get(viewIncKey).Result(); err == nil {
		if viewInc, _ := strconv.Atoi(viewIncStr); viewInc != 0 {
			post.ViewCount += uint(viewInc)
		}
	}

	publisherId := post.PublisherId

	user, err := model.GetUserById(int(publisherId))
	if err != nil {
		fmt.Println("GetUserById error: ", err)
		return nil, false, false, err
	}

	isLiked, err := model.CheckIsLiked(uint(id), usrId)
	if err != nil {
		fmt.Println("IsLiked error: ", err)
		return nil, false, false, err
	}
	isFavorited := false
	isFavorited, err = model.CheckFavoriteExist(usrId, uint(id))
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

	return postDetail, isLiked, isFavorited, nil

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

func GetPosts(page int, pageSize int) ([]*model.ApiPostDetail, int64, error) {

	cacheKey := fmt.Sprintf("dormgo:posts:page:%d:size:%d", page, pageSize)
	client := myredis.GetClient()
	val, err := client.Get(cacheKey).Result()
	if err == nil {
		// 缓存命中 (Hit)
		var cachedData CachedPostsData
		if jsonErr := json.Unmarshal([]byte(val), &cachedData); jsonErr == nil {
			fmt.Printf("🚀 Cache HIT: %s\n", cacheKey)
			return cachedData.List, cachedData.Total, nil
		}
	}
	fmt.Printf("🐢 Cache MISS: %s, Querying DB...\n", cacheKey)

	posts, total, err := model.GetPosts(page, pageSize)
	if err != nil {
		fmt.Println("GetPosts error: ", err)
		return nil, 0, err
	}
	postDetails := make([]*model.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := model.GetUserById(int(post.PublisherId))
		if err != nil {
			fmt.Println("GetUserById error: ", err)
			continue
		}
		p := post
		postDetails = append(postDetails, &model.ApiPostDetail{
			PublisherName:   user.Username,
			PublisherAvator: user.Avatar,
			DgPost:          p,
		})
	}
	saveData := CachedPostsData{
		List:  postDetails,
		Total: total,
	}
	jsonBytes, _ := json.Marshal(saveData)
	client.Set(cacheKey, string(jsonBytes), time.Minute*10)

	return postDetails, total, nil
}
