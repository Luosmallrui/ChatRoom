package controller

import (
	"chatroom/config"
	"chatroom/middleware"
	"chatroom/pkg/timeutil"
	"strings"

	//"chatroom/pkg/core"
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/service"
	"chatroom/types"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Redis        *redis.Client
	Session      *cache.JwtTokenStorage
	Config       *config.Config
	UserService  service.IUserService
	UsersRepo    *dao.Users
	OrganizeRepo *dao.Organize
}

func (u *User) RegisterRouter(r gin.IRouter) {
	authorize := middleware.Auth(u.Config.Jwt.Secret, "admin", u.Session)
	g := r.Group("/api/v1/user")
	g.Use(authorize)
	g.GET("/detail", context.HandlerFunc(u.Detail))
	g.GET("/setting", context.HandlerFunc(u.Setting))
	g.POST("/update", context.HandlerFunc(u.ChangeDetail)) // 修改用户信息
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(ctx *context.Context) error {
	in := &types.UserDetailUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if in.Birthday != "" {
		if !timeutil.IsDateFormat(in.Birthday) {
			return ctx.InvalidParams("birthday 格式错误")
		}
	}

	uid := ctx.UserId()
	_, err := u.UsersRepo.UpdateById(ctx.Ctx(), ctx.UserId(), map[string]any{
		"nickname": strings.TrimSpace(strings.Replace(in.Nickname, " ", "", -1)),
		"avatar":   in.Avatar,
		"gender":   in.Gender,
		"motto":    in.Motto,
		"birthday": in.Birthday,
	})

	if err != nil {
		return ctx.ErrorBusiness("个人信息修改失败！")
	}

	_ = u.UsersRepo.ClearTableCache(ctx.Ctx(), uid)

	return ctx.Success(nil, "个人信息修改成功！")
}

// Detail 个人用户信息
func (u *User) Detail(ctx *context.Context) error {
	user, err := u.UsersRepo.FindByIdWithCache(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(user)
}

// Setting 用户设置
func (u *User) Setting(ctx *context.Context) error {

	uid := ctx.UserId()

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
