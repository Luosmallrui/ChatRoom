package dao

import (
	"chatroom/model"
	"gorm.io/gorm"
)

type Admin struct {
	Repo[model.Admin]
}

func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{Repo: NewRepo[model.Admin](db)}
}
