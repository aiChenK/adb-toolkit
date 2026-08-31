package controllers

import (
	"adb-toolkit/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ToolController struct{}

func (c *ToolController) Ip(ctx *gin.Context) {
	ip, err := core.GetLocalIp()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success":    false,
			"errMessage": err.Error(),
			"data":       "127.0.0.1",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"errMessage": "",
		"data":       ip,
	})
}
