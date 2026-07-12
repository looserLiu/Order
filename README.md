# 智慧记账应用

一个功能完整的个人/家庭记账应用，支持多账户、多分类、预算管理、报表分析、资产管理等功能。

## 技术栈

### 前端

- React 18 + TypeScript
- Vite (构建工具)
- TailwindCSS (样式)
- Zustand (状态管理)
- React Query (数据获取)
- Recharts (图表)

### 后端

- Go 1.21+
- Gin (Web 框架)
- GORM (ORM)
- PostgreSQL (数据库)
- Redis (缓存)
- JWT (认证)

### 基础设施

- Docker + Docker Compose
- Nginx (反向代理)
- MinIO (对象存储)

## 项目结构

```
.
├── backend/                 # 后端服务
│   ├── cmd/api/            # API 服务入口
│   ├── internal/
│   │   ├── config/         # 配置管理
│   │   ├── handlers/       # HTTP 处理器
│   │   ├── models/         # 数据模型
│   │   ├── repositories/   # 数据访问层
│   │   ├── services/       # 业务逻辑层
│   │   └── middleware/     # 中间件
│   ├── pkg/                # 公共包
│   ├── docs/               # API 文档
│   └── test/               # 单元测试
├── frontend/               # 前端应用
│   ├── src/
│   │   ├── components/     # 公共组件
│   │   ├── pages/          # 页面组件
│   │   ├── services/       # API 服务
│   │   ├── stores/         # 状态管理
│   │   └── __tests__/      # 单元测试
│   └── vite.config.ts
├── docs/                   # 项目文档
├── docker-compose.yml      # Docker 编排
└── ARCHITECTURE.md         # 架构设计文档
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15 (或 Docker)
- Redis 7 (或 Docker)

### 后端开发

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息

# 运行数据库迁移
go run cmd/api/main.go migrate

# 启动开发服务器
go run cmd/api/main.go
```

### 前端开发

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置 API 地址

# 启动开发服务器
npm run dev
```

### Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
```

## 环境变量

### 后端 (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=ledger_user
DB_PASSWORD=ledger_pass
DB_NAME=ledger_db
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key
API_PORT=8080
```

### 前端 (.env)

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=智慧记账
```

## API 文档

API 文档使用 OpenAPI 3.0 规范，位于 [`backend/docs/api.yaml`](backend/docs/api.yaml)。

可以使用以下工具查看：

- Swagger UI
- Postman (导入 YAML 文件)
- Redoc

## 测试

### 后端测试

```bash
cd backend
go test -v ./test/...
```

### 前端测试

```bash
cd frontend
npm run test
```

## 代码规范

### 后端

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行静态检查

### 前端

- 遵循 ESLint 规范
- 使用 Prettier 格式化代码
- 组件使用函数式组件 + Hooks

## 架构设计

详细的架构设计请参考 [ARCHITECTURE.md](ARCHITECTURE.md)。

### 分层架构

```
┌─────────────────────────────────────────┐
│            HTTP Handler 层                │
│  (Gin Routes, 参数验证, 响应封装)         │
└───────────────────┬─────────────────────┘
                    │
┌───────────────────▼─────────────────────┐
│            Service 业务层                 │
│  (业务逻辑, 事务管理, 权限验证)           │
└───────────────────┬─────────────────────┘
                    │
┌───────────────────▼─────────────────────┐
│         Repository 数据访问层             │
│  (CRUD 操作, 查询构建)                   │
└───────────────────┬─────────────────────┘
                    │
┌───────────────────▼─────────────────────┐
│              Model 数据模型               │
│  (GORM Models, 数据库映射)               │
└─────────────────────────────────────────┘
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件
