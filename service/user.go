package service

import (
	"chatroom/dao"
	"chatroom/model"
	"context"
	"errors"
	"time"

	"chatroom/pkg/encrypt"
	"gorm.io/gorm"
)

var _ IUserService = (*UserService)(nil)

type IUserService interface {
	Register(ctx context.Context, opt *UserRegisterOpt) (*model.Users, error)
	Login(mobile string, password string) (*model.Users, error)
	Forget(opt *UserForgetOpt) (bool, error)
	UpdatePassword(uid int, oldPassword string, password string) error
}

type UserService struct {
	UsersRepo *dao.Users
}

type UserRegisterOpt struct {
	Nickname string
	Mobile   string
	Password string
	Platform string
}

// Register 注册用户
func (s *UserService) Register(ctx context.Context, opt *UserRegisterOpt) (*model.Users, error) {
	if s.UsersRepo.IsMobileExist(ctx, opt.Mobile) {
		return nil, errors.New("账号已存在! ")
	}

	user := &model.Users{
		Mobile:   opt.Mobile,
		Nickname: opt.Nickname,
		Avatar:   "",
		Gender:   model.UsersGenderDefault,
		//Password:  encrypt.HashPassword(opt.Password),
		Motto:     "",
		Email:     "",
		Birthday:  "",
		IsRobot:   model.UsersGenderDefault,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.UsersRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 登录处理
func (s *UserService) Login(mobile string, password string) (*model.Users, error) {

	user, err := s.UsersRepo.FindByMobile(context.Background(), mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("登录账号不存在! ")
		}

		return nil, err
	}

	if !encrypt.VerifyPassword(user.Password, password) {
		return nil, errors.New("登录密码填写错误! ")
	}

	return user, nil
}

// UserForgetOpt ForgetRequest 账号找回接口验证
type UserForgetOpt struct {
	Mobile   string
	Password string
	SmsCode  string
}

// Forget 账号找回
func (s *UserService) Forget(opt *UserForgetOpt) (bool, error) {

	user, err := s.UsersRepo.FindByMobile(context.Background(), opt.Mobile)
	if err != nil || user.Id == 0 {
		return false, errors.New("账号不存在! ")
	}

	affected, err := s.UsersRepo.UpdateById(context.TODO(), user.Id, map[string]any{
		"password": encrypt.HashPassword(opt.Password),
	})

	return affected > 0, err
}

// UpdatePassword 修改用户密码
func (s *UserService) UpdatePassword(uid int, oldPassword string, password string) error {

	user, err := s.UsersRepo.FindById(context.TODO(), uid)
	if err != nil {
		return errors.New("用户不存在！")
	}

	if !encrypt.VerifyPassword(user.Password, oldPassword) {
		return errors.New("密码验证不正确！")
	}

	_, err = s.UsersRepo.UpdateById(context.TODO(), user.Id, map[string]any{
		"password": encrypt.HashPassword(password),
	})

	return err
}
