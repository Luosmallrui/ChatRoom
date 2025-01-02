package service

import "chatroom/dao"

type OrganizeService struct {
	*dao.Source
	Repo *dao.Organize
}
