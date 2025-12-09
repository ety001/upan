# VPJ - 网络闪存

一个基于 Go + Gin + React + Vite 的文件上传/下载服务。

## 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: React + TypeScript + Vite + Tailwind CSS 3
- **容器化**: Docker + Docker Compose

## 功能特性

- 文件上传（支持最大文件大小限制）
- 通过 6 位提取码下载文件
- 自动清理过期文件
- 响应式 UI 设计

## 项目结构

```
vpj/
├── backend/              # Go 后端
│   ├── cmd/server/      # 应用入口
│   └── internal/        # 内部包
│       ├── config/       # 配置管理
│       ├── handlers/     # HTTP 处理器
│       ├── models/       # 数据模型
│       └── tasks/        # 后台任务
├── frontend/             # React 前端
│   └── src/
│       ├── components/   # React 组件
│       └── api/          # API 客户端
├── docker-compose.yml    # Docker Compose 配置
└── .env.example         # 环境变量示例
```

## 环境变量

创建 `.env` 文件（参考 `.env.example`）：

```env
PORT=8080
FILE_MAX_SIZE=30          # 文件大小限制（MB）
FILE_EXPIRE_TIME=6        # 文件过期时间（小时）
STORAGE_PATH=./storage/files
DB_PATH=./storage/db/database.db
```

## 开发

### 后端开发

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

## 使用 Docker

### 构建和运行

项目已配置为单容器部署，使用 supervisor 管理前后端进程：

```bash
# 构建并启动服务
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 查看 supervisor 状态
docker-compose exec app supervisorctl status

# 停止服务
docker-compose down
```

### 访问

- 应用: http://localhost
- API: http://localhost/api

### 单容器部署

前后端已打包为一个 Docker 镜像，使用 supervisor 管理：
- 后端服务运行在容器内的 8080 端口
- 前端由 nginx 提供服务，并代理后端 API
- 所有进程由 supervisor 统一管理，自动重启

## API 接口

### 获取配置

```
GET /api/config
```

响应:
```json
{
  "file_max_size": 30,
  "file_expire_time": 6
}
```

### 上传文件

```
POST /api/upload
Content-Type: multipart/form-data

Form Data:
  o: <file>
```

响应:
```json
{
  "status": true,
  "code": "abc123",
  "expired_at": 1234567890
}
```

### 下载文件

```
GET /api/file/:code
```

## 许可证

MIT

