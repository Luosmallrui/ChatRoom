package socket

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Session 客户端管理实例
var Session *session

var once sync.Once

// session 渠道客户端
type session struct {
	Chat    *Channel // 默认分组   聊天频道
	Example *Channel // 案例分组  示例频道

	channels map[string]*Channel //频道映射表
	// 可自行注册其它渠道...
}

// 2. 频道查找
func (s *session) Channel(name string) (*Channel, bool) {
	val, ok := s.channels[name]
	return val, ok
}

func Initialize(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	once.Do(func() {
		InitAck()
		initialize(ctx, eg, fn)
	})
}

func initialize(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	Session = &session{
		Chat:     NewChannel("chat", make(chan *SenderContent, 1024*20)),
		Example:  NewChannel("example", make(chan *SenderContent, 100)),
		channels: map[string]*Channel{},
	}

	Session.channels["chat"] = Session.Chat
	Session.channels["example"] = Session.Example

	// 延时启动守护协程

	// 3. 启动各个守护协程
	time.AfterFunc(3*time.Second, func() {

		// 健康检查协程
		eg.Go(func() error {
			defer fn("health exit")
			return health.Start(ctx)
		})

		// ACK确认协程
		eg.Go(func() error {
			defer fn("ack exit")
			return ack.Start(ctx)
		})

		// 聊天频道协程

		eg.Go(func() error {
			defer fn("chat exit")
			return Session.Chat.Start(ctx)
		})

		//eg.Go(func() error {
		//	defer fn("example exit")
		//	return Session.Example.Start(ctx)
		//})
	})
}
