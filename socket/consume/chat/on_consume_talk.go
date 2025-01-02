package chat

import (
	"chatroom/model"
	"chatroom/pkg/core"
	"chatroom/pkg/logger"
	"chatroom/socket"
	"chatroom/types"
	"context"
	"encoding/json"
	"fmt"
	"go-chat/internal/service"
	"time"
)

// 聊天消息事件
func (h *Handler) onConsumeTalk(ctx context.Context, body []byte) {
	var in types.SubEventImMessagePayload
	if err := json.Unmarshal(body, &in); err != nil {
		fmt.Println("Err SubEventImMessagePayload===>", err)
		logger.Errorf("[ChatSubscribe] onConsumeTalk Unmarshal err: %s", err.Error())
		return
	}

	if in.TalkMode == types.ChatPrivateMode {
		h.onConsumeTalkPrivateMessage(ctx, in)
	} else if in.TalkMode == types.ChatGroupMode {
		h.onConsumeTalkGroupMessage(ctx, in)
	}
}

// 私有消息(点对点消息)
func (h *Handler) onConsumeTalkPrivateMessage(ctx context.Context, in types.SubEventImMessagePayload) {
	message := model.TalkUserMessage{}
	if err := json.Unmarshal([]byte(in.Message), &message); err != nil {
		fmt.Println("[ChatSubscribe] onConsumeTalkPrivateMessage Unmarshal err:", err.Error())
		return
	}

	// 没在线则不推送
	clientIds, _ := h.ClientConnectService.GetUidFromClientIds(ctx, core.GetServerId(), socket.Session.Chat.Name(), message.UserId)
	if len(clientIds) == 0 {
		return
	}

	var extra any
	if err := json.Unmarshal([]byte(message.Extra), &extra); err != nil {
		return
	}

	var quote any
	if err := json.Unmarshal([]byte(message.Quote), &quote); err != nil {
		return
	}

	body := types.ImMessagePayloadBody{
		MsgId:     message.MsgId,
		Sequence:  int(message.Sequence),
		MsgType:   message.MsgType,
		UserId:    message.FromId,
		Nickname:  "",
		Avatar:    "",
		IsRevoked: model.No,
		SendTime:  message.CreatedAt.Format(time.DateTime),
		Extra:     extra,
		Quote:     quote,
	}

	if body.UserId > 0 {
		user, err := h.UserRepo.FindByIdWithCache(ctx, message.FromId)
		if err != nil {
			return
		}

		body.Nickname = user.Nickname
		body.Avatar = user.Avatar
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetAck(true)

	c.SetMessage(types.PushEventImMessage, types.ImMessagePayload{
		TalkMode: types.ChatPrivateMode,
		ToFromId: message.ToFromId,
		FromId:   message.FromId,
		Body:     body,
	})

	socket.Session.Chat.Write(c)
}

// 群消息
func (h *Handler) onConsumeTalkGroupMessage(ctx context.Context, in types.SubEventImMessagePayload) {
	message := model.TalkGroupMessage{}
	if err := json.Unmarshal([]byte(in.Message), &message); err != nil {
		return
	}

	clientIds := h.RoomStorage.GetClientIDAll(int32(message.GroupId))

	if len(clientIds) == 0 {
		return
	}

	var extra any
	if err := json.Unmarshal([]byte(message.Extra), &extra); err != nil {
		return
	}

	var quote any
	if err := json.Unmarshal([]byte(message.Quote), &quote); err != nil {
		return
	}

	data := service.TalkRecord{
		MsgId:     message.MsgId,
		Sequence:  int(message.Sequence),
		MsgType:   message.MsgType,
		UserId:    message.FromId,
		Nickname:  "",
		Avatar:    "",
		IsRevoked: model.No,
		SendTime:  message.SendTime.Format(time.DateTime),
		Extra:     extra,
		Quote:     quote,
	}

	if data.UserId > 0 {
		user, err := h.UserRepo.FindByIdWithCache(ctx, message.FromId)
		if err != nil {
			return
		}

		data.Nickname = user.Nickname
		data.Avatar = user.Avatar
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetAck(true)

	c.SetMessage(types.PushEventImMessage, types.ImMessagePayload{
		TalkMode: types.ChatGroupMode,
		ToFromId: message.GroupId,
		FromId:   message.FromId,
		Body:     data,
	})

	socket.Session.Chat.Write(c)
}
