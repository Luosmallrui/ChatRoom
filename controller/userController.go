package controller

import (
	//"chatroom/pkg/core"
	"chatroom/context"
	"chatroom/dao"
	"chatroom/service"
	"chatroom/types"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserController struct {
	Redis        *redis.Client
	UserService  service.IUserService
	UsersRepo    *dao.Users
	OrganizeRepo *dao.Organize
}

func (u *UserController) RegisterRouter(r gin.IRouter) {
	g := r.Group("/api/v1/user")
	g.GET("/list", context.HandlerFunc(u.Detail))
	g.GET("/setting", context.HandlerFunc(u.Setting))
}

// Detail 个人用户信息
func (u *UserController) Detail(ctx *context.Context) error {
	user, err := u.UsersRepo.FindByIdWithCache(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(user)
}

// Setting 用户设置
func (u *UserController) Setting(ctx *context.Context) error {

	//uid := ctx.UserId()
	uid := 3

	user, err := u.UsersRepo.FindByIdWithCache(ctx.Ctx(), uid)
	if err != nil {
		return ctx.Error(err.Error())
	}

	isOk, err := u.OrganizeRepo.IsQiyeMember(ctx.Ctx(), uid)
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(&types.UserSettingResponse{
		UserInfo: &types.UserInfo{
			Uid:      int32(user.Id),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Motto:    user.Motto,
			Gender:   int32(user.Gender),
			IsQiye:   isOk,
			Mobile:   user.Mobile,
			Email:    user.Email,
		},
		Setting: &types.ConfigInfo{},
	})
}
