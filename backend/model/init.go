package model

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() (err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mytsql.hos"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		println("数据库连接失败")
	} else {
		println("数据库连接成功")
	}
	//创建用户表
	err = db.AutoMigrate(&DgUser{})
	if err != nil {
		println("创建用户表失败")
	} else {
		println("创建用户表成功")
	}
	//创建post表
	err = db.AutoMigrate(&DgPost{})
	if err != nil {
		println("创建post表失败")
	} else {
		println("创建post表成功")
	}
	//创建宿舍表
	err = db.AutoMigrate(&DgDorm{})
	if err != nil {
		println("创建dorm表失败")
	} else {
		println("创建dorm表成功")
	}
	//创建类型表
	err = db.AutoMigrate(&DgType{})
	if err != nil {
		println("创建type表失败")
	} else {
		println("创建type表成功")
	}
	//创建admin表
	err = db.AutoMigrate(&DgAdmin{})
	if err != nil {
		println("创建admin表失败")
	} else {
		println("创建admin表成功")
	}
	err = db.AutoMigrate(&DgComment{})
	if err != nil {
		println("创建Comment表失败")
	} else {
		println("创建Comment表成功")
	}
	err = db.AutoMigrate(&DgImages{})
	if err != nil {
		println("创建Images表失败")
	} else {
		println("创建Images表成功")
	}
	err = db.AutoMigrate(&DgPostLike{})
	if err != nil {
		println("创建PostLike表失败")
	} else {
		println("创建PostLike表成功")
	}
	//创建cmtlike表
	err = db.AutoMigrate(&DgCmtLike{})
	if err != nil {
		println("创建cmtlike表失败")
	} else {
		println("创建cmtlike表成功")
	}

	return
}
