# 赞助功能快速开始

## 一、准备工作

### 1. 准备二维码图片

#### 支付宝收款码
1. 打开支付宝APP → 点击右上角 "+" → "收钱"
2. 保存二维码图片
3. 上传到图床（推荐GitHub）
4. 获取图片URL

#### 微信收款码
1. 打开微信APP → "我" → "服务" → "收付款" → "二维码收款"
2. 保存二维码图片
3. 上传到图床
4. 获取图片URL

### 2. 图床推荐

**GitHub（推荐）**
```bash
# 1. 创建图片仓库
mkdir sponsor-images
cd sponsor-images
git init

# 2. 添加二维码图片
cp ~/Downloads/alipay-qr.png ./
cp ~/Downloads/wechat-qr.png ./

# 3. 推送到GitHub
git add .
git commit -m "Add sponsor QR codes"
git remote add origin https://github.com/yourname/sponsor-images.git
git push -u origin main

# 4. 获取图片URL
# 支付宝: https://raw.githubusercontent.com/yourname/sponsor-images/main/alipay-qr.png
# 微信: https://raw.githubusercontent.com/yourname/sponsor-images/main/wechat-qr.png
```

## 二、配置赞助功能

### 1. 修改配置文件

编辑 `config/config.toml`，添加以下内容：

```toml
[sponsor]
# 启用赞助功能
enabled = true

# 支付宝二维码（替换为您的实际URL）
alipayQRCode = "https://raw.githubusercontent.com/yourname/sponsor-images/main/alipay-qr.png"

# 微信二维码（替换为您的实际URL）
wechatQRCode = "https://raw.githubusercontent.com/yourname/sponsor-images/main/wechat-qr.png"

# USDT钱包地址
usdtAddress = "TNfSYG6F2vkiibd6J6mhhHNWDgWgNdF5hN"

# 赞助说明
sponsorText = "感谢您的支持！"
```

### 2. 重启服务

```bash
# 停止服务
pkill ghproxy

# 启动服务
./ghproxy -c config/config.toml
```

## 三、访问测试

### 1. 访问赞助页面

浏览器打开：`http://localhost:8080/sponsor`

### 2. 检查功能

- [ ] 页面正常显示
- [ ] 支付宝二维码可见
- [ ] 微信二维码可见
- [ ] USDT地址正确
- [ ] 赞助者列表显示
- [ ] 返回首页按钮可用

### 3. 测试扫码

使用手机支付宝/微信扫描二维码，确认可以正常识别收款码。

## 四、生产环境部署

### 1. 使用HTTPS

确保您的域名已配置SSL证书：

```nginx
# Nginx配置示例
server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
    }
}
```

### 2. 确保图片URL使用HTTPS

```toml
[sponsor]
alipayQRCode = "https://raw.githubusercontent.com/..."  # 使用HTTPS
wechatQRCode = "https://raw.githubusercontent.com/..."   # 使用HTTPS
```

### 3. 添加到主页（可选）

修改前端代码，在主页添加赞助入口：

```html
<a href="/sponsor" class="btn btn-primary">💝 赞助支持</a>
```

## 五、常见问题

### Q1: 二维码不显示？

**解决方案：**
1. 检查图片URL是否正确
2. 确认图片可以公开访问
3. 查看浏览器控制台是否有错误
4. 测试图片URL是否可以正常访问

### Q2: 配置不生效？

**解决方案：**
1. 确认配置文件路径正确
2. 检查TOML语法是否正确
3. 重启服务
4. 查看日志是否有错误信息

### Q3: 移动端显示异常？

**解决方案：**
1. 清除浏览器缓存
2. 检查CSS是否正确加载
3. 确认使用了响应式设计

## 六、完整示例

### 最小配置

```toml
[sponsor]
enabled = true
alipayQRCode = "https://example.com/alipay.png"
wechatQRCode = "https://example.com/wechat.png"
```

### 完整配置

```toml
[sponsor]
enabled = true
alipayQRCode = "https://raw.githubusercontent.com/yourname/sponsor-images/main/alipay-qr.png"
wechatQRCode = "https://raw.githubusercontent.com/yourname/sponsor-images/main/wechat-qr.png"
usdtAddress = "TNfSYG6F2vkiibd6J6mhhHNWDgWgNdF5hN"
sponsorText = "您的赞助将用于服务器维护和功能开发，感谢您的支持！"
```

## 七、下一步

- [ ] 在主页添加赞助入口
- [ ] 自定义赞助者列表
- [ ] 添加更多支付方式
- [ ] 集成到管理后台

## 需要帮助？

- 查看详细文档：[docs/sponsor.md](sponsor.md)
- 提交Issue：[GitHub Issues](https://github.com/WJQSERVER-STUDIO/ghproxy/issues)
- 加入讨论：[Telegram群组](https://t.me/ghproxy_go)
