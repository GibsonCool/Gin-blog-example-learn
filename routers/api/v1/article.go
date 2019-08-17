package v1

import (
	"Gin-blog-example/pkg/app"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/util"
	"Gin-blog-example/service/article_service"
	"Gin-blog-example/service/tag_service"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 获取单个文章
// @Description 通过文章 id 获取文章信息
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles/{id} [get]
func GetArticle(ctx *gin.Context) {
	appG := app.Gin{C: ctx}

	id := com.StrTo(ctx.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID 必须大于0")

	if valid.HasErrors() {
		msg := app.MarkErrors(valid.Errors)
		appG.ResponseMsg(http.StatusBadRequest, e.InvalidParams, msg, nil)
		return
	}

	articleService := article_service.ArticleService{ID: id}
	existByID, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCheckExistArticleFail, nil)
		return
	}

	if !existByID {
		appG.Response(http.StatusOK, e.ErrorNotExistArticle, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorGetArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// @Summary 获取多个文章
// @Description 可通过可选参数获取符合条件的文章列表
// @Produce json
// @Param tag_id body int false "TagId"
// @Param state body int false "State"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles [get]
func GetArticleList(ctx *gin.Context) {
	appG := app.Gin{C: ctx}

	valid := validation.Validation{}

	state := -1
	if arg := ctx.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	}

	tagId := -1
	if arg := ctx.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签 ID 必须大于0")
	}

	if valid.HasErrors() {
		msg := app.MarkErrors(valid.Errors)
		appG.ResponseMsg(http.StatusBadRequest, e.InvalidParams, msg, nil)
		return
	}

	articleService := article_service.ArticleService{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(ctx),
		PageSize: setting.AppSetting.PageSize,
	}

	count, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCountArticleFail, nil)
		return
	}
	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorGetArticlesFail, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = count
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 新增文章
// @Description 可通过可选参数 tag_id 获取同标签下的所有文章信息
// @Produce json
// @Param tag_id body int true "TagId"
// @Param state body int true "State"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param cover_image_url body string true "CoverImageUrl"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles [post]
func AddArticle(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
		form AddArticleForm
	)

	httpCode, errCode, msg := app.BindAndValid(ctx, &form)
	if errCode != e.SUCCESS {
		appG.ResponseMsg(httpCode, errCode, msg, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	existByTagID, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorExistTagFail, nil)
		return
	}

	if !existByTagID {
		appG.Response(http.StatusOK, e.ErrorNotExistTag, nil)
		return
	}

	articleService := article_service.ArticleService{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	if err = articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAddArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 修改文章信息
// @Description 根据文章 id 修改文章属性信息
// @Produce json
// @Param id path int true "ID"
// @Param tag_id body int false "TagId"
// @Param state body int false "State"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string false "ModifiedBy"
// @Param cover_image_url body string false "CoverImageUrl"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles/{id} [put]
func EditArticle(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
		form = EditArticleForm{ID: com.StrTo(ctx.Param("id")).MustInt()}
	)

	httpCode, errCode, msg := app.BindAndValid(ctx, &form)

	if errCode != e.SUCCESS {
		appG.ResponseMsg(httpCode, errCode, msg, nil)
		return
	}

	articleService := article_service.ArticleService{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}

	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCheckExistArticleFail, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistArticle, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorExistTagFail, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistTag, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorEditArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// @Summary 删除文章
// @Description 根据文章id 删除对应文章信息
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(ctx *gin.Context) {
	appG := app.Gin{C: ctx}
	id := com.StrTo(ctx.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message(e.IDMustGreaterThanZero)

	if valid.HasErrors() {
		msg := app.MarkErrors(valid.Errors)
		appG.ResponseMsg(http.StatusOK, e.InvalidParams, msg, nil)
	}

	articleService := article_service.ArticleService{ID: id}

	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCheckExistArticleFail, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistArticle, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorDeleteArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
