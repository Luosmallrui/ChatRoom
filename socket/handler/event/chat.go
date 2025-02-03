package event

import (
	"chatroom/dao"
	"chatroom/pkg/business"
	"chatroom/pkg/core/socket"
	"chatroom/pkg/jsonutil"
	"chatroom/service"
	"chatroom/socket/handler/event/chat"
	"chatroom/types"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tidwall/gjson"
)

type ChatEvent struct {
	Redis           *redis.Client
	GroupMemberRepo *dao.GroupMember
	MemberService   service.IGroupMemberService
	Handler         *chat.Handler
	RoomStorage     *socket.RoomStorage
	PushMessage     *business.PushMessage
}

// OnOpen 连接成功回调事件
func (c *ChatEvent) OnOpen(client socket.IClient) {
	ctx := context.TODO()

	now := time.Now()

	// 客户端加入群房间
	for _, groupId := range c.GroupMemberRepo.GetUserGroupIds(ctx, client.Uid()) {
		_ = c.RoomStorage.Insert(int32(groupId), client.Cid(), now.Unix())
	}

	// 推送上线消息
	_ = c.PushMessage.Push(ctx, types.ImTopicChat, &types.SubscribeMessage{
		Event: types.SubEventContactStatus,
		Payload: jsonutil.Encode(types.SubEventContactStatusPayload{
			Status: 1,
			UserId: client.Uid(),
		}),
	})
}

// OnMessage 消息回调事件
func (c *ChatEvent) OnMessage(client socket.IClient, message []byte) {
	res := gjson.GetBytes(message, "event")
	if !res.Exists() {
		return
	}

	// 获取事件名
	event := res.String()
	if event != "" {
		// 触发事件
		c.Handler.Call(context.TODO(), client, event, message)
	}
}

// OnClose 连接关闭回调事件
func (c *ChatEvent) OnClose(client socket.IClient, code int, text string) {
	ctx := context.TODO()

	now := time.Now()

	// 客户端退出群房间
	for _, groupId := range c.GroupMemberRepo.GetUserGroupIds(ctx, client.Uid()) {
		_ = c.RoomStorage.Delete(int32(groupId), client.Cid(), now.Unix())
	}

	// 推送下线消息
	_ = c.PushMessage.Push(ctx, types.ImTopicChat, &types.SubscribeMessage{
		Event: types.SubEventContactStatus,
		Payload: jsonutil.Encode(types.SubEventContactStatusPayload{
			Status: 2,
			UserId: client.Uid(),
		}),
	})
}
