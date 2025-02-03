package controller

import (
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/service"
	"chatroom/types"
	"github.com/gin-gonic/gin"
)

type Group struct {
	RedisLock         *cache.RedisLock
	Repo              *dao.Source
	UsersRepo         *dao.Users
	GroupRepo         *dao.Group
	GroupMemberRepo   *dao.GroupMember
	GroupApplyStorage *cache.GroupApplyStorage
	//GroupNoticeRepo    *dao.GroupNotice
	TalkSessionRepo    *dao.TalkSession
	GroupService       service.IGroupService
	GroupMemberService service.IGroupMemberService
	TalkSessionService service.ITalkSessionService
	UserService        service.IUserService
	ContactService     service.IContactService
	//Message            message.IService
}

func (g *Group) RegisterRouter(r gin.IRouter) {
	c := r.Group("/api/v1/group")
	c.GET("/list", context.HandlerFunc(g.List))                   // 获取好友列表
	c.GET("/apply/unread", context.HandlerFunc(g.ApplyUnreadNum)) // 入群申请未读
}

func (g *Group) ApplyUnreadNum(ctx *context.Context) error {
	return ctx.Success(map[string]any{
		"unread_num": g.GroupApplyStorage.Get(ctx.Ctx(), ctx.UserId()),
	})
}

func (g *Group) List(ctx *context.Context) error {
	items, err := g.GroupService.List(ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	resp := &types.GroupListResponse{}
	resp.Items = make([]*types.GroupItem, 0)

	for _, item := range items {
		resp.Items = append(resp.Items, &types.GroupItem{
			GroupID:   int32(item.Id),
			GroupName: item.GroupName,
			Avatar:    item.Avatar,
			Profile:   item.Profile,
			Leader:    int32(item.Leader),
			CreatorID: int32(item.CreatorId),
		})
	}

	return ctx.Success(resp)
}
