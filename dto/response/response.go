package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 API 响应格式
type Response struct {
	Success    bool        `json:"success"`
	ErrMessage string      `json:"errMessage"`
	Data       interface{} `json:"data"`
}

// Success 成功响应 (HTTP 200)
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Success:    true,
		ErrMessage: "",
		Data:       data,
	})
}

// Fail 失败响应
func Fail(ctx *gin.Context, statusCode int, errMsg string) {
	ctx.JSON(statusCode, Response{
		Success:    false,
		ErrMessage: errMsg,
		Data:       nil,
	})
}

// FailWithData 带附带数据的失败响应
func FailWithData(ctx *gin.Context, statusCode int, errMsg string, data interface{}) {
	ctx.JSON(statusCode, Response{
		Success:    false,
		ErrMessage: errMsg,
		Data:       data,
	})
}
