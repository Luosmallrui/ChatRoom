// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"chatroom/config"
	"chatroom/controller"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/pkg/business"
	"chatroom/pkg/client"
	"chatroom/pkg/core"
	"chatroom/service"
	"chatroom/service/message"
	"github.com/google/wire"
)

// Injectors from wire.go:

func NewHttpInjector(conf *config.Config) *core.AppProvider {
	engine := core.NewGinServer()
	redisClient := client.NewRedisClient(conf)
	db := client.NewMySQLClient(conf)
	users := dao.NewUsers(db, redisClient)
	userService := &service.UserService{
		UsersRepo: users,
	}
	organize := dao.NewOrganize(db)
	user := &controller.User{
		Redis:        redisClient,
		UserService:  userService,
		UsersRepo:    users,
		OrganizeRepo: organize,
	}
	admin := dao.NewAdmin(db)
	jwtTokenStorage := dao.NewTokenSessionStorage(redisClient)
	captchaStorage := dao.NewCaptchaStorage(redisClient)
	captcha := dao.NewBase64Captcha(captchaStorage)
	auth := &controller.Auth{
		Config:          conf,
		AdminRepo:       admin,
		UserRepo:        users,
		JwtTokenStorage: jwtTokenStorage,
		ICaptcha:        captcha,
		UserService:     userService,
	}
	redisLock := cache.NewRedisLock(redisClient)
	cacheJwtTokenStorage := cache.NewTokenSessionStorage(redisClient)
	messageStorage := cache.NewMessageStorage(redisClient)
	serverStorage := cache.NewSidStorage(redisClient)
	clientStorage := cache.NewClientStorage(redisClient, conf, serverStorage)
	unreadStorage := cache.NewUnreadStorage(redisClient)
	contactRemark := cache.NewContactRemark(redisClient)
	relation := cache.NewRelation(redisClient)
	contact := dao.NewContact(db, contactRemark, relation)
	group := dao.NewGroup(db)
	source := dao.NewSource(db, redisClient)
	groupMember := dao.NewGroupMember(db, relation)
	talkService := &service.TalkService{
		Source:          source,
		GroupMemberRepo: groupMember,
	}
	talkSession := dao.NewTalkSession(db)
	talkSessionService := &service.TalkSessionService{
		Source:          source,
		TalkSessionRepo: talkSession,
	}
	sequence := cache.NewSequence(redisClient)
	daoSequence := dao.NewSequence(db, sequence)
	pushMessage := &business.PushMessage{
		Redis: redisClient,
	}
	groupService := &service.GroupService{
		Source:          source,
		GroupRepo:       group,
		GroupMemberRepo: groupMember,
		Relation:        relation,
		Sequence:        daoSequence,
		PushMessage:     pushMessage,
	}
	authService := &service.AuthService{
		OrganizeRepo:    organize,
		ContactRepo:     contact,
		GroupRepo:       group,
		GroupMemberRepo: groupMember,
	}
	contactService := &service.ContactService{
		Source:      source,
		ContactRepo: contact,
	}
	clientConnectService := &service.ClientConnectService{
		Storage: clientStorage,
	}
	session := &controller.Session{
		RedisLock:            redisLock,
		Session:              cacheJwtTokenStorage,
		Config:               conf,
		MessageStorage:       messageStorage,
		ClientStorage:        clientStorage,
		UnreadStorage:        unreadStorage,
		ContactRemark:        contactRemark,
		ContactRepo:          contact,
		UsersRepo:            users,
		GroupRepo:            group,
		TalkService:          talkService,
		TalkSessionService:   talkSessionService,
		UserService:          userService,
		GroupService:         groupService,
		AuthService:          authService,
		ContactService:       contactService,
		ClientConnectService: clientConnectService,
	}
	contactGroup := dao.NewContactGroup(db)
	contactGroupService := &service.ContactGroupService{
		Source:           source,
		ContactGroupRepo: contactGroup,
	}
	contactApplyService := &service.ContactApplyService{
		Source:      source,
		PushMessage: pushMessage,
	}
	fileUpload := dao.NewFileUpload(db)
	vote := cache.NewVote(redisClient)
	groupVote := dao.NewGroupVote(db, vote)
	iFilesystem := config.NewFilesystem(conf)
	robot := dao.NewRobot(db)
	messageService := &message.Service{
		Source:              source,
		GroupMemberRepo:     groupMember,
		SplitUploadRepo:     fileUpload,
		TalkRecordsVoteRepo: groupVote,
		UsersRepo:           users,
		Filesystem:          iFilesystem,
		UnreadStorage:       unreadStorage,
		MessageStorage:      messageStorage,
		ServerStorage:       serverStorage,
		ClientStorage:       clientStorage,
		Sequence:            daoSequence,
		RobotRepo:           robot,
		PushMessage:         pushMessage,
	}
	controllerContact := &controller.Contact{
		Session:             cacheJwtTokenStorage,
		Config:              conf,
		ClientStorage:       clientStorage,
		ContactRepo:         contact,
		UsersRepo:           users,
		OrganizeRepo:        organize,
		ContactGroupRepo:    contactGroup,
		TalkSessionRepo:     talkSession,
		ContactService:      contactService,
		UserService:         userService,
		TalkListService:     talkSessionService,
		ContactGroupService: contactGroupService,
		ContactApplyService: contactApplyService,
		MessageService:      messageService,
	}
	groupMemberService := &service.GroupMemberService{
		Source:          source,
		GroupMemberRepo: groupMember,
	}
	controllerGroup := &controller.Group{
		RedisLock:          redisLock,
		Repo:               source,
		UsersRepo:          users,
		GroupRepo:          group,
		GroupMemberRepo:    groupMember,
		TalkSessionRepo:    talkSession,
		GroupService:       groupService,
		GroupMemberService: groupMemberService,
		TalkSessionService: talkSessionService,
		UserService:        userService,
		ContactService:     contactService,
	}
	controllers := &controller.Controllers{
		User:    user,
		Auth:    auth,
		Session: session,
		Contact: controllerContact,
		Group:   controllerGroup,
	}
	appProvider := &core.AppProvider{
		Config:      conf,
		Engine:      engine,
		Controllers: controllers,
	}
	return appProvider
}

// wire.go:

var providerSet = wire.NewSet(client.NewMySQLClient, client.NewRedisClient, config.NewFilesystem)
