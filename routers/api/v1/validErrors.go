package v1

import "github.com/astaxie/beego/validation"

//将表单验证的错误信息，转换为string
func ValidErrorsToStr(errs []*validation.Error) (msg string) {
	for _, err := range errs {
		msg += err.Message + ","
	}
	return
}
