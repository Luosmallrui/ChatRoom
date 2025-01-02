package service

import (
	"chatroom/dao"
)

type PositionService struct {
	*dao.Source
	Repo *dao.Position
}
