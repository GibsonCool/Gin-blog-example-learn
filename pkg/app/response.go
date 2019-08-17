package app

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errorCode int, data interface{}) {
	g.C.JSON(httpCode, models.BaseResp{
		Code: errorCode,
		Msg:  e.GetMsg(errorCode),
		Data: data,
	})
}

func (g *Gin) ResponseMsg(httpCode, errorCode int, msg string, data interface{}) {
	g.C.JSON(httpCode, models.BaseResp{
		Code: errorCode,
		Msg:  getMsg(errorCode, msg),
		Data: data,
	})
}

func getMsg(errorCode int, msg string) string {
	if msg != "" {
		return msg
	}
	return e.GetMsg(errorCode)
}
