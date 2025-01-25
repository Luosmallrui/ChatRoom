package controller

import (
	"chatroom/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// Controller 接口定义所有 controller 必须实现的方法
type Controller interface {
	RegisterRouter(r gin.IRouter)
}

// Controllers 存储所有的 controller
type Controllers struct {
	User    *UserController
	Auth    *AuthController
	Session *SessionController
	//Chat    *ChatController
	//Room    *RoomController
	// ... 添加其他 controller
}

// RegisterRouters 注册所有路由
func (c *Controllers) RegisterRouters(r gin.IRouter) {
	r.Use(CORSMiddleware())
	c.User.RegisterRouter(r)
	c.Auth.RegisterRouter(r)
	c.Session.RegisterRouter(r)
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置 CORS 头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// 对于 OPTIONS 请求，直接返回 204
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
