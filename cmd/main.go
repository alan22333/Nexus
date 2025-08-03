package main

import (
	"Nuxus/configs"
	"Nuxus/internal/dao"
	"Nuxus/internal/routers"
	"Nuxus/internal/tasks"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

func init() {
	configs.Init()
	dao.InitMysql()
	dao.InitRedis()
}

func main() {

	// ====Cron goroutine===
	c := cron.New(cron.WithSeconds()) // 6个字段

	// 每小时的第0分0秒执行一次同步任务
	// 格式：秒 分 时 日 月 周
	_, err := c.AddFunc("0 */2 * * * *", tasks.SyncViewCountsToDB)
	if err != nil {
		log.Fatalf("添加定时任务失败: %v", err)
	}
	c.Start()
	log.Println("定时任务启动成功")
	defer c.Stop()

	// =====Gin server=====
	r := routers.SetupRouter()

	err = r.Run(fmt.Sprintf("127.0.0.1:%d", configs.Conf.Server.Port))
	if err != nil {
		log.Fatal("server start failed !")
	}

}
