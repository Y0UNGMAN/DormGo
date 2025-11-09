# 前后端搭建过程：
## 后端：
1.初始化一个go模块:go mod init github.com/Y0UNGMAN/DormGo/backend
2.安装gin框架：go get -u github.com/gin-gonic/gin
3.安装orm 以及 mysql 数据库相关： go get -u gorm.io/gorm
                                go get -u "gorm.io/driver/mysql"
本地创建数据库：dormgo_db
用户：root    密码：123456
4.安装跨域工具： go get github.com/gin-contrib/cors
5.启动后端： go run main.go    
## 前端：
1 安装 Node.js
2 创建vue项目：npm create vite@latest frontend
3 进入前端项目：cd frontend  安装依赖：npm install
4 安装路由插件：npm install vue-router@4
  安装pinia状态管理： npm install pinia
  安装axios请求库：   npm install axios
5 前端启动：先进入前端目录:cd frontend   启动：npm run dev

## 运行注意
clone下来在本地运行时，不需要“初始化go模块”以及“初始化go模块” ， 需要安装其他配置，再运行
