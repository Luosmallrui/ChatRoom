package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"chatroom/pkg/core"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

const (
	_MsgEventPing = "ping"
	_MsgEventPong = "pong"
	_MsgEventAck  = "ack"
)

type IClient interface {
	Cid() int64                       // 客户端ID
	Uid() int                         // 客户端关联用户ID
	Close(code int, text string)      // 关闭客户端
	Write(data *ClientResponse) error // 写入数据
	Channel() IChannel                // 获取客户端所属渠道
}

type IStorage interface {
	Bind(ctx context.Context, sid string, channel string, cid int64, uid int) error
	UnBind(ctx context.Context, sid string, channel string, cid int64) error
}

// Client WebSocket 客户端连接信息
type Client struct {
	conn     IConn                // 客户端连接
	cid      int64                // 客户端ID/客户端唯一标识
	uid      int                  // 用户ID
	lastTime int64                // 客户端最后心跳时间/心跳检测
	closed   int32                // 客户端是否关闭连接
	channel  IChannel             // 渠道分组
	storage  IStorage             // 缓存服务
	event    IEvent               // 回调方法
	outChan  chan *ClientResponse // 发送通道
}

type ClientOption struct {
	Uid         int         // 用户识别ID
	Channel     IChannel    // 渠道信息
	Storage     IStorage    // 自定义缓存组件，用于绑定用户与客户端的关系
	IdGenerator IdGenerator // 客户端ID生成器(唯一ID), 默认使用雪花算法
	Buffer      int         // 缓冲区大小根据业务，自行调整
}

type ClientResponse struct {
	IsAck   bool   `json:"-"`                 // 是否需要 ack 回调
	Ackid   string `json:"ackid,omitempty"`   // ACK ID
	Event   string `json:"event"`             // 事件名
	Content any    `json:"payload,omitempty"` // 事件内容
	Retry   int    `json:"-"`                 // 重试次数（0 默认不重试）
}

// NewClient 初始化客户端信息
func NewClient(conn IConn, option *ClientOption, event IEvent) error {
	if option.Buffer <= 0 {
		option.Buffer = 10
	}

	if event == nil {
		panic("event is nil")
	}

	client := &Client{
		conn:     conn,
		uid:      option.Uid,
		lastTime: time.Now().Unix(),
		channel:  option.Channel,
		storage:  option.Storage,
		outChan:  make(chan *ClientResponse, option.Buffer),
		event:    event,
	}

	if option.IdGenerator != nil {
		client.cid = option.IdGenerator.IdGen()
	} else {
		client.cid = defaultIdGenerator.IdGen()
	}

	// 设置客户端连接关闭回调事件
	conn.SetCloseHandler(client.hookClose)

	// 绑定客户端映射关系
	if client.storage != nil {
		err := client.storage.Bind(context.Background(), core.GetServerId(), client.channel.Name(), client.cid, client.uid)
		if err != nil {
			log.Println("[ERROR] bind client err: ", err.Error())
			return err
		}
	}
	//消息 redis publish 留到了Redis上 现在每个服务器 两台 机器 都会去解析推送上来的数据

	// 注册客户端
	client.channel.addClient(client)

	// 触发自定义的 Open 事件
	client.event.Open(client)

	// 注册心跳管理
	health.insert(client)

	return client.init()
}

// Channel  Name
func (c *Client) Channel() IChannel {
	return c.channel
}

// Cid 获取客户端ID
func (c *Client) Cid() int64 {
	return c.cid
}

// Uid 获取客户端关联的用户ID
func (c *Client) Uid() int {
	return c.uid
}

// Close 关闭客户端连接
func (c *Client) Close(code int, message string) {
	defer func() {
		if err := c.conn.Close(); err != nil {
			log.Printf("connection closed failed: %s \n", err.Error())
		}
	}()

	// 触发客户端关闭回调事件
	if err := c.hookClose(code, message); err != nil {
		log.Printf("[%s-%d-%d] client close err: %s \n", c.channel.Name(), c.cid, c.uid, err.Error())
	}
}

func (c *Client) Closed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

// Write 客户端写入数据
func (c *Client) Write(data *ClientResponse) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR] [%s-%d-%d] chan write err: %v \n", c.channel.Name(), c.cid, c.uid, err)
		}
	}()

	if c.Closed() {
		return fmt.Errorf("connection has been closed")
	}

	if data.IsAck && data.Ackid == "" {
		data.Ackid = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	c.outChan <- data //ack的content

	return nil
}

// 循环接收客户端推送信息
func (c *Client) loopAccept() {
	defer c.Close(1000, "loop accept closed")

	for {
		data, err := c.conn.Read()
		if err != nil {
			break
		}

		c.lastTime = time.Now().Unix()

		c.handleMessage(data)
	}
}

// 循环推送客户端信息
func (c *Client) loopWrite() {
	for data := range c.outChan {
		if c.Closed() {
			return
		}

		bt, err := json.Marshal(data)
		if err != nil {
			log.Printf("[ERROR] client json marshal err: %v \n", err)
			break
		}

		if err := c.conn.Write(bt); err != nil {
			log.Printf("[ERROR] [%s-%d-%d] client write err: %v \n", c.channel.Name(), c.cid, c.uid, err)
			return
		}

		if data.IsAck && data.Retry > 0 {
			data.Retry--

			ackBufferContent := &AckBufferContent{}
			ackBufferContent.cid = c.cid
			ackBufferContent.uid = int64(c.uid)
			ackBufferContent.channel = c.channel.Name()
			ackBufferContent.response = data

			ack.insert(data.Ackid, ackBufferContent)
		}
	}
}

// 初始化连接
func (c *Client) init() error {

	// 推送心跳检测配置
	_ = c.Write(&ClientResponse{
		Event: "connect",
		Content: map[string]any{
			"ping_interval": heartbeatInterval,
			"ping_timeout":  heartbeatTimeout,
		}},
	)

	// 启动协程处理推送信息
	go c.loopWrite()

	go c.loopAccept()

	return nil
}

func (c *Client) hookClose(code int, text string) error {

	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}

	close(c.outChan)

	c.event.Close(c, code, text)

	if c.storage != nil {
		err := c.storage.UnBind(context.Background(), core.GetServerId(), c.channel.Name(), c.cid)
		if err != nil {
			log.Println("[ERROR] unbind client err: ", err.Error())
			return err
		}
	}

	health.delete(c)

	c.channel.delClient(c)

	return nil
}

func (c *Client) handleMessage(data []byte) {

	event, err := c.validate(data)
	if err != nil {
		log.Printf("[ERROR] validate err: %s \n", err.Error())
		return
	}

	switch event {
	case _MsgEventPing:
		_ = c.Write(&ClientResponse{Event: _MsgEventPong})
	case _MsgEventPong:
	case _MsgEventAck:
		res := gjson.GetBytes(data, "ackid")
		if res.Exists() && res.String() != "" {
			ack.delete(res.String())
		}

	default: // 触发消息回调
		c.event.Message(c, data)
	}
}

func (c *Client) validate(data []byte) (string, error) {
	if !gjson.ValidBytes(data) {
		return "", fmt.Errorf("invalid json")
	}

	res := gjson.GetBytes(data, "event")
	if !res.Exists() {
		return "", fmt.Errorf("event is empty")
	}

	return res.String(), nil
}
