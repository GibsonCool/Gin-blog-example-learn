package jwt

import (
	"Gin-blog-example/pkg/app"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// jwt 用户信息校验中间件
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := e.SUCCESS
		token := ctx.Query("token")

		if token == "" {
			code = e.InvalidParamsEmptyToken
		} else {
			//解析token
			claims, err := util.ParseToken(token)
			if err != nil {
				//验证失败
				logging.Warn(err)
				code = e.ErrorAuthCheckTokenFail
			} else if time.Now().Unix() > claims.ExpiresAt {
				// token 超时
				code = e.ErrorAuthCheckTokenTimeout
			}
		}

		if code != e.SUCCESS {
			//token校验失败,设置返回内容
			appG := app.Gin{C: ctx}
			appG.Response(http.StatusUnauthorized, code, nil)

			//校验失败，停止调用后续的处理
			ctx.Abort()
			return
		}

		//校验通过继续下一个处理程序
		ctx.Next()
	}
}
