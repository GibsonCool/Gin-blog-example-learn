package e

//定义错误码
const (
	UnknowError             = -1111
	SUCCESS                 = 200
	ERROR                   = 500
	InvalidParams           = 400
	InvalidParamsEmptyToken = 4001

	ErrorExistTag        = 10001
	ErrorNotExistTag     = 10002
	ErrorNotExistArticle = 10003

	ErrorAuthCheckTokenFail    = 20001
	ErrorAuthCheckTokenTimeout = 20002
	ErrorAuthToken             = 20003
	ErrorAuth                  = 20004

	//保存图片失败
	ErrorUploadSaveImageFail = 30001
	//检查图片失败
	ErrorUploadCheckImageFail = 30002
	//校验图片错误，图片格式或大小
	ErrorUploadCheckImageFormat = 30003
)
