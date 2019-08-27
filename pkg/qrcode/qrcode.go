package qrcode

import (
	"Gin-blog-example/pkg/file"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
)

type Qrcode struct {
	Url    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	ExtTag = ".jpg"
)

func (q *Qrcode) GetQrCodeExt() string {
	return q.Ext
}

func (q *Qrcode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.Url) + q.GetQrCodeExt()
	if file.CheckNotExist(src) == true {
		return false
	}
	return true
}

// 根据属性生成二维码
func (q *Qrcode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.Url) + q.GetQrCodeExt()
	src := path + name

	//判断下文件是不是还没有存在
	if file.CheckNotExist(src) == true {
		//生成二维码
		code, e := qr.Encode(q.Url, q.Level, q.Mode)
		if e != nil {
			return "", "", e
		}

		// 按照宽高裁剪
		code, e = barcode.Scale(code, q.Width, q.Height)
		if e != nil {
			return "", "", e
		}

		//创建二维码图片文件
		f, e := file.MustOpen(name, path)
		if e != nil {
			return "", "", e
		}
		defer f.Close()

		// 将图像写入文件中
		// Encode 方法会将 图像（二维码）以 JPEG 4：2：0 基线格式写入文件
		// 第三个参数nil 会使用默认值75 。取值区间在1-100 值越大图片质量越高
		e = jpeg.Encode(f, code, nil)
		if e != nil {
			return "", "", e
		}
	}

	return name, path, nil
}

func NewQrcode(url string, width int, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *Qrcode {
	return &Qrcode{Url: url, Width: width, Height: height, Level: level, Mode: mode, Ext: ExtTag}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}
