package v1

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取多个文章标签
func GetTags(ctx *gin.Context) {
	// ctx.Query可用于获取?name=test&state=1这类URL参数，而c.DefaultQuery则支持设置一个默认值
	//name :=ctx.DefaultQuery("name","test")
	name := ctx.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state = -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(ctx), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}

//新增文章标签
func AddTag(ctx *gin.Context) {
	name := ctx.Query("name")
	state := com.StrTo(ctx.DefaultQuery("state", "0")).MustInt()
	createdBy := ctx.Query("created_by")

	//beego 的表单验证
	valid := validation.Validation{}

	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "name").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许为0或1")

	code := e.InvalidParams

	//表单参数如果没有错误
	if !valid.HasErrors() {
		//判断 tag 是否被创建过
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ErrorExistTag
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

//修改文章标签
func EditTag(ctx *gin.Context) {

}

//删除文章标签
func DeleteTag(ctx *gin.Context) {

}
