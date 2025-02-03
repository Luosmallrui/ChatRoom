package handler

import (
	"chatroom/config"
	"chatroom/socket"
)

type Handler struct {
	Chat        *ChatChannel
	Example     *ExampleChannel
	Config      *config.Config
	RoomStorage *socket.RoomStorage
}
