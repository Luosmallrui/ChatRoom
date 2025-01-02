package dao

import (
	"chatroom/model"
	"context"

	"gorm.io/gorm"
)

type Department struct {
	Repo[model.OrganizeDept]
}

func NewDepartment(db *gorm.DB) *Department {
	return &Department{Repo: NewRepo[model.OrganizeDept](db)}
}

func (d *Department) List(ctx context.Context) ([]*model.OrganizeDept, error) {
	return d.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("is_deleted = 1").Order("parent_id asc,order_num asc")
	})
}
