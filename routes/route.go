package routes

import (
	"adb-toolkit/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoute(engine *gin.Engine) {
	// 添加跨域中间件
	engine.Use(CORSMiddleware())

	adbCtrl := &controllers.AdbController{}
	toolCtrl := &controllers.ToolController{}

	engine.POST("/adb", adbCtrl.Index)
	engine.GET("/devices", adbCtrl.Devices)
	engine.GET("/ip", toolCtrl.Ip)
}

// CORSMiddleware 是跨域中间件函数
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
