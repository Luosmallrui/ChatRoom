package chat

import (
	"chatroom/pkg/core"
	"chatroom/pkg/core/socket"
	"chatroom/pkg/logger"
	"chatroom/types"
	"context"
	"encoding/json"
)

// 键盘输入事件消息
func (h *Handler) onConsumeTalkKeyboard(ctx context.Context, body []byte) {
	var in types.SubEventImMessageKeyboardPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: %s", err.Error())
		return
	}

	ids, _ := h.ClientConnectService.GetUidFromClientIds(ctx, core.GetServerId(), socket.Session.Chat.Name(), in.ToFromId)
	if len(ids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(ids...)
	c.SetMessage(types.PushEventImMessageKeyboard, types.ImMessageKeyboardPayload{
		FromId:   in.ToFromId,
		ToFromId: in.ToFromId,
	})

	socket.Session.Chat.Write(c)
}
