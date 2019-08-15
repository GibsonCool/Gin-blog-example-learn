package main

import (
	"Gin-blog-example/models"
	"github.com/robfig/cron"
	"log"
)

/*
	使用 cron 一个简单的定时任务调度管理
*/
func main() {
	log.Println("Starting.....")

	c := cron.New()
	c.AddFunc("0/10 * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})

	c.AddFunc("0/15 * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()

	//创建一个定时器
	//timer := time.NewTimer(time.Second * 10)

	//主要用于阻塞主程序，cron 的定时任务是开协程去做的，
	//如果不阻塞主程序，定时任务还没开始就主程序退出了
	for {
		select {
		//case <-timer.C:
		//	timer.Reset(time.Second*10)
		}
	}
}
