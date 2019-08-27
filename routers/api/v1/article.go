package v1

import (
	"Gin-blog-example/pkg/app"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/export"
	"Gin-blog-example/pkg/qrcode"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/util"
	"Gin-blog-example/service/article_service"
	"Gin-blog-example/service/tag_service"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	QrcodeUrl = "https://github.com/GibsonCool"
)

// @Summary 获取单个文章
// @Description 通过文章 id 获取文章信息
// @Produce json
// @Param id path int true "ID"
// @Param token query string true "token"
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
// @Accept	mpfd
// @Produce json
// @Param tag_id query int false "TagId"
// @Param state query int false "State"
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles [get]
func GetArticleList(ctx *gin.Context) {
	appG := app.Gin{C: ctx}

	valid := validation.Validation{}

	state := -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	}

	tagId := -1
	if arg := ctx.Query("tag_id"); arg != "" {
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
	TagID         int    `json:"tagID"  valid:"Required;Min(1)"`
	Title         string `json:"title"  valid:"Required;MaxSize(100)"`
	Desc          string `json:"desc"  valid:"Required;MaxSize(255)"`
	Content       string `json:"content"  valid:"Required;MaxSize(65535)"`
	CreatedBy     string `json:"createdBy" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `json:"coverImageUrl"  valid:"Required;MaxSize(255)"`
	State         int    `json:"state"  valid:"Range(0,1)"`
}

// @Summary 新增文章
// @Description 可通过可选参数 tag_id 获取同标签下的所有文章信息
// @Produce json
// @Param article body v1.AddArticleForm true "{json内容}"
// @Param token query string true "token"
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
	ID            int    `form:"id" valid:"Min(1)"`
	TagID         int    `form:"tag_id" valid:"Min(1)"`
	Title         string `form:"title" valid:"MaxSize(100)"`
	Desc          string `form:"desc" valid:"MaxSize(255)"`
	Content       string `form:"content" valid:"MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 修改文章信息
// @Description 根据文章 id 修改文章属性信息
// @Accept	mpfd
// @Produce json
// @Param id path int true "ID"
// @Param tag_id formData int false "TagId"
// @Param state formData int false "State"
// @Param title formData string false "Title"
// @Param desc formData string false "Desc"
// @Param content formData string false "Content"
// @Param modified_by formData string false "ModifiedBy"
// @Param cover_image_url formData string false "CoverImageUrl"
// @Param token query string true "token"
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
// @Param token query string true "token"
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

// @Summary 导出文章
// @Description 导出所有文章为 .xlsx 文件
// @Accept	mpfd
// @Produce json
// @Param title formData string false "文章标题"
// @Param state formData int false "State"
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles/export [post]
func ExportArticle(ctx *gin.Context) {
	appG := app.Gin{C: ctx}
	title := ctx.PostForm("title")
	state := -1
	if arg := ctx.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	articleService := article_service.ArticleService{
		Title:    title,
		State:    state,
		PageNum:  util.GetPage(ctx),
		PageSize: setting.AppSetting.PageSize,
	}

	filename, err := articleService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorExportArticleFail, err.Error())
		return
	}
	data := map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_sava_url": export.GetExcelFullPath() + filename,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary 生成海报二维码
// @Description  生成二维码图片
// @Produce json
// @Param token query string true "token"
// @Success 200 {object} models.BaseResp
// @Failure 500 {object} models.BaseResp
// @Router /api/v1/articles/poster/generate [post]
func GenerateArticlePoster(ctx *gin.Context) {
	appG := app.Gin{C: ctx}

	articleService := &article_service.ArticleService{}

	qr := qrcode.NewQrcode(QrcodeUrl, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.Url) + qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, articleService, qr)

	articlePosterBg := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := articlePosterBg.Generate()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGenArticlePosterFail, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}
