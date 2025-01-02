package service

import (
	"chatroom/config"
	"chatroom/dao"
	"chatroom/pkg/ipaddress"
	"context"
)

var _ IIpAddressService = (*IpAddressService)(nil)

type IIpAddressService interface {
	FindAddress(ip string) (string, error)
}

type IpAddressService struct {
	*dao.Source
	Config          *config.Config
	IpAddressClient *ipaddress.Client
}

func (i *IpAddressService) FindAddress(ip string) (string, error) {
	if val, err := i.getCache(ip); err == nil {
		return val, nil
	}

	address, err := i.IpAddressClient.GetIpInfo(context.Background(), ip)
	if err != nil {
		return "", err
	}

	_ = i.setCache(ip, address)

	return address, nil
}

func (i *IpAddressService) getCache(ip string) (string, error) {
	return i.Source.Redis().HGet(context.TODO(), "hash:ip_info", ip).Result()
}

func (i *IpAddressService) setCache(ip string, value string) error {
	return i.Source.Redis().HSet(context.TODO(), "hash:ip_info", ip, value).Err()
}
