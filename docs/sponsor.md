# 赞助功能使用说明

## 功能概述

GHProxy 现已支持赞助功能，可以在独立的赞助页面展示支付宝、微信收款二维码以及USDT钱包地址，方便用户扫码赞助。

## 配置说明

### 1. 启用赞助功能

在配置文件 `config.toml` 中添加或修改以下配置：

```toml
[sponsor]
# 启用赞助功能
enabled = true

# 支付宝二维码图片URL
alipayQRCode = "https://example.com/images/alipay-qr.png"

# 微信二维码图片URL
wechatQRCode = "https://example.com/images/wechat-qr.png"

# USDT钱包地址（TRC20）
usdtAddress = "TNfSYG6F2vkiibd6J6mhhHNWDgWgNdF5hN"

# 赞助说明文字
sponsorText = "如果您觉得本项目对您有帮助，欢迎赞助支持！"
```

### 2. 准备二维码图片

#### 支付宝二维码

1. 打开支付宝APP
2. 点击右上角 "+" -> "收钱"
3. 保存收款二维码图片
4. 将图片上传到图床（推荐使用 GitHub、Imgur、阿里云OSS等）
5. 获取图片URL并填入配置文件的 `alipayQRCode` 字段

#### 微信二维码

1. 打开微信APP
2. 点击 "我" -> "服务" -> "收付款" -> "二维码收款"
3. 保存收款二维码图片
4. 将图片上传到图床
5. 获取图片URL并填入配置文件的 `wechatQRCode` 字段

### 3. 访问赞助页面

启用赞助功能后，访问以下地址查看赞助页面：

```
http://your-domain:port/sponsor
```

例如：`http://localhost:8080/sponsor`

## 页面特性

- **响应式设计**：适配桌面和移动设备
- **美观的UI**：渐变背景、卡片式布局
- **多种支付方式**：支持支付宝、微信、USDT
- **赞助明细展示**：滚动显示赞助者列表
- **一键返回**：快速返回首页

## 自定义赞助者列表

如需添加更多赞助者，可以修改 `sponsor/sponsor.go` 文件中的 `generateSponsorHTML` 函数，在捐赠列表部分添加新的赞助者：

```html
<div class="donor-item">
    <span class="donor-name">赞助者姓名</span>
    <span class="donor-amount">赞助金额</span>
</div>
```

## 图床推荐

以下是一些常用的图床服务：

1. **GitHub** - 免费，稳定，推荐
   - 创建一个专门的仓库存放图片
   - 使用 GitHub 的 raw 链接

2. **Imgur** - 免费，简单易用
   - 上传后获取直接链接

3. **阿里云OSS / 腾讯云COS** - 付费，速度快
   - 适合生产环境使用

4. **Cloudflare R2** - 免费额度大
   - 适合大流量网站

## 注意事项

1. **图片格式**：建议使用 PNG 或 JPG 格式的二维码图片
2. **图片大小**：建议二维码图片大小在 200x200 像素左右
3. **安全性**：不要在二维码图片中包含敏感信息
4. **HTTPS**：生产环境建议使用 HTTPS，确保图片URL也是 HTTPS
5. **测试**：配置完成后，务必测试二维码是否可以正常扫码支付

## 示例配置

完整的赞助配置示例：

```toml
[sponsor]
enabled = true
alipayQRCode = "https://raw.githubusercontent.com/yourname/yourrepo/main/images/alipay-qr.png"
wechatQRCode = "https://raw.githubusercontent.com/yourname/yourrepo/main/images/wechat-qr.png"
usdtAddress = "TNfSYG6F2vkiibd6J6mhhHNWDgWgNdF5hN"
sponsorText = "您的赞助将用于服务器维护和功能开发，感谢您的支持！"
```

## 常见问题

### Q: 二维码图片无法显示？

A: 请检查：
- 图片URL是否正确
- 图片是否可以公开访问
- 是否存在跨域问题
- 图片格式是否正确

### Q: 如何在主页添加赞助入口？

A: 由于前端页面是嵌入式的，您需要：
1. 修改前端仓库 [GHProxy-Frontend](https://github.com/WJQSERVER-STUDIO/GHProxy-Frontend)
2. 在主页添加赞助链接：`<a href="/sponsor">赞助支持</a>`
3. 重新构建前端并嵌入到Go项目中

### Q: 可以只启用部分支付方式吗？

A: 可以。如果某个支付方式的配置项为空，该支付方式将不会显示在页面上。例如：
- 只使用支付宝：只填写 `alipayQRCode`，其他留空
- 只使用USDT：只填写 `usdtAddress`，二维码URL留空

## 技术实现

- **后端**：Go语言 + Touka框架
- **前端**：Bootstrap 5 + 原生HTML/CSS
- **路由**：`GET /sponsor`
- **配置**：支持热更新，无需重启服务

## 更新日志

- **v1.0.0** - 初始版本
  - 支持支付宝、微信二维码展示
  - 支持USDT钱包地址展示
  - 支持赞助者列表滚动显示
  - 响应式设计，适配移动端
