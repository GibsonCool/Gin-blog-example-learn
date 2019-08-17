package app

import (
	"Gin-blog-example/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindAndValid(ctx *gin.Context, form interface{}) (int, int, string) {
	err := ctx.Bind(form)

	if err != nil {
		return http.StatusBadRequest, e.InvalidParams, ""
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR, ""
	}

	if !check {
		msg := MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.InvalidParams, msg
	}

	return http.StatusOK, e.SUCCESS, ""
}
