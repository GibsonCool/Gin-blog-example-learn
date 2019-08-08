package main

import (
	"Gin-blog-example/pkg/e"
	"Gin-blog-example/pkg/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(e.SUCCESS, gin.H{
			"msg": "test",
		})
	})

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
