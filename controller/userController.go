package controller

import (
	//"chatroom/pkg/core"
	"chatroom/context"
	"chatroom/dao"
	"chatroom/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserController struct {
	Redis       *redis.Client
	UserService service.IUserService
	UsersRepo   *dao.Users
}

func (u *UserController) RegisterRouter(r gin.IRouter) {
	g := r.Group("/user")
	g.GET("/list", context.HandlerFunc(u.Detail))
}

// Detail 个人用户信息
func (u *UserController) Detail(ctx *context.Context) error {
	user, err := u.UsersRepo.FindByIdWithCache(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(user)
}
