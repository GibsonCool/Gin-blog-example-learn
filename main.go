package main

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/routers"
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {

	//将多模块的初始化函数放到启动流程中，可自由控制先后顺序
	setting.Setup()
	logging.Setup()
	models.Setup()

	//简单版本的
	//router := routers.InitRouter()
	//
	//server := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", setting.HttpPort),
	//	Handler:        router,
	//	ReadTimeout:    setting.ReadTimeout,
	//	WriteTimeout:   setting.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//server.RegisterOnShutdown(func() {
	//	logging.Info("程序关闭。。。。")
	//})
	//
	//server.ListenAndServe()

	// 搭配 endless 热启动版本
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Acutual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err:%v", err)
	}
}
