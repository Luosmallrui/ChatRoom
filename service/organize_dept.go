package service

import "chatroom/dao"

type DeptService struct {
	*dao.Source
	Repo *dao.Department
}
