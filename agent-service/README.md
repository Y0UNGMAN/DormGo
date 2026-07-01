# DormGo Agent Service

Agent 服务负责调用大模型做任务规划，并通过 Go 内部工具执行搜索、待确认报名和确认报名。Python 不直接连接 MySQL 或 Redis。

启动：

```bash
python3 agent-service/main.py
```

启动时会读取 Go 后端配置：

```text
backend/config.yaml
api.ApiUrl
api.ApiKey
api.ApiModel
```

Go 后端默认会请求：

```text
http://127.0.0.1:8090/chat
```

当前支持：

- 普通帖子搜索：`search_posts`
- 可报名帖子搜索：`search_signup_posts`
- 为候选帖子生成待确认报名动作：`prepare_signup`
- 用户确认后调用 Go 完成报名：`confirm_signup`

检索链路：

```text
用户自然语言
-> LLM Planner 生成结构化计划
-> Go 工具宽召回候选帖子
-> LLM Reranker 对候选帖子语义打分和排序
-> Agent 返回相关结果或执行待确认动作
```

注意：代码不维护“学习=自习/看电影=拼票”这类固定扩展词表。候选相关性由 LLM Reranker 基于帖子标题、正文和用户查询判断。

执行边界：

- Python Agent 只负责规划和调用 Go 工具。
- Go 后端负责 JWT、内部工具鉴权、Redis 待确认动作、MySQL 查询和报名事务。
