package chat

import (
	"chatroom/dao"
	"chatroom/socket"
	"context"
	"log"

	"chatroom/config"
	"chatroom/service"
	"chatroom/types"
)

var handlers map[string]func(ctx context.Context, data []byte)

type Handler struct {
	Config               *config.Config
	OrganizeRepo         *dao.Organize
	UserRepo             *dao.Users
	Source               *dao.Source
	TalkRecordsService   service.ITalkRecordService
	ContactService       service.IContactService
	ClientConnectService service.IClientConnectService
	RoomStorage          *socket.RoomStorage
}

func (h *Handler) init() {
	handlers = make(map[string]func(ctx context.Context, data []byte))

	handlers[types.SubEventImMessage] = h.onConsumeTalk
	handlers[types.SubEventImMessageKeyboard] = h.onConsumeTalkKeyboard
	handlers[types.SubEventImMessageRevoke] = h.onConsumeTalkRevoke
	handlers[types.SubEventContactStatus] = h.onConsumeContactStatus
	handlers[types.SubEventContactApply] = h.onConsumeContactApply
	handlers[types.SubEventGroupJoin] = h.onConsumeGroupJoin
	handlers[types.SubEventGroupApply] = h.onConsumeGroupApply
}

func (h *Handler) Call(ctx context.Context, event string, data []byte) {
	if handlers == nil {
		h.init()
	}

	if call, ok := handlers[event]; ok {
		call(ctx, data)
	} else {
		log.Printf("consume chat event: [%s]未注册回调事件\n", event)
	}
}
