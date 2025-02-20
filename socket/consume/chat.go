package consume

import (
	"context"

	"chatroom/socket/consume/chat"
)

type ChatSubscribe struct {
	handler *chat.Handler
}

func NewChatSubscribe(handel *chat.Handler) *ChatSubscribe {
	return &ChatSubscribe{handler: handel}
}

// Call 触发回调事件
func (s *ChatSubscribe) Call(event string, data []byte) {
	s.handler.Call(context.TODO(), event, data)
}

//send -> kafka 作为一个生产者把消息放到了消息队列里面
// 水 面包 n种商品

//
