package controller

import (
	"chatroom/config"
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/middleware"
	"chatroom/pkg/filesystem"
	"chatroom/service"
	"chatroom/types"
	"github.com/gin-gonic/gin"
)

type Emoticon struct {
	Session         *cache.JwtTokenStorage
	Config          *config.Config
	RedisLock       *cache.RedisLock
	EmoticonRepo    *dao.Emoticon
	EmoticonService service.IEmoticonService
	Filesystem      filesystem.IFilesystem
}

func (c *Emoticon) RegisterRouter(r gin.IRouter) {
	authorize := middleware.Auth(c.Config.Jwt.Secret, "admin", c.Session)
	r.Use(authorize)
	emoticon := r.Group("/api/v1/emoticon").Use(authorize)
	emoticon.GET("/customize/list", context.HandlerFunc(c.List)) // 表情包列表
}

// List 收藏列表
func (c *Emoticon) List(ctx *context.Context) error {
	resp := &types.EmoticonListResponse{
		Items: make([]*types.EmoticonItem, 0),
	}

	items, err := c.EmoticonRepo.GetCustomizeList(ctx.UserId())
	if err != nil {
		return ctx.Error(err.Error())
	}

	for _, item := range items {
		resp.Items = append(resp.Items, &types.EmoticonItem{
			EmoticonID: int32(item.Id),
			URL:        item.Url,
		})
	}

	return ctx.Success(resp)
}
