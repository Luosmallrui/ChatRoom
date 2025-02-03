package handler

import (
	"chatroom/config"
	"chatroom/pkg/core/socket"
)

type Handler struct {
	Chat        *ChatChannel
	Example     *ExampleChannel
	Config      *config.Config
	RoomStorage *socket.RoomStorage
}
