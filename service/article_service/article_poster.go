package article_service

import (
	"Gin-blog-example/pkg/file"
	"Gin-blog-example/pkg/qrcode"
	"Gin-blog-example/pkg/setting"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
)

type ArticlePoster struct {
	PosterName string
	*ArticleService
	Qr *qrcode.Qrcode
}

func NewArticlePoster(posterName string, articleService *ArticleService, qr *qrcode.Qrcode) *ArticlePoster {
	return &ArticlePoster{PosterName: posterName, ArticleService: articleService, Qr: qr}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) CheckMergedImage(path string) bool {
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}
	return true
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	file, e := file.MustOpen(a.PosterName, path)
	if e != nil {
		return nil, e
	}
	return file, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

func NewArticlePosterBg(name string, articlePoster *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{Name: name, ArticlePoster: articlePoster, Rect: rect, Pt: pt}
}

// 绘制二维码，合并海报
func (a *ArticlePosterBg) Generate() (string, string, error) {
	// 1.获取二维码吗存储路径
	fullPath := qrcode.GetQrCodeFullPath()
	// 2.生成二维码图像
	fileName, path, e := a.Qr.Encode(fullPath)
	if e != nil {
		return "", "", e
	}

	// 3.判断海报是否已经存在
	if !a.CheckMergedImage(path) {
		// 4.如果不存在，则生成待合并的图像 mergedF
		mergedF, e := a.OpenMergedImage(path)
		if e != nil {
			return "", "", e
		}
		defer mergedF.Close()

		// 5.打开事先存放的背景图片 bgF
		bgF, e := file.MustOpen(a.Name, path)
		defer bgF.Close()

		// 6.打开（2）中生成的二维码图像 qrF
		qrF, e := file.MustOpen(fileName, path)
		if e != nil {
			return "", "", e
		}
		defer qrF.Close()

		// 7.将 brF 和 qrF 文件解码为返回 image.Image
		bgImage, e := jpeg.Decode(bgF)
		if e != nil {
			return "", "", e
		}

		qrImage, e := jpeg.Decode(qrF)
		if e != nil {
			return "", "", e
		}

		// 8.创建一个新的 RGBA 图像 jpg
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

		// 9.在 jpg 图像上绘制 背景图 bgF
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		// 10.在（9）的基础上在指定的 point 上绘制二维码图像 qrF
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		// 11.将绘制完成的图像 以 JPEG 4：2：0 基线格式写入文件
		//jpeg.Encode(mergedF, jpg, &jpeg.Options{Quality: jpeg.DefaultQuality})

		// 11.绘制文字
		e = a.DrawPoster(&DrawText{
			JPG:    jpg,
			Merged: mergedF,

			Title: "Golang Gin 系列学习项目",
			X0:    80,
			Y0:    160,
			Size0: 55,

			SubTitle: "---DoubleX",
			X1:       320,
			Y1:       220,
			Size1:    30,
		}, "HanyiSentyCrayon.ttf", "Radam.ttf")

		if e != nil {
			return "", "", e
		}

	}
	return fileName, path, nil
}

/*
	绘制文字
*/
type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

// 绘制文字
func (a *ArticlePosterBg) DrawPoster(d *DrawText, titleFontName, subTitleFontName string) error {
	titleFontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + titleFontName
	subTitleFontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + subTitleFontName

	// 读取字体资源文件
	titlefontSourceBytes, err := ioutil.ReadFile(titleFontSource)
	if err != nil {
		return err
	}

	subTitlefontSourceBytes, err := ioutil.ReadFile(subTitleFontSource)
	if err != nil {
		return err
	}

	//使用 freeType 字体驱动库解析字体资源文件
	titleFont, err := freetype.ParseFont(titlefontSourceBytes)
	if err != nil {
		return err
	}

	subTitleFont, err := freetype.ParseFont(subTitlefontSourceBytes)
	if err != nil {
		return err
	}

	fc := freetype.NewContext()
	fc.SetDPI(72)              //设置分辨率
	fc.SetFont(titleFont)      //设置位置的字体
	fc.SetFontSize(d.Size0)    //设置字体打下
	fc.SetClip(d.JPG.Bounds()) // 设置剪裁区域范围以进行绘制
	fc.SetDst(d.JPG)           //设置目标图像
	fc.SetSrc(image.Black)     //设置绘制操作的原图像，通常为 image.Uniform。这里有点类似Android图像操作的画板的概念

	//设置坐标绘制文字
	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFont(subTitleFont)
	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}
	return nil
}
