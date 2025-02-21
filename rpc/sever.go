package rpc

import (
	"chatroom/config"
	"chatroom/pkg/core/socket"
	"chatroom/rpc/kitex_gen/connect"
	c "chatroom/rpc/kitex_gen/connect/connectionservice"
	"context"
	"fmt"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

// StartRpcServer 启动 Kitex RPC 服务
func StartRpcServer(conf *config.Config) error {
	// 从配置中获取端口号
	listenAddr := fmt.Sprintf(":%d", conf.Server.Rpc) // 动态设置端口

	// 解析字符串地址为 *net.TCPAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	// 创建 Kitex 服务，指定监听地址
	server := c.NewServer(
		new(ConnectionServiceImpl),
		server.WithServiceAddr(tcpAddr), // 使用解析后的地址
	)
	// 启动 Kitex 服务
	if err := server.Run(); err != nil {
		log.Printf("Kitex server failed: %v", err)
		return err
	}
	return nil
}

// ConnectionServiceImpl 是 RPC 服务的实现
type ConnectionServiceImpl struct{}

func (s *ConnectionServiceImpl) GetConnectionDetail(ctx context.Context) (*connect.ConnectionResponse, error) {
	chatCount := socket.Session.Chat.Count()
	exampleCount := socket.Session.Example.Count()
	roomNum := 22 // 替换为 handle.RoomStorage.GetRoomNum()
	return &connect.ConnectionResponse{
		Chat:    int32(chatCount),
		Example: int32(exampleCount),
		Num:     int32(roomNum),
	}, nil
}
