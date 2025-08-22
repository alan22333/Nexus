package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

func main() {
	// 通过Wire初始化整个应用
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// 启动定时任务
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("0 */2 * * * *", app.SyncTask.SyncViewCountsToDB)
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	c.Start()
	defer c.Stop()

	// 启动Web服务
	router := app.Router.SetupRouter()
	log.Printf("Server starting on :%d", app.Config.Server.Port)
	if err := router.Run(fmt.Sprintf(":%d", app.Config.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
