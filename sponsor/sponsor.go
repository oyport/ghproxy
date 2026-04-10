package sponsor

import (
	"ghproxy/config"

	"github.com/infinite-iroha/touka"
)

// SponsorPage 返回赞助页面的HTML内容
func SponsorPage(cfg *config.Config) touka.HandlerFunc {
	return func(c *touka.Context) {
		// 如果未启用赞助功能，返回404
		if !cfg.Sponsor.Enabled {
			c.String(404, "Sponsor page is not enabled")
			return
		}

		// 生成赞助页面HTML
		html := generateSponsorHTML(cfg)
		c.HTML(200, html)
	}
}

// generateSponsorHTML 生成赞助页面的HTML内容
func generateSponsorHTML(cfg *config.Config) string {
	siteName := cfg.Website.SiteName
	if siteName == "" {
		siteName = "GHProxy"
	}

	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>赞助支持 - ` + siteName + `</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        body {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px 0;
        }
        .sponsor-container {
            max-width: 900px;
            margin: 0 auto;
        }
        .sponsor-card {
            background: white;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            padding: 40px;
            margin-bottom: 30px;
        }
        .sponsor-title {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
            font-weight: 600;
        }
        .qr-code-wrapper {
            text-align: center;
            margin-bottom: 30px;
        }
        .qr-code-item {
            display: inline-block;
            margin: 0 20px;
            text-align: center;
        }
        .qr-code-img {
            width: 200px;
            height: 200px;
            border: 3px solid #667eea;
            border-radius: 10px;
            margin-bottom: 15px;
            object-fit: contain;
            background: #f8f9fa;
        }
        .qr-code-label {
            font-size: 16px;
            font-weight: 600;
            color: #555;
        }
        .usdt-section {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 10px;
            margin-top: 30px;
        }
        .usdt-address {
            font-family: 'Courier New', monospace;
            background: #e9ecef;
            padding: 15px;
            border-radius: 5px;
            word-break: break-all;
            margin-top: 10px;
        }
        .sponsor-text {
            text-align: center;
            color: #666;
            font-size: 16px;
            margin-bottom: 30px;
            line-height: 1.6;
        }
        .donor-list {
            margin-top: 40px;
        }
        .donor-list-title {
            color: #333;
            font-weight: 600;
            margin-bottom: 20px;
        }
        .donor-item {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 10px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .donor-name {
            font-weight: 600;
            color: #555;
        }
        .donor-amount {
            color: #667eea;
            font-weight: 600;
        }
        .back-btn {
            display: block;
            width: 200px;
            margin: 30px auto 0;
            text-align: center;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 30px;
            border-radius: 25px;
            text-decoration: none;
            font-weight: 600;
            transition: transform 0.3s, box-shadow 0.3s;
        }
        .back-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
            color: white;
        }
        .scroll-container {
            max-height: 300px;
            overflow-y: auto;
            padding-right: 10px;
        }
        .scroll-container::-webkit-scrollbar {
            width: 8px;
        }
        .scroll-container::-webkit-scrollbar-track {
            background: #f1f1f1;
            border-radius: 4px;
        }
        .scroll-container::-webkit-scrollbar-thumb {
            background: #667eea;
            border-radius: 4px;
        }
        .scroll-container::-webkit-scrollbar-thumb:hover {
            background: #764ba2;
        }
    </style>
</head>
<body>
    <div class="container sponsor-container">
        <div class="sponsor-card">
            <h1 class="sponsor-title">💝 赞助支持</h1>
            
            <p class="sponsor-text">` + cfg.Sponsor.SponsorText + `</p>
            
            <div class="qr-code-wrapper">
`

	// 添加支付宝二维码
	if cfg.Sponsor.AlipayQRCode != "" {
		html += `
                <div class="qr-code-item">
                    <img src="` + cfg.Sponsor.AlipayQRCode + `" alt="支付宝二维码" class="qr-code-img">
                    <div class="qr-code-label">支付宝</div>
                </div>
`
	}

	// 添加微信二维码
	if cfg.Sponsor.WechatQRCode != "" {
		html += `
                <div class="qr-code-item">
                    <img src="` + cfg.Sponsor.WechatQRCode + `" alt="微信二维码" class="qr-code-img">
                    <div class="qr-code-label">微信</div>
                </div>
`
	}

	html += `
            </div>
`

	// 添加USDT地址
	if cfg.Sponsor.USDTAddress != "" {
		html += `
            <div class="usdt-section">
                <h5><strong>USDT (TRC20) 钱包地址：</strong></h5>
                <div class="usdt-address">` + cfg.Sponsor.USDTAddress + `</div>
            </div>
`
	}

	// 添加捐赠列表
	html += `
            <div class="donor-list">
                <h4 class="donor-list-title">🏆 感谢以下赞助者</h4>
                <div class="scroll-container">
                    <div class="donor-item">
                        <span class="donor-name">starry</span>
                        <span class="donor-amount">8 USDT (TRC20)</span>
                    </div>
                    <!-- 更多赞助者可以在这里添加 -->
                </div>
            </div>
            
            <a href="/" class="back-btn">← 返回首页</a>
        </div>
    </div>
    
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>
`

	return html
}
