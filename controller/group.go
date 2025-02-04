package controller

import (
	"chatroom/config"
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/middleware"
	"chatroom/model"
	"chatroom/pkg/sliceutil"
	"chatroom/pkg/timeutil"
	"chatroom/service"
	"chatroom/types"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"slices"
)

type Group struct {
	Session            *cache.JwtTokenStorage
	Config             *config.Config
	RedisLock          *cache.RedisLock
	Repo               *dao.Source
	UsersRepo          *dao.Users
	GroupRepo          *dao.Group
	GroupMemberRepo    *dao.GroupMember
	GroupApplyStorage  *cache.GroupApplyStorage
	GroupNoticeRepo    *dao.GroupNotice
	TalkSessionRepo    *dao.TalkSession
	GroupService       service.IGroupService
	GroupMemberService service.IGroupMemberService
	TalkSessionService service.ITalkSessionService
	UserService        service.IUserService
	ContactService     service.IContactService
}

func (g *Group) RegisterRouter(r gin.IRouter) {
	authorize := middleware.Auth(g.Config.Jwt.Secret, "admin", g.Session)
	r.Use(authorize)
	c := r.Group("/api/v1/group")
	c.GET("/list", context.HandlerFunc(g.List))                   // 群组列表
	c.GET("/detail", context.HandlerFunc(g.Detail))               // 群组详情
	c.GET("/apply/unread", context.HandlerFunc(g.ApplyUnreadNum)) // 入群申请未读

	c.GET("/invite-list", context.HandlerFunc(g.GetInviteFriends)) // 获取待邀请入群好友列表

	c.POST("/create", context.HandlerFunc(g.Create))      // 创建群组
	c.GET("/member/list", context.HandlerFunc(g.Members)) // 群成员列表

}

// Detail 获取群组信息
func (g *Group) Detail(ctx *context.Context) error {
	in := &types.GroupDetailRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	groupInfo, err := g.GroupRepo.FindById(ctx.Ctx(), int(in.GroupID))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if groupInfo.Id == 0 {
		return ctx.ErrorBusiness("数据不存在")
	}

	resp := &types.GroupDetailResponse{
		GroupID:   int32(groupInfo.Id),
		GroupName: groupInfo.Name,
		Profile:   groupInfo.Profile,
		Avatar:    groupInfo.Avatar,
		CreatedAt: timeutil.FormatDatetime(groupInfo.CreatedAt),
		IsManager: uid == groupInfo.CreatorId,
		IsDisturb: 0,
		IsMute:    int32(groupInfo.IsMute),
		IsOvert:   int32(groupInfo.IsOvert),
		VisitCard: g.GroupMemberRepo.GetMemberRemark(ctx.Ctx(), int(in.GroupID), uid),
		Notice: &types.Notice{
			Content:        "",
			CreatedAt:      "",
			UpdatedAt:      "",
			ModifyUserName: "",
		},
	}

	notice, err := g.GroupNoticeRepo.GetLatestNotice(ctx.Ctx(), int(in.GroupID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if notice != nil {
		resp.Notice = &types.Notice{
			Content:        notice.Content,
			CreatedAt:      timeutil.FormatDatetime(notice.CreatedAt),
			UpdatedAt:      timeutil.FormatDatetime(notice.UpdatedAt),
			ModifyUserName: "admin",
		}
	}

	if g.TalkSessionRepo.IsDisturb(uid, groupInfo.Id, 2) {
		resp.IsDisturb = 1
	}

	return ctx.Success(resp)
}

// Members 获取群成员列表
func (g *Group) Members(ctx *context.Context) error {
	in := &types.GroupMemberListRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	group, err := g.GroupRepo.FindById(ctx.Ctx(), int(in.GroupID))
	if err != nil {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	if group != nil && group.IsDismiss == model.Yes {
		return ctx.Success([]any{})
	}

	if !g.GroupMemberRepo.IsMember(ctx.Ctx(), int(in.GroupID), ctx.UserId(), false) {
		return ctx.ErrorBusiness("非群成员无权查看成员列表！")
	}

	list := g.GroupMemberRepo.GetMembers(ctx.Ctx(), int(in.GroupID))

	items := make([]*types.GroupMemberListResponseItem, 0)
	for _, item := range list {
		items = append(items, &types.GroupMemberListResponseItem{
			UserID:   int32(item.UserId),
			Nickname: item.Nickname,
			Avatar:   item.Avatar,
			Gender:   int32(item.Gender),
			Leader:   int32(item.Leader),
			IsMute:   int32(item.IsMute),
			Remark:   item.UserCard,
			Motto:    item.Motto,
		})
	}

	slices.SortFunc(items, func(a, b *types.GroupMemberListResponseItem) int {
		return int(a.Leader - b.Leader)
	})

	return ctx.Success(&types.GroupMemberListResponse{Items: items})
}

// Create 创建群聊分组
func (g *Group) Create(ctx *context.Context) error {
	in := &types.GroupCreateRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uids := make([]int, 0)
	for _, id := range sliceutil.Unique(in.UserIds) {
		uids = append(uids, int(id))
	}

	if len(uids) < 2 {
		return ctx.InvalidParams("创建群聊失败，至少需要两个用户！")
	}

	if len(uids)+1 > model.GroupMemberMaxNum {
		return ctx.InvalidParams(fmt.Sprintf("群成员数量已达到%d上限！", model.GroupMemberMaxNum))
	}

	gid, err := g.GroupService.Create(ctx.Ctx(), &service.GroupCreateOpt{
		UserId:    ctx.UserId(),
		Name:      in.Name,
		MemberIds: uids,
	})

	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！" + err.Error())
	}

	return ctx.Success(&types.GroupCreateResponse{GroupID: int32(gid)})
}

func (g *Group) GetInviteFriends(ctx *context.Context) error {
	in := &types.GetInviteFriendsRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	items, err := g.ContactService.List(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	data := make([]*types.GetInviteFriendsResponseItem, 0)
	if in.GroupID <= 0 {
		for _, item := range items {
			data = append(data, &types.GetInviteFriendsResponseItem{
				UserID:   int32(item.Id),
				Nickname: item.Nickname,
				Avatar:   item.Avatar,
				Gender:   int32(item.Gender),
				Remark:   item.Remark,
			})
		}

		return ctx.Success(&types.GetInviteFriendsResponse{
			Items: data,
		})
	}

	mids := g.GroupMemberRepo.GetMemberIds(ctx.Ctx(), int(in.GroupID))
	if len(mids) == 0 {
		return ctx.Success(&types.GetInviteFriendsResponse{
			Items: data,
		})
	}

	for i := 0; i < len(items); i++ {
		if !slices.Contains(mids, items[i].Id) {
			data = append(data, &types.GetInviteFriendsResponseItem{
				UserID:   int32(items[i].Id),
				Nickname: items[i].Nickname,
				Avatar:   items[i].Avatar,
				Gender:   int32(items[i].Gender),
				Remark:   items[i].Remark,
			})
		}
	}

	return ctx.Success(&types.GetInviteFriendsResponse{
		Items: data,
	})
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
