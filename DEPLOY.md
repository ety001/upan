# 部署说明

## 单容器部署

项目已配置为单容器部署，前后端打包在一个 Docker 镜像中，使用 supervisor 管理进程。

## 快速开始

### 1. 构建镜像

```bash
docker build -t vpj-app .
```

### 2. 运行容器

```bash
docker run -d \
  --name vpj-app \
  -p 80:80 \
  -v $(pwd)/storage/files:/app/storage/files \
  -v $(pwd)/storage/db:/app/storage/db \
  -e FILE_MAX_SIZE=30 \
  -e FILE_EXPIRE_TIME=6 \
  vpj-app
```

### 3. 使用 Docker Compose

```bash
# 启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 查看进程状态
docker-compose exec app supervisorctl status

# 停止
docker-compose down
```

## 容器内进程管理

容器使用 supervisor 管理以下进程：

1. **backend**: Go 后端服务（端口 8080）
2. **nginx**: 前端服务和 API 代理（端口 80）

### Supervisor 命令

```bash
# 进入容器
docker-compose exec app sh

# 查看进程状态
supervisorctl status

# 重启后端服务
supervisorctl restart backend

# 重启 nginx
supervisorctl restart nginx

# 查看日志
tail -f /var/log/supervisor/backend.out.log
tail -f /var/log/supervisor/nginx.out.log
```

## 环境变量

可通过环境变量配置：

- `PORT`: 后端服务端口（默认 8080）
- `FILE_MAX_SIZE`: 文件大小限制，单位 MB（默认 30）
- `FILE_EXPIRE_TIME`: 文件过期时间，单位小时（默认 6）
- `STORAGE_PATH`: 文件存储路径（默认 /app/storage/files）
- `DB_PATH`: 数据库文件路径（默认 /app/storage/db/database.db）

## 数据持久化

建议挂载以下目录到宿主机：

- `/app/storage/files`: 上传的文件
- `/app/storage/db`: SQLite 数据库

## 健康检查

```bash
# 检查应用是否运行
curl http://localhost/api/config

# 检查健康状态
curl http://localhost/health
```

## 故障排查

### 查看所有日志

```bash
docker-compose logs -f app
```

### 查看 supervisor 日志

```bash
docker-compose exec app cat /var/log/supervisor/supervisord.log
```

### 重启所有服务

```bash
docker-compose exec app supervisorctl restart all
```

