/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package main

import (
	"easy-chat/apps/im/ws/internal/config"
	"easy-chat/apps/im/ws/internal/handler"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)

	srv := websocket.NerServer(c.ListenOn, websocket.WithServerAuthentication(handler.NewJwtAuth(
		ctx)))
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn, "......")
	srv.Start()

}
