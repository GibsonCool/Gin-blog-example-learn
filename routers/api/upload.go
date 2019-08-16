package api

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(ctx *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	//获取上传文件内容
	file, image, err := ctx.Request.FormFile("image")

	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		ctx.JSON(http.StatusOK, models.BaseResp{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: data,
		})
		return
	}

	if image == nil {
		code = e.InvalidParams
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ErrorUploadCheckImageFormat
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ErrorUploadCheckImageFail
			} else if err := ctx.SaveUploadedFile(image, src); err != nil {
				logging.Warn(err)
				code = e.ErrorUploadCheckImageFail
			} else {
				data["imageUrl"] = upload.GetImageFullUrl(imageName)
				data["imageSaveUrl"] = savePath + imageName
			}
		}
	}

	ctx.JSON(http.StatusOK, models.BaseResp{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: data,
	})
}
