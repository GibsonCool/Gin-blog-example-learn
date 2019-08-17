package app

import (
	"Gin-blog-example/pkg/logging"
	"github.com/astaxie/beego/validation"
)

//将表单验证的错误信息，转换为string
func MarkErrors(errors []*validation.Error) (msg string) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
		msg += err.Message + ","
	}
	return
}
