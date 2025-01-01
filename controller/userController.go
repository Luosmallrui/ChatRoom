package controller

import (
	//"chatroom/pkg/core"
	"chatroom/context"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) RegisterRouter(r gin.IRouter) {
	g := r.Group("/user") // middwares.AuthMiddleware()
	g.GET("/list", context.HandlerFunc(u.LoginRouter))
}

func (u *UserController) LoginRouter(ctx *context.Context) error {
	// 实现你的业务逻辑
	return ctx.Success("ok")
}
