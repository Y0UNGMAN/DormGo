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
	err = db.AutoMigrate(&dgUser{})
	if err != nil {
		println("创建用户表失败")
	} else {
		println("创建用户表成功")
	}
	//创建post表
	err = db.AutoMigrate(&dgPost{})
	if err != nil {
		println("创建post表失败")
	} else {
		println("创建post表成功")
	}
	err = db.AutoMigrate(&dgDorm{})
	if err != nil {
		println("创建dorm表失败")
	} else {
		println("创建dorm表成功")
	}
	err = db.AutoMigrate(&dgType{})
	if err != nil {
		println("创建type表失败")
	} else {
		println("创建type表成功")
	}
	DB = db

}
