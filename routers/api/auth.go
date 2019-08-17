package api

import (
	"Gin-blog-example/pkg/app"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/util"
	"Gin-blog-example/service/auth_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary 获取token
// @Description 通过账号信息获取token用于其他接口鉴权使用
// @Param username query string true "UserName"
// @Param password query string true "PassWord"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /auth [get]
func GetAuth(ctx *gin.Context) {
	appG := app.Gin{C: ctx}
	valid := validation.Validation{}

	username := ctx.Query("username")
	password := ctx.Query("password")

	a := auth{username, password}
	ok, _ := valid.Valid(&a)

	if !ok {
		msg := app.MarkErrors(valid.Errors)
		appG.ResponseMsg(http.StatusOK, e.InvalidParams, msg, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAuthCheckTokenFail, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ErrorAuth, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAuthToken, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
