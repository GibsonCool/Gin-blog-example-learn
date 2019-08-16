package e

// 定义错误码对应 msg
var MsgFlags = map[int]string{
	SUCCESS:                     "ok",
	ERROR:                       "fail",
	InvalidParams:               "请求参数错误",
	InvalidParamsEmptyToken:     "请求参数错误,token为空",
	ErrorExistTag:               "已存在该标签名称",
	ErrorNotExistTag:            "该标签不存在",
	ErrorNotExistArticle:        "该文章不存在",
	ErrorAuthCheckTokenFail:     "Token鉴权失败",
	ErrorAuthCheckTokenTimeout:  "Token已超时",
	ErrorAuthToken:              "Token生成失败",
	ErrorAuth:                   "Token错误",
	UnknowError:                 "未定义错误",
	ErrorUploadSaveImageFail:    "保存图片失败",
	ErrorUploadCheckImageFail:   "检查图片失败",
	ErrorUploadCheckImageFormat: "校验图片错误，图片格式或大小不符合",
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
)
