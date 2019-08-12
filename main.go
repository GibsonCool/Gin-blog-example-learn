package main

import (
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/routers"
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {

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
	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Acutual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err:%v", err)
	}
}
