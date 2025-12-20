package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Y0UNGMAN/DormGo/backend/logic"
	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
	"github.com/gin-gonic/gin"
)

// backend/controller/admin_controller.go

// ================= 1. 数据统计 (首页) =================

func AdminStatistics(c *gin.Context) {
	var totalUsers, totalPosts, totalComments, totalViolations, totalDorms int64

	// 1. 基础数据统计 (查询真实数据库)
	// 使用 model.DB 对各个模型进行计数
	model.DB.Model(&model.DgUser{}).Count(&totalUsers)
	model.DB.Model(&model.DgPost{}).Count(&totalPosts)
	model.DB.Model(&model.DgComment{}).Count(&totalComments)
	model.DB.Model(&model.DgViolation{}).Count(&totalViolations)
	model.DB.Model(&model.DgDorm{}).Count(&totalDorms)

	// 2. 帖子分类统计 (关键修复部分)
	var types []model.DgType
	// 查询所有存在的板块类型
	if err := model.DB.Find(&types).Error; err != nil {
		fmt.Println("查询板块类型失败:", err)
	}

	var typeNames []string
	var typeCounts []int64

	for _, t := range types {
		var count int64

		// 【关键修改】:
		// 根据 backend/model/dg_post.go 的定义: TypeId uint `gorm:"column:typeid" ...`
		// 这里的 Where 条件必须显式指定 "typeid = ?"。
		// 如果不指定或写成 "type_id"，GORM 会找不到列导致统计结果为 0。
		model.DB.Model(&model.DgPost{}).Where("typeid = ?", t.TypeId).Count(&count)

		typeNames = append(typeNames, t.TypeName)
		typeCounts = append(typeCounts, count)
	}

	// 3. 返回数据给前端
	c.JSON(http.StatusOK, gin.H{
		// 顶部卡片数据
		"total_users":      totalUsers,
		"total_dorms":      totalDorms,
		"total_contents":   totalPosts + totalComments,
		"total_violations": totalViolations,

		// 饼图数据 (内容构成)
		"pie_data": []gin.H{
			{"value": totalPosts, "name": "帖子"},
			{"value": totalComments, "name": "评论"},
			{"value": totalViolations, "name": "违规记录"},
		},

		// 柱状图数据 (帖子板块分布)
		// type_names: ["失物招领", "二手交易", ...]
		// type_values: [12, 5, ...]
		"type_names":  typeNames,
		"type_values": typeCounts,
	})
}

// ================= 2. 用户管理 =================

// GetUserList 获取用户列表（修复后）
func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	var users []model.DgUser
	var total int64

	query := model.DB.Model(&model.DgUser{}).Preload("Dorm")

	// 搜索逻辑：学号 或 用户名(Username)
	if keyword != "" {
		query = query.Where("username LIKE ? OR studentid LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	// 【关键修改】: 修改排序为 ID 从小到大 (id asc)
	err := query.Order("id asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
		return
	}

	var list []gin.H
	for _, u := range users {
		dormName := "未知"
		if u.Dorm.DormName != "" {
			dormName = u.Dorm.DormName
		}

		list = append(list, gin.H{
			"id":           u.ID,
			"student_id":   u.StudentId, // 学号
			"nickname":     u.Username,  // 用户名
			"avatar":       u.Avatar,
			"dorm_name":    dormName,
			"status":       u.Status,      // 【关键修复】返回真实状态 (1:正常, 2:封禁)
			"credit_score": u.CreditScore, // 使用真实信用分
			"created_at":   u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// BanUserHandler 封禁用户接口
func BanUserHandler(c *gin.Context) {
	var req struct {
		UserID int64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	// 调用 Logic 层 (使用独立函数调用，无需实例化结构体)
	if err := logic.BanUser(req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户封禁成功",
	})
}
func UpdateUserStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "状态已更新"})
}

func BatchUpdateUserStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "批量操作成功"})
}

func UpdateUserCredit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "信用分已调整"})
}

// ================= 3. 违规处理 =================

func GetViolationList(c *gin.Context) {
	var list []model.DgViolation
	var total int64

	userKeyword := c.Query("userKeyword")
	status := c.Query("status")
	vType := c.Query("type")
	reason := c.Query("reason")

	query := model.DB.Model(&model.DgViolation{}).Preload("User")

	if userKeyword != "" {
		query = query.Joins("JOIN dg_users ON dg_users.id = dg_violations.user_id").
			Where("dg_users.username LIKE ? OR dg_users.studentid LIKE ?", "%"+userKeyword+"%", "%"+userKeyword+"%")
	}
	if status != "" {
		query = query.Where("dg_violations.status = ?", status)
	}
	if vType != "" {
		query = query.Where("dg_violations.type = ?", vType)
	}
	if reason != "" {
		query = query.Where("dg_violations.reason LIKE ?", "%"+reason+"%")
	}

	query.Count(&total)
	query.Order("created_at desc").Find(&list)

	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

func ProcessViolation(c *gin.Context) {
	var req struct {
		Punishment string `json:"punishment"`
		Notes      string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}

	model.DB.Model(&model.DgViolation{}).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
		"status":     "processed",
		"punishment": req.Punishment,
		"notes":      req.Notes,
	})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "处理成功"})
}

// ================= 4. 内容审核 (帖子) =================

func GetContentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	var posts []model.DgPost
	var total int64

	query := model.DB.Model(&model.DgPost{}).Preload("Type").Preload("Dorm")

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status == "canceled" {
		// 1. 如果前端选了“违规/撤销”Tab，只查违规的
		query = query.Where("status = ?", "canceled")
	} else {
		// 2. 如果前端选了“全部”Tab (status=all)，只查【非违规】的
		// 这样“全部”里就不会出现已经删除的帖子了，操作完会自动消失
		query = query.Where("status != ?", "canceled")
	}

	query.Count(&total)
	query.Order("is_pinned desc, created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&posts)

	var list []gin.H
	for _, p := range posts {
		userName := "管理员"
		userDorm := "官方"

		var user model.DgUser
		if err := model.DB.First(&user, p.PublisherId).Error; err == nil {
			userName = user.Username
			userDorm = p.Dorm.DormName
		}

		pStatus := p.Status
		if pStatus == "" {
			pStatus = "normal"
		}

		list = append(list, gin.H{
			"id":             p.ID,
			"title":          p.Title,
			"content":        p.Content,
			"type_name":      p.Type.TypeName,
			"publisher_name": userName,
			"dorm_name":      userDorm,
			"status":         pStatus,
			"is_pinned":      p.IsPinned,
			"reward_type":    "无偿",
			"created_at":     p.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

func AuditContent(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	c.ShouldBindJSON(&req)
	model.DB.Model(&model.DgPost{}).Where("id = ?", id).Update("status", req.Status)

	cacheKey := fmt.Sprintf("dormgo:posts:page:%d:size:%d", 1, 10)
	myredis.GetClient().Del(cacheKey)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功"})
}

func BatchAuditContent(c *gin.Context) {
	var req struct {
		IDs    []uint `json:"ids"`
		Status string `json:"status"`
	}
	c.ShouldBindJSON(&req)
	model.DB.Model(&model.DgPost{}).Where("id IN ?", req.IDs).Update("status", req.Status)
	cacheKey := fmt.Sprintf("dormgo:posts:page:%d:size:%d", 1, 10)
	myredis.GetClient().Del(cacheKey)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "批量操作成功"})
}

// TogglePostPin 置顶/取消置顶
func TogglePostPin(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		IsPinned bool `json:"is_pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	model.DB.Model(&model.DgPost{}).Where("id = ?", id).Update("is_pinned", req.IsPinned)
	cacheKey := fmt.Sprintf("dormgo:posts:page:%d:size:%d", 1, 10)
	myredis.GetClient().Del(cacheKey)
	fmt.Printf("🗑️ [管理员置顶操作] 缓存已清除: %s\n", cacheKey)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功"})
}

// AdminCreatePost 管理员发帖
func AdminCreatePost(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	typeId, _ := strconv.Atoi(c.PostForm("typeid"))
	isPinned := c.PostForm("is_pinned") == "true"
	if title == "" || content == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "标题和内容不能为空"})
		return
	}
	var postImages []model.DgImages
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["images"]
		for i, file := range files {
			url, _ := utils.UploadFile(file, "AdminPost")
			postImages = append(postImages, model.DgImages{ImageURL: url, Order: uint(i)})
		}
	}

	officialUserID := uint(10)
	post := &model.DgPost{
		Title:       title,
		Content:     content,
		PublisherId: officialUserID,
		DormId:      1,
		TypeId:      uint(typeId),
		Status:      "normal",
		IsPinned:    isPinned,
		Images:      postImages,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := model.CreatePost(post); err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "发布失败: " + err.Error()})
		return
	}
	cacheKey := fmt.Sprintf("dormgo:posts:page:%d:size:%d", 1, 10)
	myredis.GetClient().Del(cacheKey)
	fmt.Printf("🗑️ [管理员] 缓存已清除: %s\n", cacheKey)

	c.JSON(200, gin.H{"code": 200, "msg": "发布成功"})
}

// ================= 5. 管理员个人中心 =================

func GetAdminProfile(c *gin.Context) {
	adminId, _ := c.Get("UserID")

	var admin model.DgAdmin
	if err := model.DB.First(&admin, adminId).Error; err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "管理员不存在"})
		return
	}

	// 简介默认为空时显示默认文案，但不存入库
	intro := admin.Intro
	if intro == "" {
		intro = "系统超级管理员"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":         admin.ID,
			"nickname":   admin.Username,
			"avatar":     admin.Avatar,
			"student_id": "ADMIN",
			"dorm_name":  "管理中心",
			"intro":      intro, // 返回真实简介
		},
	})
}

func UpdateAdminProfile(c *gin.Context) {
	adminId, _ := c.Get("UserID")

	// 1. 获取文本参数 (multipart/form-data)
	nickname := c.PostForm("nickname")
	intro := c.PostForm("intro")

	updates := map[string]interface{}{}

	if nickname != "" {
		updates["username"] = nickname
	}
	// 允许简介为空，或者更新为新值
	updates["intro"] = intro

	// 2. 处理头像文件上传
	file, err := c.FormFile("avatar")
	if err == nil {
		// 如果有文件上传，则进行上传处理
		// utils.UploadFile 需确保存在，通常在 utils/oss.go 中
		url, uploadErr := utils.UploadFile(file, "admin_avatar")
		if uploadErr != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "头像上传失败"})
			return
		}
		updates["avatar"] = url
	}

	// 3. 更新数据库
	if err := model.DB.Model(&model.DgAdmin{}).Where("id = ?", adminId).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// ================= 6. 楼栋 & 通知 & 配置 =================

// GetDormListAdmin 获取楼栋列表
func GetDormListAdmin(c *gin.Context) {
	var dorms []model.DgDorm
	model.DB.Find(&dorms)

	var list []gin.H
	for _, d := range dorms {
		var count int64
		model.DB.Model(&model.DgUser{}).Where("dormid = ?", d.DormId).Count(&count)
		// 【关键修改】这里统一返回 "dormname"
		list = append(list, gin.H{
			"id":          d.DormId,
			"dormname":    d.DormName, // 修正：使用 dormname
			"description": "学生公寓",
			"user_count":  count,
			"created_at":  time.Now(),
		})
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
}

func CreateDorm(c *gin.Context) {
	var dorm model.DgDorm
	// 注意：前端传来的JSON需包含 dormname
	if err := c.ShouldBindJSON(&dorm); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	model.DB.Create(&dorm)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功"})
}

func UpdateDorm(c *gin.Context) {
	id := c.Param("id")
	var dorm model.DgDorm
	c.ShouldBindJSON(&dorm)
	model.DB.Model(&model.DgDorm{}).Where("dormid = ?", id).Updates(dorm)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

func DeleteDorm(c *gin.Context) {
	id := c.Param("id")
	model.DB.Delete(&model.DgDorm{}, "dormid = ?", id)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

func GetNoticeList(c *gin.Context) {
	var notices []model.DgNotice
	var total int64
	model.DB.Model(&model.DgNotice{}).Count(&total)
	model.DB.Order("created_at desc").Find(&notices)
	c.JSON(http.StatusOK, gin.H{"list": notices, "total": total})
}

func CreateNotice(c *gin.Context) {
	var n model.DgNotice
	c.ShouldBindJSON(&n)
	n.CreatedAt = time.Now()
	model.DB.Create(&n)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "发布成功"})
}

func DeleteNotice(c *gin.Context) {
	model.DB.Delete(&model.DgNotice{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

// 敏感词 & 类型配置
func GetSensitiveWords(c *gin.Context) {
	var words []model.DgSensitiveWord
	model.DB.Find(&words)
	c.JSON(http.StatusOK, gin.H{"list": words})
}
func AddSensitiveWord(c *gin.Context) {
	var w model.DgSensitiveWord
	c.ShouldBindJSON(&w)
	model.DB.Create(&w)
	utils.WordFilter.AddWord(w.Word)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "添加成功"})
}
func DeleteSensitiveWord(c *gin.Context) {
	model.DB.Delete(&model.DgSensitiveWord{}, c.Param("id"))
	utils.UpdateFilter()
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}
func GetPostTypesAdmin(c *gin.Context) {
	var types []model.DgType
	if err := model.DB.Find(&types).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	var list []gin.H
	for _, t := range types {
		list = append(list, gin.H{"id": t.TypeId, "name": t.TypeName})
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
}

// CreatePostType 创建分类
func CreatePostType(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	newType := model.DgType{
		TypeName: req.Name,
	}
	if err := model.DB.Create(&newType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功"})
}

// UpdatePostType 更新分类名称
func UpdatePostType(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 更新数据库
	if err := model.DB.Model(&model.DgType{}).Where("typeid = ?", id).Update("typename", req.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// DeletePostType 删除分类
func DeletePostType(c *gin.Context) {
	id := c.Param("id")

	// 物理删除
	if err := model.DB.Delete(&model.DgType{}, "typeid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

func SendSystemNotification(c *gin.Context) {
	// 1. 定义请求参数
	var req struct {
		TargetType  string `json:"target_type"`  // "all" (全体) 或 "specific" (指定)
		TargetValue string `json:"target_value"` // 学号 (如果选了指定)
		Content     string `json:"content"`      // 通知内容
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	noticeTitle := "系统通知"
	historyNotice := model.DgNotice{
		Title:     noticeTitle, // 默认标题
		Content:   req.Content,
		Target:    req.TargetType,
		TargetID:  req.TargetValue,
		Publisher: "管理员", // 或者从 Token 获取当前管理员名字
		CreatedAt: time.Now(),
	}
	if err := model.DB.Create(&historyNotice).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "保存历史记录失败"})
		return
	}

	// 2. 设定官方发送者 ID (确保数据库 dg_users 表里有 id=10 的用户，头像设为官方Logo)
	const OfficialSenderID = 10

	// 3. 准备通知模板
	baseNotification := model.DgNotification{
		SenderID:  OfficialSenderID,
		Type:      "system", // 【关键】标记为系统通知，前端据此区分展示样式
		Content:   req.Content,
		PostID:    0, // 系统通知不关联帖子，设为 0
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	// 4. 根据类型发送
	if req.TargetType == "all" {
		// --- 群发逻辑 ---
		var users []model.DgUser
		// 只查 ID，性能优化
		if err := model.DB.Select("id").Find(&users).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "查询用户失败"})
			return
		}

		var notifications []model.DgNotification
		for _, u := range users {
			if u.ID == OfficialSenderID {
				continue
			} // 不发给自己

			note := baseNotification
			note.ReceiverID = u.ID
			notifications = append(notifications, note)
		}

		// 批量插入 (GORM V2 API)
		if len(notifications) > 0 {
			// 分批写入，防止一次性插入过多导致数据库报错
			if err := model.DB.CreateInBatches(notifications, 100).Error; err != nil {
				c.JSON(500, gin.H{"code": 500, "msg": "群发失败"})
				return
			}
		}

	} else if req.TargetType == "specific" {
		// --- 单发逻辑 ---
		var targetUser model.DgUser
		if err := model.DB.Where("studentid = ?", req.TargetValue).First(&targetUser).Error; err != nil {
			c.JSON(404, gin.H{"code": 404, "msg": "未找到该学号用户"})
			return
		}

		note := baseNotification
		note.ReceiverID = targetUser.ID
		if err := model.DB.Create(&note).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "发送失败"})
			return
		}
	} else {
		c.JSON(400, gin.H{"code": 400, "msg": "无效的发送类型"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "通知已推送到用户中心"})
}

func UpdateSensitiveWord(c *gin.Context)  { c.JSON(200, gin.H{"code": 200}) }
func UpdatePostTypeStatus(c *gin.Context) { c.JSON(200, gin.H{"code": 200}) }
