package controller

import (
	"fmt"
	"strconv"
	"time"

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
	isLimited, _ := strconv.ParseBool(c.PostForm("is_limited"))
	maxEnrollment, _ := strconv.Atoi(c.PostForm("max_enrollment"))
	deadlineStr := c.PostForm("deadline")
	var deadline *time.Time
	if isLimited && deadlineStr != "" {
		// 尝试解析时间
		if t, err := time.Parse(time.RFC3339, deadlineStr); err == nil {
			deadline = &t
		} else {
			fmt.Println("时间解析失败:", err)
			// 如果解析失败，视业务需求决定是报错还是忽略，这里暂时忽略
		}
	}
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
		Title:         title,
		Content:       content,
		PublisherId:   uint(publisherid),
		DormId:        uint(dormid),
		TypeId:        uint(typeid),
		Images:        postImages,
		IsLimited:     isLimited,
		MaxEnrollment: maxEnrollment,
		Deadline:      deadline,
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
