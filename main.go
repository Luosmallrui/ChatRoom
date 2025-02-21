package main

import (
	"chatroom/config"
	"chatroom/pkg/core"
	"chatroom/rpc"
	"chatroom/socket"
	"github.com/urfave/cli/v2"
	"log"
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
			// 初始化并启动 RPC 服务
			go func() {
				if err := rpc.StartRpcServer(conf); err != nil {
					log.Fatalf("Failed to start RPC server: %v", err)
				}
			}()
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
