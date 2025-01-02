package dao

import (
	"chatroom/model"
	"context"

	"gorm.io/gorm"
)

type TalkGroupMessage struct {
	Repo[model.TalkGroupMessage]
}

func NewTalkRecordGroup(db *gorm.DB) *TalkGroupMessage {
	return &TalkGroupMessage{Repo: NewRepo[model.TalkGroupMessage](db)}
}

func (t *TalkGroupMessage) FindByMsgId(ctx context.Context, msgId string) (*model.TalkGroupMessage, error) {
	return t.FindByWhere(ctx, "msg_id =?", msgId)
}
