package main

import (
	"adb-toolkit/core"
	"adb-toolkit/routes"
	"embed"
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/*
var Static embed.FS

func main() {
	portFlag := flag.String("port", "8088", "HTTP 服务监听端口")
	noBrowserFlag := flag.Bool("no-browser", false, "是否禁用自动打开浏览器")
	flag.Parse()

	// 环境检测
	core.CheckAdb()

	// 生产模式隐藏 gin 调试日志
	gin.SetMode(gin.ReleaseMode)

	// 创建路由引擎
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册业务路由
	routes.InitRoute(r)

	// 静态资源服务 (内嵌前端产物)
	r.Use(routes.Serve("/", routes.EmbedFolder(Static, "web/dist")))
	r.NoRoute(func(c *gin.Context) {
		data, err := Static.ReadFile("web/dist/index.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	serverURL := fmt.Sprintf("http://localhost:%s", *portFlag)

	fmt.Println("==================================================")
	fmt.Println(" 🛠️  ADB Toolkit 可视化调试工具箱已启动")
	fmt.Printf(" 🌐 访问地址: %s\n", serverURL)
	fmt.Println("==================================================")

	if !*noBrowserFlag {
		go func() {
			core.OpenBrowser(serverURL)
		}()
	}

	if err := r.Run(":" + *portFlag); err != nil {
		core.PrintErr(fmt.Sprintf("❌ 服务启动失败: %v", err))
	}
}
