package model

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DgUser struct {
	ID          uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	UserID      uint      `gorm:"column:user_id;" json:"userid"`
	StudentId   string    `gorm:"column:studentid;type:varchar(20);" json:"studentid"`
	Username    string    `gorm:"column:username;type:varchar(50);" json:"username"`
	Password    string    `gorm:"column:password;type:varchar(255);" json:"password"`
	Avatar      string    `gorm:"column:avatar;type:varchar(255);" json:"avatarurl"`
	Intro       string    `gorm:"column:intro;type:varchar(255);" json:"intro"`
	Email       string    `gorm:"column:email;type:varchar(255);" json:"email"`
	DormId      uint      `gorm:"column:dormid" json:"dormid"`
	Dorm        DgDorm    `gorm:"foreignkey:DormId;references:DormId" json:"dorm"`
	Status      int       `gorm:"column:status;default:1" json:"status"`               // 1:正常 2:封禁
	CreditScore int       `gorm:"column:credit_score;default:100" json:"credit_score"` // 信用分
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// 注册请求参数
type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	StudentId  string `json:"student_id" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
	DormId     uint   `json:"dorm_id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 根据id获取用户信息
func GetUserById(id int) (*DgUser, error) {
	var user DgUser
	err := DB.Preload("Dorm").Where("id=?", id).First(&user).Error
	if err != nil {
		fmt.Println("get user detail error: ", err)
		return nil, err
	}
	return &user, nil
}

func CheckUserExist(username string) (bool, error) {
	var count int64
	err := DB.Model(&DgUser{}).Where("username=?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertUser(user *DgUser) error {
	//对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func Login(username, password string) (*DgUser, error) {
	var user DgUser
	err := DB.Preload("Dorm").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err // 数据库连接等其他错误
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 【新增】验证账号状态
	// 假设 1 为正常状态，非 1 (如 2) 则禁止登录
	if user.Status != 1 {
		return nil, errors.New("账号已封禁，请联系管理员")
	}

	// 登录成功，返回用户对象
	return &user, nil
}

func UserIdToId(userId uint) (uint, error) {
	var user DgUser
	err := DB.Model(&DgUser{}).Where("user_id=?", userId).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
