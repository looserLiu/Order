# 智慧记账应用系统架构设计

## 一、系统架构总览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              客户端层                                        │
├─────────────────────┬─────────────────────┬───────────────────────────────┤
│    PC Web 端        │    移动 App 端       │      PWA/小程序(可选)          │
│  React + TS         │  React Native       │                                │
│  TailwindCSS        │  Expo                │                                │
└─────────┬───────────┴──────────┬───────────┴───────────────┬───────────────┘
          │                      │                             │
          └──────────────────────┼─────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API 网关层                                      │
│                    Nginx / Traefik (负载均衡、SSL)                          │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              应用服务层 (Go Microservices)                   │
├──────────────┬──────────────┬──────────────┬──────────────┬────────────────┤
│  用户服务     │  账目服务     │  报表服务     │  预算服务     │  资产服务       │
│  user-svc    │  ledger-svc  │  report-svc  │  budget-svc  │  asset-svc     │
├──────────────┴──────────────┴──────────────┴──────────────┴────────────────┤
│                              公共组件层                                      │
│         ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌─────────────┐   │
│         │  Auth      │  │  Notify    │  │  Search    │  │  Export     │   │
│         │  Service   │  │  Service   │  │  Service   │  │  Service    │   │
│         └────────────┘  └────────────┘  └────────────┘  └─────────────┘   │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              数据存储层                                      │
├──────────────────┬──────────────────┬───────────────────┬─────────────────┤
│   PostgreSQL     │    Redis          │   MinIO           │   Elasticsearch │
│   (主数据库)     │    (缓存/Session) │   (文件存储)       │   (搜索分析)     │
└──────────────────┴──────────────────┴───────────────────┴─────────────────┘
```

## 二、技术栈选型

### 前端
- **PC端**: React 18 + TypeScript + Vite + TailwindCSS + React Query
- **移动端**: React Native + Expo
- **状态管理**: Zustand / Redux Toolkit
- **图表**: Recharts / ECharts
- **UI组件**: Headless UI / Radix UI

### 后端
- **语言**: Go 1.21+
- **框架**: Gin / Fiber
- **ORM**: GORM
- **认证**: JWT + OAuth2
- **缓存**: Redis
- **任务队列**: Asynq (Redis-based)
- **配置**: Viper

### 基础设施
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **对象存储**: MinIO (本地兼容S3)
- **搜索**: Elasticsearch 8
- **容器**: Docker + Docker Compose

## 三、核心功能模块

### 3.1 用户模块
- [x] 用户注册/登录 (邮箱、手机号)
- [x] 第三方登录 (Google、Apple、WeChat)
- [x] 忘记密码/重置密码
- [x] 个人资料管理
- [x] 多设备登录管理
- [x] 用户偏好设置 (主题、语言、货币)

### 3.2 账目管理模块
- [x] 记账流水 (收入/支出/转账)
- [x] 多账户支持 (银行卡、信用卡、现金、支付宝、微信等)
- [x] 分类管理 (自定义分类、层级分类)
- [x] 标签系统
- [x] 商家/收款方管理
- [x] 备注/附件上传
- [x] 批量操作 (批量删除、批量修改)
- [x] 记账提醒 (定时提醒、自动提醒)
- [x] 语音记账
- [x] OCR 拍照识别 (发票/小票)

### 3.3 资产管理模块
- [x] 资产账户管理
- [x] 资产变动记录
- [x] 资产统计报表
- [x] 债务管理 (借出/借入)
- [x] 投资账户 (股票、基金、理财)

### 3.4 预算管理模块
- [x] 月度/年度预算
- [x] 分类预算
- [x] 预算提醒
- [x] 超支预警
- [x] 预算执行分析

### 3.5 报表分析模块
- [x] 日/周/月/年报表
- [x] 收支趋势图
- [x] 分类占比饼图
- [x] 收入支出对比
- [x] 多维度分析 (按时间/分类/账户/商家)
- [x] 自定义报表
- [x] 数据导出 (Excel、CSV、PDF)

### 3.6 高级功能
- [x] 数据同步 (多设备)
- [x] 数据导入 (支持其他记账软件数据迁移)
- [x] 数据备份与恢复
- [x] 周期记账 (重复账单)
- [x] 多人协作 (家庭/团队记账)
- [x] 权限管理
- [x] 消息通知

## 四、数据库设计

### 4.1 核心表结构

```sql
-- 用户表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    nickname VARCHAR(100),
    avatar_url VARCHAR(500),
    currency VARCHAR(10) DEFAULT 'CNY',
    timezone VARCHAR(50) DEFAULT 'Asia/Shanghai',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 账户表
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL, -- bank, credit, cash, wechat, alipay, investment
    balance DECIMAL(15,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'CNY',
    icon VARCHAR(100),
    color VARCHAR(20),
    is_default BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 分类表
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    parent_id UUID REFERENCES categories(id),
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(100),
    color VARCHAR(20),
    type VARCHAR(20) NOT NULL, -- income, expense, transfer
    sort_order INT DEFAULT 0,
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 记账流水表
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    account_id UUID REFERENCES accounts(id),
    target_account_id UUID REFERENCES accounts(id),
    category_id UUID REFERENCES categories(id),
    type VARCHAR(20) NOT NULL, -- income, expense, transfer
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'CNY',
    exchange_rate DECIMAL(10,4) DEFAULT 1,
    tags UUID[], -- array of tag ids
    merchant VARCHAR(255),
    note TEXT,
    attachment_urls TEXT[],
    bill_date DATE NOT NULL,
    is_recurring BOOLEAN DEFAULT false,
    recurring_rule JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 预算表
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    category_id UUID REFERENCES categories(id),
    amount DECIMAL(15,2) NOT NULL,
    period VARCHAR(20) NOT NULL, -- monthly, yearly
    start_date DATE NOT NULL,
    end_date DATE,
    alert_threshold DECIMAL(5,2) DEFAULT 0.8,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 资产变动表
CREATE TABLE asset_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    account_id UUID REFERENCES accounts(id),
    asset_type VARCHAR(20) NOT NULL, -- debt_owed, debt_owing, investment
    related_user VARCHAR(255), -- for debts
    amount DECIMAL(15,2) NOT NULL,
    interest_rate DECIMAL(5,4),
    start_date DATE,
    end_date DATE,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 标签表
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(50) NOT NULL,
    color VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 设备表
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    device_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(100),
    device_type VARCHAR(20), -- ios, android, web, desktop
    push_token VARCHAR(500),
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 团队/家庭成员表
CREATE TABLE family_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID NOT NULL,
    user_id UUID REFERENCES users(id),
    role VARCHAR(20) NOT NULL, -- owner, admin, member
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 团队/家庭表
CREATE TABLE families (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    owner_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 团队成员账目表 (共享账本)
CREATE TABLE family_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID REFERENCES families(id),
    user_id UUID REFERENCES users(id),
    account_id UUID REFERENCES accounts(id),
    category_id UUID REFERENCES categories(id),
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    note TEXT,
    bill_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## 五、API 接口设计

### 5.1 认证接口
```
POST   /api/v1/auth/register        # 用户注册
POST   /api/v1/auth/login           # 用户登录
POST   /api/v1/auth/logout          # 登出
POST   /api/v1/auth/refresh         # 刷新Token
POST   /api/v1/auth/forgot-password # 忘记密码
POST   /api/v1/auth/reset-password  # 重置密码
POST   /api/v1/auth/oauth/{provider} # 第三方登录
```

### 5.2 用户接口
```
GET    /api/v1/users/me             # 获取当前用户信息
PUT    /api/v1/users/me             # 更新用户信息
DELETE /api/v1/users/me             # 删除账户
GET    /api/v1/users/devices        # 获取设备列表
DELETE /api/v1/users/devices/:id    # 删除设备
```

### 5.3 账户接口
```
GET    /api/v1/accounts             # 获取账户列表
POST   /api/v1/accounts             # 创建账户
GET    /api/v1/accounts/:id         # 获取账户详情
PUT    /api/v1/accounts/:id         # 更新账户
DELETE /api/v1/accounts/:id         # 删除账户
GET    /api/v1/accounts/:id/balance # 获取账户余额
```

### 5.4 分类接口
```
GET    /api/v1/categories           # 获取分类列表
POST   /api/v1/categories           # 创建分类
PUT    /api/v1/categories/:id       # 更新分类
DELETE /api/v1/categories/:id       # 删除分类
GET    /api/v1/categories/tree      # 获取分类树形结构
```

### 5.5 账目接口
```
GET    /api/v1/transactions         # 获取账目列表 (支持分页、筛选)
POST   /api/v1/transactions         # 创建账目
GET    /api/v1/transactions/:id     # 获取账目详情
PUT    /api/v1/transactions/:id     # 更新账目
DELETE /api/v1/transactions/:id     # 删除账目
POST   /api/v1/transactions/batch   # 批量操作
POST   /api/v1/transactions/ocr     # OCR识别
POST   /api/v1/transactions/voice   # 语音记账
GET    /api/v1/transactions/export  # 导出账目
```

### 5.6 预算接口
```
GET    /api/v1/budgets              # 获取预算列表
POST   /api/v1/budgets              # 创建预算
PUT    /api/v1/budgets/:id          # 更新预算
DELETE /api/v1/budgets/:id          # 删除预算
GET    /api/v1/budgets/:id/progress # 获取预算执行进度
```

### 5.7 报表接口
```
GET    /api/v1/reports/summary      # 收支摘要
GET    /api/v1/reports/trend        # 趋势分析
GET    /api/v1/reports/category     # 分类统计
GET    /api/v1/reports/account      # 账户统计
GET    /api/v1/reports/compare      # 对比分析
GET    /api/v1/reports/chart        # 图表数据
GET    /api/v1/reports/export       # 导出报表
```

### 5.8 资产接口
```
GET    /api/v1/assets               # 获取资产列表
POST   /api/v1/assets               # 创建资产记录
PUT    /api/v1/assets/:id           # 更新资产记录
DELETE /api/v1/assets/:id           # 删除资产记录
GET    /api/v1/assets/summary       # 资产汇总
```

## 六、前端项目结构

```
frontend/
├── apps/
│   ├── web/                    # PC Web端
│   │   ├── src/
│   │   │   ├── components/    # 公共组件
│   │   │   ├── features/      # 功能模块
│   │   │   ├── hooks/         # 自定义Hooks
│   │   │   ├── layouts/       # 布局组件
│   │   │   ├── pages/         # 页面组件
│   │   │   ├── services/      # API服务
│   │   │   ├── stores/        # 状态管理
│   │   │   ├── types/         # TypeScript类型
│   │   │   ├── utils/         # 工具函数
│   │   │   └── App.tsx
│   │   └── vite.config.ts
│   │
│   └── mobile/                 # 移动App端 (React Native)
│       ├── src/
│       │   ├── components/
│       │   ├── screens/
│       │   ├── navigation/
│       │   ├── services/
│       │   ├── stores/
│       │   └── utils/
│       └── app.json
│
└── packages/
    ├── ui/                     # 共享UI组件库
    ├── api-client/             # API客户端
    ├── types/                  # 共享类型定义
    └── utils/                  # 共享工具函数
```

## 七、后端项目结构

```
backend/
├── cmd/
│   ├── api/                    # HTTP API服务
│   │   └── main.go
│   └── worker/                 # 异步任务服务
│       └── main.go
│
├── internal/
│   ├── config/                 # 配置管理
│   ├── middleware/            # 中间件
│   ├── models/                 # 数据模型
│   ├── handlers/               # HTTP处理器
│   ├── services/              # 业务逻辑
│   ├── repositories/          # 数据访问层
│   ├── cache/                 # 缓存层
│   └── utils/                  # 工具函数
│
├── pkg/                        # 公共包
│   ├── auth/                  # 认证
│   ├── errors/                # 错误处理
│   ├── response/              # 响应封装
│   └── validator/             # 参数验证
│
├── migrations/                 # 数据库迁移
├── api/                        # OpenAPI/Swagger文档
└── docker/
    └── Dockerfile
```

## 八、Docker 部署配置

### 8.1 docker-compose.yml

```yaml
version: '3.9'

services:
  # PostgreSQL 数据库
  db:
    image: postgres:15-alpine
    container_name: ledger_db
    environment:
      POSTGRES_USER: ${DB_USER:-ledger_user}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-ledger_pass}
      POSTGRES_DB: ${DB_NAME:-ledger_db}
    volumes:
      - db-data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ledger_user"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: ledger_redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries:5

  # MinIO 对象存储
  minio:
    image: minio/minio:latest
    container_name: ledger_minio
    environment:
      MINIO_ROOT_USER: ${MINIO_USER:-minioadmin}
      MINIO_ROOT_PASSWORD: ${MINIO_PASSWORD:-minioadmin}
    volumes:
      - minio-data:/data
    ports:
      - "9000:9000"
      - "9001:9001"
    command: server /data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  # Elasticsearch (可选，用于搜索)
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    container_name: ledger_es
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    volumes:
      - es-data:/usr/share/elasticsearch/data
    ports:
      - "9200:9200"
    healthcheck:
      test: ["CMD-SHELL", "curl -f http://localhost:9200/_cluster/health || exit 1"]
      interval: 30s
      timeout: 10s
      retries: 5

  # 后端API服务
  backend:
    build:
      context: ./backend
      dockerfile: docker/Dockerfile
    container_name: ledger_backend
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: ${DB_USER:-ledger_user}
      DB_PASSWORD: ${DB_PASSWORD:-ledger_pass}
      DB_NAME: ${DB_NAME:-ledger_db}
      REDIS_HOST: redis
      REDIS_PORT: 6379
      MINIO_ENDPOINT: minio:9000
      MINIO_USER: ${MINIO_USER:-minioadmin}
      MINIO_PASSWORD: ${MINIO_PASSWORD:-minioadmin}
      JWT_SECRET: ${JWT_SECRET:-your-secret-key}
      API_PORT: 8080
    ports:
      - "8000:8080"
    networks:
      - ledgernet
    restart: unless-stopped

  # 前端Web服务
  frontend:
    build:
      context: ./frontend/apps/web
      dockerfile: ../../docker/Dockerfile
    container_name: ledger_frontend
    depends_on:
      - backend
    environment:
      VITE_BACKEND_URL: http://backend:8080
      VITE_API_BASE_URL: /api
    ports:
      - "5173:80"
    networks:
      - ledgernet
    restart: unless-stopped

  # Nginx 反向代理 (可选)
  nginx:
    image: nginx:alpine
    container_name: ledger_nginx
    depends_on:
      - backend
      - frontend
    volumes:
      - ./docker/nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "80:80"
      - "443:443"
    networks:
      - ledgernet
    restart: unless-stopped

volumes:
  db-data:
  redis-data:
  minio-data:
  es-data:

networks:
  ledgernet:
    driver: bridge
```

## 九、核心业务流程

### 9.1 记账流程
```
用户输入 → 账户选择 → 金额输入 → 分类选择 → (可选)标签 → (可选)商家 → (可选)备注 → (可选)附件 → 保存 → 更新账户余额 → 返回结果
```

### 9.2 OCR识别流程
```
拍照/选择图片 → 上传到MinIO → 调用OCR服务 → 提取关键信息(金额、商家、日期) → 自动填充表单 → 用户确认 → 保存
```

### 9.3 报表生成流程
```
选择时间范围 → 选择筛选条件 → 查询数据库 → 数据聚合计算 → 生成图表数据 → 返回前端展示
```

## 十、安全考虑

- [x] JWT Token 认证
- [x] 密码加密存储 (bcrypt)
- [x] HTTPS 传输加密
- [x] API 频率限制
- [x] SQL 注入防护
- [x] XSS 防护
- [x] CORS 配置
- [x] 数据备份策略

## 十一、开发阶段划分

### Phase 1: 核心功能 (MVP)
- 用户注册/登录
- 账户管理
- 分类管理
- 记账流水CRUD
- 基础报表

### Phase 2: 增强功能
- 预算管理
- 数据导出
- 标签系统
- 附件上传
- 周期记账

### Phase 3: 高级功能
- OCR识别
- 语音记账
- 多人协作
- 资产/债务管理
- 数据导入

### Phase 4: 优化与扩展
- 性能优化
- 移动端适配
- 数据分析
- 国际化

---

**文档版本**: v1.0
**创建日期**: 2026-02-28
**技术栈**: React + TypeScript + Vite + TailwindCSS + Go + PostgreSQL + Redis + Docker
