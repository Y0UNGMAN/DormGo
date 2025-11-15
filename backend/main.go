package main

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/redis"
	"github.com/Y0UNGMAN/DormGo/backend/router"
	"github.com/Y0UNGMAN/DormGo/backend/settings"
)

func main() {
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
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	//5.注册路由
	r := router.App()
	//6.启动服务

	r.Run(":8080")
}
