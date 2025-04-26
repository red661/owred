package middleware

import (
	"hoyang/ownsa/data/response"
	"hoyang/ownsa/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

// ErrorHandler 是一个处理错误的中间件，用于捕获并处理运行时错误
func ErrorHandler(c *gin.Context, err interface{}) {
	// 将传入的错误包裹，并添加堆栈信息，便于调试
	goErr := errors.Wrap(err, 2)

	// 生成一个随机字符串作为错误标识
	randStr, err := utils.GenerateRandomString(8)

	// 构建返回的错误响应结构体
	webResponse := response.Response{
		Code:    500,           // 错误代码：500 表示服务器内部错误
		Success: false,         // 请求未成功
		Message: goErr.Error(), // 错误信息：包含堆栈跟踪的错误信息
		Data:    randStr,       // 返回一个随机字符串，用作错误标识，方便后续跟踪
	}

	// 中止当前请求的处理并返回 JSON 格式的错误响应
	c.AbortWithStatusJSON(http.StatusOK, webResponse)
}
