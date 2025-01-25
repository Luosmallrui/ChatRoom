package controller

import (
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/model"
	"chatroom/service"
	"chatroom/types"
	"github.com/gin-gonic/gin"
)

type Contact struct {
	ClientStorage       *cache.ClientStorage
	ContactRepo         *dao.Contact
	UsersRepo           *dao.Users
	OrganizeRepo        *dao.Organize
	ContactGroupRepo    *dao.ContactGroup
	TalkSessionRepo     *dao.TalkSession
	ContactService      service.IContactService
	UserService         service.IUserService
	TalkListService     service.ITalkSessionService
	ContactGroupService service.IContactGroupService
	//Message         message2.IService
}

func (u *Contact) RegisterRouter(r gin.IRouter) {
	c := r.Group("/api/v1/contact")
	c.GET("/list", context.HandlerFunc(u.List)) // 获取好友列表

	c.GET("/group/list", context.HandlerFunc(u.GroupList)) // 联系人分组列表
}

func (u *Contact) List(ctx *context.Context) error {
	userId := ctx.UserId()
	userId = 3

	list, err := u.ContactService.List(ctx.Ctx(), userId)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	items := make([]*types.ContactItem, 0, len(list))
	for _, item := range list {
		items = append(items,
			&types.ContactItem{
				UserID:   int32(item.Id),
				Nickname: item.Nickname,
				Gender:   int32(item.Gender),
				Motto:    item.Motto,
				Avatar:   item.Avatar,
				Remark:   item.Remark,
				GroupID:  int32(item.GroupId),
			})
	}

	return ctx.Success(&types.ContactListResponse{Items: items})
}

// GroupList 联系人分组列表
func (u *Contact) GroupList(ctx *context.Context) error {

	uid := ctx.UserId()
	uid = 3

	items := make([]*types.ContactGroupItem, 0)

	count, err := u.ContactRepo.FindCount(ctx.Ctx(), "user_id = ? and status = ?", uid, model.Yes)
	if err != nil {
		return ctx.Error(err.Error())
	}

	items = append(items, &types.ContactGroupItem{
		Name:  "全部",
		Count: int32(count),
	})

	group, err := u.ContactGroupService.GetUserGroup(ctx.Ctx(), uid)
	if err != nil {
		return ctx.Error(err.Error())
	}

	for _, v := range group {
		items = append(items, &types.ContactGroupItem{
			ID:    int32(v.Id),
			Name:  v.Name,
			Count: int32(v.Num),
			Sort:  int32(v.Sort),
		})
	}

	return ctx.Success(&types.ContactGroupListResponse{Items: items})
}
