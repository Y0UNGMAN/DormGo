package logic

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
)

func CreatePost(post *model.DgPost) (err error) {

	//1 保存到数据库
	err = model.CreatePost(post)
	//2 返回
	if err != nil {
		fmt.Println("create post error: ", err)
		return err
	}
	return nil
}
