package main

import (
	"Nuxus/configs"
	"Nuxus/internal/dao"
	"Nuxus/internal/routers"
	"fmt"
)

func init() {
	configs.Init()
	dao.InitMysql()
	dao.InitRedis()
}

func main() {
	r := routers.SetupRouter()

	r.Run(fmt.Sprintf("127.0.0.1:%d", configs.Conf.Server.Port))
}
