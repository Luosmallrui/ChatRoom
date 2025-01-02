package dao

import (
	"context"

	"chatroom/model"
	"gorm.io/gorm"
)

type Organize struct {
	Repo[model.Organize]
}

func NewOrganize(db *gorm.DB) *Organize {
	return &Organize{Repo: NewRepo[model.Organize](db)}
}

type UserInfo struct {
	UserId     int    `json:"user_id"`
	Nickname   string `json:"nickname"`
	Gender     int    `json:"gender"`
	Department string `json:"department"`
	Position   string `json:"position"`
}

func (o *Organize) List() ([]*UserInfo, error) {

	tx := o.Repo.Db.Table("organize")
	tx.Select([]string{
		"organize.user_id", "organize.dept_id as department", "organize.position_id as position",
		"users.nickname", "users.gender",
	})
	tx.Joins("left join users on users.id = organize.user_id")

	items := make([]*UserInfo, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// IsQiyeMember 判断是否是企业成员
func (o *Organize) IsQiyeMember(ctx context.Context, uid ...int) (bool, error) {

	count, err := o.Repo.FindCount(ctx, "user_id in ?", uid)
	if err != nil {
		return false, err
	}

	return int(count) == len(uid), nil
}

func (o *Organize) GetMemberIds(ctx context.Context) ([]int64, error) {
	var ids []int64
	if err := o.Repo.Db.WithContext(ctx).Table("organize").Pluck("user_id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}
