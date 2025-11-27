package controller

import (
	"fmt"
	"strconv"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
	"github.com/gin-gonic/gin"
)

func PostCreate(c *gin.Context) {

	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "title or content is empty",
		})
		return
	}
	publisherid, _ := strconv.Atoi(c.PostForm("publisherid"))
	dormid, _ := strconv.Atoi(c.PostForm("dormid"))
	typeid, _ := strconv.Atoi(c.PostForm("typeid"))

	var postImages []model.DgImages
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["images"]
		for i, file := range files {
			url, secerr := utils.UploadFile(file, "PostImg")
			if secerr != nil {
				fmt.Println("图片上传失败:", err)
				c.JSON(500, gin.H{"code": 500, "msg": "图片上传失败"})
				return
			}
			postImages = append(postImages, model.DgImages{
				ImageURL: url,
				Order:    uint(i),
			})
		}
	}

	post := &model.DgPost{
		Title:       title,
		Content:     content,
		PublisherId: uint(publisherid),
		DormId:      uint(dormid),
		TypeId:      uint(typeid),
		Images:      postImages,
	}

	err = logic.CreatePost(post)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "创建帖子失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Create Success",
	})

}
