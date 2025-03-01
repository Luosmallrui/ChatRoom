package controller

import (
	"chatroom/config"
	"chatroom/context"
	"chatroom/dao"
	"chatroom/dao/cache"
	pb "chatroom/kitex_gen/connect"
	connect "chatroom/kitex_gen/connect/connectionservice"
	"chatroom/middleware"
	"chatroom/model"
	"chatroom/pkg/encrypt"
	"chatroom/pkg/strutil"
	"chatroom/pkg/timeutil"
	"chatroom/service"
	"chatroom/types"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"log"
	"strings"
)

type Session struct {
	RedisLock            *cache.RedisLock
	Session              *cache.JwtTokenStorage
	Config               *config.Config
	MessageStorage       *cache.MessageStorage
	ClientStorage        *cache.ClientStorage
	UnreadStorage        *cache.UnreadStorage
	ContactRemark        *cache.ContactRemark
	ContactRepo          *dao.Contact
	UsersRepo            *dao.Users
	GroupRepo            *dao.Group
	TalkService          service.ITalkService
	TalkRecordsService   service.ITalkRecordService
	TalkSessionService   service.ITalkSessionService
	UserService          service.IUserService
	GroupService         service.IGroupService
	AuthService          service.IAuthService
	ContactService       service.IContactService
	ClientConnectService service.IClientConnectService
}

func (c *Session) RegisterRouter(r gin.IRouter) {
	authorize := middleware.Auth(c.Config.Jwt.Secret, "admin", c.Session)
	//r.Use(authorize)
	talk := r.Group("/api/v1/talk")
	talk.Use(authorize)
	talk.GET("/list", context.HandlerFunc(c.List))          // 会话列表
	talk.POST("/create", context.HandlerFunc(c.Create))     // 创建会话
	talk.POST("/delete", context.HandlerFunc(c.Delete))     // 删除会话
	talk.POST("/topping", context.HandlerFunc(c.Top))       // 置顶会话
	talk.POST("/disturb", context.HandlerFunc(c.Disturb))   // 会话免打扰
	talk.GET("/records", context.HandlerFunc(c.GetRecords)) // 会话面板记录
	//talk.GET("/history-records", context.HandlerFunc(c.SearchHistoryRecords)) // 历史会话记录
	//talk.GET("/forward-records", context.HandlerFunc(c.GetForwardRecords))    // 会话转发记录
	//talk.GET("/file-download", context.HandlerFunc(c.Download))               // 下载文件
	talk.POST("/clear-unread", context.HandlerFunc(c.ClearUnreadMessage)) // 清除会话未读数

	connect := r.Group("/api/v1/connect")
	connect.GET("/detail", context.HandlerFunc(c.Detail))
}

func (c *Session) Detail(ctx *context.Context) error {
	client, err := connect.NewClient("ConnectionService",
		client.WithHostPorts(fmt.Sprintf(":%d", c.Config.Server.Rpc))) // 服务名与服务端保持一致
	if err != nil {
		return err
	}

	req := &pb.EmptyRequest{}
	resp, err := client.GetConnectionDetail(ctx.Ctx(), req)
	if err != nil {
		log.Fatalf("Failed to call GetConnectionDetail: %v", err)
	}

	// 输出结果
	log.Printf("Connection Detail: Chat=%d, Example=%d, Num=%d", resp.Chat, resp.Example, resp.Num)

	return ctx.Success(map[string]any{
		"detail": resp,
	})
	return nil
}

// GetRecords 获取会话记录
func (c *Session) GetRecords(ctx *context.Context) error {
	in := &types.GetTalkRecordsRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if in.TalkMode == types.ChatGroupMode {
		err := c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
			TalkType: in.TalkMode,
			UserId:   uid,
			ToFromId: in.ToFromId,
		})

		if err != nil {
			items := make([]types.TalkRecord, 0)
			items = append(items, types.TalkRecord{
				MsgId:    strutil.NewMsgId(),
				Sequence: 1,
				MsgType:  types.ChatMsgSysText,
				Extra: model.TalkRecordExtraText{
					Content: "暂无权限查看群消息",
				},
				CreatedAt: timeutil.DateTime(),
			})

			return ctx.Success(map[string]any{
				"cursor": 1,
				"items":  items,
			})
		}
	}

	records, err := c.TalkRecordsService.FindAllTalkRecords(ctx.Ctx(), &service.FindAllTalkRecordsOpt{
		TalkType:   in.TalkMode,
		UserId:     uid,
		ReceiverId: in.ToFromId,
		Cursor:     in.Cursor,
		Limit:      in.Limit,
	})

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	for i, record := range records {
		if record.IsRevoked == model.Yes {
			records[i].Extra = make(map[string]any)
		}
	}

	return ctx.Success(map[string]any{
		"cursor": cursor,
		"items":  records,
	})
}

// Create 创建会话列表
func (c *Session) Create(ctx *context.Context) error {
	var (
		in    = &types.TalkSessionCreateRequest{}
		uid   = ctx.UserId()
		agent = strings.TrimSpace(ctx.Context.GetHeader("user-agent"))
	)

	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if agent != "" {
		agent = encrypt.Md5(agent)
	}

	// 判断对方是否是自己
	if in.TalkMode == types.ChatPrivateMode && int(in.ToFromID) == ctx.UserId() {
		return ctx.ErrorBusiness("创建失败")
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, in.ToFromID, in.TalkMode, agent)
	if !c.RedisLock.Lock(ctx.Ctx(), key, 10) {
		return ctx.ErrorBusiness("创建失败")
	}
	err := c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType: int(in.TalkMode),
		UserId:   uid,
		ToFromId: int(in.ToFromID),
	})

	if err != nil {
		fmt.Println(err)
		return ctx.ErrorBusiness("暂无权限！")
	}

	result, err := c.TalkSessionService.Create(ctx.Ctx(), &service.TalkSessionCreateOpt{
		UserId:     uid,
		TalkType:   int(in.TalkMode),
		ReceiverId: int(in.ToFromID),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	item := &types.TalkSessionItem{
		ID:        int32(result.Id),
		TalkMode:  int32(result.TalkMode),
		ToFromID:  int32(result.ToFromId),
		IsTop:     int32(result.IsTop),
		IsDisturb: int32(result.IsDisturb),
		IsOnline:  model.No,
		IsRobot:   int32(result.IsRobot),
		Name:      "",
		Avatar:    "",
		Remark:    "",
		UnreadNum: 0,
		MsgText:   "",
		UpdatedAt: timeutil.DateTime(),
	}

	if item.TalkMode == types.ChatPrivateMode {
		item.UnreadNum = int32(c.UnreadStorage.Get(ctx.Ctx(), uid, 1, int(in.ToFromID)))

		item.Remark = c.ContactRepo.GetFriendRemark(ctx.Ctx(), uid, int(in.ToFromID))
		if user, err := c.UsersRepo.FindById(ctx.Ctx(), result.ToFromId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else if result.TalkMode == types.ChatGroupMode {
		if group, err := c.GroupRepo.FindById(ctx.Ctx(), int(in.ToFromID)); err == nil {
			item.Name = group.Name
			item.Avatar = group.Avatar
		}
	}

	// 查询缓存消息
	if msg, err := c.MessageStorage.Get(ctx.Ctx(), result.TalkMode, uid, result.ToFromId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	return ctx.Success(&types.TalkSessionCreateResponse{
		ID:        item.ID,
		TalkMode:  item.TalkMode,
		ToFromID:  item.ToFromID,
		IsTop:     item.IsTop,
		IsDisturb: item.IsDisturb,
		IsOnline:  item.IsOnline,
		IsRobot:   item.IsRobot,
		Name:      item.Name,
		Avatar:    item.Avatar,
		Remark:    item.Remark,
		UnreadNum: item.UnreadNum,
		MsgText:   item.MsgText,
		UpdatedAt: item.UpdatedAt,
	})
}

// Delete 删除列表
func (c *Session) Delete(ctx *context.Context) error {
	in := &types.TalkSessionDeleteRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Delete(ctx.Ctx(), ctx.UserId(), int(in.TalkMode), int(in.ToFromID)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&types.TalkSessionDeleteResponse{})
}

// Top 置顶列表
func (c *Session) Top(ctx *context.Context) error {
	in := &types.TalkSessionTopRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Top(ctx.Ctx(), &service.TalkSessionTopOpt{
		UserId:   ctx.UserId(),
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromID),
		Action:   int(in.Action),
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&types.TalkSessionTopResponse{})
}

// Disturb 会话免打扰
func (c *Session) Disturb(ctx *context.Context) error {
	in := &types.TalkSessionDisturbRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Disturb(ctx.Ctx(), &service.TalkSessionDisturbOpt{
		UserId:   ctx.UserId(),
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromID),
		Action:   int(in.Action),
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&types.TalkSessionDisturbResponse{})
}

// List 会话列表
func (c *Session) List(ctx *context.Context) error {
	uid := ctx.UserId()

	data, err := c.TalkSessionService.List(ctx.Ctx(), uid)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	fmt.Println(len(data))
	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkMode == 1 {
			friends = append(friends, item.ToFromId)
		}
	}

	// 获取好友备注
	remarks, _ := c.ContactRepo.Remarks(ctx.Ctx(), uid, friends)

	items := make([]*types.TalkSessionItem, 0)
	for _, item := range data {
		value := &types.TalkSessionItem{
			ID:        int32(item.Id),
			TalkMode:  int32(item.TalkMode),
			ToFromID:  int32(item.ToFromId),
			IsTop:     int32(item.IsTop),
			IsDisturb: int32(item.IsDisturb),
			IsRobot:   int32(item.IsRobot),
			IsOnline:  2,
			Avatar:    item.Avatar,
			MsgText:   "...",
			UpdatedAt: timeutil.FormatDatetime(item.UpdatedAt),
			UnreadNum: int32(c.UnreadStorage.Get(ctx.Ctx(), uid, item.TalkMode, item.ToFromId)),
		}

		if item.TalkMode == 1 {
			isOnline, _ := c.ClientConnectService.IsUidOnline(ctx.Ctx(), types.ImChannelChat, int(value.ToFromID))

			value.Name = item.Nickname
			value.Avatar = item.Avatar
			value.Remark = remarks[item.ToFromId]
			value.IsOnline = lo.Ternary[int32](isOnline, 1, 2)
		} else {
			value.Name = item.GroupName
			value.Avatar = item.GroupAvatar
		}

		// 查询缓存消息
		if msg, err := c.MessageStorage.Get(ctx.Ctx(), item.TalkMode, uid, item.ToFromId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		items = append(items, value)
	}

	return ctx.Success(&types.TalkSessionListResponse{Items: items})
}

func (c *Session) ClearUnreadMessage(ctx *context.Context) error {
	in := &types.TalkSessionClearUnreadNumRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	c.UnreadStorage.Reset(ctx.Ctx(), ctx.UserId(), int(in.TalkMode), int(in.ToFromID))

	return ctx.Success(&types.TalkSessionClearUnreadNumResponse{})
}
