package main

import (
	"chatroom/config"
	"chatroom/pkg/core"
	"chatroom/socket"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
)

func NewServerCommand() core.Command {
	return core.Command{
		Name:  "server",
		Usage: "Start both HTTP and WebSocket servers",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			var wg sync.WaitGroup
			errChan := make(chan error, 2)

			// 启动 HTTP 服务
			wg.Add(1)
			go func() {
				defer wg.Done()
				httpApp := NewHttpInjector(conf)
				httpApp.RegisterRoutes()
				if err := core.Run(ctx, httpApp); err != nil {
					errChan <- err
				}
			}()

			// 启动 WebSocket 服务
			wg.Add(1)
			go func() {
				defer wg.Done()
				socketApp := NewSocketInjector(conf)
				if err := socket.Run(ctx, socketApp); err != nil {
					errChan <- err
				}
			}()

			// 等待错误或完成
			go func() {
				wg.Wait()
				close(errChan)
			}()

			// 返回第一个发生的错误（如果有）
			for err := range errChan {
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func main() {
	app := core.NewApp("v1.0.0")
	app.Register(NewServerCommand)
	os.Args = append(os.Args, "server")
	app.Run()
}
