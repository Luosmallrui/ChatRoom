package dao

import (
	"chatroom/model"
	"gorm.io/gorm"
)

type ContactGroup struct {
	Repo[model.ContactGroup]
}

func NewContactGroup(db *gorm.DB) *ContactGroup {
	return &ContactGroup{Repo: NewRepo[model.ContactGroup](db)}
}
