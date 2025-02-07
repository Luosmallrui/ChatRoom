package message

import (
	"chatroom/dao/cache"
	"chatroom/model"
	"chatroom/pkg/jsonutil"
	"chatroom/pkg/logger"
	"chatroom/pkg/strutil"
	"chatroom/types"
	"context"
	"time"
)

func (s *Service) CreatePrivateMessage(ctx context.Context, option CreatePrivateMessageOption) error {
	var (
		orgMsgId      = strutil.NewMsgId()                // 生成原始消息ID
		items         = make([]*model.TalkUserMessage, 0) // 准备消息记录数组
		quoteJsonText = "{}"                              // 默认引用消息为空
		now           = time.Now()                        // 当前时间戳
	)
	// 处理引用消息（如果有）：
	if option.QuoteId != "" {
		quoteRecord := &model.TalkUserMessage{}
		if err := s.Db().First(quoteRecord, "msg_id = ?", option.QuoteId).Error; err != nil {
			return err
		}

		user := &model.Users{}
		if err := s.Db().First(user, "id = ?", quoteRecord.FromId).Error; err != nil {
			return err
		}

		queue := &model.Quote{
			QuoteId: option.QuoteId,
			MsgType: 1,
		}

		queue.Nickname = user.Nickname
		queue.Content = s.getTextMessage(quoteRecord.MsgType, quoteRecord.Extra)
		quoteJsonText = jsonutil.Encode(queue)
	}

	items = append(items, &model.TalkUserMessage{
		MsgId:     strutil.NewMsgId(),
		Sequence:  s.Sequence.Get(ctx, option.FromId, true),
		MsgType:   option.MsgType,
		UserId:    option.FromId,
		ToFromId:  option.ToFromId,
		FromId:    option.FromId,
		Extra:     option.Extra,
		Quote:     quoteJsonText,
		OrgMsgId:  orgMsgId,
		SendTime:  now,
		IsRevoked: model.No,
		IsDeleted: model.No,
	})

	// 创建消息记录（双向存储）：
	items = append(items, &model.TalkUserMessage{
		MsgId:     strutil.NewMsgId(),
		Sequence:  s.Sequence.Get(ctx, option.ToFromId, true),
		MsgType:   option.MsgType,
		UserId:    option.ToFromId,
		ToFromId:  option.FromId,
		FromId:    option.FromId,
		Extra:     option.Extra,
		Quote:     quoteJsonText,
		OrgMsgId:  orgMsgId,
		SendTime:  now,
		IsRevoked: model.No,
		IsDeleted: model.No,
	})
	//数据库存发送的消息
	if err := s.Db().WithContext(ctx).Create(items).Error; err != nil {
		return err
	}

	// 消息推送与状态更新：
	pipe := s.Source.Redis().Pipeline()
	for _, item := range items {
		// 1. 构建订阅消息
		content := &types.SubscribeMessage{
			Event: types.SubEventImMessage,
			Payload: jsonutil.Encode(types.SubEventImMessagePayload{
				TalkMode: types.ChatPrivateMode,
				Message:  jsonutil.Encode(item),
			}),
		}

		// 2. 发布消息到Redis
		pipe.Publish(ctx, types.ImTopicChat, jsonutil.Encode(content))

		// 3. 更新未读消息计数（仅接收者）
		if item.UserId != option.FromId {
			s.UnreadStorage.PipeIncr(ctx, pipe, item.UserId, types.ChatPrivateMode, item.ToFromId)
		}

		// 4. 更新最后一条消息缓存
		_ = s.MessageStorage.Set(ctx, types.ChatPrivateMode, item.UserId, item.ToFromId, &cache.LastCacheMessage{
			Content:  s.getTextMessage(item.MsgType, option.Extra),
			Datetime: item.CreatedAt.Format(time.DateTime),
		})
	}

	_, _ = pipe.Exec(ctx)

	return nil
}

func (s *Service) CreateToUserPrivateMessage(ctx context.Context, data *model.TalkUserMessage) error {
	if data.MsgId == "" {
		data.MsgId = strutil.NewMsgId()
	}

	if data.OrgMsgId == "" {
		data.OrgMsgId = data.MsgId
	}

	if data.Sequence <= 0 {
		data.Sequence = s.Sequence.Get(ctx, data.UserId, true)
	}

	if data.Quote == "" {
		data.Quote = "{}"
	}

	if data.SendTime.IsZero() {
		data.SendTime = time.Now()
	}

	data.IsRevoked = model.No
	data.IsDeleted = model.No

	if err := s.Db().WithContext(ctx).Create(data).Error; err != nil {
		return err
	}

	err := s.PushMessage.Push(ctx, types.ImTopicChat, &types.SubscribeMessage{
		Event: types.SubEventImMessage,
		Payload: jsonutil.Encode(types.SubEventImMessagePayload{
			TalkMode: types.ChatPrivateMode,
			Message:  jsonutil.Encode(data),
		}),
	})
	if err != nil {
		logger.Errorf("SendToUserPrivateLetter redis push err:%s", err.Error())
	}

	s.UnreadStorage.Incr(ctx, data.UserId, types.ChatPrivateMode, data.ToFromId)

	// 更新最后一条消息
	_ = s.MessageStorage.Set(ctx, types.ChatPrivateMode, data.UserId, data.ToFromId, &cache.LastCacheMessage{
		Content:  s.getTextMessage(data.MsgType, data.Extra),
		Datetime: data.CreatedAt.Format(time.DateTime),
	})

	return nil
}

func (s *Service) CreatePrivateSysMessage(ctx context.Context, option CreatePrivateSysMessageOption) error {
	return s.CreateToUserPrivateMessage(ctx, &model.TalkUserMessage{
		MsgId:    strutil.NewMsgId(),
		Sequence: s.Sequence.Get(ctx, option.FromId, true),
		MsgType:  types.ChatMsgSysText,
		UserId:   option.FromId,
		ToFromId: option.ToFromId,
		FromId:   0,
		Extra: jsonutil.Encode(model.TalkRecordExtraText{
			Content: option.Content,
		}),
		Quote:    "{}",
		SendTime: time.Now(),
	})
}
