package routers

import (
	"Gin-blog-example/pkg/setting"
	v1 "Gin-blog-example/routers/api/v1"
	"github.com/gin-gonic/gin"
)

//抽离路由规则配置
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	// 测试接口
	//r.GET("/test", func(c *gin.Context) {
	//	c.JSON(e.SUCCESS, gin.H{
	//		"msg": "test",
	//	})
	//})

	apiv1 := r.Group("/api/v1")
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
	}

	return r
}
