package api

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

//生成用户token
func GetAuth(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")

	valid := validation.Validation{}

	a := auth{username, password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.InvalidParams
	msg := ""

	if ok {

		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ErrorAuthToken
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ErrorAuth
		}
		msg = e.GetMsg(code)
	} else {
		msg = util.ValidErrorsToStr(valid.Errors)
	}

	ctx.JSON(http.StatusOK, models.BaseResp{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
