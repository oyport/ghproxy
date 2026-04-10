# GHProxy Docker 部署指南

## 📋 目录

- [快速开始](#快速开始)
- [部署方式](#部署方式)
- [配置说明](#配置说明)
- [环境变量](#环境变量)
- [数据持久化](#数据持久化)
- [健康检查](#健康检查)
- [性能优化](#性能优化)
- [故障排查](#故障排查)
- [高级配置](#高级配置)

## 🚀 快速开始

### 方式一：使用 Docker Compose（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/WJQSERVER-STUDIO/ghproxy.git
cd ghproxy

# 2. 创建必要的目录
mkdir -p ghproxy/log ghproxy/config

# 3. 下载配置文件
wget -O ghproxy/config/config.toml https://raw.githubusercontent.com/WJQSERVER-STUDIO/ghproxy/main/config/config.toml
wget -O ghproxy/config/blacklist.json https://raw.githubusercontent.com/WJQSERVER-STUDIO/ghproxy/main/config/blacklist.json
wget -O ghproxy/config/whitelist.json https://raw.githubusercontent.com/WJQSERVER-STUDIO/ghproxy/main/config/whitelist.json

# 4. 启动服务
cd docker/compose
docker-compose up -d

# 5. 查看日志
docker-compose logs -f

# 6. 检查状态
docker-compose ps
```

### 方式二：使用 Docker 命令

```bash
# 1. 拉取镜像
docker pull wjqserver/ghproxy:latest

# 2. 创建目录
mkdir -p /data/ghproxy/log /data/ghproxy/config

# 3. 运行容器
docker run -d \
  --name ghproxy \
  --restart always \
  -p 7210:8080 \
  -v /data/ghproxy/log:/data/ghproxy/log \
  -v /data/ghproxy/config:/data/ghproxy/config \
  -e TZ=Asia/Shanghai \
  wjqserver/ghproxy:latest

# 4. 查看日志
docker logs -f ghproxy
```

## 📦 部署方式

### 1. 使用预构建镜像（推荐）

```bash
# 拉取最新版本
docker pull wjqserver/ghproxy:latest

# 拉取特定版本
docker pull wjqserver/ghproxy:v4.3.4

# 拉取开发版本
docker pull wjqserver/ghproxy:dev
```

### 2. 自行构建镜像

```bash
# 克隆项目
git clone https://github.com/WJQSERVER-STUDIO/ghproxy.git
cd ghproxy

# 构建Release版本
docker build -f docker/dockerfile/release/Dockerfile -t ghproxy:latest .

# 构建Dev版本
docker build -f docker/dockerfile/dev/Dockerfile -t ghproxy:dev .

# 运行自定义镜像
docker run -d -p 7210:8080 --name ghproxy ghproxy:latest
```

### 3. 使用 Docker Compose 完整配置

创建 `docker-compose.yml`：

```yaml
version: '3.9'

services:
  ghproxy:
    image: 'wjqserver/ghproxy:latest'
    container_name: ghproxy
    restart: always
    
    ports:
      - '7210:8080'
    
    volumes:
      - ./ghproxy/log:/data/ghproxy/log
      - ./ghproxy/config:/data/ghproxy/config
    
    environment:
      - TZ=Asia/Shanghai
    
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/healthcheck"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s
    
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 128M
    
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    
    networks:
      - ghproxy-network

networks:
  ghproxy-network:
    driver: bridge
```

## ⚙️ 配置说明

### 1. 基础配置

编辑 `ghproxy/config/config.toml`：

```toml
[server]
host = "0.0.0.0"
port = 8080
sizeLimit = 125  # MB
cors = "*"
debug = false

[pages]
mode = "internal"
theme = "design"  # 选择主题

[log]
logFilePath = "/data/ghproxy/log/ghproxy.log"
maxLogSize = 5  # MB
level = "info"
```

### 2. 启用管理后台

```toml
[admin]
enabled = true
username = "admin"
password = "your_secure_password"  # 请修改
pathPrefix = "/admin"
```

### 3. Docker代理配置

```toml
[docker]
enabled = true
target = "dockerhub"  # 或 "ghcr"
auth = false
```

### 4. 速率限制

```toml
[rateLimit]
enabled = true
ratePerMinute = 180
burst = 5
```

## 🌍 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| TZ | 时区设置 | Asia/Shanghai |

**使用示例**：

```yaml
environment:
  - TZ=America/New_York
```

## 💾 数据持久化

### 重要目录

| 容器路径 | 说明 | 建议挂载 |
|---------|------|---------|
| `/data/ghproxy/log` | 日志文件 | ✅ 必须挂载 |
| `/data/ghproxy/config` | 配置文件 | ✅ 必须挂载 |
| `/data/www` | 静态文件 | ⚪ 可选 |

### 挂载示例

```yaml
volumes:
  # 日志持久化
  - ./ghproxy/log:/data/ghproxy/log
  
  # 配置持久化
  - ./ghproxy/config:/data/ghproxy/config
  
  # 自定义静态文件（可选）
  - ./ghproxy/www:/data/www
```

## 🏥 健康检查

### Dockerfile 中的健康检查

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/healthcheck || exit 1
```

### Docker Compose 中的健康检查

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/api/healthcheck"]
  interval: 30s      # 检查间隔
  timeout: 10s       # 超时时间
  retries: 3         # 重试次数
  start_period: 5s   # 启动等待时间
```

### 手动检查

```bash
# 检查容器健康状态
docker inspect --format='{{.State.Health.Status}}' ghproxy

# 检查API健康状态
curl http://localhost:7210/api/healthcheck
```

## ⚡ 性能优化

### 1. 资源限制

```yaml
deploy:
  resources:
    limits:
      cpus: '2'        # 最大CPU使用
      memory: 512M     # 最大内存使用
    reservations:
      cpus: '0.5'      # 保留CPU
      memory: 128M     # 保留内存
```

### 2. 日志轮转

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"   # 单个日志文件最大10MB
    max-file: "3"     # 保留3个日志文件
```

### 3. 网络优化

```yaml
networks:
  ghproxy-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### 4. 连接池配置

编辑 `config.toml`：

```toml
[httpc]
mode = "advanced"
maxIdleConns = 100
maxIdleConnsPerHost = 60
maxConnsPerHost = 0
```

## 🔧 故障排查

### 1. 容器无法启动

```bash
# 查看容器日志
docker logs ghproxy

# 查看容器状态
docker ps -a

# 查看容器详细信息
docker inspect ghproxy
```

### 2. 配置文件问题

```bash
# 检查配置文件是否存在
ls -la ghproxy/config/

# 验证配置文件格式
cat ghproxy/config/config.toml

# 重新下载配置文件
wget -O ghproxy/config/config.toml \
  https://raw.githubusercontent.com/WJQSERVER-STUDIO/ghproxy/main/config/config.toml
```

### 3. 权限问题

```bash
# 修改目录权限
chmod -R 755 ghproxy/log ghproxy/config

# 修改所有者
chown -R 1000:1000 ghproxy/log ghproxy/config
```

### 4. 端口冲突

```bash
# 检查端口占用
netstat -tulpn | grep 7210

# 修改端口映射
# 编辑 docker-compose.yml，修改 ports 配置
ports:
  - '8080:8080'  # 改为其他端口
```

### 5. 网络问题

```bash
# 检查容器网络
docker network ls
docker network inspect ghproxy-network

# 重建网络
docker-compose down
docker-compose up -d
```

## 🎓 高级配置

### 1. 使用自定义配置

```yaml
volumes:
  - ./ghproxy/config:/data/ghproxy/config:ro  # 只读挂载
```

### 2. 多实例部署

```yaml
version: '3.9'

services:
  ghproxy-1:
    image: 'wjqserver/ghproxy:latest'
    ports:
      - '7210:8080'
    # ... 其他配置

  ghproxy-2:
    image: 'wjqserver/ghproxy:latest'
    ports:
      - '7211:8080'
    # ... 其他配置
```

### 3. 使用环境变量配置

创建 `.env` 文件：

```bash
GHPROXY_PORT=7210
GHPROXY_THEME=design
GHPROXY_LOG_LEVEL=info
```

在 `docker-compose.yml` 中使用：

```yaml
services:
  ghproxy:
    ports:
      - '${GHPROXY_PORT}:8080'
```

### 4. 自动重启策略

```yaml
restart: always  # 总是重启

# 或使用更精细的控制
restart: on-failure:5  # 失败时重启，最多5次
```

### 5. 监控和日志

```yaml
# 使用Prometheus监控（需要额外配置）
ports:
  - '7210:8080'
  - '9090:9090'  # Prometheus端口

# 集成外部日志系统
logging:
  driver: "syslog"
  options:
    syslog-address: "tcp://192.168.0.42:514"
```

## 📊 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f

# 进入容器
docker exec -it ghproxy sh

# 更新镜像
docker-compose pull
docker-compose up -d

# 清理旧镜像
docker image prune -f

# 查看资源使用
docker stats ghproxy

# 备份配置
tar -czf ghproxy-config-backup.tar.gz ghproxy/config/

# 恢复配置
tar -xzf ghproxy-config-backup.tar.gz
```

## 🔐 安全建议

1. **修改默认密码**：启用管理后台时务必修改默认密码
2. **使用HTTPS**：建议配合反向代理使用HTTPS
3. **限制资源**：设置合理的CPU和内存限制
4. **定期更新**：定期更新镜像版本
5. **备份数据**：定期备份配置文件和日志
6. **网络隔离**：使用Docker网络隔离服务

## 📝 示例：完整生产环境配置

```yaml
version: '3.9'

services:
  ghproxy:
    image: 'wjqserver/ghproxy:latest'
    container_name: ghproxy
    restart: always
    
    ports:
      - '7210:8080'
    
    volumes:
      - ./ghproxy/log:/data/ghproxy/log
      - ./ghproxy/config:/data/ghproxy/config
    
    environment:
      - TZ=Asia/Shanghai
    
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/healthcheck"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s
    
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 128M
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 120s
    
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    
    networks:
      - ghproxy-network
    
    labels:
      - "com.ghproxy.description=GitHub Proxy Server"
      - "com.ghproxy.version=4.3.4"

networks:
  ghproxy-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

## 🆘 获取帮助

- **项目地址**：https://github.com/WJQSERVER-STUDIO/ghproxy
- **问题反馈**：https://github.com/WJQSERVER-STUDIO/ghproxy/issues
- **Telegram群组**：https://t.me/ghproxy_go

## 📄 许可证

本项目使用双重许可：
- WJQserver Studio License 2.1
- Mozilla Public License Version 2.0

---

**最后更新**：2024年
