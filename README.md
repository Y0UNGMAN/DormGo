# DormGo 寝友 Go

DormGo 是一个面向校园宿舍生活的互助社区系统，支持帖子发布、宿舍楼筛选、评论点赞、报名、私信通知、管理员审核和 AI 任务助手。项目采用 Vue 3 + Go Gin + MySQL + Redis，并引入 Python Agent 服务完成自然语言任务规划、帖子语义重排和安全工具调用。

## 功能概览

- 校园互助帖子：发布帖子、图片上传、分类筛选、宿舍楼筛选、收藏、点赞、评论。
- 报名协作：支持限时/限额报名，报名动作进入后端事务，避免绕过业务规则。
- 消息通知：私信会话、未读数、报名通知和系统通知。
- 管理后台：用户管理、内容审核、楼栋管理、类型配置、敏感词管理、违规处理、统计页面。
- 内容安全：敏感词过滤与 AI 内容审核结合，降低违规内容进入社区的风险。
- AI 任务助手：用户可以用自然语言查找帖子、筛选可报名活动，并在二次确认后完成报名。

## Agent 架构

DormGo 的 Agent 不是单纯聊天机器人，而是一个受控的工具调用系统：

```text
Vue Agent 页面
  ↓
Go /api/v1/agent/chat
  ↓
Python Agent Service
  ├─ LLM Planner：将自然语言转换成结构化任务计划
  ├─ Go Tool Client：调用 Go 内部工具
  ├─ LLM Reranker：对 Go 召回的帖子进行语义重排
  └─ Session Context：读取和写回最近候选与对话上下文
  ↓
Go Internal Tools
  ├─ search_posts
  ├─ search_signup_posts
  ├─ prepare_signup
  ├─ prepare_signup_from_candidate
  ├─ confirm_signup
  ├─ get_session
  └─ save_session
```

关键边界：

- Python Agent 不直接连接 MySQL 或 Redis。
- 所有报名、权限、事务和状态写入都由 Go 后端执行。
- 带副作用的动作必须经过待确认 action，再由用户确认。
- 帖子检索采用“数据库宽召回 + LLM 语义重排”，不是固定词表匹配。

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 前端 | Vue 3, TypeScript, Pinia, Vue Router, Axios, Element Plus |
| 后端 | Go, Gin, GORM, JWT |
| 数据存储 | MySQL, Redis |
| AI 能力 | DeepSeek/OpenAI-compatible Chat API, Python Agent Service |
| 文件存储 | 阿里云 OSS，本地静态资源回退 |

## 项目结构

```text
DormGo
├── backend/                 # Go Gin 后端
│   ├── controller/           # HTTP 控制器
│   ├── logic/                # 业务逻辑与 Agent 编排
│   ├── model/                # GORM 模型与数据库访问
│   ├── router/               # 路由注册
│   ├── ai/                   # AI 调用与内容审核
│   └── config.yaml           # 本地开发配置
├── frontend/                # Vue 3 前端
│   └── src/views/            # 用户端与管理端页面
├── agent-service/           # Python Agent 服务
│   ├── main.py               # Agent HTTP 服务入口
│   ├── planner.py            # LLM Planner
│   ├── reranker.py           # LLM 语义重排
│   ├── tools.py              # Go 内部工具客户端
│   └── config.py             # 读取 Go 配置中的模型参数
└── mysql/                   # 数据库初始化 SQL
```

## 本地运行

### 1. 准备数据库

创建 MySQL 数据库：

```sql
CREATE DATABASE dormgo_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

导入初始化 SQL：

```bash
mysql -u root -p dormgo_db < mysql/dormgo_db.sql
```

### 2. 启动后端

配置 `backend/config.yaml` 中的 MySQL、Redis、OSS 和模型 API 参数，然后启动：

```bash
cd backend
go run main.go
```

默认地址：

```text
http://127.0.0.1:8080
```

### 3. 启动 Agent 服务

Agent 服务会读取 `backend/config.yaml` 里的：

```text
api.ApiUrl
api.ApiKey
api.ApiModel
```

启动：

```bash
python3 agent-service/main.py
```

默认地址：

```text
http://127.0.0.1:8090
```

### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

## 测试与检查

Python Agent 单元测试：

```bash
python3 -m unittest discover agent-service
```

前端构建：

```bash
cd frontend
npm run build
```

后端测试：

```bash
cd backend
go test ./...
```

## 当前状态

已实现核心社区功能、管理后台、内容审核、消息通知与 AI 任务助手。Agent 当前聚焦于帖子检索、语义重排、候选上下文保存和安全报名确认；校园规章问答、知识库 RAG、搜索引擎接入可在后续作为独立模块扩展。
