# GHProxy 管理后台使用说明

## 功能概述

GHProxy 管理后台提供了完整的Web界面来管理所有配置项，包括：

### 核心配置
- **网站配置**：网站名称、描述、关键词、备案信息、统计代码等
- **服务器配置**：监听地址、端口、文件大小限制、CORS设置等
- **认证配置**：认证方式、Token设置等
- **主题设置**：前端主题切换、页面模式设置

### 功能配置
- **Docker配置**：Docker代理设置、目标仓库选择
- **速率限制**：请求速率控制、突发容量设置
- **白名单配置**：仓库白名单管理
- **黑名单配置**：仓库黑名单管理
- **IP过滤配置**：IP访问控制、白名单/黑名单设置

### 高级配置
- **HTTP客户端**：连接池配置、性能优化
- **Git克隆配置**：缓存模式、Smart-Git设置
- **Shell配置**：脚本嵌套加速、API重写
- **日志配置**：日志文件、级别、大小限制
- **出站代理**：SOCKS5/HTTP代理设置
- **赞助配置**：收款二维码、钱包地址

### 系统监控
- **系统状态**：查看系统运行状态和配置概览

## 启用管理后台

### 1. 修改配置文件

编辑 `config/config.toml` 文件，找到 `[admin]` 部分：

```toml
[admin]
enabled = true  # 启用管理后台
username = "admin"  # 管理员用户名
password = "your_secure_password"  # 设置强密码
pathPrefix = "/admin"  # 管理后台路径
```

### 2. 重启服务

```bash
# 停止服务
pkill ghproxy

# 启动服务
./ghproxy -c config/config.toml
```

## 访问管理后台

启动服务后，访问以下地址：

```
http://your-domain:8080/admin/
```

使用配置文件中设置的用户名和密码进行基本认证登录。

## 配置项说明

### 网站配置 (Website)

| 配置项 | 说明 | 示例 |
|--------|------|------|
| siteName | 网站名称 | "我的GHProxy" |
| siteDescription | 网站描述 | "快速可靠的GitHub代理服务" |
| siteKeywords | SEO关键词 | "github,proxy,download" |
| icpNumber | ICP备案号 | "京ICP备12345678号" |
| analyticsCode | 统计代码 | Google Analytics或百度统计代码 |
| footerText | 页脚文本 | "© 2024 My Site" |
| customCSS | 自定义CSS | 额外的样式代码 |
| customJS | 自定义JavaScript | 额外的脚本代码 |
| contactEmail | 联系邮箱 | "admin@example.com" |
| githubUrl | GitHub链接 | "https://github.com/user/repo" |
| twitterUrl | Twitter链接 | "https://twitter.com/username" |

### 服务器配置 (Server)

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| host | 监听地址 | "0.0.0.0" |
| port | 监听端口 | 8080 |
| sizeLimit | 文件大小限制(MB) | 125 |
| memLimit | 内存限制(MB) | 0 |
| cors | CORS设置 | "*" |
| debug | 调试模式 | false |

### 主题设置 (Pages)

| 配置项 | 说明 | 可选值 |
|--------|------|--------|
| mode | 页面模式 | "internal"(内置), "external"(外部) |
| theme | 主题选择 | bootstrap, nebula, design, metro, classic, mino, hub, free |
| staticDir | 外部静态文件目录 | "/data/www" |

### Docker配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| enabled | 启用Docker代理 | false |
| target | 目标仓库 | "dockerhub", "ghcr" |
| auth | 启用认证 | false |

### 速率限制配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| enabled | 启用速率限制 | false |
| ratePerMinute | 每分钟请求数 | 180 |
| burst | 突发容量 | 5 |

### 白名单/黑名单配置

| 配置项 | 说明 | 示例 |
|--------|------|--------|
| enabled | 启用功能 | true/false |
| whitelistFile | 白名单文件路径 | "/data/ghproxy/config/whitelist.json" |
| blacklistFile | 黑名单文件路径 | "/data/ghproxy/config/blacklist.json" |

**文件格式示例**：
```json
[
  "owner/repo",
  "owner2/*"
]
```

### IP过滤配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| enabled | 启用IP过滤 | false |
| enableAllowList | 启用白名单 | false |
| enableBlockList | 启用黑名单 | false |
| ipFilterFile | IP过滤文件路径 | "/data/ghproxy/config/ipfilter.json" |

### HTTP客户端配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| mode | 模式 | "auto", "advanced" |
| maxIdleConns | 最大空闲连接数 | 100 |
| maxIdleConnsPerHost | 每主机最大空闲连接数 | 60 |
| maxConnsPerHost | 每主机最大连接数 | 0 |
| useCustomRawHeaders | 使用自定义Raw Headers | false |

### Git克隆配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| mode | 模式 | "bypass"(直通), "cache"(缓存) |
| smartGitAddr | Smart-Git地址 | "http://127.0.0.1:8080" |
| ForceH2C | 强制H2C | false |

### Shell配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| editor | 启用编辑器 | true |
| rewriteAPI | 重写API | false |

### 日志配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| logFilePath | 日志文件路径 | "/data/ghproxy/log/ghproxy.log" |
| maxLogSize | 最大日志大小(MB) | 5 |
| level | 日志级别 | "debug", "info", "warn", "error" |

### 出站代理配置

| 配置项 | 说明 | 示例 |
|--------|------|--------|
| enabled | 启用出站代理 | true/false |
| url | 代理URL | "socks5://127.0.0.1:1080" 或 "http://127.0.0.1:7890" |

### 赞助配置

| 配置项 | 说明 | 示例 |
|--------|------|--------|
| enabled | 启用赞助功能 | true/false |
| alipayQRCode | 支付宝二维码URL | "https://example.com/alipay.png" |
| wechatQRCode | 微信二维码URL | "https://example.com/wechat.png" |
| usdtAddress | USDT钱包地址 | "TNfSYG6F2vkiibd6J6mhhHNWDgWgNdF5hN" |
| sponsorText | 赞助说明文字 | "感谢您的支持！" |

## API接口

管理后台提供完整的RESTful API接口：

### 完整配置
```
GET  /admin/api/config
POST /admin/api/config
```

### 各模块配置
```
GET  /admin/api/config/website
POST /admin/api/config/website

GET  /admin/api/config/server
POST /admin/api/config/server

GET  /admin/api/config/auth
POST /admin/api/config/auth

GET  /admin/api/config/pages
POST /admin/api/config/pages

GET  /admin/api/config/docker
POST /admin/api/config/docker

GET  /admin/api/config/ratelimit
POST /admin/api/config/ratelimit

GET  /admin/api/config/whitelist
POST /admin/api/config/whitelist

GET  /admin/api/config/blacklist
POST /admin/api/config/blacklist

GET  /admin/api/config/ipfilter
POST /admin/api/config/ipfilter

GET  /admin/api/config/httpc
POST /admin/api/config/httpc

GET  /admin/api/config/gitclone
POST /admin/api/config/gitclone

GET  /admin/api/config/shell
POST /admin/api/config/shell

GET  /admin/api/config/log
POST /admin/api/config/log

GET  /admin/api/config/outbound
POST /admin/api/config/outbound

GET  /admin/api/config/sponsor
POST /admin/api/config/sponsor
```

### 系统操作
```
POST /admin/api/config/reload  # 重载配置
GET  /admin/api/status         # 获取系统状态
```

## 安全建议

1. **修改默认密码**：务必修改默认的管理员密码
2. **使用HTTPS**：建议在生产环境中使用HTTPS
3. **限制访问**：可以通过IP过滤功能限制管理后台的访问IP
4. **定期备份**：定期备份配置文件
5. **监控日志**：监控管理后台的访问日志

## 故障排查

### 无法访问管理后台

1. 检查 `admin.enabled` 是否设置为 `true`
2. 检查服务是否正常启动
3. 检查防火墙是否开放端口
4. 查看日志文件确认错误信息

### 登录失败

1. 确认用户名和密码正确
2. 检查浏览器是否支持基本认证
3. 清除浏览器缓存重试

### 配置保存失败

1. 检查配置文件权限
2. 检查磁盘空间
3. 查看日志文件确认错误信息

### 配置未生效

1. 部分配置需要重启服务才能生效（如端口、内存限制等）
2. 使用"重载配置"功能重新加载配置
3. 检查配置文件格式是否正确

## 使用示例

### 示例1：添加备案信息

1. 访问管理后台
2. 点击左侧菜单"网站配置"
3. 在"ICP备案号"输入框填写备案号，如：`京ICP备12345678号`
4. 点击"保存配置"按钮

### 示例2：添加统计代码

1. 访问管理后台
2. 点击左侧菜单"网站配置"
3. 在"统计代码"文本框粘贴统计代码，例如：

```html
<!-- Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=GA_MEASUREMENT_ID"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_MEASUREMENT_ID');
</script>
```

4. 点击"保存配置"按钮

### 示例3：切换主题

1. 访问管理后台
2. 点击左侧菜单"主题设置"
3. 在"主题选择"下拉框选择喜欢的主题（如：design）
4. 点击"保存配置"按钮
5. 刷新前端页面查看效果

### 示例4：配置Docker代理

1. 访问管理后台
2. 点击左侧菜单"Docker配置"
3. 勾选"启用Docker代理"
4. 选择目标仓库（DockerHub或GitHub Container Registry）
5. 点击"保存配置"按钮

### 示例5：设置白名单

1. 访问管理后台
2. 点击左侧菜单"白名单配置"
3. 勾选"启用白名单"
4. 设置白名单文件路径（如：`/data/ghproxy/config/whitelist.json`）
5. 点击"保存配置"按钮
6. 编辑白名单文件，添加允许的仓库：

```json
[
  "WJQSERVER-STUDIO/ghproxy",
  "golang/go",
  "kubernetes/*"
]
```

### 示例6：配置出站代理

1. 访问管理后台
2. 点击左侧菜单"出站代理"
3. 勾选"启用出站代理"
4. 填写代理URL（如：`socks5://127.0.0.1:1080`）
5. 点击"保存配置"按钮

## 配置管理最佳实践

### 1. 配置备份
定期备份配置文件：
```bash
cp /data/ghproxy/config/config.toml /data/ghproxy/config/config.toml.backup
```

### 2. 配置验证
修改配置后，检查服务是否正常运行：
```bash
curl http://localhost:8080/api/healthcheck
```

### 3. 日志监控
定期查看日志文件，监控异常：
```bash
tail -f /data/ghproxy/log/ghproxy.log
```

### 4. 性能优化
根据实际负载调整配置：
- HTTP客户端连接池参数
- 速率限制参数
- 内存限制

## 技术支持

如有问题，请访问：https://github.com/WJQSERVER-STUDIO/ghproxy

## 更新日志

### v2.0 (当前版本)
- ✅ 新增6个配置管理模块
- ✅ 支持15个配置模块的完整管理
- ✅ 新增30+个API端点
- ✅ 完善前端管理界面
- ✅ 支持实时配置保存

### v1.0
- 基础配置管理功能
- 网站配置、服务器配置、认证配置
- 系统状态查看
