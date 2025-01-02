package dao

import (
	"chatroom/model"
	"context"

	"gorm.io/gorm"
)

type Position struct {
	Repo[model.OrganizePost]
}

func NewPosition(db *gorm.DB) *Position {
	return &Position{Repo: NewRepo[model.OrganizePost](db)}
}

func (p *Position) List(ctx context.Context) ([]*model.OrganizePost, error) {
	return p.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("status = 1").Order("sort asc")
	})
}
