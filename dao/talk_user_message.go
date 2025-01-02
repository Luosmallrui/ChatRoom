package dao

import (
	"chatroom/model"
	"context"

	"gorm.io/gorm"
)

type TalkUserMessage struct {
	Repo[model.TalkUserMessage]
}

func NewTalkRecordFriend(db *gorm.DB) *TalkUserMessage {
	return &TalkUserMessage{Repo: NewRepo[model.TalkUserMessage](db)}
}

func (t *TalkUserMessage) FindByMsgId(ctx context.Context, msgId string) (*model.TalkUserMessage, error) {
	return t.FindByWhere(ctx, "msg_id = ?", msgId)
}
