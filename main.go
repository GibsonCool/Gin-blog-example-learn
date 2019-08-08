package main

import (
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/routers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	server.RegisterOnShutdown(func() {
		log.Println("程序关闭。。。。")
	})

	server.ListenAndServe()
}
