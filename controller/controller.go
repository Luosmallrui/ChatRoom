package controller

import (
	"chatroom/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Controller 接口定义所有 controller 必须实现的方法
type Controller interface {
	RegisterRouter(r gin.IRouter)
}

// Controllers 存储所有的 controller
type Controllers struct {
	User *UserController
	Auth *AuthController
	//Chat    *ChatController
	//Room    *RoomController
	// ... 添加其他 controller
}

// RegisterRouters 注册所有路由
func (c *Controllers) RegisterRouters(r gin.IRouter) {
	c.User.RegisterRouter(r)
	c.Auth.RegisterRouter(r)
	//c.Chat.RegisterRouter(r)
	//c.Room.RegisterRouter(r)
	// ... 注册其他 controller 的路由
}

type Deps struct {
	Redis *redis.Client
	//// 还可以添加其他依赖，例如 logger、数据库连接等
	//Logger *log.Logger
	UserService service.IUserService
}
