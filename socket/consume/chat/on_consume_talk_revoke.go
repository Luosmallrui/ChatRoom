package chat

import (
	"chatroom/model"
	"chatroom/pkg/core"
	"chatroom/pkg/core/socket"
	"chatroom/pkg/logger"
	"chatroom/types"
	"context"
	"encoding/json"
	"fmt"
)

type ConsumeTalkRevoke struct {
	MsgId string `json:"msg_id"`
}

// 撤销聊天消息
func (h *Handler) onConsumeTalkRevoke(ctx context.Context, body []byte) {
	var in ConsumeTalkRevoke
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: %s", err.Error())
		return
	}

	var record model.TalkRecord
	if err := h.Source.Db().First(&record, "msg_id = ?", in.MsgId).Error; err != nil {
		return
	}

	var clientIds []int64
	if record.TalkType == types.ChatPrivateMode {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids, _ := h.ClientConnectService.GetUidFromClientIds(ctx, core.GetServerId(), socket.Session.Chat.Name(), uid)
			clientIds = append(clientIds, ids...)
		}
	} else if record.TalkType == types.ChatGroupMode {
		clientIds = h.RoomStorage.GetClientIDAll(int32(record.ReceiverId))
	}

	if len(clientIds) == 0 {
		return
	}

	var user model.Users
	if err := h.Source.Db().WithContext(ctx).Select("id,nickname").First(&user, record.UserId).Error; err != nil {
		return
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(types.PushEventImMessageRevoke, map[string]any{
		"talk_type":   record.TalkType,
		"sender_id":   record.UserId,
		"receiver_id": record.ReceiverId,
		"msg_id":      record.MsgId,
		"text":        fmt.Sprintf("%s: 撤回了一条消息", user.Nickname),
	})

	socket.Session.Chat.Write(c)
}
