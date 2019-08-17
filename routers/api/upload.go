package api

import (
	"Gin-blog-example/pkg/app"
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(ctx *gin.Context) {
	appG := app.Gin{C: ctx}

	//获取上传文件内容
	file, image, err := ctx.Request.FormFile("image")

	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if image == nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()

	src := fullPath + imageName
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, e.ErrorUploadCheckImageFormat, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ErrorUploadCheckImageFail, nil)
		return
	}

	if err := ctx.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ErrorUploadSaveImageFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
	})
}
