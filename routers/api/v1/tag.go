package v1

import (
	"Gin-blog-example/pkg/app"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/export"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/util"
	"Gin-blog-example/service/tag_service"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 获取多个标签
// @Description 可选参数 tagName。获取标签列表
// @Produce json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags [get]
func GetTags(ctx *gin.Context) {
	appG := app.Gin{C: ctx}

	// ctx.Query可用于获取?name=test&state=1这类URL参数，而c.DefaultQuery则支持设置一个默认值
	//name :=ctx.DefaultQuery("name","test")
	name := ctx.Query("name")
	state := -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(ctx),
		PageSize: setting.AppSetting.PageSize,
	}

	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorGetTagsFail, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCountTagFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})

}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 新增标签
// @Produce json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query string true "CreatedBy"
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags [post]
func AddTag(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
		form AddTagForm
	)

	httpCode, errCode, msg := app.BindAndValid(ctx, &form)
	if errCode != e.SUCCESS {
		appG.ResponseMsg(httpCode, errCode, msg, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorExistTagFail, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ErrorExistTag, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAddTagFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 修改标签信息
// @Description 根据标签 id 修改标签属性信息
// @Produce json
// @Param id path int true "ID"
// @Param name query string false "name"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags/{id} [put]
func EditTag(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
		form = EditTagForm{ID: com.StrTo(ctx.Param("id")).MustInt()}
	)

	httpCode, errCode, msg := app.BindAndValid(ctx, &form)
	if errCode != e.SUCCESS {
		appG.ResponseMsg(httpCode, errCode, msg, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorExistTagFail, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistTag, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorEditTagFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除标签
// @Description 根据标签id 删除对应标签信息
// @Produce json
// @Param id path int true "ID"
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(ctx *gin.Context) {
	appG := app.Gin{C: ctx}
	valid := validation.Validation{}
	id := com.StrTo(ctx.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		msg := app.MarkErrors(valid.Errors)
		appG.ResponseMsg(http.StatusBadRequest, e.InvalidParams, msg, nil)
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorExistTagFail, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistTag, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorDeleteTagFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 导出标签
// @Description 导出所有标签为 .xlsx 文件
// @Accept	mpfd
// @Produce json
// @Param name formData string false "Name"
// @Param state formData int false "State"
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/tags/export [post]
func ExportTag(ctx *gin.Context) {
	appG := app.Gin{C: ctx}

	name := ctx.PostForm("name")
	state := -1
	if arg := ctx.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	tageService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tageService.ExportByExcelize()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorExportTagFail, err.Error())
		return
	}

	data := map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_sava_url": export.GetExcelFullPath() + filename,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

//@Summary 导入标签文件
//@Accept	mpfd
//@Produce  json
//@Param file formData file true "Excel file"
//@Param token query string true "token"
//@Success 200 {object} models.BaseResp
//@Failure 500 {object} models.BaseResp
//@Router /api/v1/tags/import [post]
func ImportTag(ctx *gin.Context) {
	appG := app.Gin{C: ctx}
	file, _, err := ctx.Request.FormFile("file")

	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ErrorImportTagFail, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
