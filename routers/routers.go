package routers

import (
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/setting"
	"github.com/gin-gonic/gin"
)

//抽离路由规则配置
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(e.SUCCESS, gin.H{
			"msg": "test",
		})
	})

	return r
}
