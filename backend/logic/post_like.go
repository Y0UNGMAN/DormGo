package logic

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
)

func PostLike(likeDetail *model.DgPostLike) error {
	//点赞表中增加记录
	err := model.PostLike(likeDetail)
	if err != nil {
		return err
	}

	//post表中点赞数加一
	err = model.AddLikeCount(likeDetail.PostID)
	if err != nil {
		secerr := model.CancelPostLike(likeDetail.PostID, likeDetail.LikerID)
		if secerr != nil {
			return fmt.Errorf("操作失败： %v , 且回滚失败： %v ", err, secerr)
		}
		return err
	}
	return nil

}

func CancelPostLike(postid uint, likerid uint) error {
	//点赞表中删去一行记录
	err := model.CancelPostLike(postid, likerid)
	if err != nil {
		return err
	}

	//post表中点赞数减一
	err = model.DelLikeCount(postid)
	if err != nil {
		secerr := model.AddLikeCount(postid)
		if secerr != nil {
			return fmt.Errorf("操作失败： %v , 且回滚失败： %v ", err, secerr)
		}
		return err
	}
	return nil
}
