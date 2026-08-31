package controllers

import (
	"adb-toolkit/core"
	"adb-toolkit/dto/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ToolController struct{}

// Ip 获取本机局域网 IP
func (c *ToolController) Ip(ctx *gin.Context) {
	ip, err := core.GetLocalIp()
	if err != nil {
		response.FailWithData(ctx, http.StatusOK, err.Error(), "127.0.0.1")
		return
	}

	response.Success(ctx, ip)
}
