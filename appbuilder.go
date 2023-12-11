package main

import (
	"appBuilder/internal/service/snowflakenode"
	"appBuilder/model"
	"flag"
	"fmt"

	"appBuilder/internal/config"
	"appBuilder/internal/svc"
	"appBuilder/internal/userhandler"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/appbuilder-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	c.UserApi.MaxBytes = 1024 * 1024 * 100 //实体大小限制
	server := rest.MustNewServer(rest.RestConf(c.UserApi))
	defer server.Stop()

	ctx, err := svc.NewServiceContext(c, "userApi")
	if err != nil {
		var n any = fmt.Sprintf("创建服务上下文失败：%s", err.Error())
		panic(n)
	}
	userhandler.RegisterHandlers(server, ctx)

	// 创建雪花节点
	snowflakenode.NewNode(1)

	ctx.DB.AutoMigrate(&model.User{})

	fmt.Printf("Starting server at %s:%d...\n", c.UserApi.Host, c.UserApi.Port)
	server.Start()
}
