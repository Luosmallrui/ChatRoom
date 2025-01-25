package controller

import (
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/service"
	"chatroom/types"
)

type Contact struct {
	ClientStorage   *cache.ClientStorage
	ContactRepo     *dao.Contact
	UsersRepo       *dao.Users
	OrganizeRepo    *dao.Organize
	TalkSessionRepo *dao.TalkSession
	ContactService  service.IContactService
	UserService     service.IUserService
	TalkListService service.ITalkSessionService
	//Message         message2.IService
}

func (c *Contact) List(ctx *context.Context) error {
	userId := ctx.UserId()
	userId = 3

	list, err := c.ContactService.List(ctx.Ctx(), userId)
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
