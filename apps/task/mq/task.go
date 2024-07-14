/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package main

import (
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/apps/task/mq/internal/handler"
	"easy-chat/apps/task/mq/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)

	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()

	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ......")
	serviceGroup.Start()

}
