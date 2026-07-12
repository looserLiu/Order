# 部署指南

本文档介绍如何将智慧记账应用部署到生产环境。

## 目录

- [环境要求](#环境要求)
- [Docker 部署](#docker-部署)
- [手动部署](#手动部署)
- [Kubernetes 部署](#kubernetes-部署)
- [CI/CD 配置](#cicd-配置)
- [监控与日志](#监控与日志)

## 环境要求

### 最低配置

- CPU: 2 核
- 内存: 4 GB
- 磁盘: 20 GB SSD
- 带宽: 5 Mbps

### 软件依赖

- Docker 20.10+
- Docker Compose 2.0+
- PostgreSQL 15
- Redis 7
- Nginx 1.20+

## Docker 部署

### 1. 克隆代码

```bash
git clone https://github.com/looserLiu/order.git
cd order
```

### 2. 配置环境变量

```bash
# 后端配置
cp backend/.env.example backend/.env
vim backend/.env
```

```env
# backend/.env
DB_HOST=db
DB_PORT=5432
DB_USER=ledger_user
DB_PASSWORD=your_secure_password
DB_NAME=ledger_db
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your_jwt_secret_key
API_PORT=8080
MINIO_ENDPOINT=minio:9000
MINIO_USER=minioadmin
MINIO_PASSWORD=your_minio_password
```

```bash
# 前端配置
cp frontend/.env.example frontend/.env
vim frontend/.env
```

```env
# frontend/.env
VITE_API_BASE_URL=/api/v1
VITE_APP_TITLE=智慧记账
```

### 3. 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 4. 数据库迁移

```bash
# 进入后端容器执行迁移
docker-compose exec backend go run cmd/api/main.go migrate
```

### 5. 验证部署

```bash
# 健康检查
curl http://localhost:8080/health

# 访问前端
open http://localhost:5173
```

## 手动部署

### 后端部署

#### 1. 构建二进制文件

```bash
cd backend
go build -o bin/server cmd/api/main.go
```

#### 2. 配置生产环境

```bash
export GIN_MODE=release
export DB_HOST=your_db_host
export DB_PASSWORD=your_db_password
export JWT_SECRET=your_jwt_secret
```

#### 3. 使用 systemd 管理

创建 `/etc/systemd/system/ledger-backend.service`:

```ini
[Unit]
Description=Ledger Backend Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=ledger
WorkingDirectory=/opt/ledger/backend
ExecStart=/opt/ledger/backend/bin/server
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable ledger-backend
sudo systemctl start ledger-backend
```

### 前端部署

#### 1. 构建静态文件

```bash
cd frontend
npm run build
```

#### 2. 配置 Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /opt/ledger/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# 测试配置
sudo nginx -t

# 重新加载
sudo nginx -s reload
```

## Kubernetes 部署

### 1. 创建命名空间

```bash
kubectl create namespace ledger
```

### 2. 创建 Secret

```bash
kubectl create secret generic ledger-secrets \
  --namespace ledger \
  --from-literal=db-password=your_db_password \
  --from-literal=jwt-secret=your_jwt_secret \
  --from-literal=minio-password=your_minio_password
```

### 3. 部署后端

`k8s/backend-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ledger-backend
  namespace: ledger
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ledger-backend
  template:
    metadata:
      labels:
        app: ledger-backend
    spec:
      containers:
        - name: backend
          image: ledger/backend:latest
          ports:
            - containerPort: 8080
          envFrom:
            - secretRef:
                name: ledger-secrets
          resources:
            requests:
              cpu: 100m
              memory: 256Mi
            limits:
              cpu: 500m
              memory: 512Mi
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: ledger-backend
  namespace: ledger
spec:
  selector:
    app: ledger-backend
  ports:
    - port: 8080
      targetPort: 8080
  type: ClusterIP
```

```bash
kubectl apply -f k8s/backend-deployment.yaml
```

### 4. 部署前端

`k8s/frontend-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ledger-frontend
  namespace: ledger
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ledger-frontend
  template:
    metadata:
      labels:
        app: ledger-frontend
    spec:
      containers:
        - name: frontend
          image: ledger/frontend:latest
          ports:
            - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: ledger-frontend
  namespace: ledger
spec:
  selector:
    app: ledger-frontend
  ports:
    - port: 80
      targetPort: 80
  type: ClusterIP
```

```bash
kubectl apply -f k8s/frontend-deployment.yaml
```

### 5. 配置 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ledger-ingress
  namespace: ledger
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  rules:
    - host: your-domain.com
      http:
        paths:
          - path: /api/(.*)
            pathType: Prefix
            backend:
              service:
                name: ledger-backend
                port:
                  number: 8080
          - path: /(.*)
            pathType: Prefix
            backend:
              service:
                name: ledger-frontend
                port:
                  number: 80
```

```bash
kubectl apply -f k8s/ingress.yaml
```

## CI/CD 配置

### GitHub Actions

`.github/workflows/ci.yml`:

```yaml
name: CI/CD

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  backend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: "1.21"
      - name: Run tests
        run: cd backend && go test -v ./test/...

  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: "18"
      - name: Install dependencies
        run: cd frontend && npm ci
      - name: Run tests
        run: cd frontend && npm run test:run

  build:
    needs: [backend-test, frontend-test]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build backend image
        run: docker build -t ledger/backend ./backend
      - name: Build frontend image
        run: docker build -t ledger/frontend ./frontend
```

## 监控与日志

### 日志收集

```bash
# 查看后端日志
docker-compose logs -f backend

# 查看前端日志
docker-compose logs -f frontend
```

### 健康检查

```bash
# 后端健康检查
curl http://localhost:8080/health

# 数据库健康检查
docker-compose exec db pg_isready -U ledger_user
```

### 性能监控

推荐使用以下工具：

- Prometheus + Grafana (指标监控)
- ELK Stack (日志分析)
- Sentry (错误追踪)

## 备份与恢复

### 数据库备份

```bash
# 创建备份
docker-compose exec db pg_dump -U ledger_user ledger_db > backup_$(date +%Y%m%d).sql

# 恢复备份
docker-compose exec -T db psql -U ledger_user ledger_db < backup_20240101.sql
```

### 文件备份

```bash
# 备份上传文件
docker-compose exec minio mc mirror /data /backup/minio
```

## 故障排查

### 常见问题

1. **数据库连接失败**
   - 检查 `DB_HOST` 和 `DB_PORT` 配置
   - 确认数据库服务已启动
   - 检查防火墙设置

2. **JWT 验证失败**
   - 确认 `JWT_SECRET` 配置一致
   - 检查 Token 是否过期

3. **前端无法访问 API**
   - 检查 Nginx 反向代理配置
   - 确认 `VITE_API_BASE_URL` 配置正确

## 安全建议

- 使用 HTTPS (配置 SSL 证书)
- 定期更新依赖包
- 使用强密码和密钥
- 启用防火墙
- 定期备份数据
- 限制 API 访问频率
