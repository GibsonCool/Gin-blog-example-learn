package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

/*
	文件操作相关
*/

//获取文件大小
func GetSize(file multipart.File) (int, error) {
	content, e := ioutil.ReadAll(file)
	return len(content), e
}

//获取文件后缀
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

//检查文件是否存在
func CheckNotExist(src string) bool {
	_, e := os.Stat(src)
	return os.IsNotExist(e)
}

//检查文件权限
func CheckPermission(src string) bool {
	_, e := os.Stat(src)
	return os.IsPermission(e)
}

//新建文件夹
func MkDir(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

//如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		return MkDir(src)
	}
	return nil
}

//打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}
