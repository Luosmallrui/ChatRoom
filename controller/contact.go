package controller

import (
	"chatroom/config"
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/middleware"
	"chatroom/model"
	"chatroom/pkg/timeutil"
	"chatroom/service"
	"chatroom/service/message"
	"chatroom/types"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Contact struct {
	Session             *cache.JwtTokenStorage
	Config              *config.Config
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
	ContactApplyService service.IContactApplyService
	MessageService      message.IService
	//Message         message2.IService
}

func (u *Contact) RegisterRouter(r gin.IRouter) {
	authorize := middleware.Auth(u.Config.Jwt.Secret, "admin", u.Session)
	r.Use(authorize)
	c := r.Group("/api/v1/contact")
	c.GET("/list", context.HandlerFunc(u.List))                   // 获取好友列表
	c.GET("/search", context.HandlerFunc(u.Search))               //查找好友
	c.GET("/detail", context.HandlerFunc(u.Detail))               //用户详情信息
	c.POST("/online-status", context.HandlerFunc(u.OnlineStatus)) //联系人在线状态

	c.GET("/group/list", context.HandlerFunc(u.GroupList)) // 联系人分组列表

	// 联系人申请相关
	c.GET("/apply/records", context.HandlerFunc(u.ContactApplyList))  // 联系人申请列表
	c.POST("/apply/create", context.HandlerFunc(u.Create))            // 添加联系人申请
	c.POST("/apply/accept", context.HandlerFunc(u.Accept))            // 同意人申请列表
	c.GET("/apply/unread-num", context.HandlerFunc(u.ApplyUnreadNum)) // 联系人申请未读数
	//c.POST("/apply/decline", context.HandlerFunc(handler.V1.ContactApply.Decline)) // 拒绝人申请列表
}

// ApplyUnreadNum 获取好友申请未读数
func (u *Contact) ApplyUnreadNum(ctx *context.Context) error {
	return ctx.Success(map[string]any{
		"unread_num": u.ContactApplyService.GetApplyUnreadNum(ctx.Ctx(), ctx.UserId()),
	})
}

// ContactApplyList 获取联系人申请列表
func (u *Contact) ContactApplyList(ctx *context.Context) error {

	list, err := u.ContactApplyService.List(ctx.Ctx(), 4)
	if err != nil {
		return ctx.Error(err.Error())
	}
	items := make([]*types.ContactApplyListItem, 0, len(list))
	for _, item := range list {
		items = append(items, &types.ContactApplyListItem{
			ID:        int32(item.Id),
			UserID:    int32(item.UserId),
			FriendID:  int32(item.FriendId),
			Remark:    item.Remark,
			Nickname:  item.Nickname,
			Avatar:    item.Avatar,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}
	u.ContactApplyService.ClearApplyUnreadNum(ctx.Ctx(), ctx.UserId())
	return ctx.Success(&types.ContactApplyListResponse{Items: items})
}

// Accept 同意联系人添加申请
func (u *Contact) Accept(ctx *context.Context) error {
	in := &types.ContactApplyAcceptRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	fmt.Println(uid)
	applyInfo, err := u.ContactApplyService.Accept(ctx.Ctx(), &service.ContactApplyAcceptOpt{
		Remarks: in.Remark,
		ApplyId: int(in.ApplyID),
		UserId:  uid,
	})

	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	_ = u.MessageService.CreatePrivateSysMessage(ctx.Ctx(), message.CreatePrivateSysMessageOption{
		FromId:   uid,
		ToFromId: applyInfo.UserId,
		Content:  "你们已成为好友，可以开始聊天咯！",
	})

	_ = u.MessageService.CreatePrivateSysMessage(ctx.Ctx(), message.CreatePrivateSysMessageOption{
		FromId:   applyInfo.UserId,
		ToFromId: uid,
		Content:  "你们已成为好友，可以开始聊天咯！",
	})

	return ctx.Success(&types.ContactApplyAcceptResponse{})
}

// Create 创建联系人申请
func (u *Contact) Create(ctx *context.Context) error {
	in := &types.ContactApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if u.ContactRepo.IsFriend(ctx.Ctx(), uid, int(in.UserID), false) {
		return ctx.Success(nil)
	}

	if err := u.ContactApplyService.Create(ctx.Ctx(), &service.ContactApplyCreateOpt{
		UserId:   ctx.UserId(),
		Remarks:  in.Remark,
		FriendId: int(in.UserID),
	}); err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&types.ContactApplyCreateResponse{})
}

func (u *Contact) Search(ctx *context.Context) error {
	in := &types.ContactSearchRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := u.UsersRepo.FindByMobile(ctx.Ctx(), in.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.ErrorBusiness("用户不存在！")
		}

		return ctx.ErrorBusiness(err.Error())
	}
	return ctx.Success(&types.ContactSearchResponse{
		UserID:   int32(user.Id),
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
	})
}

func (u *Contact) List(ctx *context.Context) error {
	userId := ctx.UserId()
	fmt.Println(userId)

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

// Detail 联系人详情信息
func (u *Contact) Detail(ctx *context.Context) error {
	in := &types.ContactDetailRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	user, err := u.UsersRepo.FindByIdWithCache(ctx.Ctx(), int(in.UserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.ErrorBusiness("用户不存在！")
		}

		return ctx.ErrorBusiness(err.Error())
	}

	data := types.ContactDetailResponse{
		UserID:   int32(user.Id),
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
		Email:    user.Email,
		FriendInfo: &types.FriendInfo{
			IsFriend: "N",
			GroupID:  0,
			Remark:   "",
		},
	}

	if uid != user.Id {
		contact, err := u.ContactRepo.FindByWhere(ctx.Ctx(), "user_id = ? and friend_id = ?", uid, user.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil && contact.Status == 1 {
			if u.ContactRepo.IsFriend(ctx.Ctx(), uid, user.Id, false) {
				data.FriendInfo.IsFriend = "Y"
				data.FriendInfo.GroupID = int32(contact.GroupId)
				data.FriendInfo.Remark = contact.Remark
			}
		} else {
			isOk, _ := u.OrganizeRepo.IsQiyeMember(ctx.Ctx(), uid, user.Id)
			if isOk {
				data.FriendInfo.IsFriend = "Y"
			}
		}
	}

	return ctx.Success(&data)
}

// OnlineStatus 获取联系人在线状态
func (u *Contact) OnlineStatus(ctx *context.Context) error {
	in := &types.ContactOnlineStatusRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	resp := &types.ContactOnlineStatusResponse{
		OnlineStatus: 1,
	}

	if u.ClientStorage.IsOnline(ctx.Ctx(), types.ImChannelChat, fmt.Sprintf("%d", in.UserID)) {
		resp.OnlineStatus = 2
	}

	return ctx.Success(resp)
}
