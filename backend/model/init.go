package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Database() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/dormgo_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		println("数据库连接失败")
	} else {
		println("数据库连接成功")
	}
	//创建用户表
	err = db.AutoMigrate(&dg_user{})
	if err != nil {
		println("创建用户表失败")
	} else {
		println("创建用户表成功")
	}
	//创建post表
	err = db.AutoMigrate(&dg_post{})
	if err != nil {
		println("创建post表失败")
	} else {
		println("创建post表成功")
	}
	DB = db

}
