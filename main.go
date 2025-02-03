package main

import (
	"chatroom/config"
	"chatroom/pkg/core"
	"chatroom/socket"
	"github.com/urfave/cli/v2"
	"os"
)

func NewHttpCommand() core.Command {
	return core.Command{
		Name:  "Gin",
		Usage: "Gin Command - Gin Web Sever",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			//logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "http")
			// 初始化依赖注入
			app := NewHttpInjector(conf)
			// 注册路由
			app.RegisterRoutes()
			// 启动 HTTP 服务
			return core.Run(ctx, app)
		},
	}
}

func NewWebSocketCommand() core.Command {
	return core.Command{
		Name:  "websocket",
		Usage: "websocket Command",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			//logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "http")
			// 初始化依赖注入
			app := NewSocketInjector(conf)
			// 启动 HTTP 服务
			return socket.Run(ctx, app)
		},
	}
}

func main() {
	app := core.NewApp("v1.0.0")
	app.Register(NewHttpCommand)
	app.Register(NewWebSocketCommand)
	os.Args = append(os.Args, "websocket")
	app.Run()

}
