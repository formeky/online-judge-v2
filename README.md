# Online Judge V2

基于 Go 的在线评测系统，支持多语言代码提交、自动编译与沙箱评测。

## 系统架构

采用生产者-消费者架构，包含两个独立组件：

- **Server**：提供 RESTful API，处理题目管理、代码提交，并将评测任务投递到消息队列
- **Worker**：从消息队列消费任务，在沙箱中执行代码并返回评测结果

```
用户提交代码 -> Server API -> RocketMQ -> Worker -> 沙箱执行 -> 结果写入数据库
```

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.21 |
| Web 框架 | Gin |
| 数据库 | MySQL 8.0 + GORM |
| 消息队列 | Apache RocketMQ 5.1.4 |
| 配置管理 | Viper |
| 日志 | Zap |
| 认证 | JWT |

## 支持的评测语言

| 语言 | 编译器/运行时 |
|------|-------------|
| C | gcc (C11) |
| C++ | g++ (C++17) |
| Java | OpenJDK 17 |
| Python | Python 3 |
| Go | Go 1.21+ |

## 目录结构

```
online-judge-v2/
├── cmd/
│   ├── server/main.go          # API 服务入口
│   └── worker/main.go          # 评测 Worker 入口
├── configs/
│   └── config.yaml             # 配置文件
├── deployments/
│   ├── Dockerfile.server
│   ├── Dockerfile.worker
│   └── docker-compose.yml      # 基础设施编排
├── internal/
│   ├── api/                    # HTTP API 层（路由、Handler、中间件）
│   ├── config/                 # 配置加载
│   ├── database/               # 数据库连接与迁移
│   ├── judge/                  # 评测核心引擎
│   ├── model/                  # 数据模型
│   ├── mq/                     # 消息队列生产者/消费者
│   ├── repository/             # 数据访问层
│   └── sandbox/                # 沙箱（Cgroup + Namespace，仅 Linux）
└── testdata/                   # 示例测试数据
```

## 快速开始

### 前置要求

- Go 1.21+
- Docker & Docker Compose
- Worker 需运行在 Linux 系统，并安装 gcc、g++、javac、python3、go

### 本地开发

```bash
# 1. 整理依赖
make tidy

# 2. 启动基础设施（MySQL + RocketMQ）
make docker-up

# 3. 启动 API 服务
make run-server

# 4. 新终端启动评测 Worker
make run-worker
```

### 编译运行

```bash
make build-server
make build-worker

./bin/server
./bin/worker
```

### Docker 部署

```bash
docker build -f deployments/Dockerfile.server -t oj-server:latest .
docker build -f deployments/Dockerfile.worker -t oj-worker:latest .

docker run -d --name oj-server -p 8080:8080 -v $(pwd)/configs:/configs oj-server:latest
# Worker 沙箱需要特权模式
docker run -d --name oj-worker --privileged -v $(pwd)/configs:/configs oj-worker:latest
```

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make tidy` | 整理依赖 |
| `make build-server` | 编译 API 服务到 `bin/server` |
| `make build-worker` | 编译 Worker 到 `bin/worker` |
| `make run-server` | 运行 API 服务 |
| `make run-worker` | 运行评测 Worker |
| `make docker-up` | 启动基础设施容器 |
| `make docker-down` | 停止容器 |
| `make docker-logs` | 查看容器日志 |

## API

### 路由

```
GET  /health

/api/v1/problems
├── GET    /                    # 题目列表（分页）
├── GET    /:id                 # 题目详情
├── POST   /                    # 创建题目 [Admin]
├── PUT    /:id                 # 更新题目 [Admin]
├── DELETE /:id                 # 删除题目 [Admin]
└── POST   /:id/testcases      # 添加测试用例 [Admin]

/api/v1/submissions [需登录]
├── POST   /                    # 提交代码
├── GET    /                    # 提交列表（分页）
└── GET    /:id                 # 提交详情
```

### 示例

```bash
# 健康检查
curl http://localhost:8080/health

# 提交代码
curl -X POST http://localhost:8080/api/v1/submissions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "problem_id": 1,
    "language": "python",
    "code": "a, b = map(int, input().split())\nprint(a + b)"
  }'
```

## 评测状态

| 状态 | 说明 |
|------|------|
| `pending` | 等待评测 |
| `compiling` | 编译中 |
| `running` | 运行中 |
| `accepted` | 通过 |
| `wrong_answer` | 答案错误 |
| `time_limit_exceeded` | 超时 |
| `memory_limit_exceeded` | 内存超限 |
| `runtime_error` | 运行时错误 |
| `compile_error` | 编译错误 |
| `system_error` | 系统错误 |

## 配置说明

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"

database:
  dsn: "oj:ojpasswd@tcp(127.0.0.1:3306)/oj?charset=utf8mb4&parseTime=True&loc=Local"

rocketmq:
  name_servers:
    - "127.0.0.1:9876"
  topic: "judge_submission"
  consumer_group: "judge_worker_group"

judge:
  work_dir: "/tmp/oj-judge"
  worker_concurrency: 4
  compile_timeout: 30

jwt:
  secret: "change-me-in-production"
  expire_hours: 24
```

> 生产环境请务必修改 `jwt.secret`。
