package utils

import (
	"log"
	"os"
)

// ErrorPanic 用于检查错误并在发生错误时触发 panic。
// 如果传入的错误不为 nil，则会抛出一个 panic，通常用于无法恢复的错误场景。
// 该函数可用于简化错误检查的代码，并确保错误被立即报告。
// 参数:
//
//	err: 需要检查的错误对象。如果错误不为 nil，则触发 panic。
func ErrorPanic(err error) {
	if err != nil {
		panic(err) // 如果错误存在，则抛出 panic
	}
}

// NewLogger 创建一个新的日志记录器
func NewLogger() *log.Logger {
	return log.New(os.Stdout, "\r\n", log.LstdFlags)
}
