package upload

import (
	"Gin-blog-example/pkg/file"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/util"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

//获取图片完整访问 url
func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

//获取图片路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

//获取图片名称
func GetImageName(name string) string {
	//获取后缀
	ext := file.GetExt(name)
	//删除后缀 ext 得到文件名称,如果文件不是已 ext 结尾的。返回原始串
	fileName := strings.TrimSuffix(name, ext)
	//对文件名称进行 md5 取摘要
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

// 获取图片完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// 检查 img 是否是配置中支持的格式后缀
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

// 检查是否超出最大限制
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)

	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

//检查图片
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd error : %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err :%v", err)
	}

	permission := file.CheckPermission(src)

	if permission == true {
		return fmt.Errorf("file.CheckPermission Permission denied src :%s", src)
	}
	return nil
}
