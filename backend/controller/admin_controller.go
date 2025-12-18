package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	myredis "github.com/Y0UNGMAN/DormGo/backend/redis"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
	"github.com/gin-gonic/gin"
)

// ================= 1. 数据统计 (首页) =================

func AdminStatistics(c *gin.Context) {
	var totalUsers, totalPosts, totalComments, totalViolations int64

	// 统计真实数据库数据
	model.DB.Model(&model.DgUser{}).Count(&totalUsers)
	model.DB.Model(&model.DgPost{}).Count(&totalPosts)
	model.DB.Model(&model.DgComment{}).Count(&totalComments)
	model.DB.Model(&model.DgViolation{}).Count(&totalViolations)

	// 模拟趋势数据
	c.JSON(http.StatusOK, gin.H{
		"total_users":      totalUsers,
		"total_dorms":      12, // 可查库
		"total_contents":   totalPosts + totalComments,
		"total_violations": totalViolations,

		// 图表数据 (模拟)
		"dates":            []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"},
		"users_trend_data": []int{12, 18, 25, 30, 45, 50, 60},
		"posts_trend_data": []int{5, 8, 12, 15, 20, 18, 25},
		"pie_data": []gin.H{
			{"value": totalPosts, "name": "帖子"},
			{"value": totalComments, "name": "评论"},
			{"value": totalViolations, "name": "违规"},
		},
	})
}

// ================= 2. 用户管理 =================

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	var users []model.DgUser
	var total int64

	query := model.DB.Model(&model.DgUser{}).Preload("Dorm")

	// 搜索逻辑：学号 或 昵称
	if keyword != "" {
		query = query.Where("username LIKE ? OR studentid LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
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
			"student_id":   u.StudentId, // 统一使用学号
			"nickname":     u.Username,
			"avatar":       u.Avatar,
			"dorm_name":    dormName,
			"credit_score": 5.0,      // 默认信用分
			"status":       "normal", // 默认状态
		})
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
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

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":         admin.ID,
			"nickname":   admin.Username,
			"avatar":     admin.Avatar,
			"student_id": "ADMIN", // 管理员无学号，显示 ADMIN
			"dorm_name":  "管理中心",
			"intro":      "系统超级管理员",
		},
	})
}

func UpdateAdminProfile(c *gin.Context) {
	adminId, _ := c.Get("UserID")
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}

	model.DB.Model(&model.DgAdmin{}).Where("id = ?", adminId).Updates(map[string]interface{}{
		"username": req.Nickname,
		"avatar":   req.Avatar,
	})

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
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "添加成功"})
}
func DeleteSensitiveWord(c *gin.Context) {
	model.DB.Delete(&model.DgSensitiveWord{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}
func GetPostTypesAdmin(c *gin.Context) {
	var types []model.DgType
	model.DB.Find(&types)
	var list []gin.H
	for _, t := range types {
		list = append(list, gin.H{"id": t.TypeId, "name": t.TypeName})
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
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

func UpdateBasicConfig(c *gin.Context)    { c.JSON(200, gin.H{"code": 200}) }
func UpdateSecurityConfig(c *gin.Context) { c.JSON(200, gin.H{"code": 200}) }
func UpdateSensitiveWord(c *gin.Context)  { c.JSON(200, gin.H{"code": 200}) }
func CreatePostType(c *gin.Context)       { c.JSON(200, gin.H{"code": 200}) }
func UpdatePostType(c *gin.Context)       { c.JSON(200, gin.H{"code": 200}) }
func UpdatePostTypeStatus(c *gin.Context) { c.JSON(200, gin.H{"code": 200}) }
func DeletePostType(c *gin.Context)       { c.JSON(200, gin.H{"code": 200}) }
