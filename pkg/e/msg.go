package e

// 定义错误码对应 msg
var MsgFlags = map[int]string{
	SUCCESS:                     "ok",
	ERROR:                       "fail",
	InvalidParams:               "请求参数错误",
	InvalidParamsEmptyToken:     "请求参数错误,token为空",
	ErrorExistTag:               "已存在该标签名称",
	ErrorExistTagFail:           "获取已存在标签失败",
	ErrorNotExistTag:            "该标签不存在",
	ErrorGetTagsFail:            "获取所有标签失败",
	ErrorCountTagFail:           "统计标签失败",
	ErrorAddTagFail:             "新增标签失败",
	ErrorEditTagFail:            "修改标签失败",
	ErrorDeleteTagFail:          "删除标签失败",
	ErrorExportTagFail:          "导出标签失败",
	ErrorImportTagFail:          "导入标签失败",
	ErrorNotExistArticle:        "该文章不存在",
	ErrorAddArticleFail:         "新增文章失败",
	ErrorDeleteArticleFail:      "删除文章失败",
	ErrorCheckExistArticleFail:  "检查文章是否存在失败",
	ErrorEditArticleFail:        "修改文章失败",
	ErrorCountArticleFail:       "统计文章失败",
	ErrorGetArticlesFail:        "获取多个文章失败",
	ErrorGetArticleFail:         "获取单个文章失败",
	ErrorExportArticleFail:      "导出文章失败",
	ErrorImportArticleFail:      "导入文章失败",
	ErrorGenArticlePosterFail:   "生成文章海报失败",
	ErrorAuthCheckTokenFail:     "Token鉴权失败",
	ErrorAuthCheckTokenTimeout:  "Token已超时",
	ErrorAuthToken:              "Token生成失败",
	ErrorAuth:                   "Token错误",
	ErrorUploadSaveImageFail:    "保存图片失败",
	ErrorUploadCheckImageFail:   "检查图片失败",
	ErrorUploadCheckImageFormat: "校验图片错误，图片格式或大小有问题",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}

// 定义常用错误 tips
const (
	NameNotEmpty   = "名称不能为空"
	NameMaxSize100 = "名称最长为100字符"

	CreatedManNotEmpty   = "创建人不能为空"
	CreatedManMaxSize100 = "创建人最长为100字符"

	ModifiedManNotEmpty   = "修改人不能为空"
	ModifiedManMaxSize100 = "修改人最长为100字符"

	StateMustZeroOrOne       = "状态只允许为0或1"
	TagIdMustGreaterThanZero = "标签ID必须大于0"

	IDNotEmpty            = "ID 不能为空"
	IDMustGreaterThanZero = "ID 必须大于0"

	TitleNotEmpty       = "标题不能为空"
	TitleMaxSize100     = "标题最长为100字符"
	DescNotEmpty        = "简述不能为空"
	DescMaxSize255      = "简述最长为255字符"
	ContentNotEmpty     = "内容不能为空"
	ContentMaxSize65535 = "内容最长为65535字符"

	CoverImageUrlNotEmpty   = "封面图片不能为空"
	CoverImageUrlMaxSize255 = "图片地址最长为255字符"
)
