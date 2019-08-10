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

//获取单个文章
func GetArticle(ctx *gin.Context) {

	id := com.StrTo(ctx.Param("id")).MustInt()
	valid := validation.Validation{}

	valid.Min(id, 1, "id").Message("ID 必须大于0")

	code := e.InvalidParams
	var (
		msg  string
		data interface{}
	)

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ErrorNotExistArticle
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

//获取多个文章
func GetArticleList(ctx *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	valid := validation.Validation{}

	state := -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	}

	tagId := -1
	if arg := ctx.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签 ID 必须大于0")
	}

	code := e.InvalidParams
	var msg string
	if !valid.HasErrors() {
		data["lists"] = models.GetArticleList(util.GetPage(ctx), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)

		code = e.SUCCESS
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

//新增文章
func AddArticle(ctx *gin.Context) {
	tagId := com.StrTo(ctx.Query("tag_id")).MustInt()
	title := ctx.Query("title")
	desc := ctx.Query("desc")
	content := ctx.Query("content")
	createdBy := ctx.Query("created_by")
	state := com.StrTo(ctx.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message(e.TagIdMustGreaterThanZero)
	valid.Required(title, "title").Message(e.TitleNotEmpty)
	valid.Required(desc, "desc").Message(e.DescNotEmpty)
	valid.Required(content, "content").Message(e.ContentNotEmpty)
	valid.Required(createdBy, "created_by").Message(e.CreatedManNotEmpty)
	valid.Range(state, 0, 1, "state").Message(e.StateMustZeroOrOne)

	code := e.InvalidParams
	msg := ""
	data := make(map[string]interface{})
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ErrorNotExistTag
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

//修改文章
func EditArticle(ctx *gin.Context) {
	id := com.StrTo(ctx.Param("id")).MustInt()
	tagId := com.StrTo(ctx.Query("tag_id")).MustInt()
	title := ctx.Query("title")
	desc := ctx.Query("desc")
	content := ctx.Query("content")
	modifiedBy := ctx.Query("modified_by")

	var (
		state = -1
		code  = e.InvalidParams
		msg   = ""
		valid = validation.Validation{}
	)

	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message(e.StateMustZeroOrOne)
	}

	valid.Min(id, 1, "id").Message(e.IDMustGreaterThanZero)
	valid.MaxSize(title, 100, "title").Message(e.TitleMaxSize100)
	valid.MaxSize(desc, 255, "desc").Message(e.DescMaxSize255)
	valid.MaxSize(content, 65535, "content").Message(e.ContentMaxSize65535)
	valid.Required(modifiedBy, "modified_by").Message(e.ModifiedManNotEmpty)
	valid.MaxSize(modifiedBy, 100, "modified_by").Message(e.ModifiedManMaxSize100)

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}

				if title != "" {
					data["title"] = title
				}

				if desc != "" {
					data["desc"] = desc
				}

				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				code = e.SUCCESS

			} else {
				code = e.ErrorNotExistTag
			}
		} else {
			code = e.ErrorNotExistArticle
		}

		msg = e.GetMsg(code)
	} else {
		msg = util.ValidErrorsToStr(valid.Errors)
	}

	ctx.JSON(http.StatusOK, models.BaseResp{
		Code: code,
		Msg:  msg,
		Data: make(map[string]interface{}),
	})

}

//删除文章
func DeleteArticle(ctx *gin.Context) {
	id := com.StrTo(ctx.Param("id")).MustInt()

	var (
		code  = e.InvalidParams
		valid = validation.Validation{}
		msg   = ""
	)

	valid.Min(id, 1, "id").Message(e.IDMustGreaterThanZero)

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ErrorNotExistArticle
		}
		msg = e.GetMsg(code)
	} else {
		msg = util.ValidErrorsToStr(valid.Errors)
	}

	ctx.JSON(http.StatusOK, models.BaseResp{
		Code: code,
		Msg:  msg,
		Data: make(map[string]interface{}),
	})

}
