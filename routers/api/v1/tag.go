package v1

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 获取多个标签
// @Description 可选参数 tagName。获取标签列表
// @Produce json
// @Param name body string false "Name"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags [get]
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

	data["lists"] = models.GetTags(util.GetPage(ctx), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}

// @Summary 新增标签
// @Produce json
// @Param name body int true "Name"
// @Param state body int false "State"
// @Param created_by body string true "CreatedBy"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags [post]
func AddTag(ctx *gin.Context) {
	name := ctx.Query("name")
	state := com.StrTo(ctx.DefaultQuery("state", "0")).MustInt()
	createdBy := ctx.Query("created_by")

	//beego 的表单验证
	valid := validation.Validation{}

	valid.Required(name, "name").Message(e.NameNotEmpty)
	valid.MaxSize(name, 100, "name").Message(e.NameMaxSize100)
	valid.Required(createdBy, "created_by").Message(e.CreatedManNotEmpty)
	valid.MaxSize(createdBy, 100, "name").Message(e.CreatedManMaxSize100)
	valid.Range(state, 0, 1, "state").Message(e.StateMustZeroOrOne)

	code := e.InvalidParams
	var msg string
	//表单参数如果没有错误
	if !valid.HasErrors() {
		//判断 tag 是否被创建过
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ErrorExistTag
		}
		msg = e.GetMsg(code)
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
		msg = util.ValidErrorsToStr(valid.Errors)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": make(map[string]string),
	})

}

// @Summary 修改标签信息
// @Description 根据标签 id 修改标签属性信息
// @Produce json
// @Param id path int true "ID"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags/{id} [put]
func EditTag(ctx *gin.Context) {
	id := com.StrTo(ctx.Param("id")).MustInt()
	name := ctx.Query("name")
	modifiedBy := ctx.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message(e.StateMustZeroOrOne)
	}

	valid.Required(id, "id").Message(e.IDNotEmpty)
	valid.Required(modifiedBy, "modified_by").Message(e.ModifiedManNotEmpty)
	valid.MaxSize(modifiedBy, 100, "modified_by").Message(e.ModifiedManMaxSize100)
	valid.MaxSize(name, 100, "name").Message(e.NameMaxSize100)

	code := e.InvalidParams
	var msg string
	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}

			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
			code = e.SUCCESS
		} else {
			code = e.ErrorNotExistTag
		}

		msg = e.GetMsg(code)
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
		msg = util.ValidErrorsToStr(valid.Errors)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": make(map[string]string),
	})

}

// @Summary 删除标签
// @Description 根据标签id 删除对应标签信息
// @Produce json
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(ctx *gin.Context) {
	id := com.StrTo(ctx.Param("id")).MustInt()

	valid := validation.Validation{}

	valid.Min(id, 1, "id").Message(e.IDMustGreaterThanZero)

	code := e.InvalidParams
	var msg string

	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
			code = e.SUCCESS
		} else {
			code = e.ErrorNotExistTag
		}
		msg = e.GetMsg(code)
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
		msg = util.ValidErrorsToStr(valid.Errors)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": make(map[string]string),
	})
}
