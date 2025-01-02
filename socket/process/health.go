package process

import (
	"context"
	"log"
	"time"

	"chatroom/pkg/core"
	"go-chat/internal/repository/cache"
)

type HealthSubscribe struct {
	storage *cache.ServerStorage
}

func NewHealthSubscribe(storage *cache.ServerStorage) *HealthSubscribe {
	return &HealthSubscribe{storage}
}

func (s *HealthSubscribe) Setup(ctx context.Context) error {

	log.Println("Start HealthSubscribe")

	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			if err := s.storage.Set(ctx, core.GetServerId(), time.Now().Unix()); err != nil {
				//logger.Std().Error(fmt.Sprintf("Websocket HealthSubscribe Report Err: %s", err.Error()))
			}
		}
	}
}
