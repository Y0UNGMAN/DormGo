package model

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() (err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		println("数据库连接失败")
	} else {
		println("数据库连接成功")
	}
	//创建用户表
	err = DB.AutoMigrate(&DgUser{})
	if err != nil {
		println("创建用户表失败")
	} else {
		println("创建用户表成功")
	}
	//创建post表
	err = DB.AutoMigrate(&DgPost{})
	if err != nil {
		println("创建post表失败")
	} else {
		println("创建post表成功")
	}
	//创建宿舍表
	err = DB.AutoMigrate(&DgDorm{})
	if err != nil {
		println("创建dorm表失败")
	} else {
		println("创建dorm表成功")
	}
	//创建类型表
	err = DB.AutoMigrate(&DgType{})
	if err != nil {
		println("创建type表失败")
	} else {
		println("创建type表成功")
	}
	//创建admin表
	err = DB.AutoMigrate(&DgAdmin{})
	if err != nil {
		println("创建admin表失败")
	} else {
		println("创建admin表成功")
	}
	err = DB.AutoMigrate(&DgComment{})
	if err != nil {
		println("创建Comment表失败")
	} else {
		println("创建Comment表成功")
	}
	err = DB.AutoMigrate(&DgImages{})
	if err != nil {
		println("创建Images表失败")
	} else {
		println("创建Images表成功")
	}
	err = DB.AutoMigrate(&DgPostLike{})
	if err != nil {
		println("创建PostLike表失败")
	} else {
		println("创建PostLike表成功")
	}
	//创建cmtlike表
	err = DB.AutoMigrate(&DgCmtLike{})
	if err != nil {
		println("创建cmtlike表失败")
	} else {
		println("创建cmtlike表成功")
	}

	return
}
