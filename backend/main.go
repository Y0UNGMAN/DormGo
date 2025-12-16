package main

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/router"
	"github.com/Y0UNGMAN/DormGo/backend/settings"
	"github.com/Y0UNGMAN/DormGo/backend/utils"
)

func main() {
	//初始化雪花算法
	if err := utils.Init("2025-11-01", 1); err != nil {
		fmt.Printf("Init snowflake failed, err:%v\n", err)
		return // 如果初始化失败，直接退出程序，因为后面生成不了ID
	}
	fmt.Println("雪花算法初始化成功！")
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//2.初始化日志

	//3.初始化MySql连接
	if err := model.Init(); err != nil {
		fmt.Printf("init model failed, err:%v\n", err)
		return
	}
	//4.初始化Redis连接
	//if err := redis.Init(); err != nil {
	//	fmt.Printf("init redis failed, err:%v\n", err)
	//	return
	//}
	//5.注册路由
	r := router.App()

	//6.启动服务
	r.Run(":8080")
}
