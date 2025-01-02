package dao

import (
	"chatroom/model"
	"context"

	"gorm.io/gorm"
)

type TalkSession struct {
	Repo[model.TalkSession]
}

func NewTalkSession(db *gorm.DB) *TalkSession {
	return &TalkSession{Repo: NewRepo[model.TalkSession](db)}
}

func (t *TalkSession) IsDisturb(uid int, receiverId int, talkType int) bool {
	resp, err := t.Repo.FindByWhere(context.TODO(), "user_id = ? and to_from_id = ? and talk_mode = ?", uid, receiverId, talkType)
	return err == nil && resp.IsDisturb == 1
}

func (t *TalkSession) FindBySessionId(uid int, receiverId int, talkType int) int {

	resp, err := t.Repo.FindByWhere(context.TODO(), "user_id = ? and to_from_id = ? and talk_mode = ?", uid, receiverId, talkType)
	if err != nil {
		return 0
	}

	return resp.Id
}
