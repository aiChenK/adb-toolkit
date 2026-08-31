package controllers

import (
	"adb-toolkit/core"
	"adb-toolkit/dto/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdbController struct{}

func (c *AdbController) Index(ctx *gin.Context) {
	var commandForm request.CommandForm
	if err := ctx.ShouldBind(&commandForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"errMessage": "缺少必要参数: " + err.Error(),
		})
		return
	}

	if commandForm.Op != "free" && commandForm.Ip == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"errMessage": "设备 IP 不能为空",
		})
		return
	}

	if commandForm.Op == "setProxy" && commandForm.ProxyAddr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"errMessage": "代理地址不能为空",
		})
		return
	}

	if (commandForm.Op == "clear" || commandForm.Op == "stop") && commandForm.PackageName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"errMessage": "包名不能为空",
		})
		return
	}

	if commandForm.Port == "" {
		commandForm.Port = "5555"
	}

	output, err := core.Exec(commandForm)
	if err != nil {
		errMsg := output
		if errMsg == "" {
			errMsg = err.Error()
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success":    false,
			"errMessage": errMsg,
			"data":       output,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"errMessage": "",
		"data":       output,
	})
}

// Devices 获取当前连接的所有 ADB 设备
func (c *AdbController) Devices(ctx *gin.Context) {
	devices, err := core.GetDevices()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success":    false,
			"errMessage": err.Error(),
			"data":       []core.DeviceInfo{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"errMessage": "",
		"data":       devices,
	})
}
