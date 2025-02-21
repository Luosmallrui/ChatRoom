package main

import (
	"chatroom/config"
	"chatroom/pkg/core"
	"chatroom/socket"
	"github.com/urfave/cli/v2"
)

// 创建 HTTP 服务命令
func NewHttpServerCommand() core.Command {
	return core.Command{
		Name:  "http",
		Usage: "Start the HTTP server",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			httpApp := NewHttpInjector(conf)
			httpApp.RegisterRoutes()
			if err := core.Run(ctx, httpApp); err != nil {
				return err
			}
			return nil
		},
	}
}

// 创建 WebSocket 服务命令
func NewWebSocketServerCommand() core.Command {
	return core.Command{
		Name:  "websocket",
		Usage: "Start the WebSocket server",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			socketApp := NewSocketInjector(conf)
			if err := socket.Run(ctx, socketApp); err != nil {
				return err
			}
			return nil
		},
	}
}

func main() {
	// 创建 CLI 应用
	app := core.NewApp("v1.0.0")
	app.Register(NewHttpServerCommand)
	app.Register(NewWebSocketServerCommand)

	// 运行命令
	app.Run()
}
