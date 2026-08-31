package controllers

import (
	"adb-toolkit/core"
	"adb-toolkit/dto/request"
	"adb-toolkit/dto/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdbController struct{}

// Index 执行 ADB 指令
func (c *AdbController) Index(ctx *gin.Context) {
	var commandForm request.CommandForm
	if err := ctx.ShouldBind(&commandForm); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "缺少必要参数: "+err.Error())
		return
	}

	if commandForm.Op != "free" && commandForm.Ip == "" {
		response.Fail(ctx, http.StatusBadRequest, "设备 IP 不能为空")
		return
	}

	if commandForm.Op == "setProxy" && commandForm.ProxyAddr == "" {
		response.Fail(ctx, http.StatusBadRequest, "代理地址不能为空")
		return
	}

	if (commandForm.Op == "clear" || commandForm.Op == "stop") && commandForm.PackageName == "" {
		response.Fail(ctx, http.StatusBadRequest, "包名不能为空")
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
		response.FailWithData(ctx, http.StatusOK, errMsg, output)
		return
	}

	response.Success(ctx, output)
}

// Devices 获取当前连接的所有 ADB 设备
func (c *AdbController) Devices(ctx *gin.Context) {
	devices, err := core.GetDevices()
	if err != nil {
		response.FailWithData(ctx, http.StatusOK, err.Error(), []core.DeviceInfo{})
		return
	}

	response.Success(ctx, devices)
}
