package routers

import (
	_ "Gin-blog-example/docs"
	"Gin-blog-example/middleware/jwt"
	"Gin-blog-example/pkg/export"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/upload"
	"Gin-blog-example/routers/api"
	v1 "Gin-blog-example/routers/api/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

//抽离路由规则配置
func InitRouter() *gin.Engine {
	//设置 gin mode 需要放在 New() 之前才生效
	gin.SetMode(setting.ServerSetting.RunMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//配置支持静态资源--》图片的访问
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))

	r.GET("/auth", api.GetAuth)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")

	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		apiv1.POST("/tags/export", v1.ExportTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticleList)
		//根据id获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//根据id更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//根据id删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
