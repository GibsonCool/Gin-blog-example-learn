package e

//定义错误码
const (
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

	UnknowError = -1111
)
