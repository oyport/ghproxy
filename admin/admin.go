package admin

import (
	"encoding/json"
	"ghproxy/config"
	"net/http"
	"sync"

	"github.com/fenthope/bauth"
	"github.com/infinite-iroha/touka"
)

var (
	cfg      *config.Config
	cfgFile  string
	cfgMutex sync.RWMutex
)

// InitAdmin 初始化管理模块
func InitAdmin(config *config.Config, configFile string) {
	cfg = config
	cfgFile = configFile
}

// SetupAdminRoutes 设置管理后台路由
func SetupAdminRoutes(r *touka.Engine, config *config.Config) {
	if !config.Admin.Enabled {
		return
	}

	adminRouter := r.Group(config.Admin.PathPrefix)

	// 基本认证中间件
	adminRouter.Use(bauth.BasicAuth(map[string]string{
		config.Admin.Username: config.Admin.Password,
	}, "GHProxy Admin"))

	// 管理首页
	adminRouter.GET("/", AdminDashboardHandler)

	// 配置管理API
	adminRouter.GET("/api/config", GetConfigHandler)
	adminRouter.POST("/api/config", UpdateConfigHandler)
	adminRouter.GET("/api/config/website", GetWebsiteConfigHandler)
	adminRouter.POST("/api/config/website", UpdateWebsiteConfigHandler)
	adminRouter.GET("/api/config/server", GetServerConfigHandler)
	adminRouter.POST("/api/config/server", UpdateServerConfigHandler)
	adminRouter.GET("/api/config/auth", GetAuthConfigHandler)
	adminRouter.POST("/api/config/auth", UpdateAuthConfigHandler)
	
	// 新增配置管理API
	adminRouter.GET("/api/config/pages", GetPagesConfigHandler)
	adminRouter.POST("/api/config/pages", UpdatePagesConfigHandler)
	adminRouter.GET("/api/config/sponsor", GetSponsorConfigHandler)
	adminRouter.POST("/api/config/sponsor", UpdateSponsorConfigHandler)
	adminRouter.GET("/api/config/docker", GetDockerConfigHandler)
	adminRouter.POST("/api/config/docker", UpdateDockerConfigHandler)
	adminRouter.GET("/api/config/ratelimit", GetRateLimitConfigHandler)
	adminRouter.POST("/api/config/ratelimit", UpdateRateLimitConfigHandler)
	adminRouter.GET("/api/config/whitelist", GetWhitelistConfigHandler)
	adminRouter.POST("/api/config/whitelist", UpdateWhitelistConfigHandler)
	adminRouter.GET("/api/config/blacklist", GetBlacklistConfigHandler)
	adminRouter.POST("/api/config/blacklist", UpdateBlacklistConfigHandler)
	adminRouter.GET("/api/config/httpc", GetHttpcConfigHandler)
	adminRouter.POST("/api/config/httpc", UpdateHttpcConfigHandler)
	adminRouter.GET("/api/config/gitclone", GetGitCloneConfigHandler)
	adminRouter.POST("/api/config/gitclone", UpdateGitCloneConfigHandler)
	adminRouter.GET("/api/config/shell", GetShellConfigHandler)
	adminRouter.POST("/api/config/shell", UpdateShellConfigHandler)
	adminRouter.GET("/api/config/log", GetLogConfigHandler)
	adminRouter.POST("/api/config/log", UpdateLogConfigHandler)
	adminRouter.GET("/api/config/ipfilter", GetIPFilterConfigHandler)
	adminRouter.POST("/api/config/ipfilter", UpdateIPFilterConfigHandler)
	adminRouter.GET("/api/config/outbound", GetOutboundConfigHandler)
	adminRouter.POST("/api/config/outbound", UpdateOutboundConfigHandler)

	// 配置重载
	adminRouter.POST("/api/config/reload", ReloadConfigHandler)

	// 系统状态
	adminRouter.GET("/api/status", GetSystemStatusHandler)
}

// AdminDashboardHandler 管理后台首页
func AdminDashboardHandler(c *touka.Context) {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GHProxy 管理后台</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        body { background-color: #f5f5f5; }
        .sidebar { background-color: #343a40; color: white; min-height: 100vh; }
        .sidebar a { color: #adb5bd; text-decoration: none; display: block; padding: 10px 15px; }
        .sidebar a:hover, .sidebar a.active { background-color: #495057; color: white; }
        .content { padding: 20px; }
    </style>
</head>
<body>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md-2 sidebar p-0">
                <h4 class="p-3 border-bottom">GHProxy 管理</h4>
                <a href="#" class="active" onclick="loadPage('dashboard')">仪表盘</a>
                <a href="#" onclick="loadPage('website')">网站配置</a>
                <a href="#" onclick="loadPage('pages')">主题设置</a>
                <a href="#" onclick="loadPage('server')">服务器配置</a>
                <a href="#" onclick="loadPage('auth')">认证配置</a>
                <a href="#" onclick="loadPage('sponsor')">赞助配置</a>
                <a href="#" onclick="loadPage('docker')">Docker配置</a>
                <a href="#" onclick="loadPage('ratelimit')">速率限制</a>
                <a href="#" onclick="loadPage('whitelist')">白名单配置</a>
                <a href="#" onclick="loadPage('blacklist')">黑名单配置</a>
                <a href="#" onclick="loadPage('httpc')">HTTP客户端</a>
                <a href="#" onclick="loadPage('gitclone')">Git克隆配置</a>
                <a href="#" onclick="loadPage('shell')">Shell配置</a>
                <a href="#" onclick="loadPage('log')">日志配置</a>
                <a href="#" onclick="loadPage('ipfilter')">IP过滤配置</a>
                <a href="#" onclick="loadPage('outbound')">出站代理</a>
                <a href="#" onclick="loadPage('advanced')">高级配置</a>
            </div>
            <div class="col-md-10 content">
                <div id="main-content">
                    <h2>欢迎来到 GHProxy 管理后台</h2>
                    <p>请从左侧菜单选择要管理的配置项。</p>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        function loadPage(page) {
            const content = document.getElementById('main-content');
            switch(page) {
                case 'dashboard':
                    content.innerHTML = '<h2>仪表盘</h2><p>系统状态和概览信息</p>';
                    loadSystemStatus();
                    break;
                case 'website':
                    loadWebsiteConfig();
                    break;
                case 'pages':
                    loadPagesConfig();
                    break;
                case 'server':
                    loadServerConfig();
                    break;
                case 'auth':
                    loadAuthConfig();
                    break;
                case 'sponsor':
                    loadSponsorConfig();
                    break;
                case 'docker':
                    loadDockerConfig();
                    break;
                case 'ratelimit':
                    loadRateLimitConfig();
                    break;
                case 'whitelist':
                    loadWhitelistConfig();
                    break;
                case 'blacklist':
                    loadBlacklistConfig();
                    break;
                case 'httpc':
                    loadHttpcConfig();
                    break;
                case 'gitclone':
                    loadGitCloneConfig();
                    break;
                case 'shell':
                    loadShellConfig();
                    break;
                case 'log':
                    loadLogConfig();
                    break;
                case 'ipfilter':
                    loadIPFilterConfig();
                    break;
                case 'outbound':
                    loadOutboundConfig();
                    break;
                case 'advanced':
                    content.innerHTML = '<h2>高级配置</h2><p>高级配置选项</p>';
                    break;
            }
        }

        function loadSystemStatus() {
            fetch('/admin/api/status')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>系统状态</h2>' +
                        '<div class="row">' +
                        '<div class="col-md-4"><div class="card"><div class="card-body"><h5>服务器端口</h5><p>' + data.serverPort + '</p></div></div></div>' +
                        '<div class="col-md-4"><div class="card"><div class="card-body"><h5>认证状态</h5><p>' + (data.authEnabled ? '已启用' : '未启用') + '</p></div></div></div>' +
                        '<div class="col-md-4"><div class="card"><div class="card-body"><h5>速率限制</h5><p>' + (data.rateLimitEnabled ? '已启用' : '未启用') + '</p></div></div></div>' +
                        '</div>';
                });
        }

        function loadWebsiteConfig() {
            fetch('/admin/api/config/website')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>网站配置</h2>' +
                        '<form id="website-form">' +
                        '<div class="mb-3"><label class="form-label">网站名称</label><input type="text" class="form-control" name="siteName" value="' + (data.siteName || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">网站描述</label><input type="text" class="form-control" name="siteDescription" value="' + (data.siteDescription || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">网站关键词</label><input type="text" class="form-control" name="siteKeywords" value="' + (data.siteKeywords || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">ICP备案号</label><input type="text" class="form-control" name="icpNumber" value="' + (data.icpNumber || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">统计代码</label><textarea class="form-control" name="analyticsCode" rows="3">' + (data.analyticsCode || '') + '</textarea></div>' +
                        '<div class="mb-3"><label class="form-label">页脚文本</label><input type="text" class="form-control" name="footerText" value="' + (data.footerText || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">联系邮箱</label><input type="email" class="form-control" name="contactEmail" value="' + (data.contactEmail || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">GitHub链接</label><input type="text" class="form-control" name="githubUrl" value="' + (data.githubUrl || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">Twitter链接</label><input type="text" class="form-control" name="twitterUrl" value="' + (data.twitterUrl || '') + '"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('website-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {};
                        formData.forEach((value, key) => config[key] = value);

                        fetch('/admin/api/config/website', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadServerConfig() {
            fetch('/admin/api/config/server')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>服务器配置</h2>' +
                        '<form id="server-form">' +
                        '<div class="mb-3"><label class="form-label">监听地址</label><input type="text" class="form-control" name="host" value="' + (data.host || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">监听端口</label><input type="number" class="form-control" name="port" value="' + (data.port || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">文件大小限制(MB)</label><input type="number" class="form-control" name="sizeLimit" value="' + (data.sizeLimit || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">CORS设置</label><input type="text" class="form-control" name="cors" value="' + (data.cors || '') + '"></div>' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="debug" ' + (data.debug ? 'checked' : '') + '><label class="form-check-label">调试模式</label></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('server-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            host: formData.get('host'),
                            port: parseInt(formData.get('port')),
                            sizeLimit: parseInt(formData.get('sizeLimit')),
                            cors: formData.get('cors'),
                            debug: formData.has('debug')
                        };

                        fetch('/admin/api/config/server', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('配置保存成功！需要重启服务生效。');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadAuthConfig() {
            fetch('/admin/api/config/auth')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>认证配置</h2>' +
                        '<form id="auth-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用认证</label></div>' +
                        '<div class="mb-3"><label class="form-label">认证方式</label><select class="form-control" name="method"><option value="parameters" ' + (data.method === 'parameters' ? 'selected' : '') + '>参数认证</option><option value="header" ' + (data.method === 'header' ? 'selected' : '') + '>Header认证</option></select></div>' +
                        '<div class="mb-3"><label class="form-label">认证Key</label><input type="text" class="form-control" name="key" value="' + (data.key || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">认证Token</label><input type="text" class="form-control" name="token" value="' + (data.token || '') + '"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('auth-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            method: formData.get('method'),
                            key: formData.get('key'),
                            token: formData.get('token')
                        };

                        fetch('/admin/api/config/auth', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadPagesConfig() {
            fetch('/admin/api/config/pages')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>主题设置</h2>' +
                        '<form id="pages-form">' +
                        '<div class="mb-3"><label class="form-label">页面模式</label><select class="form-control" name="mode"><option value="internal" ' + (data.mode === 'internal' ? 'selected' : '') + '>内置模式</option><option value="external" ' + (data.mode === 'external' ? 'selected' : '') + '>外部模式</option></select></div>' +
                        '<div class="mb-3"><label class="form-label">主题选择</label><select class="form-control" name="theme"><option value="bootstrap" ' + (data.theme === 'bootstrap' ? 'selected' : '') + '>Bootstrap</option><option value="nebula" ' + (data.theme === 'nebula' ? 'selected' : '') + '>Nebula</option><option value="design" ' + (data.theme === 'design' ? 'selected' : '') + '>Design</option><option value="metro" ' + (data.theme === 'metro' ? 'selected' : '') + '>Metro</option><option value="classic" ' + (data.theme === 'classic' ? 'selected' : '') + '>Classic</option><option value="mino" ' + (data.theme === 'mino' ? 'selected' : '') + '>Mino</option><option value="hub" ' + (data.theme === 'hub' ? 'selected' : '') + '>Hub</option><option value="free" ' + (data.theme === 'free' ? 'selected' : '') + '>Free</option></select></div>' +
                        '<div class="mb-3"><label class="form-label">外部静态文件目录</label><input type="text" class="form-control" name="staticDir" value="' + (data.staticDir || '') + '"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('pages-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            mode: formData.get('mode'),
                            theme: formData.get('theme'),
                            staticDir: formData.get('staticDir')
                        };

                        fetch('/admin/api/config/pages', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('主题配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadSponsorConfig() {
            fetch('/admin/api/config/sponsor')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>赞助配置</h2>' +
                        '<form id="sponsor-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用赞助功能</label></div>' +
                        '<div class="mb-3"><label class="form-label">支付宝二维码URL</label><input type="text" class="form-control" name="alipayQRCode" value="' + (data.alipayQRCode || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">微信二维码URL</label><input type="text" class="form-control" name="wechatQRCode" value="' + (data.wechatQRCode || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">USDT钱包地址</label><input type="text" class="form-control" name="usdtAddress" value="' + (data.usdtAddress || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">赞助说明文字</label><textarea class="form-control" name="sponsorText" rows="3">' + (data.sponsorText || '') + '</textarea></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('sponsor-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            alipayQRCode: formData.get('alipayQRCode'),
                            wechatQRCode: formData.get('wechatQRCode'),
                            usdtAddress: formData.get('usdtAddress'),
                            sponsorText: formData.get('sponsorText')
                        };

                        fetch('/admin/api/config/sponsor', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('赞助配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadDockerConfig() {
            fetch('/admin/api/config/docker')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>Docker配置</h2>' +
                        '<form id="docker-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用Docker代理</label></div>' +
                        '<div class="mb-3"><label class="form-label">目标仓库</label><select class="form-control" name="target"><option value="dockerhub" ' + (data.target === 'dockerhub' ? 'selected' : '') + '>DockerHub</option><option value="ghcr" ' + (data.target === 'ghcr' ? 'selected' : '') + '>GitHub Container Registry</option></select></div>' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="auth" ' + (data.auth ? 'checked' : '') + '><label class="form-check-label">启用认证</label></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('docker-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            target: formData.get('target'),
                            auth: formData.has('auth')
                        };

                        fetch('/admin/api/config/docker', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('Docker配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadRateLimitConfig() {
            fetch('/admin/api/config/ratelimit')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>速率限制配置</h2>' +
                        '<form id="ratelimit-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用速率限制</label></div>' +
                        '<div class="mb-3"><label class="form-label">每分钟请求数</label><input type="number" class="form-control" name="ratePerMinute" value="' + (data.ratePerMinute || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">突发容量</label><input type="number" class="form-control" name="burst" value="' + (data.burst || '') + '"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('ratelimit-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            ratePerMinute: parseInt(formData.get('ratePerMinute')),
                            burst: parseInt(formData.get('burst'))
                        };

                        fetch('/admin/api/config/ratelimit', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('速率限制配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadWhitelistConfig() {
            fetch('/admin/api/config/whitelist')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>白名单配置</h2>' +
                        '<div class="alert alert-info">白名单功能允许您指定允许访问的仓库列表。启用后，只有白名单中的仓库才能被访问。</div>' +
                        '<form id="whitelist-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用白名单</label></div>' +
                        '<div class="mb-3"><label class="form-label">白名单文件路径</label><input type="text" class="form-control" name="whitelistFile" value="' + (data.whitelistFile || '') + '"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>' +
                        '<div class="mt-3"><h5>白名单文件格式说明</h5><p>白名单文件为JSON格式，包含允许访问的仓库列表。例如：</p><pre>[\\n  "owner/repo",\\n  "owner2/*"\\n]</pre><p>支持通配符 * 表示匹配所有仓库</p></div>';

                    document.getElementById('whitelist-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            whitelistFile: formData.get('whitelistFile')
                        };

                        fetch('/admin/api/config/whitelist', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('白名单配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadBlacklistConfig() {
            fetch('/admin/api/config/blacklist')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>黑名单配置</h2>' +
                        '<div class="alert alert-warning">黑名单功能允许您屏蔽特定的仓库。启用后，黑名单中的仓库将无法被访问。</div>' +
                        '<form id="blacklist-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用黑名单</label></div>' +
                        '<div class="mb-3"><label class="form-label">黑名单文件路径</label><input type="text" class="form-control" name="blacklistFile" value="' + (data.blacklistFile || '') + '"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>' +
                        '<div class="mt-3"><h5>黑名单文件格式说明</h5><p>黑名单文件为JSON格式，包含需要屏蔽的仓库列表。例如：</p><pre>[\\n  "owner/bad-repo",\\n  "spammer/*"\\n]</pre><p>支持通配符 * 表示匹配所有仓库</p></div>';

                    document.getElementById('blacklist-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            blacklistFile: formData.get('blacklistFile')
                        };

                        fetch('/admin/api/config/blacklist', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('黑名单配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadHttpcConfig() {
            fetch('/admin/api/config/httpc')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>HTTP客户端配置</h2>' +
                        '<form id="httpc-form">' +
                        '<div class="mb-3"><label class="form-label">模式</label><select class="form-control" name="mode"><option value="auto" ' + (data.mode === 'auto' ? 'selected' : '') + '>自动</option><option value="advanced" ' + (data.mode === 'advanced' ? 'selected' : '') + '>高级</option></select></div>' +
                        '<div class="mb-3"><label class="form-label">最大空闲连接数</label><input type="number" class="form-control" name="maxIdleConns" value="' + (data.maxIdleConns || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">每主机最大空闲连接数</label><input type="number" class="form-control" name="maxIdleConnsPerHost" value="' + (data.maxIdleConnsPerHost || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">每主机最大连接数</label><input type="number" class="form-control" name="maxConnsPerHost" value="' + (data.maxConnsPerHost || '') + '"></div>' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="useCustomRawHeaders" ' + (data.useCustomRawHeaders ? 'checked' : '') + '><label class="form-check-label">使用自定义Raw Headers</label></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('httpc-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            mode: formData.get('mode'),
                            maxIdleConns: parseInt(formData.get('maxIdleConns')),
                            maxIdleConnsPerHost: parseInt(formData.get('maxIdleConnsPerHost')),
                            maxConnsPerHost: parseInt(formData.get('maxConnsPerHost')),
                            useCustomRawHeaders: formData.has('useCustomRawHeaders')
                        };

                        fetch('/admin/api/config/httpc', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('HTTP客户端配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadGitCloneConfig() {
            fetch('/admin/api/config/gitclone')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>Git克隆配置</h2>' +
                        '<form id="gitclone-form">' +
                        '<div class="mb-3"><label class="form-label">模式</label><select class="form-control" name="mode"><option value="bypass" ' + (data.mode === 'bypass' ? 'selected' : '') + '>直通</option><option value="cache" ' + (data.mode === 'cache' ? 'selected' : '') + '>缓存</option></select></div>' +
                        '<div class="mb-3"><label class="form-label">Smart-Git地址</label><input type="text" class="form-control" name="smartGitAddr" value="' + (data.smartGitAddr || '') + '"></div>' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="ForceH2C" ' + (data.ForceH2C ? 'checked' : '') + '><label class="form-check-label">强制H2C</label></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('gitclone-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            mode: formData.get('mode'),
                            smartGitAddr: formData.get('smartGitAddr'),
                            ForceH2C: formData.has('ForceH2C')
                        };

                        fetch('/admin/api/config/gitclone', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('Git克隆配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadShellConfig() {
            fetch('/admin/api/config/shell')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>Shell配置</h2>' +
                        '<form id="shell-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="editor" ' + (data.editor ? 'checked' : '') + '><label class="form-check-label">启用编辑器</label></div>' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="rewriteAPI" ' + (data.rewriteAPI ? 'checked' : '') + '><label class="form-check-label">重写API</label></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('shell-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            editor: formData.has('editor'),
                            rewriteAPI: formData.has('rewriteAPI')
                        };

                        fetch('/admin/api/config/shell', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('Shell配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadLogConfig() {
            fetch('/admin/api/config/log')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>日志配置</h2>' +
                        '<form id="log-form">' +
                        '<div class="mb-3"><label class="form-label">日志文件路径</label><input type="text" class="form-control" name="logFilePath" value="' + (data.logFilePath || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">最大日志大小(MB)</label><input type="number" class="form-control" name="maxLogSize" value="' + (data.maxLogSize || '') + '"></div>' +
                        '<div class="mb-3"><label class="form-label">日志级别</label><select class="form-control" name="level"><option value="debug" ' + (data.level === 'debug' ? 'selected' : '') + '>Debug</option><option value="info" ' + (data.level === 'info' ? 'selected' : '') + '>Info</option><option value="warn" ' + (data.level === 'warn' ? 'selected' : '') + '>Warn</option><option value="error" ' + (data.level === 'error' ? 'selected' : '') + '>Error</option></select></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('log-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            logFilePath: formData.get('logFilePath'),
                            maxLogSize: parseInt(formData.get('maxLogSize')),
                            level: formData.get('level')
                        };

                        fetch('/admin/api/config/log', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('日志配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadIPFilterConfig() {
            fetch('/admin/api/config/ipfilter')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>IP过滤配置</h2>' +
                        '<form id="ipfilter-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用IP过滤</label></div>' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enableAllowList" ' + (data.enableAllowList ? 'checked' : '') + '><label class="form-check-label">启用白名单</label></div>' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enableBlockList" ' + (data.enableBlockList ? 'checked' : '') + '><label class="form-check-label">启用黑名单</label></div>' +
                        '<div class="mb-3"><label class="form-label">IP过滤文件路径</label><input type="text" class="form-control" name="ipFilterFile" value="' + (data.ipFilterFile || '') + '"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('ipfilter-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            enableAllowList: formData.has('enableAllowList'),
                            enableBlockList: formData.has('enableBlockList'),
                            ipFilterFile: formData.get('ipFilterFile')
                        };

                        fetch('/admin/api/config/ipfilter', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('IP过滤配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        function loadOutboundConfig() {
            fetch('/admin/api/config/outbound')
                .then(res => res.json())
                .then(data => {
                    const content = document.getElementById('main-content');
                    content.innerHTML = '<h2>出站代理配置</h2>' +
                        '<form id="outbound-form">' +
                        '<div class="mb-3 form-check"><input type="checkbox" class="form-check-input" name="enabled" ' + (data.enabled ? 'checked' : '') + '><label class="form-check-label">启用出站代理</label></div>' +
                        '<div class="mb-3"><label class="form-label">代理URL</label><input type="text" class="form-control" name="url" value="' + (data.url || '') + '" placeholder="socks5://127.0.0.1:1080 或 http://127.0.0.1:7890"></div>' +
                        '<button type="submit" class="btn btn-primary">保存配置</button>' +
                        '</form>';

                    document.getElementById('outbound-form').onsubmit = function(e) {
                        e.preventDefault();
                        const formData = new FormData(e.target);
                        const config = {
                            enabled: formData.has('enabled'),
                            url: formData.get('url')
                        };

                        fetch('/admin/api/config/outbound', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify(config)
                        })
                        .then(res => res.json())
                        .then(data => {
                            if(data.success) {
                                alert('出站代理配置保存成功！');
                            } else {
                                alert('配置保存失败：' + data.message);
                            }
                        });
                    };
                });
        }

        // 页面加载时显示仪表盘
        loadPage('dashboard');
    </script>
</body>
</html>`
	c.String(200, html)
}

// GetConfigHandler 获取完整配置
func GetConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg)
}

// UpdateConfigHandler 更新完整配置
func UpdateConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var newConfig config.Config
	if err := json.Unmarshal(c.Request.Body, &newConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	// 保存配置到文件
	if err := newConfig.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	cfg = &newConfig
	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "配置更新成功",
	})
}

// GetWebsiteConfigHandler 获取网站配置
func GetWebsiteConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Website)
}

// UpdateWebsiteConfigHandler 更新网站配置
func UpdateWebsiteConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var websiteConfig config.WebsiteConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&websiteConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Website = websiteConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "网站配置更新成功",
	})
}

// GetServerConfigHandler 获取服务器配置
func GetServerConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Server)
}

// UpdateServerConfigHandler 更新服务器配置
func UpdateServerConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var serverConfig config.ServerConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&serverConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Server = serverConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "服务器配置更新成功，需要重启服务生效",
	})
}

// GetAuthConfigHandler 获取认证配置
func GetAuthConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Auth)
}

// UpdateAuthConfigHandler 更新认证配置
func UpdateAuthConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var authConfig config.AuthConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&authConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Auth = authConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "认证配置更新成功",
	})
}

// ReloadConfigHandler 重载配置
func ReloadConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	newConfig, err := config.LoadConfig(cfgFile)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "重载配置失败: " + err.Error(),
		})
		return
	}

	cfg = newConfig
	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "配置重载成功",
	})
}

// GetSystemStatusHandler 获取系统状态
func GetSystemStatusHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, map[string]interface{}{
		"serverPort":       cfg.Server.Port,
		"authEnabled":      cfg.Auth.Enabled,
		"rateLimitEnabled": cfg.RateLimit.Enabled,
		"dockerEnabled":    cfg.Docker.Enabled,
		"whitelistEnabled": cfg.Whitelist.Enabled,
		"blacklistEnabled": cfg.Blacklist.Enabled,
		"pagesMode":        cfg.Pages.Mode,
		"pagesTheme":       cfg.Pages.Theme,
		"adminEnabled":     cfg.Admin.Enabled,
		"sponsorEnabled":   cfg.Sponsor.Enabled,
	})
}

// GetPagesConfigHandler 获取页面配置
func GetPagesConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Pages)
}

// UpdatePagesConfigHandler 更新页面配置
func UpdatePagesConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var pagesConfig config.PagesConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&pagesConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Pages = pagesConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "页面配置更新成功",
	})
}

// GetSponsorConfigHandler 获取赞助配置
func GetSponsorConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Sponsor)
}

// UpdateSponsorConfigHandler 更新赞助配置
func UpdateSponsorConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var sponsorConfig config.SponsorConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&sponsorConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Sponsor = sponsorConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "赞助配置更新成功",
	})
}

// GetDockerConfigHandler 获取Docker配置
func GetDockerConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Docker)
}

// UpdateDockerConfigHandler 更新Docker配置
func UpdateDockerConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var dockerConfig config.DockerConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&dockerConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Docker = dockerConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Docker配置更新成功",
	})
}

// GetRateLimitConfigHandler 获取速率限制配置
func GetRateLimitConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.RateLimit)
}

// UpdateRateLimitConfigHandler 更新速率限制配置
func UpdateRateLimitConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var rateLimitConfig config.RateLimitConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&rateLimitConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.RateLimit = rateLimitConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "速率限制配置更新成功",
	})
}

// GetWhitelistConfigHandler 获取白名单配置
func GetWhitelistConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Whitelist)
}

// UpdateWhitelistConfigHandler 更新白名单配置
func UpdateWhitelistConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var whitelistConfig config.WhitelistConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&whitelistConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Whitelist = whitelistConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "白名单配置更新成功",
	})
}

// GetBlacklistConfigHandler 获取黑名单配置
func GetBlacklistConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Blacklist)
}

// UpdateBlacklistConfigHandler 更新黑名单配置
func UpdateBlacklistConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var blacklistConfig config.BlacklistConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&blacklistConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Blacklist = blacklistConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "黑名单配置更新成功",
	})
}

// GetHttpcConfigHandler 获取HTTP客户端配置
func GetHttpcConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Httpc)
}

// UpdateHttpcConfigHandler 更新HTTP客户端配置
func UpdateHttpcConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var httpcConfig config.HttpcConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&httpcConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Httpc = httpcConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "HTTP客户端配置更新成功",
	})
}

// GetGitCloneConfigHandler 获取Git克隆配置
func GetGitCloneConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.GitClone)
}

// UpdateGitCloneConfigHandler 更新Git克隆配置
func UpdateGitCloneConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var gitCloneConfig config.GitCloneConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&gitCloneConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.GitClone = gitCloneConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Git克隆配置更新成功",
	})
}

// GetShellConfigHandler 获取Shell配置
func GetShellConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Shell)
}

// UpdateShellConfigHandler 更新Shell配置
func UpdateShellConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var shellConfig config.ShellConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&shellConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Shell = shellConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Shell配置更新成功",
	})
}

// GetLogConfigHandler 获取日志配置
func GetLogConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Log)
}

// UpdateLogConfigHandler 更新日志配置
func UpdateLogConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var logConfig config.LogConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&logConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Log = logConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "日志配置更新成功",
	})
}

// GetIPFilterConfigHandler 获取IP过滤配置
func GetIPFilterConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.IPFilter)
}

// UpdateIPFilterConfigHandler 更新IP过滤配置
func UpdateIPFilterConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var ipFilterConfig config.IPFilterConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&ipFilterConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.IPFilter = ipFilterConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "IP过滤配置更新成功",
	})
}

// GetOutboundConfigHandler 获取出站代理配置
func GetOutboundConfigHandler(c *touka.Context) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	c.JSON(200, cfg.Outbound)
}

// UpdateOutboundConfigHandler 更新出站代理配置
func UpdateOutboundConfigHandler(c *touka.Context) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	var outboundConfig config.OutboundConfig
	if err := json.NewDecoder(c.Request.Body).Decode(&outboundConfig); err != nil {
		c.JSON(400, map[string]interface{}{
			"success": false,
			"message": "配置格式错误: " + err.Error(),
		})
		return
	}

	cfg.Outbound = outboundConfig

	// 保存配置到文件
	if err := cfg.WriteConfig(cfgFile); err != nil {
		c.JSON(500, map[string]interface{}{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "出站代理配置更新成功",
	})
}
