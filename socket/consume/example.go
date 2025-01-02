package consume

import (
	"chatroom/socket/handler/event/example"
	"context"
)

type ExampleSubscribe struct {
	handler *example.Handler
}

func NewExampleSubscribe(handler *example.Handler) *ExampleSubscribe {
	return &ExampleSubscribe{handler: handler}
}

// Call 触发回调事件
func (s *ExampleSubscribe) Call(event string, data []byte) {
	s.handler.Call(context.TODO(), event, data)
}
