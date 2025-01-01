package core

import (
	"chatroom/controller"
	"github.com/gin-gonic/gin"
	"os"

	"chatroom/config"
	"github.com/urfave/cli/v2"
)

type App struct {
	app *cli.App
}

type AppProvider struct {
	Config      *config.Config
	Engine      *gin.Engine
	Controllers *controller.Controllers
}

func (app *AppProvider) RegisterRoutes() {
	app.Controllers.RegisterRouters(app.Engine)
}

type Action func(ctx *cli.Context, conf *config.Config) error

type Command struct {
	Name        string
	Usage       string
	Flags       []cli.Flag
	Action      Action
	Subcommands []Command
}

func NewApp(version string) *App {
	return &App{
		app: &cli.App{
			Name:    "ChatRoom",
			Usage:   "在线聊天应用",
			Version: version,
		},
	}
}

func (c *App) Register(fn func() Command) {
	c.app.Commands = append(c.app.Commands, c.command(fn()))
}

func (c *App) command(cm Command) *cli.Command {
	cd := &cli.Command{
		Name:  cm.Name,
		Usage: cm.Usage,
		Flags: make([]cli.Flag, 0),
	}

	if len(cm.Subcommands) > 0 {
		for _, v := range cm.Subcommands {
			cd.Subcommands = append(cd.Subcommands, c.command(v))
		}
	} else {
		if len(cm.Flags) > 0 {
			cd.Flags = append(cd.Flags, cm.Flags...)
		}

		var isConfig bool

		for _, flag := range cd.Flags {
			if flag.Names()[0] == "config" {
				isConfig = true
				break
			}
		}

		if !isConfig {
			cd.Flags = append(cd.Flags, &cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "./config.yaml",
				Usage:       "配置文件路径",
				DefaultText: "./config.yaml",
			})
		}

		if cm.Action != nil {
			cd.Action = func(ctx *cli.Context) error {
				return cm.Action(ctx, config.New(ctx.String("config")))
			}
		}
	}

	return cd
}

func (c *App) Run() {
	_ = c.app.Run(os.Args)
}
